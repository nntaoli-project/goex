package huobi

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"testing"
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
