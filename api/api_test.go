package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
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

	db.SetDbUp(_testConn, "api")
}

func Handler(status int, data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(data)
	}))
}

func TestAPI_Routes(t *testing.T) {
	urls := []string{"/api/", "/api/search/?q=a", "/api/crawl/?url=x"}

	m := mux.NewRouter()
	m.StrictSlash(true)

	APIRoutes(m, _testConn)

	w := httptest.NewRecorder()

	for _, url := range urls {
		r, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err.Error())
		}

		m.ServeHTTP(w, r)

		assert.Equal(t, w.Code, 200)
	}
}

func TestAPI_RootHandler(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	APIRootHandler(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Body.String(), "API")
}

func TestAPI_CrawlHandler(t *testing.T) {
	defer db.TearDbDown(_testConn)

	data := []byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Ghandi</title>
</head>

<body>
    <p>Be the change that you wish to see in the world</p>
</body>
</html>`)
	ts := Handler(200, data)
	defer ts.Close()

	r, err := http.NewRequest("GET", fmt.Sprintf("/api/crawl?url=%s", url.QueryEscape(ts.URL)), nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APICrawlHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"response":"Successfully crawled: %s"}`, ts.URL))

	var response []interface{}
	res, err := rdb.Db(db.Database).Table(db.IndexTable).Run(_testConn.Session)
	if err != nil {
		t.Error(err.Error())
	}

	res.All(&response)
	assert.Equal(t, len(response), 4)
}

func TestAPI_CrawlHandler_BadURL(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/crawl?url=a", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APICrawlHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 500, w.Code)
}

func TestAPI_CrawlHandler_EmptyParameter(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/crawl?url=", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APICrawlHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Body.String(), `{"error":"URL parameter 'url' was empty."}`)
}

func TestAPI_SearchHandler(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/search?q=hello+world", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APISearchHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)

	resp := make(map[string]interface{})
	if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, resp["count"], 0)
	assert.Equal(t, resp["results"], nil)
}

func TestAPI_SearchHandler_EmptyParameter(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/search?q=", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APISearchHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Body.String(), `{"error":"Query parameter 'q' was empty."}`)
}

func TestAPI_SearchHandler_RaisesError(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/search?q=hello+world", nil)
	if err != nil {
		t.Error(err.Error())
	}

	// Remove index
	rdb.Db(db.Database).Table(db.IndexTable).IndexDrop("word").Exec(_testConn.Session)

	w := httptest.NewRecorder()
	h := APISearchHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 500, w.Code)

	// Re-add index
	rdb.Db(db.Database).Table(db.IndexTable).IndexCreate("word").Exec(_testConn.Session)
}
