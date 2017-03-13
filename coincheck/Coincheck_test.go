package coincheck

import (
	"github.com/nntaoli/crypto_coin_api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var api = New(http.DefaultClient, "", "")

func TestCoincheck_GetTicker(t *testing.T) {
	ticker, err := api.GetTicker(coinapi.BTC_JPY)
	assert.NoError(t, err)
	t.Log(ticker)
}

func TestCoincheck_GetDepth(t *testing.T) {
	depth, err := api.GetDepth(3, coinapi.BTC_JPY)
	assert.NoError(t, err)
	t.Log(depth)
}
