package bigone

import (
	. "github.com/nntaoli-project/goex"
	"net/http"
	"testing"
)

var (
	bo = New(http.DefaultClient, "", "")
)

func TestBigone_GetTicker(t *testing.T) {
	return
	t.Log(bo.GetTicker(ETH_BTC))
}
func TestBigone_GetDepth(t *testing.T) {
	return
	t.Log(bo.GetDepth(1, ETH_BTC))
}
func TestBigone_GetAccount(t *testing.T) {
	return
	t.Log(bo.GetAccount())
}
func TestBigone_GetUnfinishOrders(t *testing.T) {
	return
	BIG_BTC := NewCurrencyPair2("BIG_BTC")
	t.Log(bo.GetUnfinishOrders(BIG_BTC))
}
func TestBigone_GetOrderHistorys(t *testing.T) {
	return
	TCT_BTC := NewCurrencyPair2("TCT_BTC")
	t.Log(bo.GetOrderHistorys(TCT_BTC, 1, 1))
}
func TestBigone_LimitSell(t *testing.T) {
	return
	TCT_BTC := NewCurrencyPair2("TCT_BTC")
	t.Log(bo.LimitSell("322", "1", TCT_BTC))
}
func TestBigone_CancelOrder(t *testing.T) {
	return
	t.Log(bo.CancelOrder("9f352ec4-3502-4dea-bdd4-0860d22f80e3", EOS_BTC))
}
func TestBigone_GetOneOrder(t *testing.T) {
	return
	TCT_BTC := NewCurrencyPair2("TCT_BTC")

	t.Log(bo.GetOneOrder("ccfe4661-82b5-4dfa-bc55-1f4ec61b1611", TCT_BTC))
}
