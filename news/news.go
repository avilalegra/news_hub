package news

import (
	"context"
	"fmt"
	strip "github.com/grokify/html-strip-tags-go"
	"html"
	"log"
	"regexp"
	"strings"
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
	PubTime     int64
	RegUnixTime int64
}

// MatchPercent TODO: Refactor towards efficiency
func (p Preview) MatchPercent(searchExpr string) int {
	exprWords := splitWords(searchExpr)
	contWords := splitWords(strip.StripTags(html.UnescapeString(p.Title + " " + p.Description)))
	var matches []string

	for _, exprWord := range exprWords {
		regx := regexp.MustCompile(".*" + exprWord + ".*")
		for _, cword := range contWords {
			if regx.MatchString(cword) {
				matches = append(matches, exprWord)
				break
			}
		}
	}
	matchingPercent := 100 * len(matches) / len(exprWords)

	return matchingPercent
}

func splitWords(str string) (words []string) {
	words = strings.Fields(str)
	for i, w := range words {
		words[i] = strings.ToLower(strings.Trim(w, ",.;"))
	}
	return
}

type Provider interface {
	Provide(context.Context, chan<- Preview, chan<- error)
}

type Finder interface {
	FindRelated(searchExpr string) []Preview
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

func (c Cleaner) Run(ctx context.Context) {
	for running := true; running; {
		select {
		case <-c.Trigger:
			fmt.Println("running cleaner")
			limit := time.Now().Unix() - c.Ttl
			expired := c.KeeperFinder.FindBefore(limit)
			for _, preview := range expired {
				c.KeeperFinder.Remove(preview)
			}
		case <-ctx.Done():
			running = false
		}
	}
}
