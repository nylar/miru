package miru

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	rdb "github.com/dancannon/gorethink"
)

// Result holds data for a result's document and index
type Result struct {
	Document
	Index
}

// Results holds all of the results, the time taken to perform the query and the
// number of results.
type Results struct {
	Speed   float64  `json:"speed"`
	Count   int64    `json:"count"`
	Results []Result `json:"results"`
}

// RenderSpeed formats the speed of the query
func (rxs *Results) RenderSpeed() string {
	return fmt.Sprintf("%.4f seconds", rxs.Speed)
}

// RenderSpeedHTML formats the speed of the query and escapes for use in templates.
func (rxs *Results) RenderSpeedHTML() template.HTML {
	return template.HTML(rxs.RenderSpeed())
}

// RenderCount formats the number of results
func (rxs *Results) RenderCount() string {
	return fmt.Sprintf("Showing %d of %d results.", len(rxs.Results), rxs.Count)
}

// RenderCountHTML formats the number of results and escapes for use in templates.
func (rxs *Results) RenderCountHTML() template.HTML {
	return template.HTML(rxs.RenderCount())
}

// Search returns a list of Results for a given query.
func (rxs *Results) Search(query string, c *Context) error {
	start := time.Now()

	query = Normalise(query)

	keywords := rxs.ParseQuery(query)
	results, err := rdb.Db(
		c.Config.Database.Name).Table(
		c.Config.Tables.Index).GetAllByIndex(
		"word", rdb.Args(keywords)).EqJoin(
		"doc_id", rdb.Db(c.Config.Database.Name).Table(
			c.Config.Tables.Document)).Zip().OrderBy(
		rdb.Desc("count")).Run(c.Db)

	if err != nil {
		return err
	}
	results.All(&rxs.Results)

	t := time.Since(start).Seconds()
	rxs.Speed = t

	rxs.Count = int64(len(rxs.Results))

	return nil
}

// ParseQuery splits words into a list of individual words.
func (rxs *Results) ParseQuery(query string) []string {
	return strings.Split(query, " ")
}
