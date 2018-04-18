package btcmarkets

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var btcm = New(http.DefaultClient, "", "")

func TestBtcm_GetTicker(t *testing.T) {
	AUD := goex.NewCurrency("AUD", "")
	BTC_AUD := goex.NewCurrencyPair(goex.BTC, AUD)
	ticker, err := btcm.GetTicker(BTC_AUD)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
