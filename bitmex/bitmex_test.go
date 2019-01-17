package bitmex

import (
	"github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var httpProxyClient = &http.Client{
	//Transport: &http.Transport{
	//	Proxy: func(req *http.Request) (*url.URL, error) {
	//		return &url.URL{
	//			Scheme: "socks5",
	//			Host:   "127.0.0.1:1080"}, nil
	//	},
	//	Dial: (&net.Dialer{
	//		Timeout: 10 * time.Second,
	//	}).Dial,
	//},
	Timeout: 10 * time.Second,
}
var (
	proxyURL  = "socks5://127.0.0.1:1080"
	bm_key    = ""
	bm_secret = ""
)
var bm = New(httpProxyClient, bm_key, bm_secret, BaseURL)

func TestBitmex_GetDepth(t *testing.T) {
	dep, err := bm.GetFutureDepth(goex.NewCurrencyPair(goex.XBT, goex.USD), "", 2)
	assert.Nil(t, err)
	t.Log(dep)
}

func TestBitmex_GetAccount(t *testing.T) {
	//return
	acc, err := bm.GetFutureUserinfo()
	t.Log(acc, err)
}

func TestBitmex_GetUnfinishFutureOrders(t *testing.T) {
	return
	ord, err := bm.GetUnfinishFutureOrders(goex.NewCurrencyPair(goex.XBT, goex.USD), "")
	t.Log(ord, err)

}
func TestBitmex_GetFutureOrder(t *testing.T) {
	return
	ord, err := bm.GetFutureOrder("4338b74c-7012-1e68-7472-94bf50816f0e", goex.NewCurrencyPair(goex.XBT, goex.USD), "")
	t.Log(ord, err)

}
func TestBitmex_GetFutureOrdersHistory(t *testing.T) {
	return
	ord, err := bm.GetFutureOrdersHistory(goex.NewCurrencyPair(goex.XBT, goex.USD), "")
	t.Log(ord, err)
	t.Log(len(ord))

}
