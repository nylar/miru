package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nylar/miru/api"
	"github.com/nylar/miru/app"
)

func main() {
	ctx := app.NewContext()
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

	api.APIRoutes(r, ctx)

	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%s", ctx.Config.Api.Port), nil)
}
