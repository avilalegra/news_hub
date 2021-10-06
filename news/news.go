package news

type Source struct {
	Title       string
	Link        string
	Description string
	Language    string
}

type Extract struct {
	Title       string
	Link        string
	Description string
	Source      *Source
}

type NewsProvider interface {
	FetchNews() ([]*Extract, error)
}
