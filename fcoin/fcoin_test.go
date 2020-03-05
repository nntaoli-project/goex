package fcoin

import (
	"github.com/nntaoli-project/goex"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var ft = NewFCoin(&http.Client{
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
}, "dd2d899cdb694322a10954070580030b", "64f05065a6684d5182f7a382c3d23069")

func TestFCoin_GetTicker(t *testing.T) {
	return
	t.Log(ft.GetTicker(goex.NewCurrencyPair2("BTC_USDT")))
}

func TestFCoin_GetDepth(t *testing.T) {
	return
	dep, _ := ft.GetDepth(1, goex.BTC_USDT)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}

func TestFCoin_GetAccount(t *testing.T) {
	acc, _ := ft.GetAccount()
	t.Log(acc)
}

func TestFCoin_LimitBuy(t *testing.T) {
	return
	t.Log(ft.LimitBuy("0.01", "100", goex.BTC_USDT))
}

func TestFCoin_LimitSell(t *testing.T) {
	return
	t.Log(ft.LimitSell("0.01", "50", goex.ETC_USD))
}

func TestFCoin_GetOneOrder(t *testing.T) {
	return
	t.Log(ft.GetOneOrder("KRcowt_w79qxcBdooYb-RxtZ_67TFcme7eUXU8bMusg=", goex.ETC_USDT))
}

func TestFCoin_CancelOrder(t *testing.T) {
	return
	t.Log(ft.CancelOrder("-MR0CItwW-rpSFJau7bfCyUBrw9nrkLNipV9odvPlRQ=", goex.ETC_USDT))
}

//Zh4_uPspaahORdOXQhy7X-7W6rcVCzRYO9qL3HMQG6qFzWiTgSgR2xW0UqrN6G81yg4KRKI6mY7VIHoi9iSuRg==
func TestFCoin_GetUnfinishOrders(t *testing.T) {
	return
	t.Log(ft.GetUnfinishOrders(goex.ETC_USDT))
}

func TestFCoin_GetOrderHistorys(t *testing.T) {
	return
	t.Log(ft.GetOrderHistorys(goex.BTC_USDT, 1, 1))
}

func TestFCoin_AssetTransfer(t *testing.T) {
	return
	ft.AssetTransfer(goex.NewCurrency("FT", ""), "0.000945618753747253", "assets", "spot")
}

func TestFCoin_GetAssets(t *testing.T) {
	return
	acc, _ := ft.GetAssets()
	t.Log(acc)
}

func TestFCoin_GetKlineRecords(t *testing.T) {
	return
	t.Log(ft.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN, 20, 0))
}
func TestFCoin_MarketSell(t *testing.T) {
	return
	//ord, err := ft.LimitBuy("0.001", "9500", goex.BTC_USD)
	//t.Log(ord, err)
	cancelled, err := ft.CancelOrder("ICpzuGFpNc69bY59cgyo3DsjIui4J3GJjIXAEeDYATP3-DFzQEOCqu1dWkFgU8yDfSDJ-WIsn9FIQXVAzrwDzA==", goex.BTC_USD)
	t.Log(cancelled, err)
	cancelled, err = ft.CancelOrder("ICpzuGFpNc69bY59cgyo3DsjIui4J3GJjIXAEeDYATP3-DFzQEOCqu1dWkFgU8yDfSDJ-WIsn9FIQXVAzrwDzA==", goex.BTC_USD)
	t.Log(cancelled, err)

}
