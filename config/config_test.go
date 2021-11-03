package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestRawConfigLoaderReturnsConfig(t *testing.T) {
	loader := newRawConfigLoader([]byte(tsConfigs[0].yaml))
	appConfig, _ := loader()
	assert.Equal(t, tsConfigs[0].config, *appConfig)
}

func TestRawConfigLoaderReturnsErrorOnBadConfig(t *testing.T) {
	config := `
rss_news_provider:
  sources
    http://api2.rtve.es/rss/temas_noticias.xml
  delay: 5
`
	loader := newRawConfigLoader([]byte(config))
	appConfig, err := loader()

	assert.Nil(t, appConfig)
	assert.NotNil(t, err)
}

func TestFileConfigLoader(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	loader := newFileConfigLoader(file.Name())

	if _, err := file.Write([]byte(tsConfigs[0].yaml)); err != nil {
		panic(err)
	}

	appConfig, _ := loader()
	assert.Equal(t, tsConfigs[0].config, *appConfig)

	file.Truncate(0)
	file.Seek(0, 0)

	if _, err := file.Write([]byte(tsConfigs[1].yaml)); err != nil {
		panic(err)
	}

	appConfig, _ = loader()
	assert.Equal(t, tsConfigs[1].config, *appConfig)
}

func TestLoadConfigFuncUpdatesAppConfig(t *testing.T) {
	defaultLoader = newRawConfigLoader([]byte(tsConfigs[0].yaml))

	LoadConfig()

	assert.Equal(t, tsConfigs[0].config, Current)
}

func TestLoadConfigFuncReturnsErrorOnBadConfig(t *testing.T) {
	config := `
rss_news_provider:
  sources
    http://api2.rtve.es/rss/temas_noticias.xml
  delay: 5
`
	defaultLoader = newRawConfigLoader([]byte(config))
	err := LoadConfig()
	assert.NotNil(t, err)
}

func TestLoadConfigFuncNotifyChanges(t *testing.T) {
	var conf AppConfig
	defaultLoader = newRawConfigLoader([]byte(tsConfigs[0].yaml))
	go func() {
		conf = <-Subject
	}()

	LoadConfig()
	<-time.After(1 * time.Millisecond)
	assert.Equal(t, tsConfigs[0].config, conf)
}

var tsConfigs = []struct {
	yaml   string
	config AppConfig
}{
	{
		`
rss_news_provider:
  sources:
    - http://api2.rtve.es/rss/temas_noticias.xml
    - http://rss.cnn.com/rss/edition_world.rss
  delay: 5
`,
		AppConfig{
			RssNewsProvidersConfig{
				Sources: []string{
					"http://api2.rtve.es/rss/temas_noticias.xml",
					"http://rss.cnn.com/rss/edition_world.rss",
				},
				DelayInMinutes: 5,
			},
		},
	},
	{
		`
rss_news_provider:
  sources:
    - http://rss.cnn.com/rss/edition_world.rss
  delay: 1
`,
		AppConfig{
			RssNewsProvidersConfig{
				Sources: []string{
					"http://rss.cnn.com/rss/edition_world.rss",
				},
				DelayInMinutes: 1,
			},
		},
	},
}
