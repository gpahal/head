/*
Package head provides functions for parsing information in a HTML head tag.

The head.ProcessURL function takes a url string and a *http.Client, makes a
GET request to the url using the http client and returns a new *head.Object
from the returned HTML. If client is nil, a default client is used.

    object, err := head.ProcessURL("http://ogp.me", nil)

The head.ParseHTML function takes an io.Reader, reads HTML, parses the HTML
and returns a new *head.Object.

    resp, _ := http.Get("http://ogp.me")
    // ignoring the error and other response attributes (like status code)
    // for simplicity
    defer resp.Body.Close()

    object, err := head.ParseHTML(resp.Body)

*/
package head
