package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func Handler(status int, data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(data)
	}))
}

func TestCrawler_getDocument(t *testing.T) {
	data := []byte("Hello, World")
	ts := Handler(200, data)
	defer ts.Close()

	d, err := getDocument(ts.URL)
	assert.Equal(t, d, data)
	assert.NoError(t, err)
}

func TestCrawler_getDocument_BadResponse(t *testing.T) {
	_, err := getDocument("")
	assert.Error(t, err)
}

func TestCrawler_NewDocument(t *testing.T) {
	html := []byte(`<html><body><p>Hello, World!</p></body></html>`)
	doc := NewDocument(html)
	assert.IsType(t, new(goquery.Document), doc)
}

func TestCrawler_NewDocument_StripsUnwantedTags(t *testing.T) {
	html := []byte(`
<!DOCTYPE html>
<html>
<head>
	<title>Hello</title>

	<script type="text/javascript">
	alert("Hello, World");
	</script>

	<style>
	* { font-family: 'Comic Sans' }
	</style>
</head>

<body>
	<p>Hello, World!</p>
</body>
</html>`)

	doc := NewDocument(html)
	js := doc.Find("script")
	css := doc.Find("style")
	p := doc.Find("p")

	// Should be removed and thus be 0 matching nodes.
	assert.Equal(t, js.Length(), 0)
	assert.Equal(t, css.Length(), 0)

	// Everything else should be left as is.
	assert.Equal(t, p.Length(), 1)
}
