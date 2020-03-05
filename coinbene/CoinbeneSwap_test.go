package coinbene

import (
	goex "github.com/nntaoli-project/goex"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var (
	httpProxyClient = &http.Client{
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
	coinbeneSwap = NewCoinbeneSwap(goex.APIConfig{
		HttpClient:   httpProxyClient,
		Endpoint:     "",
		ApiKey:       "",
		ApiSecretKey: "",
	})
)

func TestCoinbeneSwap_GetFutureTicker(t *testing.T) {
	t.Log(coinbeneSwap.GetFutureTicker(goex.BTC_USD, goex.SWAP_CONTRACT))
}

func TestCoinbeneSwap_GetFutureDepth(t *testing.T) {
	t.Log(coinbeneSwap.GetFutureDepth(goex.BTC_USDT, goex.SWAP_CONTRACT, 2))
}

func TestCoinbeneSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(coinbeneSwap.GetFutureUserinfo())
}

func TestCoinbeneSwap_GetFuturePosition(t *testing.T) {
	t.Log(coinbeneSwap.GetFuturePosition(goex.BTC_USDT, goex.SWAP_CONTRACT))
}

func TestCoinbeneSwap_PlaceFutureOrder(t *testing.T) {
	t.Log(coinbeneSwap.PlaceFutureOrder(goex.BTC_USDT, goex.SWAP_CONTRACT, "10000", "1", goex.OPEN_BUY, 0, 10))
}

func TestCoinbeneSwap_FutureCancelOrder(t *testing.T) {
	t.Log(coinbeneSwap.FutureCancelOrder(goex.BTC_USDT, goex.SWAP_CONTRACT, "580719990266232832"))
}

func TestCoinbeneSwap_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(coinbeneSwap.GetUnfinishFutureOrders(goex.BTC_USDT, goex.SWAP_CONTRACT))
}

func TestCoinbeneSwap_GetFutureOrder(t *testing.T) {
	t.Log(coinbeneSwap.GetFutureOrder("123", goex.BTC_USDT, goex.SWAP_CONTRACT))
}
