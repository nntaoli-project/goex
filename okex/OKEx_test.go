package okex

import (
	goex "github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
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
	HttpClient:    http.DefaultClient,
	ApiKey:        "",
	ApiSecretKey:  "",
	ApiPassphrase: "",
}

var okex = NewOKEx(config2)

func TestOKExSpot_GetAccount(t *testing.T) {
	t.Log(okex.GetAccount())
}

func TestOKExSpot_LimitBuy(t *testing.T) {
	t.Log(okex.OKExSpot.LimitBuy("0.001", "9910", goex.BTC_USD)) //42152275c599444aa8ec1d33bd8003fb , 3117789364492288
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
