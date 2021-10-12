package main

import (
	"avilego.me/news_hub/persistence"
	"context"
)

func main() {
	defer func() {
		if err := persistence.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
