package news

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinder(t *testing.T) {
	for i, tData := range tsFinder {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			t.Parallel()
			previews := tData.finder.FindRelated(tData.keywords)
			assert.Equal(t, tData.expected, previews)
		})
	}
}

var tsFinder = []struct {
	keywords string
	finder   Finder
	expected []Preview
}{
	{
		"AMD",
		FinderMock{Previews[0:2]},
		Previews[0:2],
	},
	{
		"nothing",
		FinderMock{nil},
		nil,
	},
}
