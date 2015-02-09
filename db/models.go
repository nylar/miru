package db

import "fmt"

type Site struct {
	SiteID string `gorethink:"site_id"`
}

func NewSite(url string) *Site {
	site := new(Site)
	site.SiteID = url
	return site
}

type Document struct {
	DocID   string `gorethink:"doc_id"`
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

type Index struct {
	IndexID string `gorethink:"index_id"`
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
