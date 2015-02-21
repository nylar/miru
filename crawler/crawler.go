package crawler

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/nylar/miru/db"
	"github.com/nylar/miru/index"
)

var (
	UserAgent    = "Miru/1.0 (+http://www.miru.nylar.io)"
	UnwantedTags = "style, script, link, iframe, frame, embed"
)

func getDocument(url string) ([]byte, error) {
	client := &http.Client{}
	data := []byte{}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", UserAgent)

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}

	data, _ = ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	return data, nil
}

func newDocument(document []byte) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(document))
	doc.Find(UnwantedTags).Remove()

	return doc
}

func Crawl(url string, conn *db.Connection) error {
	data, err := getDocument(url)
	if err != nil {
		return err
	}

	doc := newDocument(data)
	document := db.NewDocument(url, url, ExtractTitle(doc), ExtractText(doc))
	document.Put(conn)

	i := index.Index(document.Content, document.DocID)

	i.Put(conn)

	return nil
}
