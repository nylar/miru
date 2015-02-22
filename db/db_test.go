package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb_NewConnection(t *testing.T) {
	host := os.Getenv("RETHINKDB_URL")
	db := "test"

	conn, err := NewConnection(db, host)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}

func TestDb_NewConnection_BadConfig(t *testing.T) {
	host := ""
	db := "test"

	conn, err := NewConnection(db, host)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestDb_OpenBolt(t *testing.T) {
	c := new(Connection)

	if _, err := os.Create("test.db"); err != nil {
		t.Error(err.Error())
	}

	err := c.OpenBolt("test.db")
	assert.NoError(t, err)

	c.CloseBolt()
	os.Remove("test.db")
}

func TestDB_OpenBolt_NonExistantDB(t *testing.T) {
	c := new(Connection)

	err := c.OpenBolt("")
	assert.Error(t, err)
}
