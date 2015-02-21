package crawler

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	rdb "github.com/dancannon/gorethink"
	"github.com/nylar/miru/db"
	"github.com/stretchr/testify/assert"
)

var _testConn *db.Connection

func init() {
	var err error
	_testConn, err = db.NewConnection("test", os.Getenv("RETHINKDB_URL"))
	if err != nil {
		log.Fatalln("Could not create a connection for testing. Exiting.")
	}

	db.SetDbUp(_testConn, "crawler")
}

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
	defer db.TearDbDown(_testConn)

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

	err := Crawl(ts.URL, _testConn)
	assert.NoError(t, err)

	var response []interface{}
	res, err := rdb.Db(db.Database).Table(db.IndexTable).Run(_testConn.Session)
	if err != nil {
		t.Error(err.Error())
	}

	res.All(&response)

	assert.Equal(t, len(response), 1)
}

func TestCrawler_Crawl_BadURL(t *testing.T) {
	err := Crawl("", _testConn)
	assert.Error(t, err)
}
