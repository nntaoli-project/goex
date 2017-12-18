package wex

import (
	"net/http"
	"testing"
	"github.com/nntaoli-project/GoEx"
)

var wex = New(http.DefaultClient, "", "")

func TestWex_GetTicker(t *testing.T) {
	ticker, err := wex.GetTicker(goex.BTC_USD)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
