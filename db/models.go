package db

import (
	"errors"

	rdb "github.com/dancannon/gorethink"
	"github.com/satori/go.uuid"
)

type Document struct {
	DocID   string `gorethink:"id" json:"document_id"`
	Url     string `gorethink:"url" json:"url"`
	Site    string `gorethink:"site" json:"site"`
	Title   string `gorethink:"title" json:"title"`
	Content string `gorethink:"content" json:"content"`
}

func NewDocument(url, site, title, content string) *Document {
	doc := new(Document)
	doc.DocID = uuid.NewV4().String()
	doc.Url = url
	doc.Site = site
	doc.Title = title
	doc.Content = content

	return doc
}

func (d *Document) Put(conn *Connection) error {
	res, _ := rdb.Db(Database).Table(DocumentTable).Insert(d).RunWrite(conn.Session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

type Index struct {
	IndexID string `gorethink:"id" json:"index_id"`
	DocID   string `gorethink:"doc_id" json:"document_id"`
	Word    string `gorethink:"word" json:"word"`
	Count   int64  `gorethink:"count" json:"count"`
}

func NewIndex(docID, word string, count int64) *Index {
	index := new(Index)
	index.IndexID = uuid.NewV4().String()
	index.DocID = docID
	index.Word = word
	index.Count = count

	return index
}

func (i *Index) Put(conn *Connection) error {
	res, _ := rdb.Db(Database).Table(IndexTable).Insert(i).RunWrite(conn.Session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

type Indexes []*Index

func (ixs *Indexes) Put(conn *Connection) error {
	res, _ := rdb.Db(Database).Table(IndexTable).Insert(ixs).RunWrite(conn.Session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}
