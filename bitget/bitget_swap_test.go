package bitget

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"net/url"
	"testing"
)

var bg = NewSwap(&goex.APIConfig{
	HttpClient: &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   "127.0.0.1:1080"}, nil
			},
		},
	}, //需要代理的这样配置
	Endpoint: "https://capi.bitget.io",
	ClientId: "",
	Lever:    0,
})

func TestBitgetSwap_GetFutureTicker(t *testing.T) {
	t.Log(bg.GetFutureTicker(goex.ETH_USDT, ""))
}

func TestBitgetSwap_GetServerTime(t *testing.T) {
	t.Log(bg.GetServerTime())
}

func TestBitgetSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(bg.GetFutureUserinfo(goex.ETH_USDT))
}

func TestBitgetSwap_LimitFuturesOrder(t *testing.T) {
	t.Log(bg.LimitFuturesOrder(goex.ETH_USDT, "", "350", "1", goex.CLOSE_BUY))
}

func TestBitgetSwap_GetFuturePosition(t *testing.T) {
	t.Log(bg.GetFuturePosition(goex.ETH_USDT, ""))
}

func TestBitgetSwap_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(bg.GetUnfinishFutureOrders(goex.ETH_USDT, ""))
}

func TestBitgetSwap_SetMarginLevel(t *testing.T) {
	t.Log(bg.SetMarginLevel(goex.ETH_USDT, 10, 2))
}

func TestBitgetSwap_GetMarginLevel(t *testing.T) {
	t.Log(bg.GetMarginLevel(goex.ETH_USDT))
}

func TestBitgetSwap_GetContractInfo(t *testing.T) {
	t.Log(bg.GetContractInfo(goex.ETH_USDT))
}

func TestBitgetSwap_GetFutureOrder(t *testing.T) {
	t.Log(bg.GetFutureOrder("671529783552638913", goex.ETH_USDT, ""))
}

func TestBitgetSwap_FutureCancelOrder(t *testing.T) {
	t.Log(bg.FutureCancelOrder(goex.ETH_USDT, "", "671529783552638913"))
}

func TestBitgetSwap_ModifyAutoAppendMargin(t *testing.T) {
	t.Log(bg.ModifyAutoAppendMargin(goex.ETH_USDT, 1, 1))
}
