package head

import (
	"io"

	"github.com/gpahal/head/validator"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	attrRel   = "rel"
	attrHREF  = "href"
	attrType  = "type"
	attrTitle = "title"

	attrCharset   = "charset"
	attrProperty  = "property"
	attrName      = "name"
	attrHTTPEquiv = "http-equiv"
	attrItemProp  = "itemprop"
)

// Object represents parsed HTML elements inside the <head> tag.
type Object struct {
	Title   string `json:"title"`
	Base    string `json:"base"`
	Charset string `json:"charset"`

	Links map[string][]*Link `json:"links"`
	Metas map[string]string  `json:"metas"`
}

// Link represents information attached to a <link> element.
type Link struct {
	HREF  string
	Type  string
	Title string
}

func newObject() *Object {
	return &Object{
		Links: make(map[string][]*Link),
		Metas: make(map[string]string),
	}
}

// ParseHTML takes an io.Reader, reads HTML, parses the HTML and returns a new
// *Object.
func ParseHTML(buf io.Reader) (*Object, error) {
	obj := newObject()
	isInsideHead := false
	isInsideTitle := false
	z := html.NewTokenizer(buf)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				return obj, nil
			}

			return nil, z.Err()
		}

		isStartTagToken := tt == html.StartTagToken
		if !isInsideHead && !isStartTagToken {
			continue
		}
		if tt == html.CommentToken || tt == html.DoctypeToken {
			continue
		}

		if isInsideTitle {
			isInsideTitle = false
			if tt == html.TextToken {
				titleText := string(z.Text())
				obj.Title = titleText
				continue
			}
		}

		isEndTagToken := tt == html.EndTagToken
		isSelfClosingTagToken := tt == html.SelfClosingTagToken
		if isStartTagToken || isEndTagToken || isSelfClosingTagToken {
			name, hasAttr := z.TagName()
			nameAtom := atom.Lookup(name)
			if !isInsideHead {
				if nameAtom == atom.Head {
					if isStartTagToken {
						isInsideHead = true
					} else {
						return obj, nil
					}
				}

				continue
			}

			if nameAtom == atom.Title && isStartTagToken {
				if isStartTagToken {
					isInsideTitle = true
				} else if isEndTagToken {
					isInsideTitle = false
				}

				continue
			}

			// skip if the current tag doesn't have any attributes or is an end
			// tag token
			if !hasAttr || isEndTagToken {
				continue
			}

			// base tag
			if nameAtom == atom.Base {
				var key, value []byte
				var keyString string
				for hasAttr {
					key, value, hasAttr = z.TagAttr()
					keyString = atom.String(key)
					if keyString == attrHREF {
						if href := string(value); validator.ValidateHREF(href) {
							obj.Base = href
						}
					}
				}
			}

			// link tag
			if nameAtom == atom.Link {
				var key, value []byte
				var keyString, relValue string
				link := &Link{}
				for hasAttr {
					key, value, hasAttr = z.TagAttr()
					keyString = atom.String(key)
					if keyString == attrRel {
						relValue = string(value)
					} else if keyString == attrHREF {
						if href := string(value); validator.ValidateHREF(href) {
							link.HREF = href
						}
					} else if keyString == attrType {
						// TODO: validation
						link.Type = string(value)
					} else if keyString == attrTitle {
						link.Title = string(value)
					}
				}

				if relValue != "" {
					obj.Links[relValue] = append(obj.Links[relValue], link)
				}
			}

			// meta tag
			if nameAtom == atom.Meta {
				var key, value []byte
				var keyString, propertyValue, contentValue string
				var hasCharset bool
				for hasAttr {
					key, value, hasAttr = z.TagAttr()
					keyString = atom.String(key)
					if keyString == attrCharset {
						// TODO: validation
						obj.Charset = string(value)
						hasCharset = true
						break
					} else if keyString == attrProperty ||
						keyString == attrName ||
						keyString == attrHTTPEquiv ||
						keyString == attrItemProp {
						propertyValue = string(value)
					} else if keyString == "content" {
						contentValue = string(value)
					}
				}

				if !hasCharset && propertyValue != "" {
					obj.Metas[propertyValue] = contentValue
				}
			}
		}
	}
}
