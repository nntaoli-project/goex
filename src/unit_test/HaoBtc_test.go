package unit

import (
	"testing"
	"rest/haobtc"
	"net/http"
	"rest"
	"github.com/stretchr/testify/assert"
)

func Test_GetTicker(t *testing.T) {
	api := haobtc.New(http.DefaultClient, "", "");
	_, err := api.GetTicker(rest.LTC_CNY);
	assert.Error(t, err);

	ticker, err2 := api.GetTicker(rest.BTC_CNY);
	assert.NoError(t, err2);
	t.Log(ticker);
}

func Test_GetDepth(t *testing.T) {
	api := haobtc.New(http.DefaultClient, "", "");
	_, err := api.GetDepth(5, rest.LTC_CNY);
	assert.Error(t, err);

	depth, err2 := api.GetDepth(3, rest.BTC_CNY);
	assert.NoError(t, err2);
	t.Log(depth);
}