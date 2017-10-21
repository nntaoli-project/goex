package btcc

import (
	"github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var btch = NewBTCChina(http.DefaultClient, "", "")

func TestBTCChina_GetTicker(t *testing.T) {
	ticker, err := btch.GetTicker(goex.BTC_CNY)
	assert.Nil(t, err)
	t.Log(ticker)
}

func TestBTCChina_GetDepth(t *testing.T) {
	dep, err := btch.GetDepth(2, goex.BTC_CNY)
	assert.Nil(t, err)
	t.Log(dep)
}

func TestBTCChina_GetAccount(t *testing.T) {
	acc, err := btch.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}

func TestBTCChina_LimitBuy(t *testing.T) {
	ord, err := btch.LimitBuy("0.001", "200", goex.LTC_CNY)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestBTCChina_CancelOrder(t *testing.T) {
	r, err := btch.CancelOrder("24956079", goex.LTC_CNY)
	assert.Nil(t, err)
	t.Log(r)
}

func TestBTCChina_GetOneOrder(t *testing.T) {
	order, err := btch.GetOneOrder("24956079", goex.LTC_CNY)
	assert.Nil(t, err)
	t.Log(order)
}

func TestBTCChina_GetUnfinishOrders(t *testing.T) {
	ords, err := btch.GetUnfinishOrders(goex.LTC_CNY)
	assert.Nil(t, err)
	t.Log(ords)
}
