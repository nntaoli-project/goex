package ocx

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"net/url"
	"testing"
)

var o = New(&http.Client{}, "", "")

func TestOcx_GetServerTime(t *testing.T) {
	return
	t.Log(o.GetServerTime())
}

func TestOcx_buildSigned(t *testing.T) {
	return
	method := "GET"
	path := "/api/v2/markets"
	para := url.Values{}
	para.Set("access_key", "xxx")
	para.Set("tonce", "123456789")
	para.Set("foo", "bar")
	t.Log(o.buildSigned(method, path, &para))

}

func TestOcx_LimitSell(t *testing.T) {
	//return

	o.LimitBuy("1", "1", goex.BTC_USDT)
}
