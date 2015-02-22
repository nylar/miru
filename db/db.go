package db

import (
	"github.com/boltdb/bolt"
	rdb "github.com/dancannon/gorethink"
)

var (
	Database      = "miru"
	DocumentTable = "documents"
	IndexTable    = "indexes"
)

type Connection struct {
	Session *rdb.Session
	Bolt    *bolt.DB
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

func (c *Connection) OpenBolt(db string) error {
	b, err := bolt.Open(db, 0666, nil)
	if err != nil {
		return err
	}

	c.Bolt = b
	return nil
}

func (c *Connection) CloseBolt() {
	c.Bolt.Close()
}
