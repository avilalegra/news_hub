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

const appConfFilePath = "config/app_config.yaml"

var Current AppConfig

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

func LoadConfig() error {
	conf, err := defaultLoader.LoadConfig()
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

var configChanges = make(chan AppConfig)

var Subject <-chan AppConfig = configChanges

var defaultLoader Loader

func newDefaultLoader() Loader {
	reader, err := os.Open(env.ProjDir() + "/" + appConfFilePath)
	if err != nil {
		panic(err)
	}
	return Loader{reader}
}

func init() {
	defaultLoader = newDefaultLoader()
	if err := LoadConfig(); err != nil {
		panic(err)
	}
	listenReloadSignal()
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
