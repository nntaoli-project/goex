package okex

import (
	goex "github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
var config2 = &goex.APIConfig{
	Endpoint: "https://www.okex.com",
	//HttpClient: &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: func(req *http.Request) (*url.URL, error) {
	//			return &url.URL{
	//				Scheme: "socks5",
	//				Host:   "127.0.0.1:1080"}, nil
	//		},
	//	},
	//},
	ApiKey:        "",
	ApiSecretKey:  "",
	ApiPassphrase: "",
}

var okex = NewOKEx(config2)

func TestOKExSpot_GetAccount(t *testing.T) {
	t.Log(okex.GetAccount())
}

func TestOKExSpot_LimitBuy(t *testing.T) {
	t.Log(okex.OKExSpot.LimitBuy("0.001", "9910", goex.BTC_USD))
}

func TestOKExSpot_CancelOrder(t *testing.T) {
	t.Log(okex.OKExSpot.CancelOrder("3117823911340032", goex.BTC_USD))

}

func TestOKExSpot_GetOneOrder(t *testing.T) {
	t.Log(okex.OKExSpot.GetOneOrder("42152275c599444aa8ec1d33bd8003fb", goex.BTC_USD))
}

func TestOKExSpot_GetUnfinishOrders(t *testing.T) {
	t.Log(okex.OKExSpot.GetUnfinishOrders(goex.EOS_BTC))
}

func TestOKExSpot_GetTicker(t *testing.T) {
	t.Log(okex.OKExSpot.GetTicker(goex.BTC_USD))
}

func TestOKExSpot_GetDepth(t *testing.T) {
	dep, err := okex.OKExSpot.GetDepth(2, goex.EOS_BTC)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}

func TestOKExFuture_GetFutureTicker(t *testing.T) {
	t.Log(okex.OKExFuture.GetFutureTicker(goex.BTC_USD, "BTC-USD-190927"))
	t.Log(okex.OKExFuture.GetFutureTicker(goex.BTC_USD, goex.QUARTER_CONTRACT))
}

func TestOKExFuture_GetFutureUserinfo(t *testing.T) {
	t.Log(okex.OKExFuture.GetFutureUserinfo())
}

func TestOKExFuture_GetFuturePosition(t *testing.T) {
	t.Log(okex.OKExFuture.GetFuturePosition(goex.EOS_USD, goex.QUARTER_CONTRACT))
}

func TestOKExFuture_PlaceFutureOrder(t *testing.T) {
	t.Log(okex.OKExFuture.PlaceFutureOrder(goex.EOS_USD, goex.THIS_WEEK_CONTRACT, "5.8", "1", goex.OPEN_BUY, 0, 10))
}

func TestOKExFuture_PlaceFutureOrder2(t *testing.T) {
	t.Log(okex.OKExFuture.PlaceFutureOrder2(0, &goex.FutureOrder{
		Currency:     goex.EOS_USD,
		ContractName: goex.QUARTER_CONTRACT,
		OType:        goex.OPEN_BUY,
		OrderType:    ORDINARY,
		Price:        5.9,
		Amount:       10,
		LeverRate:    10}))
}

func TestOKExFuture_FutureCancelOrder(t *testing.T) {
	t.Log(okex.OKExFuture.FutureCancelOrder(goex.EOS_USD, goex.QUARTER_CONTRACT, "e88bd3361de94512b8acaf9aa154f95a"))
}

func TestOKExFuture_GetFutureOrder(t *testing.T) {
	t.Log(okex.OKExFuture.GetFutureOrder("3145664744431616", goex.EOS_USD, goex.QUARTER_CONTRACT))
}

func TestOKExFuture_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(okex.OKExFuture.GetUnfinishFutureOrders(goex.EOS_USD, goex.QUARTER_CONTRACT))
}

func TestOKExFuture_MarketCloseAllPosition(t *testing.T) {
	t.Log(okex.OKExFuture.MarketCloseAllPosition(goex.BTC_USD, goex.THIS_WEEK_CONTRACT, goex.CLOSE_BUY))
}

func TestOKExFuture_GetRate(t *testing.T) {
	t.Log(okex.OKExFuture.GetRate())
}
