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
