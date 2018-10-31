package coin58

import (
	"testing"
	"net/http"
	"github.com/nntaoli-project/GoEx"
	"time"
)

var coin58 = New58Coin(http.DefaultClient, "", "")

func TestCoin58_GetTicker(t *testing.T) {
	t.Log(coin58.GetTicker(goex.NewCurrencyPair2("58b_btc")))
}

func TestCoin58_GetDepth(t *testing.T) {
	b := time.Now()
	dep , _ := coin58.GetDepth(1 , goex.NewCurrencyPair2("58b_ucc"))
	//dep, _ := coin58.GetDepth(2, goex.BTC_USD)
	t.Log(time.Now().Sub(b))
	t.Log(dep.BidList)
	t.Log(dep.AskList)
}

func TestCoin58_GetAccount(t *testing.T) {
	acc , err := coin58.GetAccount()
	t.Log(err)
	t.Log(acc)
}

func TestCoin58_LimitSell(t *testing.T) {
	ord , err := coin58.LimitSell("11.2" , "0.17" , goex.NewCurrencyPair2("58b_ucc"))
	t.Log(err)
	t.Log(ord)
}

func TestCoin58_GetOneOrder(t *testing.T) {
	ord , _ := coin58.GetOneOrder("23613" , goex.NewCurrencyPair2("58b_ucc")) //23478
	t.Logf("%+v" , ord)
}

func TestCoin58_GetUnfinishOrders(t *testing.T) {
	t.Log(coin58.GetUnfinishOrders(goex.NewCurrencyPair2("58b_usdt")))
}

func TestCoin58_CancelOrder(t *testing.T) {
	t.Log(coin58.CancelOrder("72" , goex.NewCurrencyPair2("58b_usdt")))
}

