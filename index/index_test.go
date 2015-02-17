package index

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/nylar/miru/db"
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
		Output db.Indexes
	}{
		{
			"I am a block of text and I am going to be indexed",
			db.Indexes{
				db.NewIndex(docID, "block", 1),
				db.NewIndex(docID, "text", 1),
				db.NewIndex(docID, "going", 1),
				db.NewIndex(docID, "indexed", 1),
			},
		},
	}

	for x, test := range tests {
		i := Index(test.Input, docID)
		id := fmt.Sprintf("%s::%s", docID, i[x].Word)

		assert.Equal(
			t,
			base64.StdEncoding.EncodeToString([]byte(id)),
			i[x].IndexID,
		)
		assert.Equal(t, len(test.Output), len(i))
	}
}

func TestIndex_RemoveDuplicates(t *testing.T) {
	indexes := Index("hello world cruel world hello world", "")
	assert.Equal(t, len(indexes), 3)
}
