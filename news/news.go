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
	FetchNews() ([]Extract, error)
}

var Register []Extract

func Update(providers ...NewsProvider) (int, []error) {
	recentNews := make([]Extract, 0, len(providers))
	errors := make([]error, 0)
	reschan := make(chan []Extract, 1)
	errchan := make(chan []error, 1)

	for _, p := range providers {
		go func(prv NewsProvider, reschan chan<- []Extract, errchan chan<- []error) {
			news, err := prv.FetchNews()
			if err != nil {
				errchan <- errors
			} else {
				reschan <- news
			}
		}(p, reschan, errchan)
	}

	for i := 0; i < len(providers); i++ {
		select {
		case news := <-reschan:
			recentNews = append(recentNews, news...)
		case errs := <-errchan:
			errors = append(errors, errs...)
		}
	}

	Register = recentNews

	return len(recentNews), errors
}

func Add(extract Extract) {
	Register = append(Register, extract)
}

func All() []Extract {
	return Register
}
