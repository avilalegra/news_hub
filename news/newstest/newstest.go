package newstest

import (
	"avilego.me/recent_news/news"
)

type Finder struct {
	Keywords string
	Previews []news.Preview
}

func (b Finder) FindBefore(unixTime int64) []news.Preview {
	panic("implement me")
}

func (b Finder) Find(keywords string) []news.Preview {
	if b.Keywords != keywords {
		return nil
	}
	return b.Previews
}
