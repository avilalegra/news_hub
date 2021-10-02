package http

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	handler := http.HandlerFunc(Search)
	params := url.Values{}
	params.Add("expr", "something")

	assert.HTTPBodyContains(t, handler, "GET", "/search", params, "hihi")
}
