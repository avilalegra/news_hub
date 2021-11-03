package config

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestConfigLoaderReturnsConfig(t *testing.T) {
	config := `
rss_news_provider:
  sources:
    - http://api2.rtve.es/rss/temas_noticias.xml
    - http://rss.cnn.com/rss/edition_world.rss
  delay: 5
`
	loader := Loader{strings.NewReader(config)}
	appConfig, _ := loader.LoadConfig()

	assert.Equal(t,
		AppConfig{
			RssNewsProvidersConfig{
				Sources: []string{
					"http://api2.rtve.es/rss/temas_noticias.xml",
					"http://rss.cnn.com/rss/edition_world.rss",
				},
				DelayInMinutes: 5,
			},
		}, *appConfig)
}

func TestConfigLoaderReturnsErrorOnBadConfig(t *testing.T) {
	config := `
rss_news_provider:
  sources
    http://api2.rtve.es/rss/temas_noticias.xml
  delay: 5
`
	loader := Loader{strings.NewReader(config)}
	appConfig, err := loader.LoadConfig()

	assert.Nil(t, appConfig)
	assert.NotNil(t, err)
}

func TestLoadConfigFuncUpdatesAppConfig(t *testing.T) {
	config := `
rss_news_provider:
  sources:
    - http://api2.rtve.es/rss/temas_noticias.xml
    - http://rss.cnn.com/rss/edition_world.rss
  delay: 5
`
	defaultLoader = Loader{strings.NewReader(config)}

	LoadConfig()

	assert.Equal(t,
		AppConfig{
			RssNewsProvidersConfig{
				Sources: []string{
					"http://api2.rtve.es/rss/temas_noticias.xml",
					"http://rss.cnn.com/rss/edition_world.rss",
				},
				DelayInMinutes: 5,
			},
		}, Current)
}

func TestLoadConfigFuncReturnsErrorOnBadConfig(t *testing.T) {
	config := `
rss_news_provider:
  sources
    http://api2.rtve.es/rss/temas_noticias.xml
  delay: 5
`
	defaultLoader = Loader{strings.NewReader(config)}

	err := LoadConfig()

	assert.NotNil(t, err)
}

func TestLoadConfigFuncNotifyChanges(t *testing.T) {
	config := `
rss_news_provider:
  sources:
    - http://api2.rtve.es/rss/temas_noticias.xml
    - http://rss.cnn.com/rss/edition_world.rss
  delay: 5
`
	defaultLoader = Loader{strings.NewReader(config)}

	var conf AppConfig
	go func() {
		conf = <-Subject
	}()

	LoadConfig()
	<-time.After(1 * time.Millisecond)
	assert.Equal(t,
		AppConfig{
			RssNewsProvidersConfig{
				Sources: []string{
					"http://api2.rtve.es/rss/temas_noticias.xml",
					"http://rss.cnn.com/rss/edition_world.rss",
				},
				DelayInMinutes: 5,
			},
		}, conf)
}
