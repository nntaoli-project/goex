package bitstamp

import (
	"github.com/nntaoli-project/goex"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

var client = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		log.Println("======")
		return nil
	},
}
var btmp = NewBitstamp(&client, "", "", "")

func TestBitstamp_GetAccount(t *testing.T) {
	acc, err := btmp.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}

func TestBitstamp_GetTicker(t *testing.T) {
	ticker, err := btmp.GetTicker(goex.BTC_USD)
	assert.Nil(t, err)
	t.Log(ticker)
}

func TestBitstamp_GetDepth(t *testing.T) {
	dep, err := btmp.GetDepth(5, goex.BTC_USD)
	assert.Nil(t, err)
	t.Log(dep.BidList)
	t.Log(dep.AskList)
}

func TestBitstamp_LimitBuy(t *testing.T) {
	ord, err := btmp.LimitBuy("55", "0.12", goex.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestBitstamp_LimitSell(t *testing.T) {
	ord, err := btmp.LimitSell("40", "0.22", goex.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestBitstamp_MarketBuy(t *testing.T) {
	ord, err := btmp.MarketBuy("1", "", goex.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestBitstamp_MarketSell(t *testing.T) {
	ord, err := btmp.MarketSell("2", "", goex.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestBitstamp_CancelOrder(t *testing.T) {
	r, err := btmp.CancelOrder("311242779", goex.XRP_USD)
	assert.Nil(t, err)
	t.Log(r)
}

func TestBitstamp_GetUnfinishOrders(t *testing.T) {
	ords, err := btmp.GetUnfinishOrders(goex.XRP_USD)
	assert.Nil(t, err)
	t.Log(ords)
}

func TestBitstamp_GetOneOrder(t *testing.T) {
	ord, err := btmp.GetOneOrder("311752078", goex.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
