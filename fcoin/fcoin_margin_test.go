package fcoin

import (
	"github.com/nntaoli-project/GoEx"
	"testing"
)

var fm = &FCoinMargin{ft}

func TestFCoinMargin_Borrow(t *testing.T) {
	//return
	t.Log(fm.Borrow(goex.BorrowParameter{
		CurrencyPair: goex.BTC_USDT,
		Currency:     goex.USDT,
		Amount:       100,
	}))
}

func TestFCoinMargin_Repayment(t *testing.T) {
	//return
	t.Log(fm.Repayment(goex.RepaymentParameter{
		BorrowParameter: goex.BorrowParameter{
			CurrencyPair: goex.BTC_USDT,
			Currency:     goex.USDT,
			Amount:       100.065,
		},
		BorrowId: "uQ7Gzird8kW0rbsC9Cu-RlcY7cGgrog23dEVugBh9JA",
	}))
}

func TestFCoinMargin_AssetTransferIn(t *testing.T) {
	//return
	t.Log(fm.AssetTransferIn(goex.USDT, "80", ASSETS, goex.BTC_USDT))
}

func TestFCoinMargin_GetMarginAccount(t *testing.T) {
	//return
	t.Log(fm.GetMarginAccount(goex.BTC_USDT))
}
