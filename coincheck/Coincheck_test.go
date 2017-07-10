package coincheck

import (
	. "github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var api = New(http.DefaultClient, "", "")

func TestCoincheck_GetTicker(t *testing.T) {
	ticker, err := api.GetTicker(BTC_JPY)
	assert.NoError(t, err)
	t.Log(ticker)
}

func TestCoincheck_GetDepth(t *testing.T) {
	depth, err := api.GetDepth(3, BTC_JPY)
	assert.NoError(t, err)
	t.Log(depth)
}
