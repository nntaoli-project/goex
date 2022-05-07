package exx

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"net/url"
	"testing"
)


var (
	api_key       = "yourAccessKey"
	api_secretkey = "yourSecretKey"
	exx           = New(http.DefaultClient, api_key, api_secretkey)
)

func TestExx_Signed(t *testing.T) {
	return
	params := url.Values{}
	exx.accessKey = "yourAccessKey"
	exx.secretKey = "yourSecretKey"
	exx.buildPostForm(&params)
	t.Log(params)
}

func TestExx_GetAccount(t *testing.T) {
	//return
	acc, err := exx.GetAccount()
	t.Log(acc, err)
	//t.Log(acc.SubAccounts[goex.BTC])
}

func TestExx_GetTicker(t *testing.T) {
	return
	ticker, err := exx.GetTicker(goex.BTC_USD)
	t.Log(ticker, err)
}

func TestExx_GetDepth(t *testing.T) {
	return
	dep, _ := exx.GetDepth(2, goex.BTC_USDT)
	t.Log(dep)
	t.Log(dep.AskList[0])
	t.Log(dep.BidList[0])
}

func TestExx_LimitSell(t *testing.T) {
	return
	ord, err := exx.LimitSell("0.001", "75000", goex.NewCurrencyPair2("BTC_QC"))
	t.Log(err)
	t.Log(ord)
}

func TestExx_LimitBuy(t *testing.T) {
	return
	ord, err := exx.LimitBuy("2", "4", goex.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ord)
}

func TestExx_CancelOrder(t *testing.T) {
	return
	r, err := exx.CancelOrder("201802014255365", goex.NewCurrencyPair2("BTC_QC"))
	t.Log(err)
	t.Log(r)
}

func TestExx_GetUnfinishOrders(t *testing.T) {
	return
	ords, err := exx.GetUnfinishOrders(goex.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ords)
}

func TestExx_GetOneOrder(t *testing.T) {
	return
	ord, err := exx.GetOneOrder("20180201341043", goex.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ord)
}
