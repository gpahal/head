package head

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"time"
)

// InvalidStatusCodeError is returned when the HTTP response status code is
// < 200 or >= 300. This is the final HTTP response after the http.Client has
// followed the redirect policy.
type InvalidStatusCodeError int

func (e InvalidStatusCodeError) Error() string {
	return fmt.Sprintf("invalid status code: %d", int(e))
}

// InvalidContentTypeError is returned when the Content-Type HTTP response
// header is invalid and cannot be parsed.
type InvalidContentTypeError string

func (e InvalidContentTypeError) Error() string {
	return fmt.Sprintf("invalid content type: %s", string(e))
}

// ContentTypeNotHTMLError is returned when the Content-Type HTTP response
// header indicates that the content is not HTML.
type ContentTypeNotHTMLError string

func (e ContentTypeNotHTMLError) Error() string {
	return fmt.Sprintf(
		"content type not html: %s (should be \"text/html\" or \"application/xhtml+xml\" or empty)", string(e))
}

type urlParseError struct {
	cause error
}

func (e *urlParseError) Error() string { return e.cause.Error() }

type httpClientError struct {
	cause error
}

func (e *httpClientError) Error() string { return e.cause.Error() }

// ProcessURL makes a GET request to the url using the http client and returns
// a new *Object from the returned HTML. If client is nil, a default
// client is used.
func ProcessURL(urlString string, client *http.Client) (*Object, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, &urlParseError{cause: err}
	}

	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	header := make(http.Header)
	header.Add("User-Agent", "head (github.com/gpahal/head)")
	resp, err := client.Do(&http.Request{
		Method:     "GET",
		URL:        parsedURL,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     header,
	})
	if err != nil {
		return nil, &httpClientError{cause: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, InvalidStatusCodeError(resp.StatusCode)
	}

	// TODO: check for content encoding as the html package supports only utf-8

	contentTypeValue := resp.Header.Get("Content-Type")
	if contentTypeValue != "" {
		mediaType, _, err := mime.ParseMediaType(contentTypeValue)
		if err != nil {
			return nil, InvalidContentTypeError(contentTypeValue)
		}
		if mediaType != "text/html" && mediaType != "application/xhtml+xml" {
			return nil, ContentTypeNotHTMLError(mediaType)
		}
	}

	return ParseHTML(resp.Body)
}
