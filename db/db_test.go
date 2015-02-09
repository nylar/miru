package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb_NewConnection(t *testing.T) {
	host := "localhost:28015"
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
