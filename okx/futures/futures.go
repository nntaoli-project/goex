package futures

import (
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

func (f *Futures) NewCurrencyPair(baseSym, quoteSym, contractAlias string) model.CurrencyPair {
	return f.currencyPairM[baseSym+quoteSym+contractAlias]
}
