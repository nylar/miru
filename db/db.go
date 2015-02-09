package db

import rdb "github.com/dancannon/gorethink"

var (
	Database      = "miru"
	SiteTable     = "sites"
	DocumentTable = "documents"
	IndexTable    = "indexes"
)

type Connection struct {
	Session *rdb.Session
}

func NewConnection(db, host string) (*Connection, error) {
	conn := new(Connection)
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address:  host,
		Database: db,
	})

	if err != nil {
		return nil, err
	}

	conn.Session = session
	return conn, nil
}
