package news

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

func Add(preview Preview) {
	register = append(register, preview)
}

func All() []Preview {
	return register
}
