# head

[![GoDoc](https://godoc.org/github.com/gpahal/head?status.svg)](https://godoc.org/github.com/gpahal/head)

A go library for parsing information in a HTML head tag.

## Installation

```sh
go get github.com/gpahal/head
```

## Usage

### Processing a URL

The head.ProcessURL function takes a url string and a \*http.Client, makes a
GET request to the url using the http client and returns a new \*head.Object
from the returned HTML. If client is nil, a default client is used.

```go
object, err := head.ProcessURL("http://ogp.me", nil)
```

### Parsing HTML

The head.ParseHTML function takes an io.Reader, reads HTML, parses the HTML
and returns a new \*head.Object.

```go
resp, _ := http.Get("http://ogp.me")
// ignoring the error and other response attributes (like status code)
// for simplicity
defer resp.Body.Close()

object, err := head.ParseHTML(resp.Body)
```

### Documentation

The complete API documentation is available on
[GoDoc](https://godoc.org/github.com/gpahal/head).


## License

Licensed under MIT license ([LICENSE](LICENSE) or [opensource.org/licenses/MIT](https://opensource.org/licenses/MIT))
