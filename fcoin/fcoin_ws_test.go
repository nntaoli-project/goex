package fcoin

import (
	"github.com/nntaoli-project/goex"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var fcws = NewFCoinWs(&http.Client{
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
	fcws.ProxyUrl("socks5://127.0.0.1:1080")
	fcws.SetCallbacks(printfTicker, printfDepth, printfTrade, nil)
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

func TestFCoinWs_GetTickerWithWs(t *testing.T) {
	return
	fcws.SubscribeTicker(goex.BTC_USDT)
	time.Sleep(time.Second * 10)
}
func TestFCoinWs_GetDepthWithWs(t *testing.T) {
	return
	fcws.SubscribeDepth(goex.BTC_USDT, 20)
	time.Sleep(time.Second * 10)
}
func TestFCoinWs_GetKLineWithWs(t *testing.T) {
	return
	fcws.SubscribeKline(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN)
	time.Sleep(time.Second * 120)
}
func TestFCoinWs_GetTradesWithWs(t *testing.T) {
	return
	fcws.SubscribeTrade(goex.BTC_USDT)
	time.Sleep(time.Second * 10)
}
func TestNewFCoinWs(t *testing.T) {
	fcws.SubscribeTrade(goex.BTC_USDT)
	//fcws.SubscribeKline(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN)
	//fcws.SubscribeDepth(goex.BTC_USDT, 20)
	fcws.SubscribeTicker(goex.BTC_USDT)
	time.Sleep(time.Minute * 10)

}
