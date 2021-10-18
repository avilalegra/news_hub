package main

import (
	"avilego.me/recent_news/container"
	"avilego.me/recent_news/persistence"
	"context"
)

func main() {
	defer func() {
		if err := persistence.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	container.GetCollector().Run()
}
