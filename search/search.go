package search

import (
	"fmt"
	"html/template"

	"github.com/nylar/miru/db"
)

type Result struct {
	db.Document
	db.Index
}

type Results struct {
	Speed   float64
	Count   int64
	Results []Result
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
