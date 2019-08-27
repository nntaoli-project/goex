package okex

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"net/url"
	"testing"
)

var config = &goex.APIConfig{
	HttpClient: &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   "127.0.0.1:1080"}, nil
			},
		},
	},
	Endpoint:      "https://www.okex.com",
	ApiKey:        "",
	ApiSecretKey:  "",
	ApiPassphrase: "",
}

var okExSwap = NewOKExSwap(config)

func TestOKExSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(okExSwap.GetFutureUserinfo())
}

func TestOKExSwap_PlaceFutureOrder(t *testing.T) {
	t.Log(okExSwap.PlaceFutureOrder(goex.XRP_USD, goex.SWAP_CONTRACT, "0.2", "10", goex.OPEN_BUY, 0, 0))
}

func TestOKExSwap_FutureCancelOrder(t *testing.T) {
	t.Log(okExSwap.FutureCancelOrder(goex.XRP_USD, goex.SWAP_CONTRACT, "309935122485305344"))
}

func TestOKExSwap_GetFutureOrder(t *testing.T) {
	t.Log(okExSwap.GetFutureOrder("309935122485305344", goex.XRP_USD, goex.SWAP_CONTRACT))
}

func TestOKExSwap_GetFuturePosition(t *testing.T) {
	t.Log(okExSwap.GetFuturePosition(goex.BTC_USD, goex.SWAP_CONTRACT))
}

func TestOKExSwap_GetFutureDepth(t *testing.T) {
	t.Log(okExSwap.GetFutureDepth(goex.LTC_USD, goex.SWAP_CONTRACT, 10))
}

func TestOKExSwap_GetFutureTicker(t *testing.T) {
	t.Log(okExSwap.GetFutureTicker(goex.BTC_USD, goex.SWAP_CONTRACT))
}

func TestOKExSwap_GetUnfinishFutureOrders(t *testing.T) {
	ords, _ := okExSwap.GetUnfinishFutureOrders(goex.XRP_USD, goex.SWAP_CONTRACT)
	for _, ord := range ords {
		t.Log(ord.OrderID2, ord.ClientOid)
	}

}

func TestOKExSwap_GetHistoricalFunding(t *testing.T) {
	for i := 1; ; i++ {
		funding, err := okExSwap.GetHistoricalFunding(goex.SWAP_CONTRACT, goex.BTC_USD, i)
		t.Log(err, len(funding))
	}
}
