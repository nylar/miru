package miru

import (
	"encoding/json"
	"net/http"
	"sort"

	rdb "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Response provides a generic interface for writing API messages back to the client.
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type queueList struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// QueueList is a sortable interface for keeping queue items in order.
type QueueList []queueList

func (ql QueueList) Len() int { return len(ql) }

func (ql QueueList) Swap(i, j int) { ql[i], ql[j] = ql[j], ql[i] }

func (ql QueueList) Less(i, j int) bool { return ql[i].Name < ql[j].Name }

// APIRoutes configures the routes for the API, cross-origin resource sharing is
// applied to each route then can be reached by external requests.
func APIRoutes(m *mux.Router, c *Context) {
	s := m.PathPrefix("/api").Subrouter()

	_c := cors.New(cors.Options{})

	s.Handle("/queue/{name}", _c.Handler(APIQueueHandler(c))).Methods("GET")
	s.Handle("/queues/", _c.Handler(APIQueuesHandler(c))).Methods("GET")
	s.Handle("/crawl", _c.Handler(APICrawlHandler(c))).Methods("GET")
	s.Handle("/search", _c.Handler(APISearchHandler(c))).Methods("GET")
	s.Handle("/sites", _c.Handler(APISitesHandler(c))).Methods("GET")
}

// APISitesHandler (GET) returns a list of sites.
func APISitesHandler(c *Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)

		type site struct {
			Site string `gorethinkdb:"site" json:"site"`
		}

		sites := []site{}

		res, err := rdb.Db(c.Config.Database.Name).Table(c.Config.Tables.Document).Pluck("site").Distinct().Run(c.Db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(Response{
				Status:  http.StatusInternalServerError,
				Message: "Could not retrieve sites",
			})
			return
		}

		res.All(&sites)

		encoder.Encode(sites)
	})
}

// APISearchHandler (GET) allows one to search the datastore. Accepts one
// parameter: 'q', which is a URL encoded string.
func APISearchHandler(c *Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		query := r.URL.Query().Get("q")

		// No 'q' parameter, quit now instead of wasting a request.
		if len(query) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{
				Status:  http.StatusBadRequest,
				Message: "Query parameter 'q' was empty.",
			})
			return
		}

		res := Results{}
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

// APICrawlHandler (GET) allows one to provide a URL to be crawled. Will recursively
// crawl in the background.
func APICrawlHandler(c *Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		q := NewQueue()
		url := r.URL.Query().Get("url")

		// No 'url' parameter, quit because no URL to crawl.
		if len(url) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{
				Status:  http.StatusBadRequest,
				Message: "URL parameter 'url' was empty.",
			})
			return
		}

		site, _ := RootURL(url)

		q.Name = site
		c.Queues.Add(q)

		if err := Crawl(url, c, q); err != nil {
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

// APIQueuesHandler (GET) returns a list of active queues.
func APIQueuesHandler(c *Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		queues := []queueList{}
		for _, q := range c.Queues.Queues {
			item := queueList{Name: q.Name, Status: q.Status}
			queues = append(queues, item)
		}
		sort.Sort(QueueList(queues))

		encoder := json.NewEncoder(w)
		encoder.Encode(queues)
	})
}

// APIQueueHandler (GET) returns a single queue.
func APIQueueHandler(c *Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		name := mux.Vars(r)["name"]

		type item struct {
			Item string `json:"item"`
			Done bool   `json:"done"`
		}

		type queue struct {
			Name   string `json:"name"`
			Status string `json:"status"`
			Items  []item `json:"items"`
		}

		// Queue not found, return Bad Request.
		_q, ok := c.Queues.Queues[name]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{
				Status:  http.StatusBadRequest,
				Message: "Name provided is not a valid queue.",
			})
			return
		}
		q := new(queue)
		q.Name = _q.Name
		q.Status = _q.Status

		for k := range _q.Manager {
			i := item{Item: k, Done: true}
			for _, _item := range _q.Items {
				if k == _item {
					i.Done = false
				}
			}
			q.Items = append(q.Items, i)
		}

		encoder.Encode(q)
	})
}
