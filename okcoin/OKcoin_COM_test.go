package okcoin

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var okcom = NewCOM(http.DefaultClient, "", "")

func TestOKCoinCOM_API_GetTicker(t *testing.T) {
	ticker, _ := okcom.GetTicker(goex.BTC_CNY)
	t.Log(ticker)
}
