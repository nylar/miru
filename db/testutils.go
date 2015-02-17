package db

import (
	"fmt"

	rdb "github.com/dancannon/gorethink"
)

func SetDbUp(conn *Connection, db string) {
	Database = fmt.Sprintf("testing_%s", db)

	rdb.DbCreate(Database).Exec(conn.Session)
	rdb.Db(Database).TableCreate(DocumentTable).Exec(conn.Session)
	rdb.Db(Database).TableCreate(IndexTable).Exec(conn.Session)
	rdb.Db(Database).Table(IndexTable).IndexCreate("word").Exec(conn.Session)
}

func TearDbDown(conn *Connection) {
	// Clear 'Document' table
	rdb.Db(Database).Table(DocumentTable).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(conn.Session)

	// Clear 'Index' table
	rdb.Db(Database).Table(IndexTable).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(conn.Session)
}
