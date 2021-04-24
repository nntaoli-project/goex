package binance

import (
	"github.com/Jameslu041/goex"
	"net/http"
	"testing"
)

var wallet *Wallet

func init() {
	wallet = NewWallet(&goex.APIConfig{
		HttpClient:   http.DefaultClient,
		ApiKey:       "",
		ApiSecretKey: "",
	})
}

func TestWallet_Transfer(t *testing.T) {
	t.Log(wallet.Transfer(goex.TransferParameter{
		Currency: "USDT",
		From:     goex.SPOT,
		To:       goex.SWAP_USDT,
		Amount:   100,
	}))
}
