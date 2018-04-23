package huobi

import (
	"github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"internal/log"
	"net/http"
	"testing"
	"time"
)

//
var hbpro = NewHuobiPro(http.DefaultClient, "", "", "")

func TestHuobiPro_GetTicker(t *testing.T) {
	ticker, err := hbpro.GetTicker(goex.XRP_BTC)
	assert.Nil(t, err)
	t.Log(ticker)
}

func TestHuobiPro_GetDepth(t *testing.T) {
	dep, err := hbpro.GetDepth(2, goex.LTC_BTC)
	assert.Nil(t, err)
	t.Log(dep)
}

func TestHuobiPro_GetAccountId(t *testing.T) {
	id, err := hbpro.GetAccountId()
	assert.Nil(t, err)
	t.Log(id)
}

func TestHuobiPro_GetAccount(t *testing.T) {
	acc, err := hbpro.GetAccount()
	assert.Nil(t, err)
	t.Log(acc.SubAccounts[goex.LTC], acc.SubAccounts[goex.BTC], acc.SubAccounts[goex.BCH])
}

func TestHuobiPro_LimitBuy(t *testing.T) {
	ord, err := hbpro.LimitBuy("", "0.09122", goex.BCC_BTC)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestHuobiPro_LimitSell(t *testing.T) {
	ord, err := hbpro.LimitSell("1", "0.212", goex.BCC_BTC)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestHuobiPro_MarketSell(t *testing.T) {
	ord, err := hbpro.MarketSell("0.1738", "0.212", goex.BCC_BTC)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestHuobiPro_MarketBuy(t *testing.T) {
	ord, err := hbpro.MarketBuy("0.02", "", goex.BCC_BTC)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestHuobiPro_GetUnfinishOrders(t *testing.T) {
	ords, err := hbpro.GetUnfinishOrders(goex.ETC_USDT)
	assert.Nil(t, err)
	t.Log(ords)
}

func TestHuobiPro_CancelOrder(t *testing.T) {
	r, err := hbpro.CancelOrder("600329873", goex.ETH_USDT)
	assert.Nil(t, err)
	t.Log(r)
	t.Log(err)
}

func TestHuobiPro_GetOneOrder(t *testing.T) {
	ord, err := hbpro.GetOneOrder("1116237737", goex.LTC_BTC)
	assert.Nil(t, err)
	t.Log(ord)
}

func TestHuobiPro_GetOrderHistorys(t *testing.T) {
	ords, err := hbpro.GetOrderHistorys(goex.NewCurrencyPair2("HT_USDT"), 1, 3)
	t.Log(err)
	t.Log(ords)
}

func TestHuobiPro_GetDepthWithWs(t *testing.T) {
	hbpro.GetDepthWithWs(goex.BTC_USDT, func(dep *goex.Depth) {
		log.Println("%+v", *dep)
	})
	time.Sleep(time.Minute)
}
