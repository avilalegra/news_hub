package persistence

import (
	_ "avilego.me/news_hub/env"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var Client *mongo.Client
var Database *mongo.Database

func init() {
	setClient()
	RecreateDb()
}

func setClient() {
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
}

func RecreateDb() {
	Database = Client.Database(os.Getenv("DbName"))
	err := Database.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
	ensureIndexes()
}

func ensureIndexes() {
	collection := Database.Collection("news_previews")
	index := mongo.IndexModel{
		Keys: bson.D{
			{"title", "text"},
			{"description", "text"},
		},
		Options: options.Index().SetWeights(bson.D{
			{"title", 9},
			{"description", 3},
		}),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), index)

	if err != nil {
		panic(err)
	}
}
