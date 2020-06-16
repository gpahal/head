package head

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type noopCloser struct {
	io.Reader
}

func (noopCloser) Close() error { return nil }

type fakeRoundTripper struct {
	contentType string
	body        string
}

func (frt *fakeRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	contentType := frt.contentType
	if contentType == "" {
		contentType = "text/html"
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     map[string][]string{"Content-Type": {contentType}},
		Body:       noopCloser{strings.NewReader(frt.body)},
	}, nil
}

func TestProcessURL(t *testing.T) {
	t.Run("T=normal", func(t *testing.T) {
		object, err := ProcessURL("http://example.com/",
			&http.Client{Transport: &fakeRoundTripper{body: htmlStringGeneral}})
		if err != nil {
			t.Errorf("expected error to be nil in processing URL: %s", err.Error())
			return
		}

		if object.Title == "" {
			t.Error("expected title not to be empty")
			return
		}
	})

	t.Run("T=content-type:invalid", func(t *testing.T) {
		_, err := ProcessURL("http://example.com/",
			&http.Client{Transport: &fakeRoundTripper{contentType: "application/json", body: "{}"}})
		if err == nil {
			t.Error("expected error to not be nil")
			return
		}
		if _, ok := err.(ContentTypeNotHTMLError); !ok {
			t.Errorf("expected error to be of type ContentTypeNotHTML, got %+v", err)
		}
	})
}
