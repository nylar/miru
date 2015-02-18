package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadConfig(t *testing.T) {
	conf, err := LoadConfig(DefaultConfig)

	assert.NoError(t, err)

	assert.Equal(t, conf.Database.Host, "localhost:28015")
	assert.Equal(t, conf.Database.Name, "miru")

	assert.Equal(t, conf.Tables.Index, "indexes")
	assert.Equal(t, conf.Tables.Document, "documents")

	assert.Equal(t, conf.Api.Port, "8036")
}

func TestConfig_LoadConfig_BadData(t *testing.T) {
	_, err := LoadConfig("[da")

	assert.Error(t, err)
}
