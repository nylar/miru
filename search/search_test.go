package search

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
