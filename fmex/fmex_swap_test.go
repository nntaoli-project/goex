package fmex

import (
	goex "github.com/nntaoli-project/GoEx"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var fm = NewFMex(&goex.APIConfig{HttpClient: &http.Client{
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
	return
	t.Log(fm.GetFutureTicker(goex.BTC_USDT, ""))
}

func TestFMexSwap_GetFutureDepth(t *testing.T) {
	return
	t.Log(fm.GetFutureDepth(goex.BTC_USDT, "", 6))
}

func TestFMexSwap_GetKlineRecords(t *testing.T) {
	t.Log(fm.GetKlineRecords("", goex.BTC_USDT, goex.KLINE_PERIOD_15MIN, 1, 1))
}
