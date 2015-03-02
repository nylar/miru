package crawler

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
