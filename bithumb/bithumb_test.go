package bithumb

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"testing"
)

var bh = New(http.DefaultClient, "", "")

func TestBithumb_GetTicker(t *testing.T) {
	ticker, err := bh.GetTicker(goex.NewCurrencyPair2("ALL_KAW"))
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func TestBithumb_GetDepth(t *testing.T) {
	dep, err := bh.GetDepth(1, goex.BTC_KRW)
	t.Log("err=>", err)
	t.Log("asks=>", dep.AskList)
	t.Log("bids=>", dep.BidList)
}
