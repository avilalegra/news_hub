package handler

import (
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/news"
	"avilego.me/recent_news/persistence"
	"github.com/stretchr/testify/assert"
	"io"
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
	content, _ := io.ReadAll(resp.Body)
	assert.Equal(t,
		`{"Count":2,"Data":{"Sources":[{"Title":"Phoronix","Link":"https://www.phoronix.com/","Description":"Linux Hardware Reviews \u0026 News","Language":"en-US"}],"Previews":[{"Title":"AMD Posts Code Enabling \"Cyan Skillfish\" Display Support Due To Different DCN2 Variant","Link":"https://www.phoronix.com/scan.php?page=news_item\u0026px=AMD-Cyan-Skillfish-DCN-2.01","Description":"Since July we've seen AMD open-source driver engineers posting code for \"Cyan Skillfish\" as an APU with Navi 1x graphics. While initial support for Cyan Skillfish was merged for Linux 5.15, it turns out the display code isn't yet wired up due to being a different DCN2 variant for its display block...","SourceLink":"https://www.phoronix.com/"},{"Title":"Linux 5.16 To Bring Initial DisplayPort 2.0 Support For AMD Radeon Driver (AMDGPU)","Link":"https://www.phoronix.com/scan.php?page=news_item\u0026px=AMDGPU-DP-2.0-Linux-5.16","Description":"A batch of feature updates was submitted today for DRM-Next of early feature work slated to come to the next version of the Linux kernel...","SourceLink":"https://www.phoronix.com/"}]}}`,
		string(content),
	)
}

func loadDbFixtures() {
	keeper := factory.Keeper()
	keeper.Store(news.Previews[0])
	keeper.Store(news.Previews[1])
}
