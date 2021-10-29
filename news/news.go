package news

import (
	"fmt"
	"log"
	"time"
)

type Source struct {
	Title       string
	Link        string
	Description string
	Language    string
}

type Preview struct {
	Title       string
	Link        string
	Description string
	Source      *Source
	RegUnixTime int64
}

type AsyncProvider interface {
	ProvideAsync(chan<- Preview, chan<- error)
}

type Finder interface {
	FindRelated(keywords string) []Preview
	FindBefore(unixTime int64) []Preview
}

type Keeper interface {
	Store(preview Preview)
	Remove(preview Preview)
}

type KeeperFinder interface {
	Keeper
	Finder
}

type Collector struct {
	Providers []AsyncProvider
	Keeper    Keeper
	Logger    *log.Logger
}

func (c Collector) Run() {
	prvChan := make(chan Preview)
	errChan := make(chan error)

	for _, p := range c.Providers {
		p.ProvideAsync(prvChan, errChan)
	}

	go func() {
		for {
			select {
			case preview := <-prvChan:
				c.Keeper.Store(preview)
			case err := <-errChan:
				c.Logger.Println(err)
			}
		}
	}()
}

//Cleaner ensures that expired news are removed
type Cleaner struct {
	KeeperFinder KeeperFinder
	Trigger      <-chan time.Time
	Ttl          int64
}

func (c Cleaner) Run() {
	go func() {
		for range c.Trigger {
			fmt.Println("running cleaner")
			limit := time.Now().Unix() - c.Ttl
			expired := c.KeeperFinder.FindBefore(limit)
			for _, preview := range expired {
				c.KeeperFinder.Remove(preview)
			}
		}
	}()
}
