package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nylar/miru/api"
	"github.com/nylar/miru/config"
	"github.com/nylar/miru/db"
)

var (
	conn *db.Connection
	conf *config.Config
)

func main() {
	var err error

	file, err := ioutil.ReadFile("config.toml")
	data := string(file)
	if err != nil {
		log.Print("Could not load the config file. Using the default settings.")
		data = config.DefaultConfig
	}

	if conf, err = config.LoadConfig(data); err != nil {
		log.Fatal("Could not load the config. Exiting.")
	}

	if conn, err = db.NewConnection(conf.Database.Name, conf.Database.Host); err != nil {
		log.Fatalln(
			"Could not connect to the database. Please check your config",
			"settings and try again.")
	}

	db.Database = conf.Database.Name
	db.IndexTable = conf.Tables.Index
	db.DocumentTable = conf.Tables.Document

	r := mux.NewRouter()
	r.StrictSlash(true)

	api.APIRoutes(r)

	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%s", conf.Api.Port), nil)
}
