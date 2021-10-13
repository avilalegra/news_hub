package news

import (
	"avilego.me/news_hub/persistence"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Add(preview Preview) error
}

type MongoRepo struct {
	db *mongo.Database
}

var DefRepo MongoRepo

func (r *MongoRepo) Add(preview Preview) error {
	col := r.db.Collection("news_previews")
	if _, err := col.InsertOne(context.TODO(), preview); err != nil {
		return err
	}
	return nil
}

func init() {
	DefRepo = MongoRepo{persistence.Database}
}
