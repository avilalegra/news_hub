package news

import (
	"github.com/grokify/html-strip-tags-go"
	"html"
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
	Link        string
	Description string
	Source      *Source
}

type Provider interface {
	FetchNews() ([]Preview, error)
}

var register []Preview

func Update(providers ...Provider) (int, []error) {
	recentNews := make([]Preview, 0, len(providers))
	errors := make([]error, 0)
	resChan := make(chan []Preview, 1)
	errChan := make(chan []error, 1)

	for _, p := range providers {
		go func(prv Provider, resChan chan<- []Preview, errChan chan<- []error) {
			previews, err := prv.FetchNews()
			if err != nil {
				errChan <- errors
			} else {
				resChan <- previews
			}
		}(p, resChan, errChan)
	}

	for i := 0; i < len(providers); i++ {
		select {
		case previews := <-resChan:
			recentNews = append(recentNews, previews...)
		case errs := <-errChan:
			errors = append(errors, errs...)
		}
	}
	register = recentNews
	return len(recentNews), errors
}

func Load(preview ...Preview) {
	register = preview
}

func All() []Preview {
	return register
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
