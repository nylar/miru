package miru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex_Stopper(t *testing.T) {
	tests := []struct {
		Input  string
		Output bool
	}{
		{
			"computer",
			false,
		},

		{
			"the",
			true,
		},

		{
			"technology",
			false,
		},

		{
			"wasn't",
			true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output, Stopper(test.Input))
	}
}

func TestIndex_Normalise(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			"Capitalised",
			"capitalis",
		},
		{
			"UPPERCASE",
			"uppercas",
		},
		{
			"lowercase",
			"lowercas",
		},
		{
			"the",
			"",
		},
		// {
		// 	"with-punctuation?!",
		// 	"with-punctuation",
		// },
		{
			"stemmed",
			"stem",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output, Normalise(test.Input))
	}
}

func TestIndex_Index(t *testing.T) {
	docID := "example.com/about/"

	tests := []struct {
		Input  string
		Output Indexes
	}{
		{
			"I am a block of text and I am going to be indexed",
			Indexes{
				NewIndex(docID, "block", 1),
				NewIndex(docID, "text", 1),
				NewIndex(docID, "going", 1),
				NewIndex(docID, "indexed", 1),
			},
		},
	}

	for x, test := range tests {
		i := Indexer(test.Input, docID)

		assert.NotEqual(
			t,
			"",
			i[x].IndexID,
		)
		assert.Equal(t, len(test.Output), len(i))
	}
}

func TestIndex_RemoveDuplicates(t *testing.T) {
	indexes := Indexer("hello world cruel world hello world", "")
	assert.Equal(t, len(indexes), 3)
}
