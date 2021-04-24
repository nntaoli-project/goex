package huobi

import (
	"github.com/Jameslu041/goex"
	"testing"
)

var wallet *Wallet

func init() {
	wallet = NewWallet(&goex.APIConfig{
		HttpClient:   httpProxyClient,
		ApiKey:       "",
		ApiSecretKey: "",
	})
}

func TestWallet_Transfer(t *testing.T) {
	t.Log(wallet.Transfer(goex.TransferParameter{
		Currency: "BTC",
		From:     goex.SWAP_USDT,
		To:       goex.SPOT,
		Amount:   11,
	}))
}
