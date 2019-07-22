package okex

import (
	"github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

//
var config2 = &goex.APIConfig{
	Endpoint: "https://www.okex.com",
	//HttpClient: &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: func(req *http.Request) (*url.URL, error) {
	//			return &url.URL{
	//				Scheme: "socks5",
	//				Host:   "127.0.0.1:1080"}, nil
	//		},
	//	},
	//}, //需要代理的这样配置
	HttpClient:    http.DefaultClient,
	ApiKey:        "",
	ApiSecretKey:  "",
	ApiPassphrase: "",
}

var okex = NewOKEx(config2) //线上请用APIBuilder构建

func TestOKExSpot_GetAccount(t *testing.T) {
	t.Log(okex.GetAccount())
}

func TestOKExSpot_BatchPlaceOrders(t *testing.T) {
	t.Log(okex.OKExSpot.BatchPlaceOrders([]goex.Order{
		goex.Order{
			Cid:       okex.UUID(),
			Currency:  goex.XRP_USD,
			Amount:    10,
			Price:     0.32,
			Side:      goex.BUY,
			Type:      "limit",
			OrderType: goex.ORDINARY,
		},
		{
			Cid:       okex.UUID(),
			Currency:  goex.EOS_USD,
			Amount:    1,
			Price:     5.2,
			Side:      goex.BUY,
			OrderType: goex.ORDINARY,
		},
		goex.Order{
			Cid:       okex.UUID(),
			Currency:  goex.XRP_USD,
			Amount:    10,
			Price:     0.33,
			Side:      goex.BUY,
			Type:      "limit",
			OrderType: goex.ORDINARY,
		}}))
}

func TestOKExSpot_LimitBuy(t *testing.T) {
	t.Log(okex.OKExSpot.LimitBuy("0.001", "9910", goex.BTC_USD))
}

func TestOKExSpot_CancelOrder(t *testing.T) {
	t.Log(okex.OKExSpot.CancelOrder("2a647e51435647708b1c840802bf70e5", goex.BTC_USD))

}

func TestOKExSpot_GetOneOrder(t *testing.T) {
	t.Log(okex.OKExSpot.GetOneOrder("42152275c599444aa8ec1d33bd8003fb", goex.BTC_USD))
}

func TestOKExSpot_GetUnfinishOrders(t *testing.T) {
	t.Log(okex.OKExSpot.GetUnfinishOrders(goex.EOS_BTC))
}

func TestOKExSpot_GetTicker(t *testing.T) {
	t.Log(okex.OKExSpot.GetTicker(goex.BTC_USD))
}

func TestOKExSpot_GetDepth(t *testing.T) {
	dep, err := okex.OKExSpot.GetDepth(2, goex.EOS_BTC)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}

func TestOKExFuture_GetFutureTicker(t *testing.T) {
	t.Log(okex.OKExFuture.GetFutureTicker(goex.BTC_USD, "BTC-USD-190927"))
	t.Log(okex.OKExFuture.GetFutureTicker(goex.BTC_USD, goex.QUARTER_CONTRACT))
	t.Log(okex.OKExFuture.GetFutureDepth(goex.BTC_USD, goex.QUARTER_CONTRACT, 2))
	t.Log(okex.OKExFuture.GetContractValue(goex.XRP_USD))
	t.Log(okex.OKExFuture.GetFutureIndex(goex.EOS_USD))
	t.Log(okex.OKExFuture.GetFutureEstimatedPrice(goex.EOS_USD))
}

func TestOKExFuture_GetFutureUserinfo(t *testing.T) {
	t.Log(okex.OKExFuture.GetFutureUserinfo())
}

func TestOKExFuture_GetFuturePosition(t *testing.T) {
	t.Log(okex.OKExFuture.GetFuturePosition(goex.EOS_USD, goex.QUARTER_CONTRACT))
}

func TestOKExFuture_PlaceFutureOrder(t *testing.T) {
	t.Log(okex.OKExFuture.PlaceFutureOrder(goex.EOS_USD, goex.THIS_WEEK_CONTRACT, "5.8", "1", goex.OPEN_BUY, 0, 10))
}

func TestOKExFuture_PlaceFutureOrder2(t *testing.T) {
	t.Log(okex.OKExFuture.PlaceFutureOrder2(0, &goex.FutureOrder{
		Currency:     goex.EOS_USD,
		ContractName: goex.QUARTER_CONTRACT,
		OType:        goex.OPEN_BUY,
		OrderType:    goex.ORDINARY,
		Price:        5.9,
		Amount:       10,
		LeverRate:    10}))
}

func TestOKExFuture_FutureCancelOrder(t *testing.T) {
	t.Log(okex.OKExFuture.FutureCancelOrder(goex.EOS_USD, goex.QUARTER_CONTRACT, "e88bd3361de94512b8acaf9aa154f95a"))
}

func TestOKExFuture_GetFutureOrder(t *testing.T) {
	t.Log(okex.OKExFuture.GetFutureOrder("3145664744431616", goex.EOS_USD, goex.QUARTER_CONTRACT))
}

func TestOKExFuture_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(okex.OKExFuture.GetUnfinishFutureOrders(goex.EOS_USD, goex.QUARTER_CONTRACT))
}

func TestOKExFuture_MarketCloseAllPosition(t *testing.T) {
	t.Log(okex.OKExFuture.MarketCloseAllPosition(goex.BTC_USD, goex.THIS_WEEK_CONTRACT, goex.CLOSE_BUY))
}

func TestOKExFuture_GetRate(t *testing.T) {
	t.Log(okex.OKExFuture.GetRate())
}

func TestOKExFuture_GetKlineRecords(t *testing.T) {
	since := time.Now().Add(-24 * time.Hour).Unix()
	kline, err := okex.OKExFuture.GetKlineRecords(goex.QUARTER_CONTRACT, goex.BTC_USD, goex.KLINE_PERIOD_4H, 0, int(since))
	assert.Nil(t, err)
	for _, k := range kline {
		t.Logf("%+v", k.Kline)
	}
}

func TestOKExWallet_GetAccount(t *testing.T) {
	t.Log(okex.OKExWallet.GetAccount())
}

func TestOKExWallet_Transfer(t *testing.T) {
	t.Log(okex.OKExWallet.Transfer(TransferParameter{
		Currency:     goex.EOS.Symbol,
		From:         SPOT,
		To:           SPOT_MARGIN,
		Amount:       20,
		InstrumentId: goex.EOS_USDT.ToLower().ToSymbol("-")}))
}

func TestOKExWallet_Withdrawal(t *testing.T) {
	t.Log(okex.OKExWallet.Withdrawal(WithdrawParameter{
		Currency:    goex.EOS.Symbol,
		Amount:      100,
		Destination: 2,
		ToAddress:   "",
		TradePwd:    "",
		Fee:         "0.01",
	}))
}

func TestOKExWallet_GetDepositAddress(t *testing.T) {
	t.Log(okex.OKExWallet.GetDepositAddress(goex.BTC))
}

func TestOKExWallet_GetWithDrawalFee(t *testing.T) {
	t.Log(okex.OKExWallet.GetWithDrawalFee(nil))
}

func TestOKExWallet_GetDepositHistory(t *testing.T) {
	t.Log(okex.OKExWallet.GetDepositHistory(&goex.BTC))
}

func TestOKExWallet_GetWithDrawalHistory(t *testing.T) {
	t.Log(okex.OKExWallet.GetWithDrawalHistory(&goex.XRP))
}

func TestOKExMargin_GetMarginAccount(t *testing.T) {
	t.Log(okex.OKExMargin.GetMarginAccount(goex.EOS_USDT))
}

func TestOKExMargin_Borrow(t *testing.T) {
	t.Log(okex.OKExMargin.Borrow(goex.BorrowParameter{
		Currency:     goex.EOS,
		CurrencyPair: goex.EOS_USDT,
		Amount:       10,
	}))
}

func TestOKExMargin_Repayment(t *testing.T) {
	t.Log(okex.OKExMargin.Repayment(goex.RepaymentParameter{
		BorrowParameter: goex.BorrowParameter{
			Currency:     goex.EOS,
			CurrencyPair: goex.EOS_USDT,
			Amount:       10},
		BorrowId: "123"}))
}

func TestOKExMargin_PlaceOrder(t *testing.T) {
	t.Log(okex.OKExMargin.PlaceOrder(&goex.Order{
		Currency:  goex.EOS_USDT,
		Amount:    0.2,
		Price:     6,
		Type:      "limit",
		OrderType: goex.ORDINARY,
		Side:      goex.SELL,
	}))
}

func TestOKExMargin_GetUnfinishOrders(t *testing.T) {
	t.Log(okex.OKExMargin.GetUnfinishOrders(goex.EOS_USDT))
}

func TestOKExMargin_CancelOrder(t *testing.T) {
	t.Log(okex.OKExMargin.CancelOrder("3174778420532224", goex.EOS_USDT))
}

func TestOKExMargin_GetOneOrder(t *testing.T) {
	t.Log(okex.OKExMargin.GetOneOrder("3174778420532224", goex.EOS_USDT))
}
