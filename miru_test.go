package miru

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gorilla/mux"

	rdb "github.com/dancannon/gorethink"
)

var (
	_ctx *Context

	_db, _index, _document string

	m = mux.NewRouter().StrictSlash(true)
)

func init() {
	ctx := NewContext()

	if err := ctx.LoadConfig("../config.toml"); err != nil {
		log.Fatalln(err.Error())
	}

	_db = "miru_test"
	_index = ctx.Config.Tables.Index
	_document = ctx.Config.Tables.Document

	ctx.Config.Database.Name = _db
	ctx.Config.Tables.Index = _index
	ctx.Config.Tables.Document = _document

	if err := ctx.Connect(os.Getenv("RETHINKDB_URL")); err != nil {
		log.Fatalln(err.Error())
	}

	_ctx = ctx

	SetUp(_ctx)
}

func Handler(status int, data []byte) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.Write(data)
		}))
}

func SetUp(c *Context) {
	rdb.DbCreate(_db).Exec(c.Db)
	rdb.Db(_db).TableCreate(_document).Exec(c.Db)
	rdb.Db(_db).TableCreate(_index).Exec(c.Db)
	rdb.Db(_db).Table(_index).IndexCreate("word").Exec(c.Db)
}

func TearDown(c *Context) {
	// Clear 'Document' table
	rdb.Db(_db).Table(_document).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(c.Db)

	// Clear 'Index' table
	rdb.Db(_db).Table(_index).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(c.Db)
}
