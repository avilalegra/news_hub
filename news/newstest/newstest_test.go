package newstest

import (
	"avilego.me/recent_news/news"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinder(t *testing.T) {
	for _, tData := range tsFinder {
		previews := tData.finder.Find(tData.keywords)
		assert.Equal(t, tData.expected, previews)
	}
}

var tsFinder = []struct {
	keywords string
	finder   Finder
	expected []news.Preview
}{
	{
		"AMD",
		Finder{"AMD", Previews[0:2]},
		Previews[0:2],
	},
	{
		"nothing",
		Finder{"AMD", Previews[0:2]},
		nil,
	},
}
