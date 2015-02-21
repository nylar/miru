package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_ExtractTitleFromTitle(t *testing.T) {
	html := []byte(`
<!DOCTYPE html>
<html>
<head>
	<title>The Title</title>
</head>
<body></body>
</html>`)

	doc := newDocument(html)

	title := ExtractTitle(doc)

	assert.Equal(t, "The Title", title)
}

func TestParser_ExtractTitleFromHeading1(t *testing.T) {
	html := []byte(`
<!DOCTYPE html>
<html>
<head></head>
<body>
	<h1>The Title</h1>
</body>
</html>`)

	doc := newDocument(html)

	title := ExtractTitle(doc)

	assert.Equal(t, "The Title", title)
}

func TestParse_ExtractTitlePrecendence(t *testing.T) {
	html := []byte(`
<!DOCTYPE html>
<html>
<head>
	<title>About Us</title>
</head>
<body>
	<h1>We rock</h1>
</body>
</html>`)

	doc := newDocument(html)

	title := ExtractTitle(doc)

	assert.Equal(t, "About Us", title)
	assert.NotEqual(t, "We Rock", title)
	assert.NotEqual(t, "", title)
}

func TestParser_ExtractTitleEmpty(t *testing.T) {
	html := []byte(`<!DOCTYPE html><html><head></head><body></body></html>`)

	doc := newDocument(html)

	title := ExtractTitle(doc)

	assert.Equal(t, "", title)
}

func TestParser_ExtractTextEmpty(t *testing.T) {
	doc := newDocument([]byte(""))
	text := ExtractText(doc)
	assert.Equal(t, "", text)
}

func TestParser_ExtractTextFromPTags(t *testing.T) {
	html := []byte(`<p>I am text one.</p><p>I am text two.</p>`)
	doc := newDocument(html)
	text := ExtractText(doc)
	assert.Equal(t, "I am text one.\nI am text two.", text)
}

func TestParser_ExtractLinks_Empty(t *testing.T) {
	doc := newDocument([]byte{})
	q := NewQueue()

	ExtractLinks(doc, q)

	assert.Equal(t, q.Len(), 0)
}

func TestParser_ExtractLinks_Valid(t *testing.T) {
	htmlSoup := []byte(`
<p>
	<a href="http://example.org/1">Link 1</a>
	<br>
	<a href="http://example.org/2">Link 2</a>
</p>`)

	doc := newDocument(htmlSoup)
	q := NewQueue()

	ExtractLinks(doc, q)

	assert.Equal(t, q.Len(), 2)
}

func TestParser_ExtractLinks_Invalid(t *testing.T) {
	// This should return an error but html.Parse doesn't seem to care.
	invalidHTML := []byte(`<html><body><aef<eqf>>>qq></body></ht>`)

	doc := newDocument(invalidHTML)
	q := NewQueue()
	ExtractLinks(doc, q)

	assert.Equal(t, q.Len(), 0)
}

func TestParser_ExtractLinks_NoDuplicates(t *testing.T) {
	htmlWithDupes := []byte(`
<p>
	<a href="http://example.org/1">Link 1</a>
	<a href="http://example.org/2">Link 1</a>
	<a href="http://example.org/3">Link 1</a>
	<a href="http://example.org/1">Link 1</a>
</p>`)

	doc := newDocument(htmlWithDupes)
	q := NewQueue()

	ExtractLinks(doc, q)

	assert.Equal(t, q.Len(), 3)
}
