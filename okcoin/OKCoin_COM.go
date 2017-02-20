package okcoin

import (
	"net/http"
)

const (
	EXCHANGE_NAME_COM = "okcoin_com"
)

type OKCoinCOM_API struct {
	OKCoinCN_API
}

func NewCOM(client *http.Client, api_key, secret_key string) *OKCoinCOM_API {
	return &OKCoinCOM_API{OKCoinCN_API{client, api_key, secret_key, "https://www.okcoin.com/api/v1/"}}
}

func (ctx *OKCoinCOM_API) GetExchangeName() string {
	return EXCHANGE_NAME_COM
}
