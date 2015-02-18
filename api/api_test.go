package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Routes(t *testing.T) {
	urls := []string{"/api/"}

	m := mux.NewRouter()
	m.StrictSlash(true)

	APIRoutes(m)

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
