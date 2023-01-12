package futures

import (
	"github.com/nntaoli-project/goex/v2/okx/common"
	"github.com/nntaoli-project/goex/v2/options"
)

type Futures struct {
	*common.OKxV5
}

func New() *Futures {
	return &Futures{OKxV5: common.New()}
}

func (f *Futures) NewPrvApi(apiOpts ...options.ApiOption) *PrvApi {
	return NewPrvApi(f.OKxV5, apiOpts...)
}
