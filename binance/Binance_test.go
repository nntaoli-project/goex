package binance

import (
	//"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var ba = New(http.DefaultClient, "", "")

func TestBinance_GetTicker(t *testing.T) {
	ticker, _ := ba.GetTicker(goex.ETH_BTC)
	t.Log(ticker)
}
func TestBinance_LimitBuy(t *testing.T) {
	order, err := ba.LimitBuy("11", "1", goex.LTC_BTC)
	t.Log(order, err)
}

func TestBinance_GetDepth(t *testing.T) {
	dep, err := ba.GetDepth(5, goex.ETH_BTC)
	t.Log(err)
	if err == nil {
		t.Log(dep.AskList)
		t.Log(dep.BidList)
	}
	dep, err = ba.GetDepth(2, goex.ETH_BTC)
	t.Log(err)
	if err == nil {
		t.Log(dep.AskList)
		t.Log(dep.BidList)
	}
}

func TestBinance_GetAccount(t *testing.T) {
	account, err := ba.GetAccount()
	t.Log(account, err)
}
