package factory

import (
	"avilego.me/recent_news/env"
	"avilego.me/recent_news/news"
	"avilego.me/recent_news/persistence"
	"avilego.me/recent_news/rss"
	"log"
	"os"
	"time"
)

var providers = []news.Provider{
	rss.NewRssNewsProvider(
		[]rss.Source{
			rss.NewSource("http://api2.rtve.es/rss/temas_noticias.xml"),
			rss.NewSource("http://rss.cnn.com/rss/edition_world.rss"),
			rss.NewSource("https://e00-elmundo.uecdn.es/elmundo/rss/espana.xml"),
			rss.NewSource("https://rss.elconfidencial.com/espana/"),
			rss.NewSource("https://e00-expansion.uecdn.es/rss/portada.xml"),
			rss.NewSource("https://www.phoronix.com/rss.php"),
			rss.NewSource("http://feeds2.feedburner.com/libertaddigital/portada"),
			rss.NewSource("https://www.lavanguardia.com/rss/home.xml"),
			rss.NewSource("https://web.gencat.cat/es/actualitat/rss.html"),
			rss.NewSource("https://www.abc.es/rss/feeds/abcPortada.xml"),
		},
		time.Tick(5*time.Minute),
	),
}

func logger() *log.Logger {
	file, _ := os.OpenFile(env.LogFile(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	return log.New(file, "", log.LstdFlags)
}

func Collector() news.Collector {
	return news.Collector{
		Providers: providers,
		Keeper:    Keeper(),
		Logger:    logger(),
	}
}

func Finder() news.Finder {
	return persistence.NewMongoKeeperFinder()
}

func Keeper() news.Keeper {
	return persistence.NewMongoKeeperFinder()
}

func Cleaner() news.Cleaner {
	return news.Cleaner{
		KeeperFinder: persistence.NewMongoKeeperFinder(),
		Trigger:      time.Tick(1 * time.Hour),
		Ttl:          int64((24 * time.Hour).Seconds()),
	}
}
