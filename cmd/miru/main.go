package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nylar/miru"
)

func main() {
	ctx := miru.NewContext()
	ctx.InitQueues()

	if err := ctx.LoadConfig("config.toml"); err != nil {
		log.Fatalln("Could not load config file.")
		return
	}

	if err := ctx.Connect(ctx.Config.Database.Host); err != nil {
		log.Fatalln("Could not connect to the database.")
		return
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	miru.APIRoutes(r, ctx)

	http.Handle("/", r)
	log.Println(fmt.Sprintf("Serving on http://localhost:%s", ctx.Config.Api.Port))
	http.ListenAndServe(fmt.Sprintf(":%s", ctx.Config.Api.Port), nil)
}
