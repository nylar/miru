package db

import (
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

func TestModels_NewSite(t *testing.T) {
	url := "example.com"
	site := NewSite(url)

	assert.IsType(t, site, new(Site))
	assert.Equal(t, site.SiteID, url)
}

func TestModels_SitePut(t *testing.T) {
	defer TearDbDown()

	site := Site{SiteID: "example.com"}

	err := site.Put(TestConn)
	assert.NoError(t, err)

	res, err := rdb.Db(Database).Table(SiteTable).Get(site.SiteID).Run(TestConn.Session)
	assert.NoError(t, err)

	var s Site
	err = res.One(&s)
	assert.NoError(t, err)

	assert.Equal(t, s.SiteID, "example.com")
}

func TestModels_SitePut_Duplicate(t *testing.T) {
	defer TearDbDown()

	site := Site{SiteID: "example.com"}
	site2 := Site{SiteID: "example.com"}

	err := site.Put(TestConn)
	assert.NoError(t, err)

	err = site2.Put(TestConn)
	assert.Error(t, err)
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

func TestModels_DocumentPut(t *testing.T) {
	defer TearDbDown()

	doc := Document{
		DocID:   "example.com/about/",
		SiteID:  "example.com",
		Title:   "About Example Inc.",
		Content: "We make examples and things.",
	}

	err := doc.Put(TestConn)
	assert.NoError(t, err)

	res, err := rdb.Db(Database).Table(DocumentTable).Get(doc.DocID).Run(TestConn.Session)
	assert.NoError(t, err)

	var d Document
	err = res.One(&d)
	assert.NoError(t, err)

	assert.Equal(t, d.DocID, "example.com/about/")
}

func TestModels_DocumentPut_Duplicate(t *testing.T) {
	defer TearDbDown()

	doc := Document{DocID: "example.com/about/"}
	doc2 := Document{DocID: "example.com/about/"}

	err := doc.Put(TestConn)
	assert.NoError(t, err)

	err = doc2.Put(TestConn)
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
	assert.Equal(t, index.IndexID, "example.com/about/::make")
}

func TestModels_IndexPut(t *testing.T) {
	defer TearDbDown()

	index := Index{
		DocID: "example.com/about/",
		Word:  "make",
		Count: 52,
	}
	index.GenerateID(index.DocID, index.Word)

	err := index.Put(TestConn)
	assert.NoError(t, err)

	res, err := rdb.Db(Database).Table(IndexTable).Get(index.IndexID).Run(TestConn.Session)
	assert.NoError(t, err)

	var i Index
	err = res.One(&i)
	assert.NoError(t, err)

	assert.Equal(t, i.IndexID, "example.com/about/::make")
}

func TestModels_IndexPut_Duplicate(t *testing.T) {
	defer TearDbDown()

	index := Index{DocID: "example.com/about/", Word: "make"}
	index2 := Index{DocID: "example.com/about/", Word: "make"}

	err := index.Put(TestConn)
	assert.NoError(t, err)

	err = index2.Put(TestConn)
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
	err := indexes.Put(TestConn)
	assert.NoError(t, err)
}

func TestModels_IndexesPut_Duplicate(t *testing.T) {
	indexes := Indexes{
		{
			IndexID: "example.com/about/::hello",
			Word:    "hello",
		},
		{
			IndexID: "example.com/about/::hello",
			Word:    "hello",
		},
	}
	err := indexes.Put(TestConn)
	assert.Error(t, err)
}
