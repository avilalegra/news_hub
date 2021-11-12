package config

import (
	"avilego.me/recent_news/env"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appConfFilePath = env.ProjDir() + "/config/app_config.yaml"

type AppConfig struct {
	RNPConfig       RssNewsProvidersConfig `yaml:"rss_news_provider"`
	CleanerConfig   CleanerConfig          `yaml:"news_cleaner"`
	LatestNewsCount positiveNumber         `yaml:"latest_news_count"`
}

func (c AppConfig) validate() (validationError error) {
	defer func() {
		if err := recover(); err != nil {
			validationError = err.(error)
		}
	}()
	mustValidate := func(validator func() error) {
		if err := validator(); err != nil {
			panic(err)
		}
	}

	mustValidate(c.RNPConfig.validate)
	mustValidate(c.CleanerConfig.validate)
	mustValidate(c.LatestNewsCount.validate)

	if validationError != nil {
		return validationError
	}
	return
}

type RssNewsProvidersConfig struct {
	Sources       []string `yaml:",flow"`
	MinutesPeriod int      `yaml:"period"`
}

func (c RssNewsProvidersConfig) validate() error {
	if c.MinutesPeriod <= 0 {
		return errors.New("invalid rss provider config: period should be a positive number")
	}
	if len(c.Sources) == 0 {
		return errors.New("invalid rss provider config: at least one source required")
	}
	return nil
}

type CleanerConfig struct {
	Ttl           int `yaml:"ttl"`
	MinutesPeriod int `yaml:"period"`
}

func (c CleanerConfig) validate() error {
	if c.Ttl <= 0 {
		return errors.New("invalid cleaner config: ttl should be a positive number")
	}
	if c.Ttl < c.MinutesPeriod {
		return errors.New("invalid cleaner config: ttl should be greater than period")
	}
	if c.MinutesPeriod <= 0 {
		return errors.New("invalid cleaner config: period should be a positive number")
	}
	return nil
}

type positiveNumber int

func (n positiveNumber) validate() error {
	if n <= 0 {
		return errors.New("invalid config: should be a positive number")
	}
	return nil
}

type parser func() (*AppConfig, error)

var defaultParser parser

func newRawConfigParser(raw []byte) parser {
	return func() (*AppConfig, error) {
		var appConfig AppConfig
		err := yaml.Unmarshal(raw, &appConfig)
		if err != nil {
			return nil, err
		}

		return &appConfig, nil
	}
}

func newFileConfigParser(filePath string) parser {
	return func() (*AppConfig, error) {
		reader, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		raw, _ := io.ReadAll(reader)
		return newRawConfigParser(raw)()
	}
}

var configChanges = make(chan AppConfig)

var Subject <-chan AppConfig = configChanges

var Current AppConfig

func LoadConfig() error {
	conf, err := defaultParser()
	if err != nil {
		return err
	}

	if err = conf.validate(); err != nil {
		return err
	}

	Current = *conf

	select {
	case configChanges <- Current:
	case <-time.After(10 * time.Millisecond):
	}

	return nil
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
	defaultParser = newFileConfigParser(appConfFilePath)
	if err := LoadConfig(); err != nil {
		panic(err)
	}
	listenReloadSignal()
}
