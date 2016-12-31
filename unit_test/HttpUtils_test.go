package unit

import (
	"testing"
	"net/http"
	"github.com/nntaoli/crypto_coin_api"
	"github.com/stretchr/testify/assert"
	"net/url"
)

func Test_HttpGet(t *testing.T) {
	reqUrl := "https://api.huobi.com/staticmarket/ticker_btc_json.js"
	respMap, err := coinapi.HttpGet(http.DefaultClient, reqUrl)
	assert.Empty(t, err)
	t.Log(respMap)
}

func Test_HttpPost(t *testing.T) {
	reqUrl := "https://api.huobi.com/staticmarket/ticker_btc_json.js"
	resp, err := coinapi.HttpPostForm(http.DefaultClient, reqUrl, url.Values{})
	assert.Empty(t, err)
	t.Log(string(resp))
}
