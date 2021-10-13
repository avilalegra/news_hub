package persistence

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	USER = "root"
	PWD  = "toor"
	HOST = "mongo"
	PORT = "27017"
)

var DbName = "todays_news"

var Client *mongo.Client

func init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", USER, PWD, HOST, PORT)))

	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	Client = client
}
