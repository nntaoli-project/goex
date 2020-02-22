package fmex

import (
	goex "github.com/nntaoli-project/GoEx"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var fm = NewFMexSwap(&goex.APIConfig{HttpClient: &http.Client{
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
}, ApiKey: "dd2d899cdb694322a10954070580030b", ApiSecretKey: "64f05065a6684d5182f7a382c3d23069"})

func init() {
	fm.SetBaseUri("https://api.testnet.fmex.com")
}

func TestFMexSwap_GetFutureTicker(t *testing.T) {
	t.Log(fm.GetFutureTicker(goex.BTC_USDT, ""))
}

func TestFMexSwap_GetFutureDepth(t *testing.T) {
	t.Log(fm.GetFutureDepth(goex.BTC_USDT, "", 6))
}

func TestFMexSwap_GetKlineRecords(t *testing.T) {
	t.Log(fm.GetKlineRecords("", goex.BTC_USDT, goex.KLINE_PERIOD_15MIN, 1, 1))
}

func TestFMexSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(fm.GetFutureUserinfo())
}

func TestFMexSwap_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(fm.GetUnfinishFutureOrders(goex.BTC_USDT, ""))
}

func TestFMexSwap_GetFuturePosition(t *testing.T) {
	t.Log(fm.GetFuturePosition(goex.BTC_USDT, ""))
}

func TestFMexSwap_GetFutureOrderHistory(t *testing.T) {
	t.Log(fm.GetFutureOrderHistory(nil, goex.BTC_USDT, ""))
}

func TestFMexSwap_GetFutureOrder(t *testing.T) {
	t.Log(fm.GetFutureOrder("466846361910", goex.BTC_USDT, ""))
}

func TestFMexSwap_GetTrades(t *testing.T) {
	t.Log(fm.GetTrades("", goex.BTC_USDT, 1))
}

func TestFMexSwap_GetFutureIndex(t *testing.T) {
	t.Log(fm.GetFutureIndex(goex.BTC_USDT))
}
