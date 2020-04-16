package bitfinex

import (
	"net/http"
	"testing"

	"github.com/nntaoli-project/goex"
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

func TestBitfinex_GetKline(t *testing.T) {
	kline, _ := bfx.GetKlineRecords(goex.BTC_USD, goex.KLINE_PERIOD_1MONTH, 10, 0)
	for _, k := range kline {
		t.Log(k)
	}
}
