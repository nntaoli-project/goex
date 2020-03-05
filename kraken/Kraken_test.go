package kraken

import (
	"github.com/nntaoli-project/goex"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var k = New(http.DefaultClient, "", "")
var BCH_XBT = goex.NewCurrencyPair(goex.BCH, goex.XBT)

func TestKraken_GetDepth(t *testing.T) {
	dep, err := k.GetDepth(2, goex.BTC_USD)
	assert.Nil(t, err)
	t.Log(dep)
}

func TestKraken_GetTicker(t *testing.T) {
	ticker, err := k.GetTicker(goex.ETC_BTC)
	assert.Nil(t, err)
	t.Log(ticker)
}

func TestKraken_GetAccount(t *testing.T) {
	acc, err := k.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}

func TestKraken_LimitSell(t *testing.T) {
	ord, err := k.LimitSell("0.01", "6900", goex.BTC_USD)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestKraken_LimitBuy(t *testing.T) {
	ord, err := k.LimitBuy("0.01", "6100", goex.NewCurrencyPair(goex.XBT, goex.USD))
	assert.Nil(t, err)
	t.Log(ord)
}

func TestKraken_GetUnfinishOrders(t *testing.T) {
	ords, err := k.GetUnfinishOrders(goex.NewCurrencyPair(goex.XBT, goex.USD))
	assert.Nil(t, err)
	t.Log(ords)
}

func TestKraken_CancelOrder(t *testing.T) {
	r, err := k.CancelOrder("O6EAJC-YAC3C-XDEEXQ", goex.NewCurrencyPair(goex.XBT, goex.USD))
	assert.Nil(t, err)
	t.Log(r)
}

func TestKraken_GetTradeBalance(t *testing.T) {
	//	k.GetTradeBalance()
}

func TestKraken_GetOneOrder(t *testing.T) {
	ord, err := k.GetOneOrder("ODCRMQ-RDEID-CY334C", goex.BTC_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
