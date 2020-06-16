package head

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

var htmlStringGeneral = `
<!DOCTYPE html>
<html>
<head>
    <title>title</title>
    <base href="http://example.com/">
    <meta charset="utf-8">
    <link rel="rel1" href="http://example.com/rel1/1" type="text/html" title="title rel1 1">
    <link rel="rel1" href="http://example.com/rel1/2" type="text/html" title="title rel1 2">
    <link rel="rel2" href="http://example.com/rel2" type="text/html" title="title rel2">
    <meta property="property1" content="property1">
    <meta property="property2" content="property2">
    <meta name="name1" content="name1">
    <meta http-equiv="http-equiv1" content="http-equiv1">
    <meta itemprop="itemprop1" content="itemprop1">
</head>
<body>
</body>
</html>
`

func TestParseHTML(t *testing.T) {
	t.Run("T=general", func(t *testing.T) {
		got, err := ParseHTML(strings.NewReader(htmlStringGeneral))
		if err != nil {
			t.Errorf("expected error to be nil in parsing HTML, got %s", err.Error())
			return
		}

		expected := &Object{
			Title:   "title",
			Base:    "http://example.com/",
			Charset: "utf-8",

			Links: map[string][]*Link{
				"rel1": {
					&Link{HREF: "http://example.com/rel1/1", Type: "text/html", Title: "title rel1 1"},
					&Link{HREF: "http://example.com/rel1/2", Type: "text/html", Title: "title rel1 2"},
				},
				"rel2": {
					&Link{HREF: "http://example.com/rel2", Type: "text/html", Title: "title rel2"},
				},
			},
			Metas: map[string]string{
				"property1":   "property1",
				"property2":   "property2",
				"name1":       "name1",
				"http-equiv1": "http-equiv1",
				"itemprop1":   "itemprop1",
			},
		}
		if diff := deep.Equal(expected, got); diff != nil {
			t.Error(diff)
		}
	})
}
