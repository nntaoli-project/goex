package hitbtc

import (
	"github.com/nntaoli-project/goex"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const (
	PubKey    = ""
	SecretKey = ""
)

var htb *Hitbtc

func init() {
	htb = New(http.DefaultClient, PubKey, SecretKey)
}

func TestHitbtc_GetSymbols(t *testing.T) {
	t.Log(htb.GetSymbols())
}

func TestHitbtc_adaptSymbolToCurrencyPair(t *testing.T) {
	t.Log(htb.adaptSymbolToCurrencyPair("DOGEBTC").String() == "DOGE_BTC")
	t.Log(htb.adaptSymbolToCurrencyPair("BTCGUSD").String() == "BTC_GUSD")
	t.Log(htb.adaptSymbolToCurrencyPair("btctusd").String() == "BTC_TUSD")
	t.Log(htb.adaptSymbolToCurrencyPair("BTCUSDC").String() == "BTC_USDC")
	t.Log(htb.adaptSymbolToCurrencyPair("ETHEOS").String() == "ETH_EOS")
}

func TestGetTicker(t *testing.T) {
	res, err := htb.GetTicker(goex.BCC_USD)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestGetAccount(t *testing.T) {
	res, err := htb.GetAccount()
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestDepth(t *testing.T) {
	res, err := htb.GetDepth(10, YCC_BTC)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestKline(t *testing.T) {
	res, err := htb.GetKline(YCC_BTC, "1M", 10, 0)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestTrades(t *testing.T) {
	res, err := htb.GetTrades(YCC_BTC, 1519862400)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestPlaceOrder(t *testing.T) {
	res, err := htb.LimitBuy("15", "0.000008", YCC_BTC)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestCancelOrder(t *testing.T) {
	res, err := htb.CancelOrder("a605f2abbcc750da9138687bb27a2835", YCC_BTC)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestGetOneOrder(t *testing.T) {
	res, err := htb.GetOneOrder("177836e71c8d57a14648d465e893efce", YCC_BTC)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestGetOrders(t *testing.T) {
	res, err := htb.GetOrderHistorys(YCC_BTC, 1, 10)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}

func TestGetUnfinishOrders(t *testing.T) {
	res, err := htb.GetUnfinishOrders(YCC_BTC)
	require := require.New(t)
	require.Nil(err)
	t.Log(res)
}
