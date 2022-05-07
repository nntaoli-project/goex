package huobi

import (
	"github.com/nntaoli-project/goex"
	"os"
	"testing"
	"time"
)

func TestNewSpotWs(t *testing.T) {
	os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:1080")
	spotWs := NewSpotWs()
	spotWs.DepthCallback(func(depth *goex.Depth) {
		t.Log("asks=", depth.AskList)
		t.Log("bids=", depth.BidList)
	})
	spotWs.TickerCallback(func(ticker *goex.Ticker) {
		t.Log(ticker)
	})
	spotWs.SubscribeTicker(goex.NewCurrencyPair2("BTC_USDT"))
	spotWs.SubscribeTicker(goex.NewCurrencyPair2("USDT_HUSD"))
	spotWs.SubscribeTicker(goex.NewCurrencyPair2("LTC_BTC"))
	spotWs.SubscribeTicker(goex.NewCurrencyPair2("EOS_ETH"))
	spotWs.SubscribeTicker(goex.NewCurrencyPair2("LTC_HT"))
	spotWs.SubscribeTicker(goex.NewCurrencyPair2("BTT_TRX"))
	//spotWs.SubscribeDepth(goex.BTC_USDT)
	time.Sleep(time.Minute)
}
