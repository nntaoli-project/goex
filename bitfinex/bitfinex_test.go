package bitfinex

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var bfx = New(http.DefaultClient, "", "")

func TestBitfinex_GetTicker(t *testing.T) {
	ticker, _ := bfx.GetTicker(goex.ETH_BTC)
	t.Log(ticker)
}

func TestBitfinex_GetDepth(t *testing.T) {
	dep, _ := bfx.GetDepth(2, goex.ETH_BTC)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}
