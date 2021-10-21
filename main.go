package main

import (
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/persistence"
	"avilego.me/recent_news/server"
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	defer func() {
		if err := persistence.Client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	factory.Collector().Run()

	log.Fatal(
		http.ListenAndServe(os.Getenv("ServerAddr"), server.NewServerHttpHandler()),
	)
}
