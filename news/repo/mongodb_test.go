package repo

import (
	"avilego.me/news_hub/news"
	"avilego.me/news_hub/persistence"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestAdd(t *testing.T) {
	err := persistence.Database.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
	var prevs []news.Preview
	prevCol := persistence.Database.Collection("news_previews")

	cursor, err := prevCol.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	err = cursor.All(context.TODO(), &prevs)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, 0, len(prevs))

	DefRepo.Add(previews[0])
	DefRepo.Add(previews[1])

	cursor, err = prevCol.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	err = cursor.All(context.TODO(), &prevs)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, previews[0:2], prevs)
}

var sources = map[string]*news.Source{
	"phoronix": {
		Title:       `Phoronix`,
		Link:        `https://www.phoronix.com/`,
		Language:    `en-US`,
		Description: `Linux Hardware Reviews & News`,
	},
	"rtve": {
		Title:       `Noticias en rtve.es`,
		Link:        `http://www.rtve.es`,
		Description: `RSS Tags`,
	},
}

func TestFindByTitle(t *testing.T) {
	persistence.Database.Drop(context.TODO())
	prevCol := persistence.Database.Collection("news_previews")
	prevCol.InsertOne(context.TODO(), previews[0])
	prevCol.InsertOne(context.TODO(), previews[1])

	for _, tData := range tsFindByTitle {
		preview := DefRepo.findByTitle(tData.Title)
		assert.Equal(t, tData.Preview, preview)
	}
}

var tsFindByTitle = []struct {
	Title   string
	Preview *news.Preview
}{
	{
		`AMD Posts Code Enabling "Cyan Skillfish" Display Support Due To Different DCN2 Variant`,
		&previews[0],
	},

	{
		`Linux 5.16 To Bring Initial DisplayPort 2.0 Support For AMD Radeon Driver (AMDGPU)`,
		&previews[1],
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

var previews = []news.Preview{
	{
		Title:       `AMD Posts Code Enabling "Cyan Skillfish" Display Support Due To Different DCN2 Variant`,
		Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMD-Cyan-Skillfish-DCN-2.01`,
		Description: `Since July we've seen AMD open-source driver engineers posting code for "Cyan Skillfish" as an APU with Navi 1x graphics. While initial support for Cyan Skillfish was merged for Linux 5.15, it turns out the display code isn't yet wired up due to being a different DCN2 variant for its display block...`,
		Source:      sources["phoronix"],
	},
	{
		Title:       `Linux 5.16 To Bring Initial DisplayPort 2.0 Support For AMD Radeon Driver (AMDGPU)`,
		Link:        `https://www.phoronix.com/scan.php?page=news_item&px=AMDGPU-DP-2.0-Linux-5.16`,
		Description: `A batch of feature updates was submitted today for DRM-Next of early feature work slated to come to the next version of the Linux kernel...`,
		Source:      sources["phoronix"],
	},
}
