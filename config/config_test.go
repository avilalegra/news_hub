package config

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestRawConfigParserWithValidConfig(t *testing.T) {
	for i, tData := range validConfigs {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			parser := newRawConfigParser([]byte(tData.yaml))
			appConfig, _ := parser()
			assert.Equal(t, tData.config, *appConfig)
		})
	}
}

func TestRawConfigParserWithInvalidYaml(t *testing.T) {
	for i, conf := range invalidYaml {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			parser := newRawConfigParser([]byte(conf))
			appConfig, err := parser()

			assert.Nil(t, appConfig)
			assert.NotNil(t, err)
		})
	}
}

func TestFileConfigParserWithValidConfig(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	for _, tData := range validConfigs {
		file.Truncate(0)
		file.Seek(0, 0)

		if _, err := file.Write([]byte(tData.yaml)); err != nil {
			panic(err)
		}

		parser := newFileConfigParser(file.Name())
		appConfig, _ := parser()
		assert.Equal(t, tData.config, *appConfig)
	}
}

func TestFileConfigParserWithInvalidYaml(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	for _, conf := range invalidYaml {
		file.Truncate(0)
		file.Seek(0, 0)

		if _, err := file.Write([]byte(conf)); err != nil {
			panic(err)
		}

		parser := newFileConfigParser(file.Name())
		appConfig, err := parser()

		assert.Nil(t, appConfig)
		assert.NotNil(t, err)
	}
}

func TestLoadConfigFuncUpdatesAppConfig(t *testing.T) {
	for i, tData := range validConfigs {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			defaultParser = newRawConfigParser([]byte(tData.yaml))
			LoadConfig()
			assert.Equal(t, tData.config, Current)
		})
	}
}

func TestLoadConfigFuncReturnsErrorOnInvalidYaml(t *testing.T) {
	for i, conf := range invalidYaml {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			defaultParser = newRawConfigParser([]byte(conf))
			err := LoadConfig()
			assert.NotNil(t, err)
		})
	}
}

func TestLoadConfigFuncNotifyChanges(t *testing.T) {
	var conf AppConfig
	defaultParser = newRawConfigParser([]byte(validConfigs[0].yaml))
	go func() {
		conf = <-Subject
	}()

	LoadConfig()
	<-time.After(1 * time.Millisecond)
	assert.Equal(t, validConfigs[0].config, conf)
}

func TestRssNewsProviderConfigValidation(t *testing.T) {
	for i, tData := range invalidRssNewsProvidersConfig {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			err := tData.conf.validate()
			assert.Equal(t, tData.err, err)
		})
	}
}

var validConfigs = []struct {
	yaml   string
	config AppConfig
}{
	{
		`
rss_news_provider:
  sources:
    - http://api2.rtve.es/rss/temas_noticias.xml
    - http://rss.cnn.com/rss/edition_world.rss
  period: 5
`,
		AppConfig{
			RssNewsProvidersConfig{
				Sources: []string{
					"http://api2.rtve.es/rss/temas_noticias.xml",
					"http://rss.cnn.com/rss/edition_world.rss",
				},
				MinutesPeriod: 5,
			},
		},
	},
	{
		`
rss_news_provider:
  sources:
    - http://rss.cnn.com/rss/edition_world.rss
  period: 1
`,
		AppConfig{
			RssNewsProvidersConfig{
				Sources: []string{
					"http://rss.cnn.com/rss/edition_world.rss",
				},
				MinutesPeriod: 1,
			},
		},
	},
}

var invalidYaml = []string{
	`
rss_news_provider:
  sources
    http://api2.rtve.es/rss/temas_noticias.xml
  period: 5
`,
}

var invalidRssNewsProvidersConfig = []struct {
	conf RssNewsProvidersConfig
	err  error
}{
	{
		RssNewsProvidersConfig{
			Sources: []string{
				"http://rss.cnn.com/rss/edition_world.rss",
			},
			MinutesPeriod: 0,
		},
		errors.New("invalid rss provider config: minutes period should be a positive number"),
	},
	{
		RssNewsProvidersConfig{
			Sources:       []string{},
			MinutesPeriod: 1,
		},
		errors.New("invalid rss provider config: at least one source required"),
	},
}
