package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nylar/miru/app"
	"github.com/nylar/miru/queue"
	"github.com/nylar/miru/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	_ctx *app.Context
	_pkg = "admin"

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
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(data)
	}))
}

func TestAdmin_AddSiteHandler_Post(t *testing.T) {
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

	post := map[string]interface{}{
		"url": ts.URL,
	}
	body, err := json.Marshal(post)
	if err != nil {
		t.Error(err.Error())
	}

	r, err := http.NewRequest("POST", "/admin/add", bytes.NewReader(body))
	if err != nil {
		t.Error(err.Error())
	}
	defer r.Body.Close()

	w := httptest.NewRecorder()
	AdminRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAdmin_AddSiteHandler_Template(t *testing.T) {
	r, err := http.NewRequest("POST", "/admin/add", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	AdminRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAdmin_NewSiteHandler_Template(t *testing.T) {
	r, err := http.NewRequest("GET", "/admin/add", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	AdminRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAdmin_QueuesHandler_Template(t *testing.T) {
	_ctx.Queues = nil
	_ctx.InitQueues()

	r, err := http.NewRequest("GET", "/admin/queues/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	AdminRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAdmin_QueueHandler_Template(t *testing.T) {
	_ctx.Queues = nil
	_ctx.InitQueues()

	q := queue.NewQueue()
	q.Name = "1"
	_ctx.Queues.Add(q)

	r, err := http.NewRequest("GET", "/admin/queue/1", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	AdminRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAdmin_QueueHandler_404(t *testing.T) {
	_ctx.Queues = nil
	_ctx.InitQueues()

	r, err := http.NewRequest("GET", "/admin/queue/x", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	AdminRoutes(m, _ctx)
	m.ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
}
