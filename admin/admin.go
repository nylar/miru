package admin

import (
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/nylar/miru/app"
	"github.com/nylar/miru/crawler"
	"github.com/nylar/miru/queue"
)

func AdminRoutes(m *mux.Router, c *app.Context) {
	s := m.PathPrefix("/admin").Subrouter()

	s.Handle("/queue/{name}", QueueHandler(c)).Methods("GET")
	s.Handle("/queues/", QueuesHandler(c)).Methods("GET")
	s.Handle("/add", NewSiteHandler(c)).Methods("GET")
	s.Handle("/add", AddSiteHandler(c)).Methods("POST")
}

func NewSiteHandler(c *app.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateBox := rice.MustFindBox("templates")

		_new := templateBox.MustString("new.html")

		t := template.New("new")
		tmpl := template.Must(t.Parse(_new))

		tmpl.Execute(w, nil)
	})
}

func AddSiteHandler(c *app.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := queue.NewQueue()
		c.Queues.Add(q)

		url := r.FormValue("url")

		ctx := struct {
			Url        string
			Successful bool
			QueueName  string
		}{
			url,
			true,
			q.Name,
		}

		if err := crawler.Crawl(url, c, q); err != nil {
			ctx.Successful = false
		}

		templateBox := rice.MustFindBox("templates")

		add := templateBox.MustString("add.html")

		t := template.New("add")
		tmpl := template.Must(t.Parse(add))

		tmpl.Execute(w, ctx)
	})
}

func QueuesHandler(c *app.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateBox := rice.MustFindBox("templates")

		_queues := templateBox.MustString("queues.html")

		t := template.New("queues")
		tmpl := template.Must(t.Parse(_queues))

		tmpl.Execute(w, nil)
	})
}

func QueueHandler(c *app.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		q, ok := c.Queues.Queues[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx := struct {
			Name string
		}{
			q.Name,
		}
		templateBox := rice.MustFindBox("templates")

		_queue := templateBox.MustString("queue.html")

		t := template.New("queue")
		tmpl := template.Must(t.Parse(_queue))

		tmpl.Execute(w, ctx)
	})
}
