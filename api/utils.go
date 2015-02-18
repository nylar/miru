package api

import (
	"encoding/json"
	"net/http"
)

func RenderJSON(data map[string]interface{}, w http.ResponseWriter) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(j)

	return nil
}
