package db

import (
	"errors"
	"fmt"

	rdb "github.com/dancannon/gorethink"
)

type Site struct {
	SiteID string `gorethink:"id"`
}

func NewSite(url string) *Site {
	site := new(Site)
	site.SiteID = url
	return site
}

func (s *Site) Put(conn *Connection) error {
	res, _ := rdb.Db(Database).Table(SiteTable).Insert(s).RunWrite(conn.Session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

type Document struct {
	DocID   string `gorethink:"id"`
	SiteID  string `gorethink:"site_id"`
	Title   string `gorethink:"title"`
	Content string `gorethink:"content"`
}

func NewDocument(id, siteID, title, content string) *Document {
	doc := new(Document)
	doc.DocID = id
	doc.SiteID = siteID
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
	i.IndexID = fmt.Sprintf("%s::%s", i.DocID, i.Word)
}

func (i *Index) Put(conn *Connection) error {
	res, _ := rdb.Db(Database).Table(IndexTable).Insert(i).RunWrite(conn.Session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

type Indexes []Index

func (ixs *Indexes) Put(conn *Connection) error {
	res, _ := rdb.Db(Database).Table(IndexTable).Insert(ixs).RunWrite(conn.Session)
	if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}
