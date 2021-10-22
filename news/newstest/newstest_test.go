package newstest

import (
	"avilego.me/recent_news/news"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinder(t *testing.T) {
	for i, tData := range tsFinder {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			t.Parallel()
			previews := tData.finder.Find(tData.keywords)
			assert.Equal(t, tData.expected, previews)
		})
	}
}

var tsFinder = []struct {
	keywords string
	finder   Finder
	expected []news.Preview
}{
	{
		"AMD",
		Finder{"AMD", news.Previews[0:2]},
		news.Previews[0:2],
	},
	{
		"nothing",
		Finder{"AMD", news.Previews[0:2]},
		nil,
	},
}
