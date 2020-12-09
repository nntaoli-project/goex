package binance

import (
	"github.com/nntaoli-project/goex"
	"log"
	"os"
	"testing"
	"time"
)

var spotWs *SpotWs

func init() {
	os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:1080")
	spotWs = NewSpotWs()
	spotWs.DepthCallback(func(depth *goex.Depth) {
		log.Println(depth)
	})
	spotWs.TickerCallback(func(ticker *goex.Ticker) {
		log.Println(ticker)
	})
}

func TestSpotWs_DepthCallback(t *testing.T) {
	spotWs.SubscribeDepth(goex.BTC_USDT)
	time.Sleep(11 * time.Minute)
}

func TestSpotWs_SubscribeTicker(t *testing.T) {
	spotWs.SubscribeTicker(goex.LTC_USDT)
	time.Sleep(30 * time.Minute)
}
