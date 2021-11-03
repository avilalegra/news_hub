package main

import (
	"avilego.me/recent_news/config"
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/handler"
	"avilego.me/recent_news/persistence"
	"context"
	"fmt"
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

	monitorConfigDependantServices()
	go factory.Cleaner().Run()

	fmt.Printf("App running at: %s\n", os.Getenv("ServerAddr"))
	fmt.Println("Mongo express running at: localhost:8081")

	log.Fatal(
		http.ListenAndServe(os.Getenv("ServerAddr"), handler.NewServerHttpHandler()),
	)
}

func monitorConfigDependantServices() {
	runServices := func(ctx context.Context) {
		go factory.Collector().Run(ctx)
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		runServices(ctx)
		for range config.Subject {
			cancel()
			runServices(ctx)
		}
	}()
}
