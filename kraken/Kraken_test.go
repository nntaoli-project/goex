package kraken

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var kk = New(http.DefaultClient, "", "")

func TestKraken_GetTicker(t *testing.T) {
	ticker, _ := kk.GetTicker(goex.LTC_BTC)
	t.Log(ticker)
}
func TestKraken_LimitSell(t *testing.T) {
	order, err := kk.LimitSell("1", "1", goex.LTC_BTC)
	t.Log(order, err)
}

func TestKraken_GetDepth(t *testing.T) {
	dep, err := kk.GetDepth(5, goex.ETH_BTC)
	t.Log(err)
	if err == nil {
		t.Log(dep.AskList)
		t.Log(dep.BidList)
	}
}

func TestKraken_GetAccount(t *testing.T) {
	account, err := kk.GetAccount()
	t.Log(account, err)
}

func TestKraken_GetUnfinishOrders(t *testing.T) {
	orders, err := kk.GetUnfinishOrders(goex.ETH_BTC)
	t.Log(orders, err)
}