package validator

import (
	"net/url"

	"github.com/asaskevich/govalidator"
)

// ParseURL parses a url string. Parsed URL is returned as the first return
// value. If there is an error while parsing, the second return value is false.
func ParseURL(urlString string) (*url.URL, bool) {
	if !govalidator.IsURL(urlString) {
		return nil, false
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, false
	}

	return parsedURL, true
}

// ValidateHREF validates a href string that should have an http or https scheme
// or be an absolute path.
func ValidateHREF(href string) bool {
	parsedURL, ok := ParseURL(href)
	if !ok || parsedURL.Scheme == "" {
		if href[0] != '/' {
			href = "/" + href
		}
		var err error
		parsedURL, err = url.ParseRequestURI(href)
		if err != nil {
			return false
		}

		return true
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return true
}
