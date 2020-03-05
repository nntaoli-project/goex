package cryptopia

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"testing"
)

var ctp = New(http.DefaultClient, "", "")

func TestCryptopia_GetTicker(t *testing.T) {
	ticker, err := ctp.GetTicker(goex.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func TestCryptopia_GetDepth(t *testing.T) {
	depCtp, err := ctp.GetDepth(goex.BTC_USDT)
	t.Log(err)
	t.Log("AskList=>", depCtp.AskList)
	t.Log("BidList=>", depCtp.BidList)
}
