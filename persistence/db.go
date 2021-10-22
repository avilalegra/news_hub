package persistence

import (
	_ "avilego.me/recent_news/env"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"os"
	"time"
)

var Client *mongo.Client
var Database *mongo.Database

func init() {
	setClient()
	setDb()
	if !dbExists() {
		RecreateDb()
	}
}

func dbExists() bool {
	dbNames, err := Client.ListDatabaseNames(context.TODO(), bson.M{"name": os.Getenv("DbName")})
	if err != nil {
		panic(err)
	}
	return len(dbNames) == 1
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

func setDb() {
	Database = Client.Database(os.Getenv("DbName"))
}

func RecreateDb() {
	err := Database.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
	ensureIndexes()
}

func ensureIndexes() {
	collection := Database.Collection("news_previews")
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"title", bsonx.String("text")},
				{"description", bsonx.String("text")},
			},
			Options: options.Index().SetWeights(bson.D{
				{"title", 9},
				{"description", 3},
			}),
		},
	}
	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		panic(err)
	}
}
