package head

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// URLHandler is an HTTP handler that responds with a serialized *Object after
// processing a URL's <head> tag. This URL is obtained from the HTTP request.
// It can provided as query parameter or in the HTTP body. GetURL function
// determines how the URL is extracted from an HTTP request.
type URLHandler struct {
	// Client is the *http.Client used to make HTTP requests.
	Client *http.Client

	// GetURL is used to extract the URL from an HTTP request. If GetURL is
	// nil, DefaultGetURL is used.
	GetURL func(r *http.Request) (string, error)

	// WriteResponse is used by URLHandler to produce an HTTP response. It
	// transforms and serializes *Object and writes it to w. If WriteResponse
	// is nil, DefaultWriteResponse is used.
	WriteResponse func(w http.ResponseWriter, obj *Object)
}

// HTMLHandler is an HTTP handler that responds with a serialized *Object after
// processing an HTML string's <head> tag. This HTML string is obtained from
// the HTTP request. It can provided as query parameter or in the HTTP body.
// GetHTML function determines how the HTML string is extracted from an HTTP
// request.
type HTMLHandler struct {
	// GetHTML is used to extract the HTML string from an HTTP request. If
	// GetHTML is nil, DefaultGetHTML is used.
	GetHTML func(r *http.Request) (string, error)

	// WriteResponse is used by HTMLHandler to produce an HTTP response. It
	// transforms and serializes *Object and writes it to w. If WriteResponse
	// is nil, DefaultWriteResponse is used.
	WriteResponse func(w http.ResponseWriter, obj *Object)
}

func (h *URLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil {
		h = &URLHandler{}
	}

	getURL := h.GetURL
	if getURL == nil {
		getURL = DefaultGetURL
	}

	u, err := getURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if u == "" {
		http.Error(w, "url not found in request", http.StatusBadRequest)
		return
	}

	obj, err := ProcessURL(u, h.Client)
	if err != nil {
		if _, ok := err.(*urlParseError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, ok := err.(*httpClientError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse := h.WriteResponse
	if writeResponse == nil {
		writeResponse = DefaultWriteResponse
	}

	writeResponse(w, obj)
}

func (h *HTMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil {
		h = &HTMLHandler{}
	}

	getHTML := h.GetHTML
	if getHTML == nil {
		getHTML = DefaultGetHTML
	}

	htmlStr, err := getHTML(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if htmlStr == "" {
		http.Error(w, "html string not found in request", http.StatusBadRequest)
		return
	}

	obj, err := ParseHTML(strings.NewReader(htmlStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeResponse := h.WriteResponse
	if writeResponse == nil {
		writeResponse = DefaultWriteResponse
	}

	writeResponse(w, obj)
}

// DefaultGetURL is used to extract the URL from an HTTP request. It first
// tries to read the `url` query paramter and then the HTTP request body.
// Whichever is present is used as the URL that needs to be processed.
func DefaultGetURL(r *http.Request) (string, error) {
	u := r.URL.Query().Get("url")
	if u != "" {
		return u, nil
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body.Close()
	return string(bs), nil
}

// DefaultGetHTML is used to extract the HTML from an HTTP request. It simply
// reads the HTTP request body and returns it.
func DefaultGetHTML(r *http.Request) (string, error) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body.Close()
	return string(bs), nil
}

// DefaultWriteResponse serializes *Object as JSON and writes it to w.
func DefaultWriteResponse(w http.ResponseWriter, obj *Object) {
	bs, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
