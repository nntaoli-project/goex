package okcoin

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	apiKey       = getEnv("GOEX_OKEX_API_KEY", "")
	apiSecretKey = getEnv("GOEX_OKEX_API_SECRET_KEY", "")
	passphrase   = getEnv("GOEX_OKEX_PASSPHRASE", "")
	endpoint     = getEnv("GOEX_OKEX_RESTFUL_URL", "https://www.okex.me/")
	authed       = len(apiKey) > 0 && len(apiSecretKey) > 0 && len(passphrase) > 0
	okexV3       = NewOKExV3(http.DefaultClient, apiKey, apiSecretKey, passphrase, endpoint)
)

func TestOKExV3_GetFutureDepth(t *testing.T) {
	size := 10
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		dep, err := okexV3.GetFutureDepth(BTC_USD, QUARTER_CONTRACT, size)
		assert.Nil(t, err)
		t.Log(dep)
	}()
	go func() {
		defer wg.Done()
		dep, err := okexV3.GetFutureDepth(BTC_USD, SWAP_CONTRACT, size)
		assert.Nil(t, err)
		t.Log(dep)
	}()
	wg.Wait()
}

func TestOKExV3_GetFutureTicker(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		ticker, err := okexV3.GetFutureTicker(BTC_USD, QUARTER_CONTRACT)
		assert.Nil(t, err)
		t.Log(ticker)
	}()
	go func() {
		defer wg.Done()
		ticker, err := okexV3.GetFutureTicker(BTC_USD, SWAP_CONTRACT)
		assert.Nil(t, err)
		t.Log(ticker)
	}()
	wg.Wait()
}

func testPlaceAndCancel(t *testing.T, currencyPair CurrencyPair, contractType string) {
	// 以100档对手价下买单然后马上撤掉
	leverage := 20
	depth := 100
	dep, err := okexV3.GetFutureDepth(currencyPair, contractType, depth)
	assert.Nil(t, err)
	price := fmt.Sprintf("%f", dep.BidList[depth-1].Price)
	symbol, err := okexV3.GetContract(currencyPair, contractType)
	assert.Nil(t, err)
	amount := symbol.getSizeIncrement()
	orderID, err := okexV3.PlaceFutureOrder(
		currencyPair, contractType, price, amount, OPEN_BUY, 0, leverage)
	assert.Nil(t, err)
	t.Log(orderID)
	order, err := okexV3.GetFutureOrder(orderID, currencyPair, contractType)
	assert.Nil(t, err)
	t.Log(order)
	cancelled, err := okexV3.FutureCancelOrder(currencyPair, contractType, orderID)
	assert.Nil(t, err)
	assert.True(t, cancelled)
	order, err = okexV3.GetFutureOrder(orderID, currencyPair, contractType)
	assert.Nil(t, err)
	t.Log(order)
}

func TestOKExV3_PlaceAndCancelFutureOrder(t *testing.T) {
	if authed {
		testPlaceAndCancel(t, EOS_USD, QUARTER_CONTRACT)
		testPlaceAndCancel(t, EOS_USD, SWAP_CONTRACT)
	} else {
		t.Log("not authed, skip test place and cancel future order")
	}
}

func testPlaceAndGetInfo(t *testing.T, currencyPair CurrencyPair, contractType string) {
	leverage := 20
	depth := 100
	dep, err := okexV3.GetFutureDepth(currencyPair, contractType, depth)
	assert.Nil(t, err)
	price := fmt.Sprintf("%f", dep.BidList[depth-1].Price)
	symbol, err := okexV3.GetContract(currencyPair, contractType)
	assert.Nil(t, err)
	amount := symbol.getSizeIncrement()
	orderID1, err := okexV3.PlaceFutureOrder(
		currencyPair, contractType, price, amount, OPEN_BUY, 0, leverage)
	assert.Nil(t, err)
	t.Log(orderID1)
	orderID2, err := okexV3.PlaceFutureOrder(
		currencyPair, contractType, price, amount, OPEN_BUY, 0, leverage)
	assert.Nil(t, err)
	t.Log(orderID2)
	// get info of order1
	order, err := okexV3.GetFutureOrder(orderID1, currencyPair, contractType)
	assert.Nil(t, err)
	t.Log(order)
	order, err = okexV3.GetFutureOrder(orderID2, currencyPair, contractType)
	assert.Nil(t, err)
	t.Log(order)
	// sleep for a while when place order,
	time.Sleep(1 * time.Second)
	// get infos of order1 and order2
	orders, err := okexV3.GetFutureOrders([]string{orderID1, orderID2}, currencyPair, contractType)
	assert.Nil(t, err)
	assert.True(t, len(orders) == 2)
	t.Log(orders)
	//cancel order1 and order2
	cancelled, err := okexV3.FutureCancelOrder(currencyPair, contractType, orderID1)
	assert.Nil(t, err)
	assert.True(t, cancelled)
	cancelled, err = okexV3.FutureCancelOrder(currencyPair, contractType, orderID2)
	assert.Nil(t, err)
	assert.True(t, cancelled)
	// sleep for a while when cancel order,
	// the order info will be missing in orders api for all states for a short time.
	time.Sleep(2 * time.Second)
	// get infos of order1 and order2 after cancelling
	orders, err = okexV3.GetFutureOrders([]string{orderID1, orderID2}, currencyPair, contractType)
	assert.Nil(t, err)
	assert.True(t, len(orders) == 2)
	t.Log(orders)
}

func TestOKEV3_PlaceAndGetOrdersInfo(t *testing.T) {
	if authed {
		testPlaceAndGetInfo(t, EOS_USD, QUARTER_CONTRACT)
		testPlaceAndGetInfo(t, EOS_USD, SWAP_CONTRACT)
	} else {
		t.Log("not authed, skip test place future order and get order info")
	}
}

func TestOKEV3_GetFutureEstimatedPrice(t *testing.T) {
	f, err := okexV3.GetFutureEstimatedPrice(EOS_USD)
	assert.Nil(t, err)
	t.Log(f)
}

func TestOKEV3_GetFee(t *testing.T) {
	f, err := okexV3.GetFee()
	assert.Nil(t, err)
	assert.True(t, f > 0)
	t.Log(f)
}

func testGetContractValue(t *testing.T, currencyPair CurrencyPair) {
	f, err := okexV3.GetContractValue(currencyPair)
	assert.Nil(t, err)
	assert.True(t, f > 0)
	t.Log(f)
}

func TestOKEV3_GetContractValue(t *testing.T) {
	testGetContractValue(t, BTC_USD)
	testGetContractValue(t, EOS_USD)
}

func isEqualDiff(klines []FutureKline, seconds int64) bool {
	miliseconds := seconds * 1000
	for i := 0; i < len(klines) - 1; i ++ {
		diff := klines[i + 1].Timestamp-klines[i].Timestamp
		if diff != miliseconds {
			return false
		}
	}
	return true
}

func testGetKlineRecords(t *testing.T, contractType string, currency CurrencyPair, maxSize int) {
	now := time.Now().UTC()
	timestamp := (now.UnixNano() - 20*int64(time.Hour)) / int64(time.Millisecond)
	size := 10
	period := KLINE_PERIOD_1MIN
	seconds := int64(60)
	kline, err := okexV3.GetKlineRecords(contractType, currency, period, size, 0)
	assert.Nil(t, err)
	assert.True(t, len(kline) == size)
	t.Log(len(kline))
	assert.True(t, isEqualDiff(kline, seconds))
	kline, err = okexV3.GetKlineRecords(contractType, currency, period, size, int(timestamp))
	assert.Nil(t, err)
	assert.True(t, len(kline) == size)
	t.Log(len(kline))
	assert.True(t, isEqualDiff(kline, seconds))
	size = maxSize
	kline, err = okexV3.GetKlineRecords(contractType, currency, period, size, 0)
	assert.Nil(t, err)
	assert.True(t, len(kline) == size)
	t.Log(len(kline))
	assert.True(t, isEqualDiff(kline, seconds))
	kline, err = okexV3.GetKlineRecords(contractType, currency, period, size, int(timestamp))
	assert.Nil(t, err)
	assert.True(t, len(kline) == size)
	t.Log(len(kline))
	assert.True(t, isEqualDiff(kline, seconds))
	size = 3 * maxSize
	kline, err = okexV3.GetKlineRecords(contractType, currency, period, size, 0)
	assert.Nil(t, err)
	assert.True(t, len(kline) == size)
	t.Log(len(kline))
	assert.True(t, isEqualDiff(kline, seconds))
	kline, err = okexV3.GetKlineRecords(contractType, currency, period, size, int(timestamp))
	assert.Nil(t, err)
	assert.True(t, len(kline) == size)
	t.Log(len(kline))
	assert.True(t, isEqualDiff(kline, seconds))
}

func TestOKEV3_GetKlineRecords(t *testing.T) {
	testGetKlineRecords(t, QUARTER_CONTRACT, EOS_USD, 300)
	testGetKlineRecords(t, SWAP_CONTRACT, EOS_USD, 200)
}

func testGetFutureIndex(t *testing.T, currencyPair CurrencyPair) {
	f, err := okexV3.GetFutureIndex(currencyPair)
	assert.Nil(t, err)
	assert.True(t, f > 0)
	t.Log(f)
}

func TestOKEV3_GetFutureIndex(t *testing.T) {
	testGetFutureIndex(t, BTC_USD)
	testGetFutureIndex(t, EOS_USD)
}

func testGetFuturePosition(t *testing.T, currencyPair CurrencyPair, contractType string) {
	ps, err := okexV3.GetFuturePosition(currencyPair, contractType)
	assert.Nil(t, err)
	t.Log(ps)
}

func TestOKEV3_GetFuturePosition(t *testing.T) {
	if authed {
		testGetFuturePosition(t, EOS_USD, QUARTER_CONTRACT)
		testGetFuturePosition(t, EOS_USD, SWAP_CONTRACT)
	} else {
		t.Log("not authed, skip test get future position")
	}
}

func testGetTrades(t *testing.T, contractType string, currencyPair CurrencyPair) {
	trades, err := okexV3.GetTrades(contractType, currencyPair, 0)
	assert.Nil(t, err)
	t.Log(trades[0])
}

func TestOKEV3_GetTrades(t *testing.T) {
	testGetTrades(t, QUARTER_CONTRACT, EOS_USD)
	testGetTrades(t, SWAP_CONTRACT, EOS_USD)
}

func testGetFutureUserinfo(t *testing.T) {
	currencies := okexV3.getAllCurrencies()
	account, err := okexV3.GetFutureUserinfo()
	assert.Nil(t, err)
	assert.True(t, len(currencies) == len(account.FutureSubAccounts))
	t.Log(account)
}

func TestOKEV3_GetFutureUserinfo(t *testing.T) {
	if authed {
		testGetFutureUserinfo(t)
	} else {
		t.Log("not authed, skip test get future userinfo")
	}
}