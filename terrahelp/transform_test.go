package terrahelp

import (
	"github.com/stretchr/testify/assert"

	"fmt"
	"io/ioutil"
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
endsWith1ExplicitNewLine
`

var endsWithoutNewLine = `some content here
and here
endsWithoutNewLine`

var hasNewLinesInMiddle = `some content here
and here


hasNewLinesInMiddle`

func TestStreamTransformable_read(t *testing.T) {
	srcScenarios := []string{
		endsWith1ExplicitNewLine,
		endsWith2ExplicitNewLines,
		endsWithoutNewLine,
		hasNewLinesInMiddle,
	}

	for _, s := range srcScenarios {
		// Given
		r := strings.NewReader(s)
		sci := NewStreamTransformable(r, nil)

		// When
		b, err := sci.read()

		// Then
		assert.NoError(t, err)
		assert.Equal(t, s, string(b))
	}

}

func TestFileTransformable_read_DirError(t *testing.T) {

	f, err := ioutil.TempDir("", "testdir")
	sci := NewFileTransformable(f, false, "")

	// When
	_, err = sci.read()

	// Then
	assert.EqualError(t, err, fmt.Sprintf("%s must be a valid file", f))

}
