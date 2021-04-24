package binance

import (
	"github.com/Jameslu041/goex"
	"log"
	"os"
	"testing"
	"time"
)

var spotWs *SpotWs

func createSpotWs() {
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
	createSpotWs()

	spotWs.SubscribeDepth(goex.BTC_USDT)
	spotWs.SubscribeTicker(goex.LTC_USDT)
	time.Sleep(11 * time.Minute)
}

func TestSpotWs_SubscribeTicker(t *testing.T) {
	createSpotWs()

	spotWs.SubscribeTicker(goex.LTC_USDT)
	time.Sleep(30 * time.Minute)
}
