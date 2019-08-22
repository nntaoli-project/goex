package fcoin

import (
	"github.com/nntaoli-project/GoEx"
	"testing"
)

var fm = &FCoinMargin{ft}

func TestFCoinMargin_GetMarginAccount(t *testing.T) {
	t.Log(fm.GetMarginAccount(goex.BTC_USDT))
}
