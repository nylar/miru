package crawler

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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

func NewDocument(document []byte) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(document))
	doc.Find(UnwantedTags).Remove()

	return doc
}
