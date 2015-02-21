package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nylar/miru/crawler"
	"github.com/nylar/miru/db"
	"github.com/nylar/miru/search"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func APIRoutes(m *mux.Router, conn *db.Connection) {
	s := m.PathPrefix("/api").Subrouter()

	s.Handle("/crawl", APICrawlHandler(conn)).Methods("GET")
	s.Handle("/search", APISearchHandler(conn)).Methods("GET")
}

func APISearchHandler(conn *db.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		query := r.URL.Query().Get("q")

		if len(query) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{
				Status:  http.StatusBadRequest,
				Message: "Query parameter 'q' was empty.",
			})
			return
		}

		res := search.Results{}
		if err := res.Search(query, conn); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(Response{
				Status:  http.StatusInternalServerError,
				Message: "Search failed.",
			})
			return
		}

		encoder.Encode(res)
	})
}

func APICrawlHandler(conn *db.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		url := r.URL.Query().Get("url")

		if len(url) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{
				Status:  http.StatusBadRequest,
				Message: "URL parameter 'url' was empty.",
			})
			return
		}

		if err := crawler.Crawl(url, conn); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(Response{
				Status:  http.StatusInternalServerError,
				Message: "Crawling failed.",
			})
			return
		}

		encoder.Encode(Response{
			Status:  http.StatusOK,
			Message: "Crawling successful"},
		)
	})
}
