package persistence

import (
	_ "avilego.me/news_hub/env"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var Client *mongo.Client
var Database *mongo.Database

func init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MongoUri")))

	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	Client = client
	Database = Client.Database(os.Getenv("DbName"))
}
