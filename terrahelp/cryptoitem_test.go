package terrahelp

import (
	"github.com/stretchr/testify/assert"

	"strings"
	"testing"
)

var endsWith2ExplicitNewLines = `some content here
and here
endsWith2ExplicitNewLines

`

// This one doesn't work at present
var endsWith1ExplicitNewLine = `some content here
and here
endsWithExplicitNewLine
`

var endsWithoutNewLine = `some content here
and here
endsWithoutNewLine`

var hasNewLinesInMiddle = `some content here
and here


hasNewLinesInMiddle`

func TestStreamCryptoItem_readFromSource(t *testing.T) {
	srcScenarios := []string{
		//endsWith1ExplicitNewLine,
		endsWith2ExplicitNewLines,
		endsWithoutNewLine,
		hasNewLinesInMiddle,
	}

	for _, s := range srcScenarios {
		// Given
		r := strings.NewReader(s)
		sci := NewStreamCryptoItem(r, nil)

		// When
		b, err := sci.readFromSource()

		// Then
		assert.NoError(t, err)
		assert.Equal(t, s, string(b))
	}

}