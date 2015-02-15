package db

import (
	"log"

	rdb "github.com/dancannon/gorethink"
)

func init() {
	var err error
	TestConn, err = NewConnection("test", "localhost:28015")
	if err != nil {
		log.Fatalln("Could not create a connection for testing. Exiting.")
	}

	Database = "testing"

	rdb.DbCreate(Database).Exec(TestConn.Session)
	rdb.Db(Database).TableCreate(SiteTable).Exec(TestConn.Session)
	rdb.Db(Database).TableCreate(DocumentTable).Exec(TestConn.Session)
	rdb.Db(Database).TableCreate(IndexTable).Exec(TestConn.Session)
	rdb.Db(Database).Table(IndexTable).IndexCreate("word").Exec(TestConn.Session)
}
