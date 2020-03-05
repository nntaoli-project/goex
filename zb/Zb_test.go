package zb

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"testing"
)

var (
	api_key       = ""
	api_secretkey = ""
	zb            = New(http.DefaultClient, api_key, api_secretkey)
)

func TestZb_GetAccount(t *testing.T) {
	acc, err := zb.GetAccount()
	t.Log(err)
	t.Log(acc.SubAccounts[goex.BTC])
}

func TestZb_GetTicker(t *testing.T) {
	ticker, _ := zb.GetTicker(goex.BCH_USD)
	t.Log(ticker)
}

func TestZb_GetDepth(t *testing.T) {
	dep, _ := zb.GetDepth(2, goex.BCH_USDT)
	t.Log(dep)
}

func TestZb_LimitSell(t *testing.T) {
	ord, err := zb.LimitSell("0.001", "75000", goex.NewCurrencyPair2("BTC_QC"))
	t.Log(err)
	t.Log(ord)
}

func TestZb_LimitBuy(t *testing.T) {
	ord, err := zb.LimitBuy("2", "4", goex.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ord)
}

func TestZb_CancelOrder(t *testing.T) {
	r, err := zb.CancelOrder("201802014255365", goex.NewCurrencyPair2("BTC_QC"))
	t.Log(err)
	t.Log(r)
}

func TestZb_GetUnfinishOrders(t *testing.T) {
	ords, err := zb.GetUnfinishOrders(goex.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ords)
}

func TestZb_GetOneOrder(t *testing.T) {
	ord, err := zb.GetOneOrder("20180201341043", goex.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ord)
}
