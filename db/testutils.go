package db

import rdb "github.com/dancannon/gorethink"

var TestConn *Connection

func TearDbDown() {
	// Clear 'Site' table
	rdb.Db(Database).Table(SiteTable).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(TestConn.Session)

	// Clear 'Document' table
	rdb.Db(Database).Table(DocumentTable).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(TestConn.Session)

	// Clear 'Index' table
	rdb.Db(Database).Table(IndexTable).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(TestConn.Session)
}
