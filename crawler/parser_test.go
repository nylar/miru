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

	doc := NewDocument(html)

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

	doc := NewDocument(html)

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

	doc := NewDocument(html)

	title := ExtractTitle(doc)

	assert.Equal(t, "About Us", title)
	assert.NotEqual(t, "We Rock", title)
	assert.NotEqual(t, "", title)
}

func TestParser_ExtractTitleEmpty(t *testing.T) {
	html := []byte(`<!DOCTYPE html><html><head></head><body></body></html>`)

	doc := NewDocument(html)

	title := ExtractTitle(doc)

	assert.Equal(t, "", title)
}

func TestParser_ExtractTextEmpty(t *testing.T) {
	doc := NewDocument([]byte(""))
	text := ExtractText(doc)
	assert.Equal(t, "", text)
}

func TestParser_ExtractTextFromPTags(t *testing.T) {
	html := []byte(`<p>I am text one.</p><p>I am text two.</p>`)
	doc := NewDocument(html)
	text := ExtractText(doc)
	assert.Equal(t, "I am text one.\nI am text two.", text)
}
