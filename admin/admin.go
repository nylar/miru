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

	s.Handle("/", AdminAddSiteHandler(conn))
}

func AdminAddSiteHandler(conn *db.Connection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := struct {
			Url        string
			Successful bool
			Posted     bool
		}{
			"",
			true,
			false,
		}

		if r.Method == "POST" {
			url := r.PostFormValue("url")
			ctx.Url = url
			ctx.Posted = true

			if err := crawler.Crawl(url, conn); err != nil {
				ctx.Successful = false
			}
		}

		templateBox, _ := rice.FindBox("templates")

		index, _ := templateBox.String("index.html")

		t := template.New("index")
		tmpl := template.Must(t.Parse(index))

		tmpl.Execute(w, ctx)
	})
}
