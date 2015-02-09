package db

import (
	"log"
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
)

var _conn *Connection

func init() {
	var err error
	_conn, err = NewConnection("test", "localhost:28015")
	if err != nil {
		log.Fatalln("Could not create a connection for testing. Exiting.")
	}

	Database = "testing"

	rdb.DbCreate(Database).Exec(_conn.Session)
	rdb.Db(Database).TableCreate(SiteTable).Exec(_conn.Session)
	rdb.Db(Database).TableCreate(DocumentTable).Exec(_conn.Session)
	rdb.Db(Database).TableCreate(IndexTable).Exec(_conn.Session)
	rdb.Db(Database).Table(IndexTable).IndexCreate("word").Exec(_conn.Session)
}

func tearDbDown() {
	// Clear 'Site' table
	rdb.Db(Database).Table(SiteTable).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(_conn.Session)

	// Clear 'Document' table
	rdb.Db(Database).Table(DocumentTable).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(_conn.Session)

	// Clear 'Index' table
	rdb.Db(Database).Table(IndexTable).Delete(rdb.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Exec(_conn.Session)
}

func TestModels_NewSite(t *testing.T) {
	url := "example.com"
	site := NewSite(url)

	assert.IsType(t, site, new(Site))
	assert.Equal(t, site.SiteID, url)
}

func TestIndexer_SitePut(t *testing.T) {
	defer tearDbDown()

	site := Site{SiteID: "example.com"}

	err := site.Put(_conn)
	assert.NoError(t, err)

	res, err := rdb.Db(Database).Table(SiteTable).Get(site.SiteID).Run(_conn.Session)
	assert.NoError(t, err)

	var s Site
	err = res.One(&s)
	assert.NoError(t, err)

	assert.Equal(t, s.SiteID, "example.com")
}

func TestIndexer_SitePut_Duplicate(t *testing.T) {
	defer tearDbDown()

	site := Site{SiteID: "example.com"}
	site2 := Site{SiteID: "example.com"}

	err := site.Put(_conn)
	assert.NoError(t, err)

	err = site2.Put(_conn)
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

func TestIndexer_DocumentPut(t *testing.T) {
	defer tearDbDown()

	doc := Document{
		DocID:   "example.com/about/",
		SiteID:  "example.com",
		Title:   "About Example Inc.",
		Content: "We make examples and things.",
	}

	err := doc.Put(_conn)
	assert.NoError(t, err)

	res, err := rdb.Db(Database).Table(DocumentTable).Get(doc.DocID).Run(_conn.Session)
	assert.NoError(t, err)

	var d Document
	err = res.One(&d)
	assert.NoError(t, err)

	assert.Equal(t, d.DocID, "example.com/about/")
}

func TestIndexer_DocumentPut_Duplicate(t *testing.T) {
	defer tearDbDown()

	doc := Document{DocID: "example.com/about/"}
	doc2 := Document{DocID: "example.com/about/"}

	err := doc.Put(_conn)
	assert.NoError(t, err)

	err = doc2.Put(_conn)
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

func TestIndexer_IndexPut(t *testing.T) {
	defer tearDbDown()

	index := Index{
		DocID: "example.com/about/",
		Word:  "make",
		Count: 52,
	}
	index.GenerateID(index.DocID, index.Word)

	err := index.Put(_conn)
	assert.NoError(t, err)

	res, err := rdb.Db(Database).Table(IndexTable).Get(index.IndexID).Run(_conn.Session)
	assert.NoError(t, err)

	var i Index
	err = res.One(&i)
	assert.NoError(t, err)

	assert.Equal(t, i.IndexID, "example.com/about/::make")
}

func TestIndexer_IndexPut_Duplicate(t *testing.T) {
	defer tearDbDown()

	index := Index{DocID: "example.com/about/", Word: "make"}
	index2 := Index{DocID: "example.com/about/", Word: "make"}

	err := index.Put(_conn)
	assert.NoError(t, err)

	err = index2.Put(_conn)
	assert.Error(t, err)
}
