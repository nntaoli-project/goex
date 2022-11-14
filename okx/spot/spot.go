package spot

import (
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/okx/common"
)

var (
	currencyPaircacheMap = make(map[string]*CurrencyPair, 6)
)

type Spot struct {
	*common.V5
}

func New() *Spot {
	v5 := common.New()
	return &Spot{v5}
}

func (s *Spot) MarketApi() IMarketRest {
	return s.V5.MarketApi()
}
