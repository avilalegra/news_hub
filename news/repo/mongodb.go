package repo

import (
	"avilego.me/news_hub/news"
	"avilego.me/news_hub/persistence"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	Db      *mongo.Database
	prevCol *mongo.Collection
}

var DefRepo MongoRepo

func (r *MongoRepo) Add(preview news.Preview) error {
	if _, err := r.prevCol.InsertOne(context.TODO(), preview); err != nil {
		return err
	}
	return nil
}

func (r *MongoRepo) Search(keywords string) []news.Preview {
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

// TODO: Create title collection index
func (r *MongoRepo) findByTitle(title string) *news.Preview {
	var preview news.Preview
	result := r.prevCol.FindOne(context.TODO(), bson.M{"title": title})
	err := result.Decode(&preview)
	if err != nil {
		return nil
	}
	return &preview
}

func NewMongoRepo(database *mongo.Database) MongoRepo {
	return MongoRepo{database, database.Collection("news_previews")}
}

func init() {
	DefRepo = NewMongoRepo(persistence.Database)
}
