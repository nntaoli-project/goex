package allcoin

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"testing"
)

var ac = New(http.DefaultClient, "", "")

func TestAllcoin_GetAccount(t *testing.T) {
	return
	t.Log(ac.GetAccount())
}
func TestAllcoin_GetUnfinishOrders(t *testing.T) {
	return
	t.Log(ac.GetUnfinishOrders(goex.ETH_BTC))
}
func TestAllcoin_GetTicker(t *testing.T) {
	return
	t.Log(ac.GetTicker(goex.ETH_BTC))
}

func TestAllcoin_GetDepth(t *testing.T) {
	return
	dep, _ := ac.GetDepth(1, goex.ETH_BTC)
	t.Log(dep)
}

func TestAllcoin_LimitBuy(t *testing.T) {
	t.Log(ac.LimitBuy("1", "0.07", goex.ETH_BTC))
}
