package okx

import (
	"github.com/nntaoli-project/goex/v2/okx/futures"
	"github.com/nntaoli-project/goex/v2/okx/spot"
)

var (
	Spot    = spot.New()
	Futures = futures.NewFutures()
)
