package chbtc

import (
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
