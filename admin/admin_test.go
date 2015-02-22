package admin

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

	db.SetDbUp(_testConn, "admin")
}

func Handler(status int, data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(data)
	}))
}

func TestAdmin_Routes(t *testing.T) {
	urls := []string{"/admin/add"}

	m := mux.NewRouter()
	m.StrictSlash(true)

	AdminRoutes(m, _testConn)

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

func TestAdmin_AddSiteHandler_Post(t *testing.T) {
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
	h := AddSiteHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAdmin_AddSiteHandler_Template(t *testing.T) {
	r, err := http.NewRequest("POST", "/admin/add", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := AddSiteHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAdmin_NewSiteHandler_Template(t *testing.T) {
	r, err := http.NewRequest("GET", "/admin/add", nil)
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()
	h := NewSiteHandler(_testConn)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}
