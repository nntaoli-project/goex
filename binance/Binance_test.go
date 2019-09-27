package binance

import (
	"github.com/nntaoli-project/GoEx"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var ba = New(&http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("socks5://127.0.0.1:1080")
			return nil, nil
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	},
	Timeout: 10 * time.Second,
}, "", "")

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
	before := time.Now().Add(-time.Hour).Unix() * 1000
	kline, _ := ba.GetKlineRecords(goex.ETH_BTC, goex.KLINE_PERIOD_5MIN, 100, int(before))
	for _, k := range kline {
		tt := time.Unix(k.Timestamp, 0)
		t.Log(tt, k.Open, k.Close, k.High, k.Low, k.Vol)
	}
}

func TestBinance_GetTradeSymbols(t *testing.T) {
	t.Log(ba.GetTradeSymbols(goex.BTC_USDT))
}
