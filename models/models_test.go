package models

import (
	"fmt"
	"log"
	"os"
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/nylar/miru/app"
	"github.com/nylar/miru/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	_ctx *app.Context
	_pkg = "models"

	_db, _index, _document string
)

func init() {
	ctx := app.NewContext()

	if err := ctx.LoadConfig("../config.toml"); err != nil {
		log.Fatalln(err.Error())
	}

	_db = fmt.Sprintf("%s_%s", ctx.Config.Database.Name, "test")
	_index = fmt.Sprintf("%s_%s", ctx.Config.Tables.Index, _pkg)
	_document = fmt.Sprintf("%s_%s", ctx.Config.Tables.Document, _pkg)

	ctx.Config.Database.Name = _db
	ctx.Config.Tables.Index = _index
	ctx.Config.Tables.Document = _document

	if err := ctx.Connect(os.Getenv("RETHINKDB_URL")); err != nil {
		log.Fatalln(err.Error())
	}

	_ctx = ctx

	testutils.SetUp(_ctx, _db, _document, _index)
}

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
	defer testutils.TearDown(_ctx, _db, _document, _index)

	doc := NewDocument(
		"example.com/about/",
		"example.com",
		"About Example Inc.",
		"We make examples and things.",
	)

	err := doc.Put(_ctx)
	assert.NoError(t, err)

	res, err := rdb.Db(_db).Table(_document).Get(doc.DocID).Run(_ctx.Db)
	assert.NoError(t, err)

	var d Document
	err = res.One(&d)
	assert.NoError(t, err)

	assert.NotEqual(t, d.DocID, "")
}

func TestModels_DocumentPut_Duplicate(t *testing.T) {
	defer testutils.TearDown(_ctx, _db, _document, _index)

	doc := NewDocument("example.com/about/", "example.com", "", "")
	doc2 := NewDocument("example.com/about/", "example.com", "", "")

	doc.DocID = "1"
	doc2.DocID = "1"

	err := doc.Put(_ctx)
	assert.NoError(t, err)

	err = doc2.Put(_ctx)
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
	defer testutils.TearDown(_ctx, _db, _document, _index)

	index := NewIndex("example.com/about/", "make", 52)

	err := index.Put(_ctx)
	assert.NoError(t, err)

	res, err := rdb.Db(_db).Table(_index).Get(index.IndexID).Run(_ctx.Db)
	assert.NoError(t, err)

	var i Index
	err = res.One(&i)
	assert.NoError(t, err)

	assert.Equal(t, i.Word, "make")
}

func TestModels_IndexPut_Duplicate(t *testing.T) {
	defer testutils.TearDown(_ctx, _db, _document, _index)

	index := NewIndex("ZXhhbXBsZS5jb20vYWJvdXQv", "make", 52)
	index2 := NewIndex("ZXhhbXBsZS5jb20vYWJvdXQv", "make", 52)

	index.IndexID = "1"
	index2.IndexID = "1"

	err := index.Put(_ctx)
	assert.NoError(t, err)

	err = index2.Put(_ctx)
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
	err := indexes.Put(_ctx)
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
	err := indexes.Put(_ctx)
	assert.Error(t, err)
}
