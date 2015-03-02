package app

import (
	"os"
	"testing"

	"github.com/nylar/miru/config"
	"github.com/stretchr/testify/assert"
)

func TestContext_NewContext(t *testing.T) {
	ctx := NewContext()

	assert.IsType(t, new(Context), ctx)
	assert.Nil(t, ctx.Config)
	assert.Nil(t, ctx.Db)
	assert.Equal(t, len(ctx.Queues.Queues), 0)
}

func TestContext_LoadConfig(t *testing.T) {
	ctx := NewContext()

	err := ctx.LoadConfig("../config.toml")
	assert.NoError(t, err)
}

func TestContext_LoadConfig_NoConfig(t *testing.T) {
	ctx := NewContext()

	err := ctx.LoadConfig(config.DefaultConfig)
	assert.NoError(t, err)
}

func TestContext_LoadConfig_InvalidConfig(t *testing.T) {
	ctx := NewContext()

	old := config.DefaultConfig

	config.DefaultConfig = "[da"

	err := ctx.LoadConfig("xx")
	assert.Error(t, err)

	config.DefaultConfig = old
}

func TestContext_Connect(t *testing.T) {
	ctx := NewContext()
	if err := ctx.LoadConfig(config.DefaultConfig); err != nil {
		t.Fatalf("Could not load config")
	}

	err := ctx.Connect(os.Getenv("RETHINKDB_URL"))
	assert.NoError(t, err)
}

func TestContext_Connect_BadDB(t *testing.T) {
	ctx := NewContext()
	if err := ctx.LoadConfig(config.DefaultConfig); err != nil {
		t.Fatal(err.Error())
	}

	err := ctx.Connect("")
	assert.Error(t, err)
}
