package spot

import (
	"errors"
	"github.com/nntaoli-project/goex/v2/model"
)

func (s *Spot) GetExchangeInfo() (map[string]model.CurrencyPair, []byte, error) {
	currencyPairM, respBody, err := s.OKxV5.GetExchangeInfo("SPOT")
	s.currencyPairM = currencyPairM
	return currencyPairM, respBody, err
}

func (s *Spot) NewCurrencyPair(baseSym, quoteSym string, opts ...model.OptionParameter) (model.CurrencyPair, error) {
	currencyPair := s.currencyPairM[baseSym+quoteSym]
	if currencyPair.Symbol == "" {
		return currencyPair, errors.New("not found currency pair")
	}
	return currencyPair, nil
}
