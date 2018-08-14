package fcoin

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var ft = NewFCoin(http.DefaultClient, "", "")

func TestFCoin_GetTicker(t *testing.T) {
	t.Log(ft.GetTicker(goex.NewCurrencyPair2("BNB_USDT")))
}

func TestFCoin_GetDepth(t *testing.T) {
	dep, _ := ft.GetDepth(1, goex.BTC_USDT)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}

func TestFCoin_GetAccount(t *testing.T) {
	acc, _ := ft.GetAccount()
	t.Log(acc)
}

func TestFCoin_LimitBuy(t *testing.T) {
	t.Log(ft.LimitBuy("0.01", "100", goex.ETC_USD))
}

func TestFCoin_LimitSell(t *testing.T) {
	t.Log(ft.LimitSell("0.01", "50", goex.ETC_USD))
}

func TestFCoin_GetOneOrder(t *testing.T) {
	t.Log(ft.GetOneOrder("KRcowt_w79qxcBdooYb-RxtZ_67TFcme7eUXU8bMusg=", goex.ETC_USDT))
}

func TestFCoin_CancelOrder(t *testing.T) {
	t.Log(ft.CancelOrder("-MR0CItwW-rpSFJau7bfCyUBrw9nrkLNipV9odvPlRQ=", goex.ETC_USDT))
}

func TestFCoin_GetUnfinishOrders(t *testing.T) {
	t.Log(ft.GetUnfinishOrders(goex.ETC_USDT))
}
