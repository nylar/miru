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
	"github.com/nylar/miru/app"
	"github.com/nylar/miru/queue"
	"github.com/nylar/miru/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	_ctx *app.Context
	_pkg = "api"

	_db, _index, _document string

	m = mux.NewRouter().StrictSlash(true)
)

func init() {
	ctx := app.NewContext()

	if err := ctx.LoadConfig("../config.toml"); err != nil {
		log.Fatalln(err.Error())
	}

	_db = fmt.Sprintf("%s_%s", ctx.Config.Database.Name, "test")
	_index = fmt.Sprintf("%s_%s", ctx.Config.Tables.Index, _pkg)
	_document = fmt.Sprintf("%s_%s", ctx.Config.Tables.Document, _pkg)

	ctx.Config.Database.Name = _db
	ctx.Config.Tables.Index = _index
	ctx.Config.Tables.Document = _document

	if err := ctx.Connect(os.Getenv("RETHINKDB_URL")); err != nil {
		log.Fatalln(err.Error())
	}

	_ctx = ctx

	testutils.SetUp(_ctx, _db, _document, _index)
}

func Handler(status int, data []byte) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.Write(data)
		}))
}

func TestAPI_CrawlHandler(t *testing.T) {
	defer testutils.TearDown(_ctx, _db, _document, _index)

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

	r, err := http.NewRequest(
		"GET",
		"/api/crawl?url="+url.QueryEscape(ts.URL),
		nil,
	)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APICrawlHandler(_ctx)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(
		t,
		"{\"status\":200,\"message\":\"Crawling successful\"}\n",
		w.Body.String(),
	)

	var response []interface{}
	res, err := rdb.Db(_db).Table(_index).Run(_ctx.Db)
	if err != nil {
		t.Error(err.Error())
	}

	res.All(&response)
	assert.Equal(t, 4, len(response))
}

func TestAPI_CrawlHandler_BadURL(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/crawl?url=a", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APICrawlHandler(_ctx)
	h.ServeHTTP(w, r)

	assert.Equal(t, w.Code, 500)
	assert.Equal(
		t,
		"{\"status\":500,\"message\":\"Crawling failed.\"}\n",
		w.Body.String(),
	)
}

func TestAPI_CrawlHandler_EmptyParameter(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/crawl?url=", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APICrawlHandler(_ctx)
	h.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
	assert.Equal(
		t,
		w.Body.String(),
		"{\"status\":400,\"message\":\"URL parameter 'url' was empty.\"}\n",
	)
}

func TestAPI_SearchHandler(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/search?q=hello+world", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := APISearchHandler(_ctx)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)

	resp := make(map[string]interface{})
	if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 0, resp["count"])
	assert.Nil(t, resp["results"])
}

func TestAPI_SearchHandler_EmptyParameter(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/search?q=", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	APIRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
	assert.Equal(
		t,
		"{\"status\":400,\"message\":\"Query parameter 'q' was empty.\"}\n",
		w.Body.String(),
	)

}

func TestAPI_SearchHandler_RaisesError(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/search?q=hello+world", nil)
	if err != nil {
		t.Error(err.Error())
	}

	// Remove index
	rdb.Db(_db).Table(_index).IndexDrop("word").Exec(_ctx.Db)

	w := httptest.NewRecorder()
	APIRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 500, w.Code)
	assert.Equal(
		t,
		"{\"status\":500,\"message\":\"Search failed.\"}\n",
		w.Body.String(),
	)

	// Re-add index
	rdb.Db(_db).Table(_index).IndexCreate("word").Exec(_ctx.Db)
}

func TestAPI_APIQueuesHandler(t *testing.T) {
	_ctx.Queues = nil
	_ctx.InitQueues()

	r, err := http.NewRequest("GET", "/api/queues/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	APIRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)

	assert.Equal(t, "{\"queues\":{}}\n", w.Body.String())
}

func TestAPI_APIQueueHandler(t *testing.T) {
	_ctx.Queues = nil
	_ctx.InitQueues()

	q := queue.NewQueue()
	q.Name = "1"
	_ctx.Queues.Add(q)

	r, err := http.NewRequest("GET", "/api/queue/"+q.Name, nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	APIRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)

	assert.Equal(
		t,
		"{\"manager\":{},\"items\":null,\"name\":\"1\"}\n",
		w.Body.String(),
	)
}

func TestAPI_APIQueueHandler_InvalidQueue(t *testing.T) {
	_ctx.Queues = nil
	_ctx.InitQueues()

	r, err := http.NewRequest("GET", "/api/queue/1", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	APIRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)

	assert.Equal(
		t,
		"{\"status\":400,\"message\":\"Name provided is not a valid queue.\"}\n",
		w.Body.String(),
	)
}
