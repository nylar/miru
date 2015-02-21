package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nylar/miru/crawler"
	"github.com/nylar/miru/db"
	"github.com/nylar/miru/search"
)

func APIRoutes(m *mux.Router, conn *db.Connection) {
	s := m.PathPrefix("/api").Subrouter()

	s.Handle("/crawl", APICrawlHandler(conn))
	s.Handle("/search", APISearchHandler(conn))
	s.HandleFunc("/", APIRootHandler)
}

func APIRootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API"))
}

func APISearchHandler(conn *db.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if len(query) == 0 {
			resp := make(map[string]interface{})
			resp["error"] = "Query parameter 'q' was empty."
			RenderJSON(resp, w)
			return
		}

		res := search.Results{}
		if err := res.Search(query, conn); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		RenderJSON(res, w)
	})
}

func APICrawlHandler(conn *db.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		resp := make(map[string]interface{})

		if len(url) == 0 {
			resp["error"] = "URL parameter 'url' was empty."
		} else {
			if err := crawler.Crawl(url, conn); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			resp["response"] = fmt.Sprintf("Successfully crawled: %s", url)
		}

		RenderJSON(resp, w)
	})
}
