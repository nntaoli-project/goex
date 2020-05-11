package aacoin

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"net/url"
	"testing"
)

var aa = New(http.DefaultClient, "xxx", "b0xxxxxx-c6xxxxxx-94xxxxxx-dxxxx")

func TestAacoin_GetExchangeName(t *testing.T) {
	return
	params := url.Values{}
	params.Set("accessKey", aa.accessKey)
	params.Set("price", aa.accessKey)
	params.Set("quantity", aa.accessKey)
	params.Set("symbol", aa.accessKey)
	params.Set("type", aa.accessKey)
	t.Log(aa.buildSigned(&params))
}
func TestAacoin_GetAccount(t *testing.T) {
	return
	t.Log(aa.GetAccount())
}
func TestAacoin_GetTicker(t *testing.T) {
	return
	t.Log(aa.GetTicker(goex.BTC_USDT))
}
func TestAacoin_GetDepth(t *testing.T) {
	return
	t.Log(aa.GetDepth(1, goex.BTC_USDT))
}
func TestAacoin_LimitSell(t *testing.T) {
	//return
	t.Log(aa.LimitSell("1", "1000000", goex.BTC_USDT))
}

func TestAacoin_LimitBuy(t *testing.T) {
	t.Log(aa.LimitBuy("1", "1", goex.BTC_USDT))
}

func TestAacoin_GetUnfinishOrders(t *testing.T) {
	t.Log(aa.GetUnfinishOrders(goex.BTC_USDT))
}
