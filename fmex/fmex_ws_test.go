package fmex

import (
	goex "github.com/nntaoli-project/GoEx"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var fmws = NewFMexWs(&http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("socks5://127.0.0.1:1080")
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	},
	Timeout: 10 * time.Second,
})

func init() {
	fmws.ProxyUrl("socks5://127.0.0.1:1080")
	fmws.SetCallbacks(printfTicker, printfDepth, printfTrade, nil)
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

func TestFMexWs_GetTickerWithWs(t *testing.T) {
	fmws.SubscribeTicker(goex.BTC_USDT)
	time.Sleep(time.Second * 10)
}
func TestFMexWs_GetDepthWithWs(t *testing.T) {
	return
	fmws.SubscribeDepth(goex.BTC_USDT, 20)
	time.Sleep(time.Second * 3)
}
func TestFMexWs_GetKLineWithWs(t *testing.T) {
	return
	fmws.SubscribeKline(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN)
	time.Sleep(time.Second * 120)
}
func TestFMexWs_GetTradesWithWs(t *testing.T) {
	return
	fmws.SubscribeTrade(goex.BTC_USDT)
	time.Sleep(time.Second * 120)
}
