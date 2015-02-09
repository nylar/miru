package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
