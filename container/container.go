package container

import (
	"avilego.me/recent_news/env"
	"avilego.me/recent_news/news"
	"avilego.me/recent_news/persistence"
	"avilego.me/recent_news/rss"
	"time"

	"log"
	"os"
)

var Providers = []news.AsyncProvider{
	rss.NewRssNewsProvider(
		[]rss.Source{
			rss.NewSource("http://api2.rtve.es/rss/temas_noticias.xml"),
			rss.NewSource("http://rss.cnn.com/rss/edition_world.rss"),
			rss.NewSource("https://www.phoronix.com/rss.php"),
		},
		time.Tick(1*time.Minute),
	),
}

func GetLogger() *log.Logger {
	file, _ := os.OpenFile(env.LogFile(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	return log.New(file, "", log.LstdFlags)
}

func GetCollector() news.Collector {
	return news.Collector{
		Providers: Providers,
		Keeper:    persistence.NewMongoKeeper(),
		Logger:    GetLogger(),
	}
}

func GetBrowser() news.Finder {
	return persistence.NewMongoFinder()
}
