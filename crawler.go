package miru

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	// UserAgent is passed on each HTTP request to identify the crawler.
	UserAgent = "Miru/1.0 (+http://www.miru.nylar.io)"
	// UnwantedTags are stripped from all HTML documents.
	UnwantedTags = "style, script, link, iframe, frame, embed"
	// ErrUnreachableURL for when the error doesn't return 200 OK.
	ErrUnreachableURL = errors.New("Url did not return a 200 OK response.")
	// ErrInvalidURL for when not a valid URL.
	ErrInvalidURL = errors.New("Url was invalid.")
	// Delay is time in between each crawl
	Delay int64 = 5
)

func newDocument(document []byte) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(document))
	doc.Find(UnwantedTags).Remove()

	return doc
}

// IndexPage is called by ProcessPages and handles dealing with individual pages
func IndexPage(c *Context, q *Queue, url, site string) error {
	req := Request(url)
	resp, err := MustGet(req)
	if err != nil {
		return err
	}

	contents := Contents(resp)

	doc := newDocument(contents)

	d := NewDoc(doc, url, site)
	d.Put(c)

	i := Indexer(d.Content, d.DocID)
	i.Put(c)

	Links(doc, q, site)
	return nil
}

// ProcessPages process all queue items and proceeds to index them.
func ProcessPages(c *Context, q *Queue, site string, delay int64) {
	for q.Len() > 0 {
		item, _ := q.Dequeue()
		IndexPage(c, q, item, site)
		time.Sleep(time.Duration(delay) * time.Second)
	}
	q.Status = "finished"
	return
}

// Crawl processes pages concurrently
func Crawl(url string, c *Context, q *Queue) error {
	site, err := RootURL(url)
	if err != nil {
		return err
	}

	// TODO: Parse robots.txt file to determine links to avoid, get sitemap
	// (if present) and check if Crawl-Delay is defined

	if err := IndexPage(c, q, url, site); err != nil {
		return err
	}

	go func(c *Context, q *Queue, site string, delay int64) {
		ProcessPages(c, q, site, Delay)
	}(c, q, site, Delay)

	return nil
}

// Request builds a new request using UserAgent as a header
func Request(url string) *http.Request {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", UserAgent)

	return request
}

// Get uses a custom request and returns a response.
func Get(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// MustGet is a strict version of Get
func MustGet(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ErrUnreachableURL
	}

	return response, nil
}

// Contents reads data from a response into a byte slice, limits to 4mb.
func Contents(resp *http.Response) []byte {
	d, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 4194304)) // Limit to 4mb
	defer resp.Body.Close()

	return d
}

// NewDoc extracts data from a page and creates a new document.
func NewDoc(doc *goquery.Document, url, site string) *Document {
	title := ExtractTitle(doc)
	content := ExtractText(doc)

	d := NewDocument(url, site, title, content)

	return d
}

// RootURL returns the domain for a given link
func RootURL(link string) (string, error) {
	_url, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return _url.Host, nil
}

// Links extracts all internal links from a page and enqueues them.
func Links(doc *goquery.Document, q *Queue, site string) {
	links := ExtractLinks(doc)
	for _, link := range links {
		link, err := ProcessURL(link, site)
		if err != nil {
			continue
		}
		q.Enqueue(link)
	}
}
