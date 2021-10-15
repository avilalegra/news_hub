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

type Provider interface {
	RunAsync(chan<- Preview, chan<- error)
}

type Browser interface {
	Search(keywords string) []Preview
}

type PrevExistsErr struct {
	PreviewTitle string
}

func (e PrevExistsErr) Error() string {
	return fmt.Sprintf("existing preview with title %s", e.PreviewTitle)
}

type Repository interface {
	Browser
	Add(preview Preview) error
}

type Collector struct {
	Providers []Provider
	Repo      Repository
	logger    *log.Logger
}

func (c Collector) Run() {
	prvChan := make(chan Preview)
	errChan := make(chan error)

	for _, p := range c.Providers {
		p.RunAsync(prvChan, errChan)
	}

	go func() {
		for {
			select {
			case preview := <-prvChan:
				c.Repo.Add(preview)
			case err := <-errChan:
				c.logger.Print(err)
			}
		}
	}()
}
