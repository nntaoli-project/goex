package gdax

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var gdax = New(http.DefaultClient, "", "")

func TestGdax_GetTicker(t *testing.T) {
	ticker, err := gdax.GetTicker(goex.BTC_USD)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func TestGdax_Get24HStats(t *testing.T) {
	stats, err := gdax.Get24HStats(goex.BTC_USD)
	t.Log("err=>", err)
	t.Log("stats=>", stats)
}

func TestGdax_GetDepth(t *testing.T) {
	dep, err := gdax.GetDepth(2, goex.BTC_USD)
	t.Log("err=>", err)
	t.Log("bids=>", dep.BidList)
	t.Log("asks=>", dep.AskList)
}
