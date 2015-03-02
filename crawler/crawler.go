package crawler

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nylar/miru/app"
	"github.com/nylar/miru/index"
	"github.com/nylar/miru/models"
	"github.com/nylar/miru/queue"
)

var (
	UserAgent    = "Miru/1.0 (+http://www.miru.nylar.io)"
	UnwantedTags = "style, script, link, iframe, frame, embed"

	UnreachableUrlError = errors.New("Url did not return a 200 OK response.")
	InvalidUrlError     = errors.New("Url was invalid.")

	Delay int64 = 5
)

func newDocument(document []byte) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(document))
	doc.Find(UnwantedTags).Remove()

	return doc
}

func IndexPage(c *app.Context, q *queue.Queue, url, site string) error {
	req := Request(url)
	resp, err := MustGet(req)
	if err != nil {
		return err
	}

	log.Println("Indexing: %s", url)

	contents := Contents(resp)

	doc := newDocument(contents)

	d := NewDoc(doc, url, site)
	d.Put(c)

	i := index.Index(d.Content, d.DocID)
	i.Put(c)

	Links(doc, q, site)
	return nil
}

func ProcessPages(c *app.Context, q *queue.Queue, site string, delay int64) {
	for q.Len() > 0 {
		item, _ := q.Dequeue()
		IndexPage(c, q, item, site)
		time.Sleep(time.Duration(delay) * time.Second)
	}
	return
}

func Crawl(url string, c *app.Context, q *queue.Queue) error {
	site, err := RootURL(url)
	if err != nil {
		return err
	}

	// TODO: Parse robots.txt file to determine links to avoid, get sitemap
	// (if present) and check if Crawl-Delay is defined

	if err := IndexPage(c, q, url, site); err != nil {
		return err
	}

	go func(c *app.Context, q *queue.Queue, site string, delay int64) {
		ProcessPages(c, q, site, Delay)
	}(c, q, site, Delay)

	return nil
}

func Request(url string) *http.Request {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", UserAgent)

	return request
}

func Get(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func MustGet(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, UnreachableUrlError
	}

	return response, nil
}

func Contents(resp *http.Response) []byte {
	d, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 4194304)) // Limit to 4mb
	defer resp.Body.Close()

	return d
}

func NewDoc(doc *goquery.Document, url, site string) *models.Document {
	title := ExtractTitle(doc)
	content := ExtractText(doc)

	d := models.NewDocument(url, site, title, content)

	return d
}

func RootURL(link string) (string, error) {
	_url, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return _url.Host, nil
}

func Links(doc *goquery.Document, q *queue.Queue, site string) {
	links := ExtractLinks(doc)
	for _, link := range links {
		link, err := ProcessURL(link, site)
		if err != nil {
			continue
		}
		q.Enqueue(link)
	}
}
