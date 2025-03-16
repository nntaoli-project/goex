package spot

import (
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/okx/common"
	"github.com/nntaoli-project/goex/v2/options"
)

type Spot struct {
	*common.OKxV5
	currencyPairM map[string]CurrencyPair
}

func New() *Spot {
	v5 := common.New()
	currencyPairCacheMap := make(map[string]CurrencyPair, 64)
	return &Spot{v5, currencyPairCacheMap}
}

func (s *Spot) NewPrvApi(apiOps ...options.ApiOption) *PrvApi {
	prv := new(PrvApi)
	prv.Prv = s.OKxV5.NewPrvApi(apiOps...)
	prv.Prv.OKxV5 = s.OKxV5
	return prv
}
