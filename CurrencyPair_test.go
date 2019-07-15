package goex

import (
	"testing"
)

func TestCurrency2_String(t *testing.T) {
	btc := NewCurrency("btc", "bitcoin")
	btc2 := Currency{"BTC", "bitcoin.org"}
	ltc := NewCurrency("ltc", "litecoin")
	t.Log(btc == BTC)
	t.Log(ltc.Desc, btc.Desc)
	t.Log(btc == btc2)
}

func TestCurrencyPair2_String(t *testing.T) {
	btc_usd := NewCurrencyPair(NewCurrency("btc", ""), NewCurrency("usd", ""))
	t.Log(btc_usd.String() == "BTC_USD")
	t.Log(btc_usd.ToLower().ToSymbol("") == "btcusd")
	t.Log(btc_usd.ToLower().String() == "btc_usd")
	t.Log(btc_usd.Reverse().String() == "USD_BTC")
	t.Log(btc_usd.Eq(BTC_USD))
}
