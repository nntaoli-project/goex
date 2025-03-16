package futures

import (
	"errors"
	"github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/okx/common"
	"github.com/nntaoli-project/goex/v2/options"
)

type Futures struct {
	*common.OKxV5
	currencyPairM map[string]model.CurrencyPair
}

func New() *Futures {
	currencyPairM := make(map[string]model.CurrencyPair, 64)
	return &Futures{OKxV5: common.New(), currencyPairM: currencyPairM}
}

func (f *Futures) NewPrvApi(apiOpts ...options.ApiOption) *PrvApi {
	return NewPrvApi(f.OKxV5, apiOpts...)
}

func (f *Futures) GetExchangeInfo() (map[string]model.CurrencyPair, []byte, error) {
	m, b, er := f.OKxV5.GetExchangeInfo("FUTURES")
	f.currencyPairM = m
	return m, b, er
}

func (f *Futures) NewCurrencyPair(baseSym, quoteSym string, opts ...model.OptionParameter) (model.CurrencyPair, error) {
	if len(opts) >= 1 && opts[0].Key == "contractAlias" {
		contractAlias := opts[0].Value
		currencyPair := f.currencyPairM[baseSym+quoteSym+contractAlias]
		if currencyPair.Symbol != "" {
			return currencyPair, nil
		}
		return currencyPair, errors.New("not found currency pair")
	}
	return model.CurrencyPair{}, errors.New("please input contract alias option parameter")
}
