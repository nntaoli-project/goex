package binance

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
	"time"
)

var ba = New(http.DefaultClient, "", "")

func TestBinance_GetTicker(t *testing.T) {
	return
	ticker, _ := ba.GetTicker(goex.LTC_BTC)
	t.Log(ticker)
}
func TestBinance_LimitSell(t *testing.T) {
	return
	order, err := ba.LimitSell("1", "1", goex.LTC_BTC)
	t.Log(order, err)
}

func TestBinance_GetDepth(t *testing.T) {
	return
	dep, err := ba.GetDepth(5, goex.ETH_BTC)
	t.Log(err)
	if err == nil {
		t.Log(dep.AskList)
		t.Log(dep.BidList)
	}
}

func TestBinance_GetAccount(t *testing.T) {
	return
	account, err := ba.GetAccount()
	t.Log(account, err)
}

func TestBinance_GetUnfinishOrders(t *testing.T) {
	return
	orders, err := ba.GetUnfinishOrders(goex.ETH_BTC)
	t.Log(orders, err)
}
func TestBinance_GetKlineRecords(t *testing.T) {

	t.Log(ba.GetKlineRecords(goex.ETH_BTC, goex.KLINE_PERIOD_1MIN, 100, int(time.Now().Add(-2*time.Hour).UnixNano())))
}
