package testutils

import (
	"github.com/nylar/miru/app"

	rdb "github.com/dancannon/gorethink"
)

func SetUp(c *app.Context, db, document, index string) {
	rdb.DbCreate(db).Exec(c.Db)
	rdb.Db(db).TableCreate(document).Exec(c.Db)
	rdb.Db(db).TableCreate(index).Exec(c.Db)
	rdb.Db(db).Table(index).IndexCreate("word").Exec(c.Db)
}

func TearDown(c *app.Context, db, document, index string) {
	// Clear 'Document' table
	rdb.Db(db).Table(document).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(c.Db)

	// Clear 'Index' table
	rdb.Db(db).Table(index).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(c.Db)
}
