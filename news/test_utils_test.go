package news

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeperFakeStore(t *testing.T) {
	fake := KeeperFake{}
	fake.Store(Previews[0])
	assert.Equal(t, Previews[0:1], fake.Previews)
}

func TestKeeperFakeRemove(t *testing.T) {
	fake := KeeperFake{Previews: Previews}
	fake.Remove(Previews[0])
	assert.Equal(t, Previews[1:], fake.Previews)
}
