package unit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"rest/okcoin"
	"net/http"
	"rest"
)

func Test_OKCoin_Future_GetTicker(t *testing.T) {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	ticker, err := api.GetFutureTicker(rest.BTC_USD, "this_week");
	assert.NoError(t, err);
	t.Log(ticker)
}

func Test_OKCoin_Future_GetDepth(t *testing.T) {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	depth, err := api.GetFutureDepth(rest.BTC_USD, "this_week");
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
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	pos , err := api.GetFuturePosition(rest.LTC_USD , "this_week");
	assert.NoError(t , err);
	t.Log(pos);
}


func Test_OKCoin_Future_CancelOrder(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	pos , err := api.FutureCancelOrder(rest.LTC_USD , "this_week" , "1991174604");
	assert.NoError(t , err);
	t.Log(pos);
}


func Test_OKCoin_Future_PlaceOrder(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	pos , err := api.PlaceFutureOrder(rest.LTC_USD , "this_week" , "5.1" , "1" , rest.OPEN_BUY , 0 , 10);
	assert.NoError(t , err);
	t.Log(pos);
}

func Test_OKCoin_Future_GetOrder(t *testing.T)  {
	api := okcoin.NewFuture(http.DefaultClient, "", "");
	order , err := api.GetFutureOrders(1991135170 , rest.LTC_USD , "this_week");
	assert.NoError(t , err);
	t.Log(order);
}