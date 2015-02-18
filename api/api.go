package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func APIRoutes(m *mux.Router) {
	s := m.PathPrefix("/api").Subrouter()

	s.HandleFunc("/", APIRootHandler)
}

func APIRootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API"))
}
