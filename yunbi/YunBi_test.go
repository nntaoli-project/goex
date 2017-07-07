package yunbi

import (
	"testing"
	"net/http"
	"github.com/nntaoli/crypto_coin_api"
)

var (
	yb = New(http.DefaultClient, "", "")
)

func TestYunBi_GetTicker(t *testing.T) {
	t.Log(yb.GetTicker(coinapi.BTS_CNY))
	t.Log(yb.GetTicker(coinapi.SC_CNY))
	t.Log(yb.GetTicker(coinapi.EOS_CNY))
}
