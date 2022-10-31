package huobi

import (
	"github.com/nntaoli-project/goex/v2/huobi/futures"
	"github.com/nntaoli-project/goex/v2/huobi/spot"
)

var (
	Spot    = spot.New()
	Futures = futures.NewUSDTFutures()
)
