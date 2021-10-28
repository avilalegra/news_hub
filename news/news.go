package news

import (
	"fmt"
	"log"
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
	Find(keywords string) []Preview
	FindBefore(unixTime int64) []Preview
}

type Keeper interface {
	Store(preview Preview) error
}

type PrevExistsErr struct {
	PreviewTitle string
}

func (e PrevExistsErr) Error() string {
	return fmt.Sprintf("existing preview with title %s", e.PreviewTitle)
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
				if err := c.Keeper.Store(preview); err != nil {
					c.Logger.Println(err)
				}
			case err := <-errChan:
				c.Logger.Println(err)
			}
		}
	}()
}
