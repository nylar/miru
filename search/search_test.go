package search

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"testing"

	rdb "github.com/dancannon/gorethink"
	"github.com/nylar/miru/app"
	"github.com/nylar/miru/index"
	"github.com/nylar/miru/models"
	"github.com/nylar/miru/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	_ctx *app.Context
	_pkg = "search"

	_db, _index, _document string
)

func init() {
	ctx := app.NewContext()

	if err := ctx.LoadConfig("../config.toml"); err != nil {
		log.Fatalln(err.Error())
	}

	_db = fmt.Sprintf("%s_%s", ctx.Config.Database.Name, "test")
	_index = fmt.Sprintf("%s_%s", ctx.Config.Tables.Index, _pkg)
	_document = fmt.Sprintf("%s_%s", ctx.Config.Tables.Document, _pkg)

	ctx.Config.Database.Name = _db
	ctx.Config.Tables.Index = _index
	ctx.Config.Tables.Document = _document

	if err := ctx.Connect(os.Getenv("RETHINKDB_URL")); err != nil {
		log.Fatalln(err.Error())
	}

	_ctx = ctx

	testutils.SetUp(_ctx, _db, _document, _index)
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
	defer testutils.TearDown(_ctx, _db, _document, _index)

	d := models.NewDocument(
		"example.com/about/",
		"example.com",
		"Examples, Examples Everywhere",
		"This is an example of some example content remember though it's just an example",
	)

	if err := d.Put(_ctx); err != nil {
		t.Log(err.Error())
	}

	i := index.Index(d.Content, d.DocID)

	if err := i.Put(_ctx); err != nil {
		t.Log(err.Error())
	}

	res := new(Results)
	err := res.Search("exampl", _ctx)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Results), 1)
}

func TestSearch_Search_NoIndexRaisesError(t *testing.T) {
	defer testutils.TearDown(_ctx, _db, _document, _index)

	d := models.NewDocument(
		"example.com/about/",
		"example.com",
		"Examples, Examples Everywhere",
		"This is an example of some example content remember though it's just an example",
	)

	if err := d.Put(_ctx); err != nil {
		t.Log(err.Error())
	}

	i := index.Index(d.Content, d.DocID)

	if err := i.Put(_ctx); err != nil {
		t.Log(err.Error())
	}

	// Remove index
	rdb.Db(_db).Table(_index).IndexDrop("word").Exec(_ctx.Db)

	res := new(Results)
	err := res.Search("exampl", _ctx)

	assert.Error(t, err)
	assert.Equal(t, len(res.Results), 0)

	// Re-add index
	rdb.Db(_db).Table(_index).IndexCreate("word").Exec(_ctx.Db)
}
