package goex

import (
	"strings"
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
	btccny := NewCurrencyPair(NewCurrency("btc", ""), NewCurrency("cny", ""))
	t.Log(strings.ToUpper(btccny.String()))
	t.Log(BTC_CNY)
}
func TestFormatPairs_String(t *testing.T) {
	p := NewCurrencyPairsFromStrings("BTC_USD,LTC_USD,LTC_BTC,ETH_USD,ETH_BTC,ETC_BTC,ETC_USD", ",")
	t.Log(p)
	t.Log(Contains(p, BTC_USD))
}
