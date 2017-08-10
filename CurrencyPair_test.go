package goex

import (
	"strings"
	"testing"
)

func TestCurrency2_String(t *testing.T) {
	btc := NewCurrency("btc", "bitcoin")
	ltc := NewCurrency("ltc", "litecoin")
	t.Log(btc)
	t.Log(ltc)
}

func TestCurrencyPair2_String(t *testing.T) {
	btccny := NewCurrencyPair(NewCurrency("btc", ""), NewCurrency("cny", ""))
	t.Log(strings.ToUpper(btccny.String()))
	t.Log(BTC_CNY)
}
