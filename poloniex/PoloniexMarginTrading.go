package poloniex

import (
	"encoding/json"
	"errors"
	. "github.com/nntaoli-project/goex"
	"log"
	"net/url"
)

type PoloniexGenericResponse struct {
	Success int    `json:"success"`
	Error   string `json:"error"`
}

type PoloniexMarginPosition struct {
	Amount            float64 `json:"amount,string"`
	Total             float64 `json:"total,string"`
	BasePrice         float64 `json:"basePrice,string"`
	LiquidiationPrice float64 `json:"liquidiationPrice"`
	ProfitLoss        float64 `json:"pl,string"`
	LendingFees       float64 `json:"lendingFees,string"`
	Type              string  `json:"type"`
}

func (poloniex *Poloniex) MarginLimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return poloniex.placeLimitOrder("marginBuy", amount, price, currency)
}

func (poloniex *Poloniex) MarginLimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return poloniex.placeLimitOrder("marginSell", amount, price, currency)
}

func (poloniex *Poloniex) GetMarginPosition(currency CurrencyPair) (*PoloniexMarginPosition, error) {
	values := url.Values{}
	values.Set("command", "getMarginPosition")
	values.Set("currencyPair", currency.AdaptUsdToUsdt().Reverse().ToSymbol("_"))
	result := PoloniexMarginPosition{}
	err := poloniex.sendAuthenticatedRequest(values, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (poloniex *Poloniex) CloseMarginPosition(currency CurrencyPair) (bool, error) {
	values := url.Values{}
	values.Set("command", "closeMarginPosition")
	values.Set("currencyPair", currency.AdaptUsdToUsdt().Reverse().ToSymbol("_"))
	result := PoloniexGenericResponse{}
	err := poloniex.sendAuthenticatedRequest(values, &result)
	if err != nil {
		return false, err
	}
	if result.Success == 0 {
		return false, errors.New(result.Error)
	}
	return true, nil
}

func (poloniex *Poloniex) sendAuthenticatedRequest(values url.Values, result interface{}) error {
	sign, _ := poloniex.buildPostForm(&values)

	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}

	resp, err := HttpPostForm2(poloniex.client, TRADE_API, values, headers)
	if err != nil {
		log.Println(err)
		return err
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return errors.New("Unable to JSON Unmarshal response.")
	}

	return nil
}
