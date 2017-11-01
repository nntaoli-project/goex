package chbtc

import (
	"github.com/nntaoli-project/GoEx"
	"net/http"
	"testing"
)

var (
	api_key       = ""
	api_secretkey = ""
	chbtc         = New(http.DefaultClient, api_key, api_secretkey)
)

func TestChbtc_GetAccount(t *testing.T) {
	acc, _ := chbtc.GetAccount()
	t.Log(acc)
}

func TestChbtc_GetTicker(t *testing.T) {
	ticker, _ := chbtc.GetTicker(goex.BCC_CNY)
	t.Log(ticker)
}

func TestChbtc_GetDepth(t *testing.T) {
	dep, _ := chbtc.GetDepth(1, goex.ETH_CNY)
	t.Log(dep)
}
