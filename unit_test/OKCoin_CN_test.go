package unit

import (
	. "../"
	"github.com/stretchr/testify/assert"
	"testing"
	"../okcoin"
	"net/http"
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

func Test_GetOnOrder(t *testing.T)  {
	api := okcoin.New(http.DefaultClient , "api_key" , "secret_key");
	order , err := api.GetOneOrder("3503527666" , BTC_CNY);
	assert.NoError(t , err);
	t.Log(order);
}


func Test_GetUnfinishOrder(t *testing.T)  {
	api := okcoin.New(http.DefaultClient , "api_key" , "secret_key");
	order , err := api.GetUnfinishOrders(BTC_CNY);
	assert.NoError(t , err);
	t.Log(order);
}

func Test_LimitBuy(t *testing.T){
	api := okcoin.New(http.DefaultClient , "api_key" , "secret_key");
	order , err := api.LimitBuy("0.01" , "3100" , BTC_CNY);
	assert.NoError(t , err);
	t.Log(order);
}

func Test_CancelOrder(t *testing.T){
	api := okcoin.New(http.DefaultClient , "api_key" , "secret_key");
	order , err := api.CancelOrder("3595088073" , BTC_CNY);
	assert.NoError(t , err);
	t.Log(order);
}

func Test_GetKlineRecords(t *testing.T)  {
	api := okcoin.New(http.DefaultClient , "" , "");
	kline , err := api.GetKlineRecords(BTC_CNY  , "1min" ,10 , -1);
	assert.NoError(t , err);
	t.Log(kline);
}

func Test_GetOrderHistorys(t *testing.T)  {
	api := okcoin.New(http.DefaultClient , "" , "");
	orders , err := api.GetOrderHistorys(LTC_CNY  , 1 , 100);
	assert.NoError(t , err);
	t.Log("size:" , len(orders));
	t.Log(orders);
}