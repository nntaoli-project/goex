package unit

import (
	. "rest"
	"github.com/stretchr/testify/assert"
	"testing"
	"rest/okcoin"
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
	t.Log("bids:", depth.AskList);
	t.Log("asks:", depth.AskList);
}

func Test_GetOnOrder(t *testing.T)  {
	//8bdf9b3e-756e-11e4-acdf-00163e0238cc
	//3A523ECEE1F5F59A23131CA030CFE116
	api := okcoin.New(http.DefaultClient , "api_key" , "secret_key");
	order , err := api.GetOneOrder("3503527666" , BTC_CNY);
	assert.NoError(t , err);
	t.Log(order);
}


func Test_GetUnfinishOrder(t *testing.T)  {
	//8bdf9b3e-756e-11e4-acdf-00163e0238cc
	//3A523ECEE1F5F59A23131CA030CFE116
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