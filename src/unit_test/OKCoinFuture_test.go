package unit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"rest/okcoin"
	"net/http"
	"rest"
)

func Test_OKCoin_Future_GetTicker(t *testing.T) {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	ticker, err := api.GetFutureTicker(rest.BTC_USD, "this_week");
	assert.NoError(t, err);
	t.Log(ticker)
}

func Test_OKCoin_Future_GetDepth(t *testing.T) {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	depth, err := api.GetFutureDepth(rest.BTC_USD, "this_week");
	assert.NoError(t, err);
	t.Log(depth)
}

