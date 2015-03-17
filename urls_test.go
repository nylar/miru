package miru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrls_ProcessURL(t *testing.T) {
	site := "example.org"

	link, err := ProcessURL("http://example.org/about/", site)
	assert.NoError(t, err)
	assert.Equal(t, "http://example.org/about/", link)

	link, err = ProcessURL("#about", site)
	assert.Error(t, err)
	assert.Equal(t, "", link)

	link, err = ProcessURL("./about/", site)
	assert.NoError(t, err)
	assert.Equal(t, "http://example.org/about/", link)

	link, err = ProcessURL("http://www.google.com/a%20b?q=c+d", site)
	assert.Error(t, err)
	assert.Equal(t, "", link)
}
