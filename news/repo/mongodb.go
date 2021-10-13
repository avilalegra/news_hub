package repo

import (
	"avilego.me/news_hub/news"
	"avilego.me/news_hub/persistence"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	Db *mongo.Database
}

var DefRepo MongoRepo

func (r *MongoRepo) Add(preview news.Preview) error {
	col := r.Db.Collection("news_previews")
	if _, err := col.InsertOne(context.TODO(), preview); err != nil {
		return err
	}
	return nil
}

func init() {
	DefRepo = MongoRepo{persistence.Database}
}
