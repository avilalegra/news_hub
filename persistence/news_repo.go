package persistence

import (
	"avilego.me/recent_news/news"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type mongoRepo struct {
	Db           *mongo.Database
	prevCol      *mongo.Collection
	timeProvider unixTimeProvider
}

func (r mongoRepo) Store(preview news.Preview) error {
	if prev := r.findByTitle(preview.Title); prev != nil {
		return news.PrevExistsErr{PreviewTitle: prev.Title}
	}
	preview.RegUnixTime = r.timeProvider()
	if _, err := r.prevCol.InsertOne(context.TODO(), preview); err != nil {
		panic(err)
	}

	return nil
}

func (r mongoRepo) Find(keywords string) []news.Preview {
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

func (r mongoRepo) FindBefore(unixTime int64) []news.Preview {
	panic("implement me")
}

func (r mongoRepo) findByTitle(title string) *news.Preview {
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

func newMongoRepo(database *mongo.Database, timeProvider unixTimeProvider) mongoRepo {
	return mongoRepo{database, database.Collection("news_previews"), timeProvider}
}

func NewMongoKeeper() news.Keeper {
	return newMongoRepo(Database, defaultTimeProvider)
}

func NewMongoFinder() news.Finder {
	return newMongoRepo(Database, defaultTimeProvider)
}

type unixTimeProvider func() int64

var defaultTimeProvider = func() int64 {
	return time.Now().Unix()
}
