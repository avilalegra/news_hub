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
	Link        string ``
	Description string
	Source      *Source
}

type AsyncProvider interface {
	ProvideAsync(chan<- Preview, chan<- error)
}

type Finder interface {
	Find(keywords string) []Preview
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
				err := c.Keeper.Store(preview)
				if err == nil {
					c.Logger.Printf("news preview added: %s\n", preview.Title)
				}
			case err := <-errChan:
				c.Logger.Printf("provider error: %s\n", err)
			}
		}
	}()
}
