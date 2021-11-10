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

func TestFindByLinkIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	repo := newMongoRepo(Database, nil)
	loadDbFixtures()
	for _, tData := range tsFindByTitle {
		preview := repo.findByLink(tData.Link)
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

func TestLatestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RecreateDb()
	loadDbFixtures()
	repo := newMongoRepo(Database, nil)

	for i := 1; i <= len(news.Previews); i++ {
		previews := repo.FindLatest(i)
		assert.Len(t, previews, i)
		assert.Equal(t, news.Previews[0:i], previews)
	}
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
		"Guías nuevas medidas",
		0,
	},
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
		0,
	},
	{
		"directo, municipio",
		2,
	},
	{
		"Display; Support. PosTing",
		1,
	},
	{
		"<ul> <li>",
		0,
	},
	{
		"linux kernel",
		1,
	},
	{
		"linux kernel covid",
		0,
	},
	{
		"covid aforo",
		1,
	},
}

var tsFindByTitle = []struct {
	Link    string
	Preview *news.Preview
}{
	{
		`https://www.phoronix.com/scan.php?page=news_item&px=AMD-Cyan-Skillfish-DCN-2.01`,
		&news.Previews[0],
	},

	{
		`https://www.phoronix.com/scan.php?page=news_item&px=AMDGPU-DP-2.0-Linux-5.16`,
		&news.Previews[1],
	},
	{
		`https://www.phoronix.com/scan.php?page=news_item&px`, //test incomplete link
		nil,
	},
	{
		"not existing link",
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
