package aex

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var acx = New(http.DefaultClient, "", "", "")

func TestAex_GetTicker(t *testing.T) {
	ticker, err := acx.GetTicker(goex.ETH_BTC)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
