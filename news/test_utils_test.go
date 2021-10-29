package news

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinderFakeFindRelated(t *testing.T) {
	fake := KeeperFinderFake{Previews: Previews}
	for _, tData := range searchTestData {
		results := fake.FindRelated(tData.keywords)
		assert.Equal(t, tData.count, len(results), tData.keywords)
	}
}

func TestFinderFakeFindBefore(t *testing.T) {
	fake := KeeperFinderFake{Previews: Previews}
	actual := fake.FindBefore(int64(456))
	expected := []Preview{Previews[0], Previews[2], Previews[3]}
	assert.Equal(t, expected, actual)
}

func TestKeeperFakeStore(t *testing.T) {
	fake := KeeperFinderFake{}
	fake.Store(Previews[0])
	assert.Equal(t, Previews[0:1], fake.Previews)
}

func TestKeeperFakeRemove(t *testing.T) {
	fake := KeeperFinderFake{Previews: Previews}
	fake.Remove(Previews[0])
	assert.Equal(t, Previews[1:], fake.Previews)
}

var searchTestData = []struct {
	keywords string
	count    int
}{
	{
		"núcleos poblacionales",
		1,
	},
	{
		"directo, municipio",
		2,
	},
	{
		"municipio, directo",
		2,
	},
	{
		"",
		0,
	},
}
