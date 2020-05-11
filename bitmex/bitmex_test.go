package bitmex

import (
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var httpProxyClient = &http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return &url.URL{
				Scheme: "socks5",
				Host:   "127.0.0.1:1080"}, nil
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	},
	Timeout: 10 * time.Second,
}

func init() {
	logger.Log.SetLevel(logger.DEBUG)
	mex = New(&goex.APIConfig{
		Endpoint:     "https://testnet.bitmex.com/",
		HttpClient:   httpProxyClient,
		ApiKey:       "",
		ApiSecretKey: ""})
}

var mex *bitmex

func TestBitmex_GetFutureDepth(t *testing.T) {
	dep, err := mex.GetFutureDepth(goex.BTC_USD, "Z19", 5)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}

func TestBitmex_GetFutureTicker(t *testing.T) {
	tk, er := mex.GetFutureTicker(goex.BTC_USD, "")
	if assert.Nil(t, er) {
		t.Logf("buy:%.8f ,sell: %.8f ,Last:%.8f , vol:%.8f", tk.Buy, tk.Sell, tk.Last, tk.Vol)
	}
}

func TestBitmex_GetIndicativeFundingRate(t *testing.T) {
	//rate, time, err := mex.GetIndicativeFundingRate("XBTUSD")
	//if assert.Nil(t, err) {
	//	t.Log(rate)
	//	t.Log(time.Local())
	//}
}

func TestBitmex_GetFutureUserinfo(t *testing.T) {
	userinfo, err := mex.GetFutureUserinfo()
	if assert.Nil(t, err) {
		t.Logf("%.8f", userinfo.FutureSubAccounts[goex.BTC].AccountRights)
		t.Logf("%.8f", userinfo.FutureSubAccounts[goex.BTC].KeepDeposit)
		t.Logf("%.8f", userinfo.FutureSubAccounts[goex.BTC].ProfitReal)
		t.Logf("%.8f", userinfo.FutureSubAccounts[goex.BTC].ProfitUnreal)
	}
}

func TestBitmex_GetFuturePosition(t *testing.T) {
	t.Log(mex.GetFuturePosition(goex.BTC_USD, ""))
}

func TestBitmex_PlaceFutureOrder(t *testing.T) {
	//{"orderID":"ae0436f4-9229-0be1-e9ea-45073a2a404a","clOrdID":"goexba0c770d9cea445eafb12b95fe220a0f"
	t.Log(mex.PlaceFutureOrder(goex.BTC_USD, goex.SWAP_CONTRACT, "9999", "2", goex.CLOSE_SELL, 0, 10))
}

func TestBitmex_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(mex.GetUnfinishFutureOrders(goex.BTC_USD, goex.SWAP_CONTRACT))
}

func TestBitmex_GetFutureOrder(t *testing.T) {
	t.Log(mex.GetFutureOrder("ae0436f4-9229-0be1-e9ea-45073a2a404a", goex.BTC_USD, goex.SWAP_CONTRACT))
}

func TestBitmex_FutureCancelOrder(t *testing.T) {
	t.Log(mex.FutureCancelOrder(goex.BTC_USD, goex.SWAP_CONTRACT, "goexfd6fd7694877448e8ae81a9cd7ecd89a"))
}
