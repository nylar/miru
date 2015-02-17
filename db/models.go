package db

import (
	"encoding/base64"
	"errors"
	"fmt"

	rdb "github.com/dancannon/gorethink"
)

type Document struct {
	DocID   string `gorethink:"id"`
	Site    string `gorethink:"site"`
	Title   string `gorethink:"title"`
	Content string `gorethink:"content"`
}

func NewDocument(url, site, title, content string) *Document {
	doc := new(Document)
	doc.Site = site
	doc.Title = title
	doc.Content = content

	doc.GenerateID(url)

	return doc
}

func (d *Document) Put(conn *Connection) error {
	res, _ := rdb.Db(Database).Table(DocumentTable).Insert(d).RunWrite(conn.Session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

func (d *Document) GenerateID(url string) {
	d.DocID = base64.StdEncoding.EncodeToString([]byte(url))
}

type Index struct {
	IndexID string `gorethink:"id"`
	DocID   string `gorethink:"doc_id"`
	Word    string `gorethink:"word"`
	Count   int64  `gorethink:"count"`
}

func NewIndex(docID, word string, count int64) *Index {
	index := new(Index)
	index.DocID = docID
	index.Word = word
	index.Count = count

	index.GenerateID(docID, word)

	return index
}

func (i *Index) GenerateID(docID, word string) {
	id := fmt.Sprintf("%s::%s", i.DocID, i.Word)
	i.IndexID = base64.StdEncoding.EncodeToString([]byte(id))
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
