package hitbtc

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nntaoli-project/GoEx"
)

const (
	PubKey    = ""
	SecretKey = ""
)

var htb *Hitbtc

func init() {
	htb = New(http.DefaultClient, PubKey, SecretKey)
	pairs, err := htb.GetSymbols()
	if err != nil {
		panic(err)
	}

	for _, pair := range pairs {
		goex.RegisterExSymbol(htb.GetExchangeName(), pair)
	}
}

func TestGetTicker(t *testing.T) {
	res, err := htb.GetTicker(YCC_BTC)
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
