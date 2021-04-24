package bitmex

import (
	"github.com/Jameslu041/goex"
	"os"
	"testing"
	"time"
)

func TestNewSwapWs(t *testing.T) {
	os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:1080")
	ws := NewSwapWs()
	ws.DepthCallback(func(depth *goex.Depth) {
		t.Log(depth)
	})
	ws.TickerCallback(func(ticker *goex.FutureTicker) {
		t.Logf("%s %v", ticker.ContractType, ticker.Ticker)
	})
	//ws.SubscribeDepth(goex.NewCurrencyPair2("LTC_USD"), goex.SWAP_CONTRACT)
	ws.SubscribeTicker(goex.LTC_USDT , goex.SWAP_CONTRACT)

	time.Sleep(5 * time.Minute)
}
