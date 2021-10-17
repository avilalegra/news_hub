package persistence

import (
	"avilego.me/news_hub/news"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	Db      *mongo.Database
	prevCol *mongo.Collection
}

var Instance MongoRepo

func (r MongoRepo) Add(preview news.Preview) error {
	if prev := r.findByTitle(preview.Title); prev != nil {
		return news.PrevExistsErr{PreviewTitle: prev.Title}
	}
	if _, err := r.prevCol.InsertOne(context.TODO(), preview); err != nil {
		panic(err)
	}

	return nil
}

func (r MongoRepo) Search(keywords string) []news.Preview {
	var previews []news.Preview
	cursor, err := r.prevCol.Find(context.TODO(), bson.M{"$text": bson.M{"$search": keywords}})
	if err != nil {
		panic(err)
	}
	err = cursor.All(context.TODO(), &previews)
	if err != nil {
		panic(err)
	}
	return previews
}

func (r MongoRepo) findByTitle(title string) *news.Preview {
	var preview news.Preview
	result := r.prevCol.FindOne(context.TODO(), bson.M{"title": title})
	if result.Err() == mongo.ErrNoDocuments {
		return nil
	}
	err := result.Decode(&preview)
	if err != nil {
		panic(err)
	}
	return &preview
}

func NewMongoRepo(database *mongo.Database) MongoRepo {
	return MongoRepo{database, database.Collection("news_previews")}
}

func init() {
	Instance = NewMongoRepo(Database)
}
