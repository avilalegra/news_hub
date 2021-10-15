package news

import (
	strip "github.com/grokify/html-strip-tags-go"
	"html"
	"log"
	"regexp"
	"strings"
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

func Search(keywords string) []*Preview {
	var matches []*Preview
	words := strings.Fields(keywords)
	for i, w := range words {
		words[i] = strings.ToLower(strings.Trim(w, ",.;"))
	}
	regx := regexp.MustCompile(strings.Join(words, " .*"))
	for _, p := range register {
		haystack := strip.StripTags(strings.ToLower(html.UnescapeString(p.Title + " " + p.Description)))
		if regx.MatchString(haystack) {
			matches = append(matches, &p)
		}
	}
	return matches
}

var register []Preview

func Load(preview ...Preview) {
	register = preview
}

func All() []Preview {
	return register
}

type Provider interface {
	RunAsync(chan<- Preview, chan<- error)
}

type Repository interface {
	Add(preview Preview)
	Search(keywords string) []Preview
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
