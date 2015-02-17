package db

import (
	"log"
	"os"
)

var _testConn *Connection

func init() {
	var err error
	_testConn, err = NewConnection("test", os.Getenv("RETHINKDB_URL"))
	if err != nil {
		log.Fatalln("Could not create a connection for testing. Exiting.")
	}

	SetDbUp(_testConn, "db")
}
