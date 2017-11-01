package coincheck

import (
	. "github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var api = New(http.DefaultClient, "", "")

func TestCoincheck_GetTicker(t *testing.T) {
	ticker, err := api.GetTicker(CurrencyPair{BTC, Currency{"JPY", ""}})
	assert.NoError(t, err)
	t.Log(ticker)
}

func TestCoincheck_GetDepth(t *testing.T) {
	depth, err := api.GetDepth(3, CurrencyPair{BTC, NewCurrency("JPY", "")})
	assert.NoError(t, err)
	t.Log(depth)
}
