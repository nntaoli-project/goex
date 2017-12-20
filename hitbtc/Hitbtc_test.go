package hitbtc

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var htb = New(http.DefaultClient, "", "")

func TestHitbtc_GetTicker(t *testing.T) {
	ticker, err := htb.GetTicker(goex.BTC_USD)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
