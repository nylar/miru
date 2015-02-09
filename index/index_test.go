package index

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
