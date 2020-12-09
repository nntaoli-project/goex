package binance

import (
	"fmt"
	"github.com/nntaoli-project/goex"
	"strings"
)

func adaptStreamToCurrencyPair(stream string) goex.CurrencyPair {
	symbol := strings.Split(stream, "@")[0]

	if strings.HasSuffix(symbol, "usdt") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_usdt", strings.TrimSuffix(symbol, "usdt")))
	}

	if strings.HasSuffix(symbol, "usd") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_usd", strings.TrimSuffix(symbol, "usd")))
	}

	if strings.HasSuffix(symbol, "btc") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_btc", strings.TrimSuffix(symbol, "btc")))
	}

	return goex.UNKNOWN_PAIR
}

func adaptSymbolToCurrencyPair(symbol string) goex.CurrencyPair {
	symbol = strings.ToUpper(symbol)

	if strings.HasSuffix(symbol, "USD") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_USD", strings.TrimSuffix(symbol, "USD")))
	}

	if strings.HasSuffix(symbol, "USDT") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_USDT", strings.TrimSuffix(symbol, "USDT")))
	}

	if strings.HasSuffix(symbol, "PAX") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_PAX", strings.TrimSuffix(symbol, "PAX")))
	}

	if strings.HasSuffix(symbol, "BTC") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_BTC", strings.TrimSuffix(symbol, "BTC")))
	}

	return goex.UNKNOWN_PAIR
}
