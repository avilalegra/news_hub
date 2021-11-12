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

type appConfig struct {
	RNPConfig       rssNewsProvidersConfig `yaml:"rss_news_provider"`
	CleanerConfig   cleanerConfig          `yaml:"news_cleaner"`
	LatestNewsCount natNumber              `yaml:"latest_news_count"`
}

func (c appConfig) validate() (validationError error) {
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

type rssNewsProvidersConfig struct {
	Sources       []string  `yaml:",flow"`
	MinutesPeriod natNumber `yaml:"period"`
}

func (c rssNewsProvidersConfig) validate() error {
	if c.MinutesPeriod.validate() != nil {
		return errors.New("invalid rss provider config: period should be a positive number")
	}
	if len(c.Sources) == 0 {
		return errors.New("invalid rss provider config: at least one source required")
	}
	return nil
}

type cleanerConfig struct {
	Ttl           natNumber `yaml:"ttl"`
	MinutesPeriod natNumber `yaml:"period"`
}

func (c cleanerConfig) validate() error {
	if c.Ttl.validate() != nil {
		return errors.New("invalid cleaner config: ttl should be a positive number")
	}
	if c.Ttl < c.MinutesPeriod {
		return errors.New("invalid cleaner config: ttl should be greater than period")
	}
	if c.MinutesPeriod.validate() != nil {
		return errors.New("invalid cleaner config: period should be a positive number")
	}
	return nil
}

type natNumber int

func (n natNumber) validate() error {
	if n <= 0 {
		return errors.New("invalid config: should be a positive number")
	}
	return nil
}

type parser func() (*appConfig, error)

var defaultParser parser

func newRawConfigParser(raw []byte) parser {
	return func() (*appConfig, error) {
		var appConfig appConfig
		err := yaml.Unmarshal(raw, &appConfig)
		if err != nil {
			return nil, err
		}

		return &appConfig, nil
	}
}

func newFileConfigParser(filePath string) parser {
	return func() (*appConfig, error) {
		reader, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		raw, _ := io.ReadAll(reader)
		return newRawConfigParser(raw)()
	}
}

var configChanges = make(chan appConfig)

var Subject <-chan appConfig = configChanges

var Current appConfig

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
