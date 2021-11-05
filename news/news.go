package news

import (
	"context"
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

type Provider interface {
	Provide(context.Context, chan<- Preview, chan<- error)
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
	Providers []Provider
	Keeper    Keeper
	Logger    *log.Logger
}

func (c Collector) Run(ctx context.Context) {
	prvChan := make(chan Preview)
	errChan := make(chan error)

	for _, p := range c.Providers {
		go p.Provide(ctx, prvChan, errChan)
	}

	for running := true; running; {
		select {
		case preview := <-prvChan:
			c.Keeper.Store(preview)
		case err := <-errChan:
			c.Logger.Println(err)
		case <-ctx.Done():
			running = false
		}
	}
}

//Cleaner ensures that expired news are removed
type Cleaner struct {
	KeeperFinder KeeperFinder
	Trigger      <-chan time.Time
	Ttl          int64
}

func (c Cleaner) Run() {
	for range c.Trigger {
		fmt.Println("running cleaner")
		limit := time.Now().Unix() - c.Ttl
		expired := c.KeeperFinder.FindBefore(limit)
		for _, preview := range expired {
			c.KeeperFinder.Remove(preview)
		}
	}
}
