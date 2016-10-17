package unit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"../okcoin"
	"net/http"
	. "../"
)

func Test_OKCoin_Future_GetTicker(t *testing.T) {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	ticker, err := api.GetFutureTicker(BTC_USD, THIS_WEEK_CONTRACT);
	assert.NoError(t, err);
	t.Log(ticker)
}

func Test_OKCoin_Future_GetDepth(t *testing.T) {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	depth, err := api.GetFutureDepth(BTC_USD, NEXT_WEEK_CONTRACT , 2);
	assert.NoError(t, err);
	t.Log(depth)
}

func Test_OKCoin_Future_GetUserInfo(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	user , err := api.GetFutureUserinfo();
	assert.NoError(t , err);
	t.Log(user);
}

func Test_OKCoin_Future_GetFuturePosition(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient,"", "");
	pos , err := api.GetFuturePosition(LTC_USD , THIS_WEEK_CONTRACT);
	assert.NoError(t , err);
	t.Log(pos);
}


func Test_OKCoin_Future_CancelOrder(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	pos , err := api.FutureCancelOrder(LTC_USD , THIS_WEEK_CONTRACT , "1991174604");
	assert.NoError(t , err);
	t.Log(pos);
}


func Test_OKCoin_Future_PlaceOrder(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient,"", "");
	pos , err := api.PlaceFutureOrder(LTC_USD , THIS_WEEK_CONTRACT , "5" , "1" , OPEN_BUY , 0 , 10);
	assert.NoError(t , err);
	t.Log(pos);
}

func Test_OKCoin_Future_GetOrder(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient,"", "");
	order , err := api.GetFutureOrders(2027122589 , LTC_USD , "this_week");
	assert.NoError(t , err);
	t.Log(order);
}

func Test_OKCoin_Future_GetUnfinishOrder(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient,  "", "");
	order , err := api.GetUnfinishFutureOrders(LTC_USD , "this_week");
	assert.NoError(t , err);
	t.Log(order);
}