package api

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils_RenderJSON(t *testing.T) {
	json := make(map[string]interface{})
	json["a"] = 1
	json["b"] = 2
	json["c"] = 3

	jsonResponse := `{"a":1,"b":2,"c":3}`

	w := httptest.NewRecorder()
	err := RenderJSON(json, w)

	assert.NoError(t, err)
	assert.Equal(t, w.Body.String(), jsonResponse)
}

func TestUtils_RenderJSON_BadData(t *testing.T) {
	json := make(map[string]interface{})
	json["a"] = make(chan int)

	w := httptest.NewRecorder()
	err := RenderJSON(json, w)

	assert.Error(t, err)
}
