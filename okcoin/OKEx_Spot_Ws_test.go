package okcoin

import (
	"github.com/nntaoli-project/GoEx"
	"testing"
	"time"
)

func TestNewOKExSpotWs(t *testing.T) {
	okSpotWs := NewOKExSpotWs()
//	okSpotWs.ProxyUrl("socks5://127.0.0.1:1080")

	okSpotWs.SetCallbacks(func(ticker *goex.Ticker) {
		t.Log(ticker)
	}, func(depth *goex.Depth) {
		t.Log(depth)
	}, func(trade *goex.Trade) {
		t.Log(trade)
	}, func(kline *goex.Kline, i int) {
		t.Log(i, kline)
	})

	okSpotWs.ErrorHandleFunc(func(err error) {
		t.Log(err)
	})
	//	t.Log(okSpotWs.SubscribeTicker(goex.BTC_USDT))
	//	t.Log(okSpotWs.SubscribeTicker(goex.BCH_USDT))
	//okSpotWs.SubscribeDepth(goex.BTC_USDT, 5)
	//okSpotWs.SubscribeTrade(goex.BTC_USDT)
	t.Log(okSpotWs.SubscribeKline(goex.BTC_USDT, goex.KLINE_PERIOD_1H))
	time.Sleep(10 * time.Second)
}
