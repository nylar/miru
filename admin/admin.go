package admin

import (
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/nylar/miru/crawler"
	"github.com/nylar/miru/db"
)

func AdminRoutes(m *mux.Router, conn *db.Connection) {
	s := m.PathPrefix("/admin").Subrouter()

	s.Handle("/add", NewSiteHandler(conn)).Methods("GET")
	s.Handle("/add", AddSiteHandler(conn)).Methods("POST")
}

func NewSiteHandler(conn *db.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateBox := rice.MustFindBox("templates")

		_new := templateBox.MustString("new.html")

		t := template.New("new")
		tmpl := template.Must(t.Parse(_new))

		tmpl.Execute(w, nil)
	})
}

func AddSiteHandler(conn *db.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")

		ctx := struct {
			Url        string
			Successful bool
		}{
			url,
			true,
		}

		if err := crawler.Crawl(url, conn); err != nil {
			ctx.Successful = false
		}

		templateBox := rice.MustFindBox("templates")

		add := templateBox.MustString("add.html")

		t := template.New("add")
		tmpl := template.Must(t.Parse(add))

		tmpl.Execute(w, ctx)
	})
}
