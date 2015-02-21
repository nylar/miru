package db

import (
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

func TestModels_NewDocument(t *testing.T) {
	source := "example.com/about/"
	url := "example.com"
	title := "About Example Inc."
	content := "We make examples and things."
	doc := NewDocument(source, url, title, content)

	assert.IsType(t, doc, new(Document))
	assert.NotEqual(t, doc.DocID, "")
	assert.Equal(t, doc.Site, url)
	assert.Equal(t, doc.Title, title)
	assert.Equal(t, doc.Content, content)
}

func TestModels_DocumentPut(t *testing.T) {
	defer TearDbDown(_testConn)

	doc := NewDocument("example.com/about/", "example.com", "About Example Inc.", "We make examples and things.")

	err := doc.Put(_testConn)
	assert.NoError(t, err)

	res, err := rdb.Db(Database).Table(DocumentTable).Get(doc.DocID).Run(_testConn.Session)
	assert.NoError(t, err)

	var d Document
	err = res.One(&d)
	assert.NoError(t, err)

	assert.NotEqual(t, d.DocID, "")
}

func TestModels_DocumentPut_Duplicate(t *testing.T) {
	defer TearDbDown(_testConn)

	doc := NewDocument("example.com/about/", "example.com", "", "")
	doc2 := NewDocument("example.com/about/", "example.com", "", "")

	doc.DocID = "1"
	doc2.DocID = "1"

	err := doc.Put(_testConn)
	assert.NoError(t, err)

	err = doc2.Put(_testConn)
	assert.Error(t, err)
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
	assert.NotEqual(t, index.IndexID, "")
}

func TestModels_IndexPut(t *testing.T) {
	defer TearDbDown(_testConn)

	index := NewIndex("example.com/about/", "make", 52)

	err := index.Put(_testConn)
	assert.NoError(t, err)

	res, err := rdb.Db(Database).Table(IndexTable).Get(index.IndexID).Run(_testConn.Session)
	assert.NoError(t, err)

	var i Index
	err = res.One(&i)
	assert.NoError(t, err)

	assert.Equal(t, i.Word, "make")
}

func TestModels_IndexPut_Duplicate(t *testing.T) {
	defer TearDbDown(_testConn)

	index := NewIndex("ZXhhbXBsZS5jb20vYWJvdXQv", "make", 52)
	index2 := NewIndex("ZXhhbXBsZS5jb20vYWJvdXQv", "make", 52)

	index.IndexID = "1"
	index2.IndexID = "1"

	err := index.Put(_testConn)
	assert.NoError(t, err)

	err = index2.Put(_testConn)
	assert.Error(t, err)
}

func TestModels_IndexesPut(t *testing.T) {
	indexes := Indexes{
		{
			IndexID: "example.com/about/::hello",
			Word:    "hello",
		},
		{
			IndexID: "example.com/about/::world",
			Word:    "world",
		},
	}
	err := indexes.Put(_testConn)
	assert.NoError(t, err)
}

func TestModels_IndexesPut_Duplicate(t *testing.T) {
	indexes := Indexes{
		{
			IndexID: "1",
			Word:    "hello",
		},
		{
			IndexID: "1",
			Word:    "hello",
		},
	}
	err := indexes.Put(_testConn)
	assert.Error(t, err)
}
