package config

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
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

func TestParseBadConfig(t *testing.T) {
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
