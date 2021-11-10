package persistence

import (
	"avilego.me/recent_news/news"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type mongoRepo struct {
	db              *mongo.Database
	prevCol         *mongo.Collection
	timeProvider    unixTimeProvider
	minMatchPercent int
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

func (r mongoRepo) FindRelated(searchExpr string) []news.Preview {
	var previews []news.Preview
	var related []news.Preview
	cursor, err := r.prevCol.Find(context.TODO(), bson.M{"$text": bson.M{"$search": searchExpr}})
	if err != nil {
		panic(err)
	}
	err = cursor.All(context.TODO(), &previews)
	if err != nil {
		panic(err)
	}

	for _, p := range previews {
		if p.MatchPercent(searchExpr) > r.minMatchPercent {
			related = append(related, p)
		}
	}
	return related
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

func (r mongoRepo) FindLatest(count int) []news.Preview {
	var previews []news.Preview
	opts := options.Find().SetSort(bson.D{{"pubtime", -1}}).SetLimit(int64(count))
	cursor, err := r.prevCol.Find(nil, bson.D{}, opts)
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
	return mongoRepo{database, database.Collection("news_previews"), timeProvider, 80}
}

func NewMongoKeeperFinder() news.KeeperFinder {
	return newMongoRepo(Database, defaultTimeProvider)
}

type unixTimeProvider func() int64

var defaultTimeProvider = func() int64 {
	return time.Now().Unix()
}
