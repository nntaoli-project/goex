package binance

import (
	"github.com/nntaoli-project/GoEx"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var bnWs = NewBinanceWs(&http.Client{
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
})

func init() {
	bnWs.ProxyUrl("socks5://127.0.0.1:1080")
	bnWs.SetCallbacks(printfTicker, printfDepth, printfTrade, nil)
}

func printfTicker(ticker *goex.Ticker) {
	log.Println("ticker:", ticker)
}
func printfDepth(depth *goex.Depth) {
	log.Println("depth:", depth)
}
func printfTrade(trade *goex.Trade) {
	log.Println("trade:", trade)
}
func printfKline(kline *goex.Kline, period int) {
	log.Println("kline:", kline)
}

func TestBinanceWs_SubscribeTicker(t *testing.T) {
	return
	bnWs.SubscribeTicker(goex.BTC_USDT)
	time.Sleep(time.Second * 10)
}

func TestBinanceWs_GetDepthWithWs(t *testing.T) {
	return
	bnWs.SubscribeDepth(goex.BTC_USDT, 20)
	time.Sleep(time.Second * 10)
}
func TestBinanceWs_GetKLineWithWs(t *testing.T) {
	return
	bnWs.SubscribeKline(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN)
	time.Sleep(time.Second * 120)
}
func TestBinanceWs_GetTradesWithWs(t *testing.T) {
	//return
	bnWs.SubscribeTrade(goex.BTC_USDT)
	time.Sleep(time.Second * 10)
}
