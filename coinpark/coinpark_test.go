package coinpark

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"testing"
)

var cpk = New(http.DefaultClient, "", "")

func TestCpk_buildSigned(t *testing.T) {
	//[{"cmd":"user/userInfo","body":{}}]
	return
	cmds := "[{\"cmd\":\"user/userInfo\",\"body\":{}}]"
	t.Log(cpk.buildSigned(cmds))
}
func TestCpk_GetAccount(t *testing.T) {
	//return
	t.Log(cpk.GetAccount())
}

func TestCpk_LimitBuy(t *testing.T) {
	return
	t.Log(cpk.LimitBuy("1", "1", goex.BTC_USDT))
}
func TestCpk_LimitSell(t *testing.T) {
	return
	t.Log(cpk.LimitSell("1", "999999", goex.BTC_USDT))
}

func TestCpk_CancelOrder(t *testing.T) {
	return
	t.Log(cpk.CancelOrder("123", goex.BTC_USDT))
}

func TestCpk_GetPairList(t *testing.T) {
	return
	t.Log(cpk.GetPairList())
}
