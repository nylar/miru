package miru

import (
	"errors"

	rdb "github.com/dancannon/gorethink"
	"github.com/satori/go.uuid"
)

// Document stores data about a page.
type Document struct {
	DocID   string `gorethink:"id" json:"document_id"`
	Url     string `gorethink:"url" json:"url"`
	Site    string `gorethink:"site" json:"site"`
	Title   string `gorethink:"title" json:"title"`
	Content string `gorethink:"content" json:"content"`
}

// NewDocument creates a new document instance
func NewDocument(url, site, title, content string) *Document {
	doc := new(Document)
	doc.DocID = uuid.NewV4().String()
	doc.Url = url
	doc.Site = site
	doc.Title = title
	doc.Content = content

	return doc
}

// Put writes a document to the datastore.
func (d *Document) Put(c *Context) error {
	res, _ := rdb.Db(c.Config.Database.Name).Table(
		c.Config.Tables.Document).Insert(d).RunWrite(c.Db)

	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

// Index stores data on a given word in a document.
type Index struct {
	IndexID string `gorethink:"id" json:"index_id"`
	DocID   string `gorethink:"doc_id" json:"document_id"`
	Word    string `gorethink:"word" json:"word"`
	Count   int64  `gorethink:"count" json:"count"`
}

// NewIndex creates a new index instance
func NewIndex(docID, word string, count int64) *Index {
	index := new(Index)
	index.IndexID = uuid.NewV4().String()
	index.DocID = docID
	index.Word = word
	index.Count = count

	return index
}

// Put writes an index to the datastore
func (i *Index) Put(c *Context) error {
	res, _ := rdb.Db(c.Config.Database.Name).Table(
		c.Config.Tables.Index).Insert(i).RunWrite(c.Db)

	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

// Indexes is a slice of index, holds all the words in a document
type Indexes []*Index

// Put writes a slice of index to the datastore.
func (ixs *Indexes) Put(c *Context) error {
	res, _ := rdb.Db(c.Config.Database.Name).Table(
		c.Config.Tables.Index).Insert(ixs).RunWrite(c.Db)

	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}
