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
	preview := news.Previews[0]
	keeper := newMongoRepo(Database, newTimeProviderMock(preview.RegUnixTime))

	keeper.Store(news.Previews[0])
	expects := getAllStoredPreviews()
	assert.Equal(t, 1, len(expects))
	assert.Equal(t, news.Previews[0], expects[0])
}

func TestStoreDuplicatesIgnoredIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	preview := news.Previews[0]
	keeper := newMongoRepo(Database, newTimeProviderMock(preview.RegUnixTime))

	keeper.Store(preview)
	keeper.Store(preview)

	assert.Equal(t, news.Previews[:1], getAllStoredPreviews())
}

func TestRegTimeSetOnStoringIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	preview := news.Previews[2]
	keeper := newMongoRepo(Database, newTimeProviderMock(123456))

	keeper.Store(preview)
	expected := getAllStoredPreviews()[0]
	assert.Equal(t, int64(123456), expected.RegUnixTime)
}

func TestFindByTitleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	repo := newMongoRepo(Database, nil)
	loadDbFixtures()
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
	repo := newMongoRepo(Database, nil)
	loadDbFixtures()
	for _, tData := range tsSearch {
		results := repo.FindRelated(tData.keywords)
		assert.Equal(t, tData.count, len(results), tData.keywords)
	}
}

func TestFindBeforeIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	loadDbFixtures()
	repo := newMongoRepo(Database, nil)
	for _, tData := range tsFindBefore {
		previews := repo.FindBefore(tData.unixTime)
		assert.Equal(t, tData.previews, previews)
	}
}

func TestRemoveIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	loadDbFixtures()
	repo := newMongoRepo(Database, nil)

	repo.Remove(news.Previews[0])
	assert.Equal(t, news.Previews[1:], getAllStoredPreviews())

	repo.Remove(news.Preview{Link: "news link"})
	assert.Equal(t, news.Previews[1:], getAllStoredPreviews())
}

func getAllStoredPreviews() (prevs []news.Preview) {
	prevCol := Database.Collection("news_previews")
	cursor, _ := prevCol.Find(context.TODO(), bson.M{})
	err := cursor.All(context.TODO(), &prevs)
	if err != nil {
		panic(err)
	}
	return
}

func loadDbFixtures() {
	repo := newMongoRepo(Database, nil)
	for _, preview := range news.Previews {
		repo.timeProvider = newTimeProviderMock(preview.RegUnixTime)
		repo.Store(preview)
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
		`Linux 5.16 To Bring Initial`, //test incomplete title
		nil,
	},
	{
		"not existing title",
		nil,
	},
}

var tsFindBefore = []struct {
	unixTime int64
	previews []news.Preview
}{
	{
		int64(123),
		news.Previews[2:4],
	},
	{
		int64(8910),
		news.Previews,
	},
	{
		int64(0),
		nil,
	},
}
