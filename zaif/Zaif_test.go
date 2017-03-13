package zaif

import (
	"github.com/nntaoli/crypto_coin_api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var api = New(http.DefaultClient, "", "")

func TestZaif_GetTicker(t *testing.T) {
	ticker, err := api.GetTicker(coinapi.BTC_JPY)
	assert.Empty(t, err)
	t.Log(ticker)
}

func TestZaif_GetDepth(t *testing.T) {
	depth, err := api.GetDepth(4, coinapi.BTC_JPY)
	assert.Empty(t, err)
	t.Log(depth)
}
