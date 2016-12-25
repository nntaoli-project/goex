package yunbi

import (
	"testing"
	"net/http"
	"github.com/nntaoli/crypto_coin_api"
	"github.com/stretchr/testify/assert"
	"github.com/nntaoli/crypto_coin_api/yunbi"
)

var (
	yb = yunbi.New(http.DefaultClient, "", "")
)

func TestYunBi_GetTicker(t *testing.T) {
	ticker, err := yb.GetTicker(coinapi.ZEC_CNY)
	assert.NoError(t, err)
	t.Log(ticker)
}

func TestYunBi_GetDepth(t *testing.T) {
	dep, err := yb.GetDepth(2, coinapi.BTC_CNY)
	assert.NoError(t, err)
	t.Log(dep)
}

func TestYunBi_GetAccount(t *testing.T) {
	acc, err := yb.GetAccount()
	assert.NoError(t, err)
	t.Log(acc)
}

func TestYunBi_LimitBuy(t *testing.T) {
	ord, err := yb.LimitBuy("0.1", "6.5", coinapi.ETC_CNY)
	assert.NoError(t, err)
	t.Log(ord)
}

func TestYunBi_LimitSell(t *testing.T) {
	ord, err := yb.LimitSell("0.1", "8.2", coinapi.ETC_CNY)
	assert.NoError(t, err)
	t.Log(ord)
}

func TestYunBi_CancelOrder(t *testing.T) {
	r, err := yb.CancelOrder("361848402", coinapi.ETC_CNY)
	assert.NoError(t, err)
	t.Log(r)
}

func TestYunBi_GetUnfinishOrders(t *testing.T) {
	orders, err := yb.GetUnfinishOrders(coinapi.ETC_CNY)
	assert.NoError(t, err)
	t.Log(orders)
}

func TestYunBi_GetOneOrder(t *testing.T) {
	ord, err := yb.GetOneOrder("361848402", coinapi.ETC_CNY)
	assert.NoError(t, err)
	t.Log(ord)
}
