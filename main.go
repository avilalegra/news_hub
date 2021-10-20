package main

import (
	"avilego.me/recent_news/persistence"
	"context"
	"time"
)

func main() {
	defer func() {
		if err := persistence.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collector().Run()

	// This is just for demo purposes
	// Be sure to start containers: docker-compose up
	// Run the project go run main.go
	// Go to localhost:8081 and look for bd recent_news_test
	// it would be populated every one minute with news
	// from rss sources (see container.go).
	// Keep in mind that you have to wait at least one minute when you run it
	// and that if you run the tests the test db is overwritten
	// Another type of news providers wil be included in the future
	<-time.After(time.Minute * 10)
}
