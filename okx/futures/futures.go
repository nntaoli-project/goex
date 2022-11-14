package futures

import (
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/okx/common"
)

type Futures struct {
	*common.V5
}

func NewFutures() *Futures {
	return &Futures{V5: common.New()}
}

func (f *Futures) MarketApi() IMarketRest {
	return f.V5.MarketApi()
}
