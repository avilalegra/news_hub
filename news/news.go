package news

import (
	strip "github.com/grokify/html-strip-tags-go"
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
	Run(chan<- Preview)
}
