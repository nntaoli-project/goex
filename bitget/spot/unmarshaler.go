package spot

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/nntaoli-project/goex/v2/model"
	"github.com/spf13/cast"
)

type RespUnmarshaler struct {
}

func (u *RespUnmarshaler) UnmarshalGetExchangeInfoResponse(data []byte) (map[string]model.CurrencyPair, error) {
	currencyPairs := make(map[string]model.CurrencyPair)

	code, _ := jsonparser.GetString(data, "code")
	if code != "00000" {
		msg, _ := jsonparser.GetString(data, "msg")
		return currencyPairs, errors.New(msg)
	}

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		currencyPair := model.CurrencyPair{}
		symbol, _ := jsonparser.GetString(value, "symbol")
		baseCoin, _ := jsonparser.GetString(value, "baseCoin")
		quoteCoin, _ := jsonparser.GetString(value, "quoteCoin")
		pricePrecision, _ := jsonparser.GetString(value, "pricePrecision")
		quantityPrecision, _ := jsonparser.GetString(value, "quantityPrecision")

		currencyPair.Symbol = symbol
		currencyPair.BaseSymbol = baseCoin
		currencyPair.QuoteSymbol = quoteCoin
		currencyPair.PricePrecision = cast.ToInt(pricePrecision)
		currencyPair.QtyPrecision = cast.ToInt(quantityPrecision)

		k := fmt.Sprintf("%s%s", currencyPair.BaseSymbol, currencyPair.QuoteSymbol)
		currencyPairs[k] = currencyPair
	}, "data")

	return currencyPairs, err
}
