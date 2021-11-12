package config

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

func TestRawConfigParserShouldParseValidConfig(t *testing.T) {
	parser := newRawConfigParser([]byte(tsYamlConfigParsing.yaml))
	appConfig, err := parser()

	assert.Nil(t, err)
	assert.Equal(t, tsYamlConfigParsing.config, *appConfig)
}

func TestRawConfigParserReturnErrorWhenInvalidYaml(t *testing.T) {
	parser := newRawConfigParser([]byte(invalidYaml))
	appConfig, err := parser()

	assert.Nil(t, appConfig)
	assert.NotNil(t, err)
}

func TestFileConfigParserShouldParseValidConfig(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	parser := newFileConfigParser(file.Name())
	if _, err := file.Write([]byte(tsYamlConfigParsing.yaml)); err != nil {
		panic(err)
	}

	appConfig, err := parser()

	assert.Nil(t, err)
	assert.Equal(t, tsYamlConfigParsing.config, *appConfig)
}

func TestFileConfigParserReturnErrorWhenInvalidYaml(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())
	parser := newFileConfigParser(file.Name())

	if _, err := file.Write([]byte(invalidYaml)); err != nil {
		panic(err)
	}

	appConfig, err := parser()

	assert.Nil(t, appConfig)
	assert.NotNil(t, err)
}

func TestLoadConfigFuncUpdatesAppConfig(t *testing.T) {
	Current = appConfig{}
	defaultParser = func() (*appConfig, error) {
		return &validConfig, nil
	}
	err := LoadConfig()
	if err != nil {
		assert.Nil(t, err)
	}
	assert.Equal(t, validConfig, Current)
}

func TestLoadConfigFuncDoesntUpdatesAppConfigWhenInvalidConfig(t *testing.T) {
	Current = appConfig{}
	defaultParser = func() (*appConfig, error) {
		return &invalidConfig, nil
	}
	LoadConfig()
	assert.Equal(t, appConfig{}, Current)
}

func TestLoadConfigFuncNotifyChangesWhenValidConfig(t *testing.T) {
	Current = appConfig{}
	defaultParser = func() (*appConfig, error) {
		return &validConfig, nil
	}
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		conf := <-Subject
		assert.Equal(t, validConfig, conf)
		wg.Done()
	}()

	LoadConfig()
	wg.Wait()
}

func TestLoadConfigFuncDoesntNotifyChangesWhenInvalidConfig(t *testing.T) {
	Current = appConfig{}
	defaultParser = func() (*appConfig, error) {
		return &invalidConfig, nil
	}
	var wg sync.WaitGroup

	go func() {
		conf := <-Subject
		assert.Equal(t, appConfig{}, conf)
		wg.Done()
	}()

	LoadConfig()
	wg.Wait()
}

func TestLoadConfigFuncReturnErrorWhenInvalidYaml(t *testing.T) {
	Current = appConfig{}
	defaultParser = newRawConfigParser([]byte(invalidYaml))
	err := LoadConfig()
	assert.NotNil(t, err)
}

func TestLoadConfigFuncReturnErrorWhenInvalidConfig(t *testing.T) {
	Current = appConfig{}
	defaultParser = func() (*appConfig, error) {
		return &invalidConfig, nil
	}
	err := LoadConfig()
	assert.NotNil(t, err)
}

func TestRssNewsProviderConfigValidation(t *testing.T) {
	for i, tData := range tsRssNewsProvidersConfigErrorsValidation {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			err := tData.conf.validate()
			assert.Equal(t, tData.err, err)
		})
	}
}

func TestCleanerConfigValidation(t *testing.T) {
	for i, tData := range tsCleanerConfigErrorsValidation {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			err := tData.conf.validate()
			assert.Equal(t, tData.err, err)
		})
	}
}

func TestPositiveNumberTypeValidation(t *testing.T) {
	count := natNumber(-1)
	assert.Equal(t, errors.New("invalid config: should be a positive number"), count.validate())
}

func TestAppConfigValidation(t *testing.T) {
	for i, tData := range tsAppConfigValidation {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			err := tData.conf.validate()
			assert.Equal(t, tData.err, err)
		})
	}
}

var validConfig = appConfig{
	rssNewsProvidersConfig{
		Sources: []string{
			"http://api2.rtve.es/rss/temas_noticias.xml",
			"http://rss.cnn.com/rss/edition_world.rss",
		},
		MinutesPeriod: 5,
	},
	cleanerConfig{30, 10},
	10,
}

var invalidConfig = appConfig{}

var tsYamlConfigParsing = struct {
	yaml   string
	config appConfig
}{
	`
rss_news_provider:
  sources:
    - http://api2.rtve.es/rss/temas_noticias.xml
    - http://rss.cnn.com/rss/edition_world.rss
  period: 5

news_cleaner:
  ttl: 30
  period: 10

latest_news_count: 10
`,
	validConfig,
}

var invalidYaml = `
rss_news_provider:
  sources
    http://api2.rtve.es/rss/temas_noticias.xml
  period: 5
`

var tsAppConfigValidation = []struct {
	conf appConfig
	err  error
}{
	{
		appConfig{
			tsRssNewsProvidersConfigErrorsValidation[1].conf,
			validConfig.CleanerConfig,
			10,
		},
		tsRssNewsProvidersConfigErrorsValidation[1].err,
	},
	{
		appConfig{
			validConfig.RNPConfig,
			tsCleanerConfigErrorsValidation[1].conf,
			30,
		},
		tsCleanerConfigErrorsValidation[1].err,
	},
	{
		validConfig,
		nil,
	},
	{
		appConfig{
			validConfig.RNPConfig,
			validConfig.CleanerConfig,
			-1,
		},
		errors.New("invalid config: should be a positive number"),
	},
}

var tsRssNewsProvidersConfigErrorsValidation = []struct {
	conf rssNewsProvidersConfig
	err  error
}{
	{
		rssNewsProvidersConfig{
			Sources: []string{
				"http://rss.cnn.com/rss/edition_world.rss",
			},
			MinutesPeriod: 10,
		},
		nil,
	},
	{
		rssNewsProvidersConfig{
			Sources: []string{
				"http://rss.cnn.com/rss/edition_world.rss",
			},
			MinutesPeriod: 0,
		},
		errors.New("invalid rss provider config: period should be a positive number"),
	},
	{
		rssNewsProvidersConfig{
			Sources:       []string{},
			MinutesPeriod: 1,
		},
		errors.New("invalid rss provider config: at least one source required"),
	},
}

var tsCleanerConfigErrorsValidation = []struct {
	conf cleanerConfig
	err  error
}{
	{
		cleanerConfig{30, 2},
		nil,
	},
	{
		cleanerConfig{0, 2},
		errors.New("invalid cleaner config: ttl should be a positive number"),
	},
	{
		cleanerConfig{10, 30},
		errors.New("invalid cleaner config: ttl should be greater than period"),
	},
	{
		cleanerConfig{30, -10},
		errors.New("invalid cleaner config: period should be a positive number"),
	},
}
