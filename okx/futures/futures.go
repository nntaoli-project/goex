package futures

import (
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/okx/common"
	. "github.com/nntaoli-project/goex/v2/options"
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

func (f *Futures) NewCrossFuturesTradeApi(apiOpt ...ApiOption) ITradeRest {
	ft := newFCrossTrade(apiOpt...)
	ft.V5 = f.V5
	return ft
}

func (f *Futures) NewIsolatedFuturesTradeApi(apiOpt ...ApiOption) ITradeRest {
	ft := newFIsolatedTrade(apiOpt...)
	ft.V5 = f.V5
	return ft
}
