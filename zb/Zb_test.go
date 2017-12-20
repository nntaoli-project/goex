package zb

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var zb = New(http.DefaultClient, "", "")

func TestZb_GetTicker(t *testing.T) {
	ticker, err := zb.GetTicker(goex.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
