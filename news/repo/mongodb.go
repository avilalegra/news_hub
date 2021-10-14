package repo

import (
	"avilego.me/news_hub/news"
	"avilego.me/news_hub/persistence"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	Db      *mongo.Database
	newsCol *mongo.Collection
}

var DefRepo MongoRepo

func (r *MongoRepo) Add(preview news.Preview) error {
	if _, err := r.newsCol.InsertOne(context.TODO(), preview); err != nil {
		return err
	}
	return nil
}

func NewMongoRepo(database *mongo.Database) MongoRepo {
	return MongoRepo{database, database.Collection("news_previews")}
}

func init() {
	DefRepo = NewMongoRepo(persistence.Database)
}
