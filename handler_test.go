package head

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLHandler(t *testing.T) {
	t.Run("T=normal", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/", strings.NewReader("http://url-example.com"))
		w := httptest.NewRecorder()
		h := &URLHandler{Client: &http.Client{Transport: &fakeRoundTripper{body: htmlStringGeneral}}}
		h.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d (body: %s)", http.StatusOK, resp.StatusCode, body)
			return
		}

		obj := &Object{}
		err := json.Unmarshal([]byte(body), obj)
		if err != nil {
			t.Errorf("expected error to be nil in json unmarshal, got %s", err)
			return
		}
		if obj.Title == "" {
			t.Error("expected title not to be empty, got empty")
			return
		}
	})

	t.Run("T=bad-request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/", nil)
		w := httptest.NewRecorder()
		h := &URLHandler{Client: &http.Client{Transport: &fakeRoundTripper{body: htmlStringGeneral}}}
		h.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code to be %d, got %d (body: %s)", http.StatusBadRequest, resp.StatusCode,
				body)
		}
	})
}

func TestHTMLHandler(t *testing.T) {
	htmlString := strings.Replace(htmlStringGeneral, "\n", "\\n", -1)
	htmlString = strings.Replace(htmlString, "\"", "\\\"", -1)

	t.Run("T=normal", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/", strings.NewReader(htmlString))
		w := httptest.NewRecorder()
		h := &HTMLHandler{}
		h.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d (body: %s)", http.StatusOK, resp.StatusCode, body)
			return
		}

		obj := &Object{}
		err := json.Unmarshal([]byte(body), obj)
		if err != nil {
			t.Errorf("expected error to be nil in json unmarshal, got %s", err)
			return
		}
		if obj.Title == "" {
			t.Error("expected title not to be empty, got empty")
			return
		}
	})

	t.Run("T=bad-request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/", strings.NewReader(""))
		w := httptest.NewRecorder()
		h := &HTMLHandler{}
		h.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code to be %d, got %d (body: %s)", http.StatusBadRequest, resp.StatusCode,
				body)
		}
	})
}
