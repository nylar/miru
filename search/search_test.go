package search

import (
	"html/template"
	"log"
	"os"
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/nylar/miru/db"
	"github.com/nylar/miru/index"
	"github.com/stretchr/testify/assert"
)

var _testConn *db.Connection

func init() {
	var err error
	_testConn, err = db.NewConnection("test", os.Getenv("RETHINKDB_URL"))
	if err != nil {
		log.Fatalln("Could not create a connection for testing. Exiting.")
	}

	db.SetDbUp(_testConn, "search")
}

func TestSearch_RenderSpeed(t *testing.T) {
	tests := []struct {
		Input  float64
		Output string
	}{
		{
			0.3432242,
			"0.3432 seconds",
		},
		{
			0.11199,
			"0.1120 seconds",
		},
	}

	for _, test := range tests {
		res := new(Results)
		res.Speed = test.Input

		assert.Equal(t, res.RenderSpeed(), test.Output)
	}
}

func TestSearch_RenderSpeedHTML(t *testing.T) {
	tests := []struct {
		Input  float64
		Output template.HTML
	}{
		{
			0.3432242,
			template.HTML("0.3432 seconds"),
		},
		{
			0.11199,
			template.HTML("0.1120 seconds"),
		},
	}

	for _, test := range tests {
		res := new(Results)
		res.Speed = test.Input

		assert.Equal(t, res.RenderSpeedHTML(), test.Output)
	}
}

func TestSearch_RenderCount(t *testing.T) {
	r1 := Result{}
	r2 := Result{}
	r3 := Result{}
	res := new(Results)
	res.Results = []Result{r1, r2, r3}

	tests := []struct {
		Input  int64
		Output string
	}{
		{
			3,
			"Showing 3 of 3 results.",
		},
		{
			50,
			"Showing 3 of 50 results.",
		},
		{
			1000,
			"Showing 3 of 1000 results.",
		},
	}

	for _, test := range tests {
		res.Count = test.Input

		assert.Equal(t, res.RenderCount(), test.Output)
	}
}

func TestSearch_RenderCountHTML(t *testing.T) {
	r1 := Result{}
	r2 := Result{}
	r3 := Result{}
	res := new(Results)
	res.Results = []Result{r1, r2, r3}

	tests := []struct {
		Input  int64
		Output template.HTML
	}{
		{
			3,
			template.HTML("Showing 3 of 3 results."),
		},
		{
			50,
			template.HTML("Showing 3 of 50 results."),
		},
		{
			1000,
			template.HTML("Showing 3 of 1000 results."),
		},
	}

	for _, test := range tests {
		res.Count = test.Input

		assert.Equal(t, res.RenderCountHTML(), test.Output)
	}
}

func TestSearch_ParseQuery(t *testing.T) {
	res := new(Results)

	tests := []struct {
		Input  string
		Output []string
	}{
		{
			"hello world",
			[]string{"hello", "world"},
		},
		{
			"doubled-barreled word",
			[]string{"doubled-barreled", "word"},
		},
		{
			"",
			[]string{""},
		},
	}

	for _, test := range tests {
		assert.Equal(t, res.ParseQuery(test.Input), test.Output)
	}
}

func TestSearch_Search(t *testing.T) {
	defer db.TearDbDown(_testConn)

	d := db.Document{
		Title:   "Examples, Examples Everywhere",
		Content: "This is an example of some example content remember though it's just an example",
		Site:    "example.com",
	}
	d.GenerateID("example.com/about/")

	if err := d.Put(_testConn); err != nil {
		t.Log(err.Error())
	}

	i := index.Index(d.Content, d.DocID)

	if err := i.Put(_testConn); err != nil {
		t.Log(err.Error())
	}

	res := new(Results)
	err := res.Search("exampl", _testConn)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Results), 1)
}

func TestSearch_Search_NoIndexRaisesError(t *testing.T) {
	defer db.TearDbDown(_testConn)

	d := db.Document{
		Title:   "Examples, Examples Everywhere",
		Content: "This is an example of some example content remember though it's just an example",
		Site:    "example.com",
	}
	d.GenerateID("example.com/about/")

	if err := d.Put(_testConn); err != nil {
		t.Log(err.Error())
	}

	i := index.Index(d.Content, d.DocID)

	if err := i.Put(_testConn); err != nil {
		t.Log(err.Error())
	}

	// Remove index
	rdb.Db(db.Database).Table(db.IndexTable).IndexDrop("word").Exec(_testConn.Session)

	res := new(Results)
	err := res.Search("exampl", _testConn)

	assert.Error(t, err)
	assert.Equal(t, len(res.Results), 0)

	// Re-add index
	rdb.Db(db.Database).Table(db.IndexTable).IndexCreate("word").Exec(_testConn.Session)
}
