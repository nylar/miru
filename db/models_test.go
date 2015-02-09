package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModels_NewSite(t *testing.T) {
	url := "example.com"
	site := NewSite(url)

	assert.IsType(t, site, new(Site))
	assert.Equal(t, site.SiteID, url)
}

func TestModels_NewDocument(t *testing.T) {
	source := "example.com/about/"
	url := "example.com"
	title := "About Example Inc."
	content := "We make examples and things."
	doc := NewDocument(source, url, title, content)

	assert.IsType(t, doc, new(Document))
	assert.Equal(t, doc.DocID, source)
	assert.Equal(t, doc.SiteID, url)
	assert.Equal(t, doc.Title, title)
	assert.Equal(t, doc.Content, content)
}

func TestModels_NewIndex(t *testing.T) {
	doc := "example.com/about/"
	word := "make"
	var count int64 = 1

	index := NewIndex(doc, word, count)

	assert.IsType(t, index, new(Index))
	assert.Equal(t, index.DocID, doc)
	assert.Equal(t, index.Word, word)
	assert.Equal(t, index.Count, count)
	assert.Equal(t, index.IndexID, "example.com/about/::make")
}
