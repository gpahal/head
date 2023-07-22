package validator_test

import (
	"net/url"
	"testing"

	"github.com/gpahal/head/validator"
)

func TestParseURL(t *testing.T) {
	cases := []struct {
		urlString string
		ok        bool
	}{
		{"http://example.com/", true},
		{"https://example.com/", true},
		{"ftp://example.com/", true},
		{"http:/example.com/", false},
	}

	for _, c := range cases {
		parsedURL, ok := validator.ParseURL(c.urlString)
		if ok != c.ok {
			t.Errorf("parseURL %s: expected ok to be %t, got %t", c.urlString, c.ok, ok)
		}
		if !ok {
			continue
		}

		expectedParsedURL := mustParseURL(c.urlString)
		if ok && !urlsEqual(parsedURL, expectedParsedURL) {
			t.Errorf("parseURL %s: expected parsed url to be %+v, got %+v", c.urlString, expectedParsedURL,
				parsedURL)
		}
	}
}

func TestValidateHREF(t *testing.T) {
	cases := []struct {
		urlString string
		ok        bool
	}{
		{"http://example.com/", true},
		{"https://example.com/", true},
		{"ftp://example.com/", false},
		{"http:/example.com/", true},
		{"/example/absolute/path", true},
		{"example/relative/path", true},
	}

	for _, c := range cases {
		ok := validator.ValidateHREF(c.urlString)
		if ok != c.ok {
			t.Errorf("validateHREF %s: expected ok to be %t, got %t", c.urlString, c.ok, ok)
		}
	}
}

func urlsEqual(a, b *url.URL) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil || *a != *b {
		return false
	}

	return true
}

func mustParseURL(rawURL string) *url.URL {
	parsedURL, _ := url.Parse(rawURL)
	return parsedURL
}
