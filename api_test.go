package miru

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

func TestAPI_CrawlHandler(t *testing.T) {
	defer TearDown(_ctx)

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

	q := NewQueue()
	q.Name = "example.com"
	_ctx.Queues.Add(q)

	r, err := http.NewRequest("GET", "/api/queues/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	APIRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)

	assert.Equal(
		t,
		"[{\"name\":\"example.com\",\"status\":\"active\"}]\n",
		w.Body.String(),
	)
}

func TestAPI_APIQueueHandler(t *testing.T) {
	_ctx.Queues = nil
	_ctx.InitQueues()

	q := NewQueue()
	q.Name = "1"
	_ctx.Queues.Add(q)

	q.Enqueue("http://1.com/contact/")
	q.Enqueue("http://1.com/about/")

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
		"{\"name\":\"1\",\"status\":\"active\",\"items\":[{\"item\""+
			":\"http://1.com/contact/\",\"done\":false},"+
			"{\"item\":\"http://1.com/about/\",\"done\":false}]}\n",
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

func TestAPI_QueueList_Sort(t *testing.T) {
	ql := QueueList{}
	ql = append(ql, queueList{Name: "b", Status: "active"})
	ql = append(ql, queueList{Name: "a", Status: "paused"})

	assert.Equal(t, "b", ql[0].Name)
	assert.Equal(t, "a", ql[1].Name)

	sort.Sort(QueueList(ql))

	assert.Equal(t, "a", ql[0].Name)
	assert.Equal(t, "b", ql[1].Name)
}
