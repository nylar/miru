package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nylar/miru/app"
	"github.com/nylar/miru/crawler"
	"github.com/nylar/miru/queue"
	"github.com/nylar/miru/search"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func APIRoutes(m *mux.Router, c *app.Context) {
	s := m.PathPrefix("/api").Subrouter()

	s.Handle("/queue/{name}", APIQueueHandler(c)).Methods("GET")
	s.Handle("/queues/", APIQueuesHandler(c)).Methods("GET")
	s.Handle("/crawl", APICrawlHandler(c)).Methods("GET")
	s.Handle("/search", APISearchHandler(c)).Methods("GET")
}

func APISearchHandler(c *app.Context) http.Handler {
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
		if err := res.Search(query, c); err != nil {
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

func APICrawlHandler(c *app.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		q := queue.NewQueue()
		c.Queues.Add(q)
		url := r.URL.Query().Get("url")

		if len(url) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{
				Status:  http.StatusBadRequest,
				Message: "URL parameter 'url' was empty.",
			})
			return
		}

		if err := crawler.Crawl(url, c, q); err != nil {
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

func APIQueuesHandler(c *app.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.Encode(c.Queues)
	})
}

func APIQueueHandler(c *app.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		name := mux.Vars(r)["name"]

		q, ok := c.Queues.Queues[name]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{
				Status:  http.StatusBadRequest,
				Message: "Name provided is not a valid queue.",
			})
			return
		}
		encoder.Encode(q)
	})
}
