package search

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	rdb "github.com/dancannon/gorethink"
	"github.com/nylar/miru/db"
	"github.com/nylar/miru/index"
)

type Result struct {
	db.Document
	db.Index
}

type Results struct {
	Speed   float64  `json:"speed"`
	Count   int64    `json:"count"`
	Results []Result `json:"results"`
}

func (rxs *Results) RenderSpeed() string {
	return fmt.Sprintf("%.4f seconds", rxs.Speed)
}

func (rxs *Results) RenderSpeedHTML() template.HTML {
	return template.HTML(rxs.RenderSpeed())
}

func (rxs *Results) RenderCount() string {
	return fmt.Sprintf("Showing %d of %d results.", len(rxs.Results), rxs.Count)
}

func (rxs *Results) RenderCountHTML() template.HTML {
	return template.HTML(rxs.RenderCount())
}

func (rxs *Results) Search(query string, conn *db.Connection) error {
	start := time.Now()

	query = index.Normalise(query)

	keywords := rxs.ParseQuery(query)
	results, err := rdb.Db(
		db.Database).Table(
		db.IndexTable).GetAllByIndex(
		"word", rdb.Args(keywords)).EqJoin(
		"doc_id", rdb.Db(db.Database).Table(
			db.DocumentTable)).Zip().OrderBy(
		rdb.Desc("count")).Run(conn.Session)

	if err != nil {
		return err
	}
	results.All(&rxs.Results)

	t := time.Since(start).Seconds()
	rxs.Speed = t

	rxs.Count = int64(len(rxs.Results))

	return nil
}

func (rxs *Results) ParseQuery(query string) []string {
	return strings.Split(query, " ")
}
