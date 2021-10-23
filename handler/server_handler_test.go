package handler

import (
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/news"
	"avilego.me/recent_news/persistence"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiSearchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	persistence.RecreateDb()
	loadDbFixtures()
	server := httptest.NewServer(NewServerHttpHandler())
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/search?keywords=AMD")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
	assert.NotEqual(t, 0, resp.ContentLength)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func loadDbFixtures() {
	keeper := factory.Keeper()
	keeper.Store(news.Previews[0])
	keeper.Store(news.Previews[1])
}
