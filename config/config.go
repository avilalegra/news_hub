package config

import (
	"avilego.me/recent_news/env"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appConfFilePath = env.ProjDir() + "/config/app_config.yaml"

var Current AppConfig

type AppConfig struct {
	RNPConfig RssNewsProvidersConfig `yaml:"rss_news_provider"`
}

type RssNewsProvidersConfig struct {
	Sources        []string `yaml:",flow"`
	DelayInMinutes int      `yaml:"delay"`
}

type loader func() (*AppConfig, error)

var defaultLoader loader

var configChanges = make(chan AppConfig)

var Subject <-chan AppConfig = configChanges

func LoadConfig() error {
	conf, err := defaultLoader()
	if err != nil {
		return err
	}

	Current = *conf

	select {
	case configChanges <- Current:
	case <-time.After(10 * time.Millisecond):
	}

	return nil
}

func newRawConfigLoader(raw []byte) loader {
	return func() (*AppConfig, error) {
		var appConfig AppConfig
		err := yaml.Unmarshal(raw, &appConfig)
		if err != nil {
			return nil, err
		}

		return &appConfig, nil
	}
}

func newFileConfigLoader(filePath string) loader {
	return func() (*AppConfig, error) {
		reader, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		raw, _ := io.ReadAll(reader)
		return newRawConfigLoader(raw)()
	}
}

func listenReloadSignal() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR1)
	go func() {
		for {
			<-s
			if err := LoadConfig(); err != nil {
				log.Println(err)
			} else {
				log.Println("app config reloaded")
			}
		}
	}()
}

func init() {
	defaultLoader = newFileConfigLoader(appConfFilePath)
	if err := LoadConfig(); err != nil {
		panic(err)
	}
	listenReloadSignal()
}
