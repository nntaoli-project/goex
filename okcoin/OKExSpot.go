package okcoin

import "net/http"

type OKExSpot struct {
	OKCoinCN_API
}

func NewOKExSpot(client *http.Client, accesskey, secretkey string) *OKExSpot {
	return &OKExSpot{
		OKCoinCN_API{client, accesskey, secretkey, "https://www.okex.com/api/v1/"}}
}

func (ctx *OKExSpot) GetExchangeName() string {
	return "okex.com"
}
