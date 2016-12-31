package unit

import (
	. "github.com/nntaoli/crypto_coin_api"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/nntaoli/crypto_coin_api/okcoin"
	"net/http"
)

var (
	_apikey = ""
	_apiscretkey = ""
	okCoinCN = okcoin.New(http.DefaultClient, _apikey, _apiscretkey)
)

func Test_OKCoin_CN(t *testing.T) {
	var api API;
	api = okcoin.New(http.DefaultClient, "", "");
	tk, err := api.GetTicker(BTC_CNY);
	if err != nil {
		t.Logf("%s", err.Error());
	}
	assert.True(t, err == nil);
	t.Logf("last:%f buy:%f sell:%f high:%f low:%f vol:%f date:%d",
		tk.Last, tk.Buy, tk.Sell, tk.High, tk.Low, tk.Vol, tk.Date);

	depth, err := api.GetDepth(3, LTC_CNY);
	assert.True(t, err == nil);
	t.Log("bids:", depth.BidList);
	t.Log("asks:", depth.AskList);
}

func Test_GetOnOrder(t *testing.T) {
	order, err := okCoinCN.GetOneOrder("354392174", LTC_CNY);
	assert.NoError(t, err);
	t.Log(order);
}

func Test_GetUnfinishOrder(t *testing.T) {
	order, err := okCoinCN.GetUnfinishOrders(BTC_CNY);
	assert.NoError(t, err);
	t.Log(order);
}

func Test_LimitBuy(t *testing.T) {
	order, err := okCoinCN.LimitBuy("1", "26", LTC_CNY);
	assert.NoError(t, err);
	t.Log(order);
}

func Test_CancelOrder(t *testing.T) {
	order, err := okCoinCN.CancelOrder("354385312", LTC_CNY);
	assert.NoError(t, err);
	t.Log(order);
}

func Test_OKCoinCN_MarketBuy(t *testing.T) {
	order, err := okCoinCN.MarketBuy("", "3", LTC_CNY)
	assert.NoError(t, err)
	t.Log(order)
}

func Test_OKCoinCN_MarketSell(t *testing.T) {
	order, err := okCoinCN.MarketSell("0.1", "3", LTC_CNY)
	assert.NoError(t, err)
	t.Log(order)
}

func Test_GetKlineRecords(t *testing.T) {
	kline, err := okCoinCN.GetKlineRecords(BTC_CNY, "1min", 10, -1);
	assert.NoError(t, err);
	t.Log(kline);
}

func Test_GetOrderHistorys(t *testing.T) {
	orders, err := okCoinCN.GetOrderHistorys(LTC_CNY, 1, 100);
	assert.NoError(t, err);
	t.Log("size:", len(orders));
	t.Log(orders);
}