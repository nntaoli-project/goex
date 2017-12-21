package cryptopia

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var ctp = New(http.DefaultClient, "", "")

func TestCryptopia_GetTicker(t *testing.T) {
	ticker, err := ctp.GetTicker(goex.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
