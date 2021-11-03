package config

import (
	"gopkg.in/yaml.v3"
	"io"
)

type AppConfig struct {
	RNPConfig RssNewsProvidersConfig `yaml:"rss_news_provider"`
}

type RssNewsProvidersConfig struct {
	Sources        []string `yaml:",flow"`
	DelayInMinutes int      `yaml:"delay"`
}

type Loader struct {
	Reader io.Reader
}

func (l Loader) LoadConfig() (*AppConfig, error) {
	raw, err := io.ReadAll(l.Reader)
	if err != nil {
		panic(err)
	}

	var appConfig AppConfig
	err = yaml.Unmarshal(raw, &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
