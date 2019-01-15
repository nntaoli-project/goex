package okex

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var config = &goex.APIConfig{
	HttpClient:    http.DefaultClient,
	ApiPassphrase: "",
	ApiSecretKey:  "",
	ApiKey:        "",
}

var okExSwap goex.FutureRestAPI = NewOKExSwap(config)

func TestOKExSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(okExSwap.GetFutureUserinfo())
}

func TestOKExSwap_PlaceFutureOrder(t *testing.T) {
	t.Log(okExSwap.PlaceFutureOrder(goex.BTC_USD, BTC_USD_SWAP, "3675.6", "10", goex.CLOSE_BUY, 0, 0))
}

func TestOKExSwap_FutureCancelOrder(t *testing.T) {
	t.Log(okExSwap.FutureCancelOrder(goex.BTC_USD, BTC_USD_SWAP, "64-a-3e3c3c359-0"))
}

func TestOKExSwap_GetFutureOrder(t *testing.T) {
	t.Log(okExSwap.GetFutureOrder("65-4-3e62a331c-0", goex.LTC_USD, LTC_USD_SWAP))
}

func TestOKExSwap_GetFuturePosition(t *testing.T) {
	t.Log(okExSwap.GetFuturePosition(goex.BTC_USD, BTC_USD_SWAP))
}

func TestOKExSwap_GetFutureDepth(t *testing.T) {
	t.Log(okExSwap.GetFutureDepth(goex.LTC_USD, LTC_USD_SWAP, 10))
}

func TestOKExSwap_GetFutureTicker(t *testing.T) {
	t.Log(okExSwap.GetFutureTicker(goex.BTC_USD, BTC_USD_SWAP))
}

func TestOKExSwap_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(okExSwap.GetUnfinishFutureOrders(goex.BTC_USD, BTC_USD_SWAP))
}
