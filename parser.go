package miru

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractTitle looks for either a title tag or h1 tag and sets that as the title
func ExtractTitle(doc *goquery.Document) string {
	title := doc.Find("title")
	if title.Length() > 0 {
		return title.First().Text()
	}
	heading := doc.Find("h1")
	if heading.Length() > 0 {
		return heading.First().Text()
	}
	return ""
}

// ExtractText returns all p tags in a page
func ExtractText(doc *goquery.Document) string {
	texts := []string{}
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text != "" {
			texts = append(texts, text)
		}
	})
	return strings.Join(texts, "\n")
}

// ExtractLinks returns all internal links from a page.
func ExtractLinks(doc *goquery.Document) []string {
	links := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// Only interested in anchors that have a href attribute.
		link, href := s.Attr("href")
		if href {
			links = append(links, link)
		}
	})
	return links
}
