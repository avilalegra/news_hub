package persistence

import (
	"avilego.me/recent_news/news"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestStorePersistDataIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	var prevs []news.Preview
	preview := news.Previews[0]
	prevCol := Database.Collection("news_previews")
	keeper := newMongoRepo(Database, newTimeProviderMock(preview.RegUnixTime))
	cursor, _ := prevCol.Find(context.TODO(), bson.M{})
	cursor.All(context.TODO(), &prevs)
	assert.Equal(t, 0, len(prevs))

	keeper.Store(news.Previews[0])

	cursor, _ = prevCol.Find(context.TODO(), bson.D{{}})
	cursor.All(context.TODO(), &prevs)
	assert.Equal(t, news.Previews[0], prevs[0])
}

func TestStoreDuplicatesReturnErrorIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	keeper := NewMongoKeeper()
	keeper.Store(news.Previews[1])
	err := keeper.Store(news.Previews[1])
	assert.ErrorIs(t, err, news.PrevExistsErr{PreviewTitle: news.Previews[1].Title})
}

func TestRegTimeSetOnStoringIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	var prevs []news.Preview
	prevCol := Database.Collection("news_previews")
	keeper := newMongoRepo(Database, newTimeProviderMock(123456))
	preview := news.Previews[2]

	keeper.Store(preview)

	cursor, _ := prevCol.Find(context.TODO(), bson.M{})
	cursor, _ = prevCol.Find(context.TODO(), bson.D{{}})
	cursor.All(context.TODO(), &prevs)

	assert.Equal(t, int64(123456), prevs[0].RegUnixTime)
}

func TestFindByTitleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	repo := newMongoRepo(Database, nil)
	prevCol := Database.Collection("news_previews")
	prevCol.InsertOne(context.TODO(), news.Previews[0])
	prevCol.InsertOne(context.TODO(), news.Previews[1])

	for _, tData := range tsFindByTitle {
		preview := repo.findByTitle(tData.Title)
		assert.Equal(t, tData.Preview, preview)
	}
}

func TestSearchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	finder := NewMongoFinder()
	prevCol := Database.Collection("news_previews")
	for _, preview := range news.Previews {
		prevCol.InsertOne(context.TODO(), preview)
	}
	for _, tData := range tsSearch {
		results := finder.Find(tData.keywords)
		assert.Equal(t, tData.count, len(results), tData.keywords)
	}
}

var newTimeProviderMock = func(time int64) unixTimeProvider {
	return func() int64 {
		return time
	}
}

var tsSearch = []struct {
	keywords string
	count    int
}{
	{
		"núcleos poblacionales",
		1,
	},
	{
		"Lava dirección confinados",
		1,
	},
	{
		"lava dirección confinados hierro",
		1,
	},
	{
		"directo, municipio",
		2,
	},
	{
		"Display; Support. PosTing",
		2,
	},
	{
		"linux kernel",
		2,
	},
	{
		"linux kernel covid",
		3,
	},
}

var tsFindByTitle = []struct {
	Title   string
	Preview *news.Preview
}{
	{
		`AMD Posts Code Enabling "Cyan Skillfish" Display Support Due To Different DCN2 Variant`,
		&news.Previews[0],
	},

	{
		`Linux 5.16 To Bring Initial DisplayPort 2.0 Support For AMD Radeon Driver (AMDGPU)`,
		&news.Previews[1],
	},
	{
		`Linux 5.16 To Bring Initial DisplayPort 2.0`,
		nil,
	},
	{
		"not existing title",
		nil,
	},
}
