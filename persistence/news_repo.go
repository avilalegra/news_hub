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

func (r mongoRepo) Store(preview news.Preview) {
	if prev := r.findByLink(preview.Link); prev != nil {
		return
	}
	preview.RegUnixTime = r.timeProvider()
	if _, err := r.prevCol.InsertOne(context.TODO(), preview); err != nil {
		panic(err)
	}
}

func (r mongoRepo) Remove(preview news.Preview) {
	if _, err := r.prevCol.DeleteOne(context.TODO(), bson.M{"link": preview.Link}); err != nil {
		panic(err)
	}
}

func (r mongoRepo) FindRelated(keywords string) []news.Preview {
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
	var previews []news.Preview
	cursor, err := r.prevCol.Find(context.TODO(), bson.M{"regunixtime": bson.M{"$lt": unixTime}})
	if err != nil {
		panic(err)
	}
	err = cursor.All(context.TODO(), &previews)
	if err != nil {
		panic(err)
	}
	return previews
}

func (r mongoRepo) findByLink(link string) *news.Preview {
	var preview news.Preview
	result := r.prevCol.FindOne(context.TODO(), bson.M{"link": link})
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

func NewMongoKeeperFinder() news.KeeperFinder {
	return newMongoRepo(Database, defaultTimeProvider)
}

type unixTimeProvider func() int64

var defaultTimeProvider = func() int64 {
	return time.Now().Unix()
}
