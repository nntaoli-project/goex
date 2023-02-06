package spot

import "github.com/nntaoli-project/goex/v2/model"

func (s *Spot) GetExchangeInfo() (map[string]model.CurrencyPair, []byte, error) {
	currencyPairM, respBody, err := s.OKxV5.GetExchangeInfo("SPOT")
	s.currencyPairM = currencyPairM
	return currencyPairM, respBody, err
}

func (s *Spot) NewCurrencyPair(baseSym, quoteSym string) model.CurrencyPair {
	return s.currencyPairM[baseSym+quoteSym]
}
