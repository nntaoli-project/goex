package huobi

import (
	"github.com/Jameslu041/goex"
	"net/http"
	"testing"
	"time"
)

var swap *HbdmSwap

func init() {
	swap = NewHbdmSwap(&goex.APIConfig{
		HttpClient:   http.DefaultClient,
		Endpoint:     "https://api.btcgateway.pro",
		ApiKey:       "",
		ApiSecretKey: "",
		Lever:        5,
	})
}

func TestHbdmSwap_GetFutureTicker(t *testing.T) {
	t.Log(swap.GetFutureTicker(goex.BTC_USD, goex.SWAP_CONTRACT))
}

func TestHbdmSwap_GetFutureDepth(t *testing.T) {
	dep, err := swap.GetFutureDepth(goex.BTC_USD, goex.SWAP_CONTRACT, 5)
	t.Log(err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}

func TestHbdmSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(swap.GetFutureUserinfo(goex.NewCurrencyPair2("DOT_USD")))
}

func TestHbdmSwap_GetFuturePosition(t *testing.T) {
	t.Log(swap.GetFuturePosition(goex.NewCurrencyPair2("DOT_USD"), goex.SWAP_CONTRACT))
}

func TestHbdmSwap_LimitFuturesOrder(t *testing.T) {
	//784115347040780289
	t.Log(swap.LimitFuturesOrder(goex.NewCurrencyPair2("DOT_USD"), goex.SWAP_CONTRACT, "6.5", "1", goex.OPEN_SELL))
}

func TestHbdmSwap_FutureCancelOrder(t *testing.T) {
	t.Log(swap.FutureCancelOrder(goex.NewCurrencyPair2("DOT_USD"), goex.SWAP_CONTRACT, "784118017750929408"))
}

func TestHbdmSwap_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(swap.GetUnfinishFutureOrders(goex.NewCurrencyPair2("DOT_USD"), goex.SWAP_CONTRACT))
}

func TestHbdmSwap_GetFutureOrder(t *testing.T) {
	t.Log(swap.GetFutureOrder("784118017750929408", goex.NewCurrencyPair2("DOT_USD"), goex.SWAP_CONTRACT))
}

func TestHbdmSwap_GetFutureOrderHistory(t *testing.T) {
	t.Log(swap.GetFutureOrderHistory(goex.NewCurrencyPair2("KSM_USD"), goex.SWAP_CONTRACT,
		goex.OptionalParameter{}.Optional("start_time", time.Now().Add(-5*24*time.Hour).Unix()*1000),
		goex.OptionalParameter{}.Optional("end_time", time.Now().Unix()*1000)))
}
