package miru

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

func TestCrawler_newDocument(t *testing.T) {
	html := []byte(`<html><body><p>Hello, World!</p></body></html>`)
	doc := newDocument(html)
	assert.IsType(t, new(goquery.Document), doc)
}

func TestCrawler_newDocument_StripsUnwantedTags(t *testing.T) {
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

	doc := newDocument(html)
	js := doc.Find("script")
	css := doc.Find("style")
	p := doc.Find("p")

	// Should be removed and thus be 0 matching nodes.
	assert.Equal(t, js.Length(), 0)
	assert.Equal(t, css.Length(), 0)

	// Everything else should be left as is.
	assert.Equal(t, p.Length(), 1)
}

func TestCrawler_Crawl(t *testing.T) {
	defer TearDown(_ctx)

	data := []byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Example</title>
</head>

<body>
    <p>Here are some examples</p>
</body>
</html>`)
	ts := Handler(200, data)
	defer ts.Close()

	err := Crawl(ts.URL, _ctx, NewQueue())
	assert.NoError(t, err)

	var response []interface{}
	res, err := rdb.Db(_db).Table(_index).Run(_ctx.Db)
	if err != nil {
		t.Error(err.Error())
	}

	res.All(&response)

	assert.Equal(t, len(response), 1)
}

func TestCrawler_Crawl_BadURL(t *testing.T) {
	err := Crawl("", _ctx, NewQueue())
	assert.Error(t, err)
}

func TestCrawler_RootURL(t *testing.T) {
	root, err := RootURL("http://example.com/about/")
	assert.Equal(t, "example.com", root)
	assert.NoError(t, err)

	root, err = RootURL("%")
	assert.Equal(t, "", root)
	assert.Error(t, err)
}

func TestCrawler_RootURL_BadURL(t *testing.T) {
	err := Crawl("%", _ctx, NewQueue())
	assert.Error(t, err)
}

func TestCrawler_Get(t *testing.T) {
	ts := Handler(200, []byte{})
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	resp, err := Get(req)
	assert.NoError(t, err)
	assert.IsType(t, new(http.Response), resp)
}

func TestCrawler_Get_EmptyRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	resp, err := Get(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCrawler_MustGet_URLNot200(t *testing.T) {
	ts := Handler(500, []byte{})
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	resp, err := MustGet(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCrawler_Contents(t *testing.T) {
	ts := Handler(200, []byte("hello world"))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	resp, err := Get(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	contents := Contents(resp)
	assert.Equal(t, []byte("hello world"), contents)
}

func TestCrawler_Links(t *testing.T) {
	site := "example.org"
	q := NewQueue()

	htmlSoup := []byte(`
<p>
    <a href="http://example.org/1">Link 1</a>
    <br>
    <a href="http://example.org/2">Link 2</a>
</p>`)

	doc := newDocument(htmlSoup)

	Links(doc, q, site)
	assert.Equal(t, 2, len(q.Items))
}

func TestCrawler_Links_InvalidUrls(t *testing.T) {
	site := "example.org"
	q := NewQueue()

	htmlSoup := []byte(`
<p>
    <a href="http://example.com/1">Link 1</a>
    <br>
    <a href="http://example.com/2">Link 2</a>
</p>`)

	doc := newDocument(htmlSoup)

	Links(doc, q, site)
	assert.Equal(t, 0, len(q.Items))
}

func TestCrawler_ProcessPages(t *testing.T) {
	defer TearDown(_ctx)

	ts := Handler(200, []byte(`
<p>
    <a href="foo">Link 1</a>
    <br>
    <a href="bar">Link 2</a>
</p>`))
	defer ts.Close()

	q := NewQueue()
	err := Crawl(ts.URL, _ctx, q)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(q.Manager))
}
