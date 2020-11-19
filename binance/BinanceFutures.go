package binance

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"net/http"
	"sort"
	"sync"
	"time"
)

type SymbolInfo struct {
	Symbol         string
	Pair           string
	ContractType   string `json:"contractType"`
	DeliveryDate   int64  `json:"deliveryDate"`
	ContractStatus string `json:"contractStatus"`
	ContractSize   int    `json:"contractSize"`
	PricePrecision int    `json:"pricePrecision"`
}

type BinanceFutures struct {
	base         Binance
	exchangeInfo *struct {
		Symbols []SymbolInfo `json:"symbols"`
	}
}

func NewBinanceFutures(config *APIConfig) *BinanceFutures {
	if config.Endpoint == "" {
		config.Endpoint = "https://dapi.binance.com"
	}

	if config.HttpClient == nil {
		config.HttpClient = http.DefaultClient
	}

	bs := &BinanceFutures{
		base: Binance{
			baseUrl:    config.Endpoint,
			accessKey:  config.ApiKey,
			apiV1:      config.Endpoint + "/dapi/v1/",
			secretKey:  config.ApiSecretKey,
			httpClient: config.HttpClient,
		},
	}

	go bs.GetExchangeInfo()

	return bs
}

func (bs *BinanceFutures) SetBaseUri(uri string) {
	bs.base.baseUrl = uri
}

func (bs *BinanceFutures) GetExchangeName() string {
	return BINANCE_FUTURES
}

func (bs *BinanceFutures) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return nil, err
	}

	ticker24hrUri := bs.base.apiV1 + "ticker/24hr?symbol=" + symbol
	tickerBookUri := bs.base.apiV1 + "ticker/bookTicker?symbol=" + symbol

	var (
		ticker24HrResp []interface{}
		tickerBookResp []interface{}
		err1           error
		err2           error
		wg             = sync.WaitGroup{}
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		ticker24HrResp, err1 = HttpGet3(bs.base.httpClient, ticker24hrUri, map[string]string{})
	}()

	go func() {
		defer wg.Done()
		tickerBookResp, err2 = HttpGet3(bs.base.httpClient, tickerBookUri, map[string]string{})
	}()

	wg.Wait()

	if err1 != nil {
		return nil, err1
	}

	if err2 != nil {
		return nil, err2
	}

	if len(ticker24HrResp) == 0 {
		return nil, errors.New("response is empty")
	}

	if len(tickerBookResp) == 0 {
		return nil, errors.New("response is empty")
	}

	ticker24HrMap := ticker24HrResp[0].(map[string]interface{})
	tickerBookMap := tickerBookResp[0].(map[string]interface{})

	var ticker Ticker
	ticker.Pair = currencyPair
	ticker.Date = ToUint64(tickerBookMap["time"])
	ticker.Last = ToFloat64(ticker24HrMap["lastPrice"])
	ticker.Buy = ToFloat64(tickerBookMap["bidPrice"])
	ticker.Sell = ToFloat64(tickerBookMap["askPrice"])
	ticker.High = ToFloat64(ticker24HrMap["highPrice"])
	ticker.Low = ToFloat64(ticker24HrMap["lowPrice"])
	ticker.Vol = ToFloat64(ticker24HrMap["volume"])

	return &ticker, nil
}

func (bs *BinanceFutures) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return nil, err
	}

	limit := 5
	if size <= 5 {
		limit = 5
	} else if size <= 10 {
		limit = 10
	} else if size <= 20 {
		limit = 20
	} else if size <= 50 {
		limit = 50
	} else if size <= 100 {
		limit = 100
	} else if size <= 500 {
		limit = 500
	} else {
		limit = 1000
	}

	depthUri := bs.base.apiV1 + "depth?symbol=%s&limit=%d"

	ret, err := HttpGet(bs.base.httpClient, fmt.Sprintf(depthUri, symbol, limit))
	if err != nil {
		return nil, err
	}
	logger.Debug(ret)

	var dep Depth

	dep.ContractType = contractType
	dep.Pair = currencyPair
	eT := int64(ret["E"].(float64))
	dep.UTime = time.Unix(0, eT*int64(time.Millisecond))

	for _, item := range ret["asks"].([]interface{}) {
		ask := item.([]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{
			Price:  ToFloat64(ask[0]),
			Amount: ToFloat64(ask[1]),
		})
	}

	for _, item := range ret["bids"].([]interface{}) {
		bid := item.([]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{
			Price:  ToFloat64(bid[0]),
			Amount: ToFloat64(bid[1]),
		})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return &dep, nil
}

func (bs *BinanceFutures) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetFutureUserinfo(currencyPair ...CurrencyPair) (*FutureAccount, error) {
	accountUri := bs.base.apiV1 + "account"
	resp, err := HttpGet5(bs.base.httpClient, accountUri, map[string]string{})
	if err != nil {
		return nil, err
	}
	logger.Debug(string(resp))
	return nil, nil
}

func (bs *BinanceFutures) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) LimitFuturesOrder(currencyPair CurrencyPair, contractType, price, amount string, openType int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) MarketFuturesOrder(currencyPair CurrencyPair, contractType, amount string, openType int) (*FutureOrder, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetFee() (float64, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetDeliveryTime() (int, int, int, int) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetKlineRecords(contractType string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

func (bs *BinanceFutures) GetExchangeInfo() {
	exchangeInfoUri := bs.base.apiV1 + "exchangeInfo"
	ret, err := HttpGet5(bs.base.httpClient, exchangeInfoUri, map[string]string{})
	if err != nil {
		logger.Error("[exchangeInfo] Http Error", err)
		return
	}

	err = json.Unmarshal(ret, &bs.exchangeInfo)
	if err != nil {
		logger.Error("json unmarshal response content error , content= ", string(ret))
		return
	}

	logger.Debug("[ExchangeInfo]", bs.exchangeInfo)
}

func (bs *BinanceFutures) adaptToSymbol(pair CurrencyPair, contractType string) (string, error) {
	if contractType == THIS_WEEK_CONTRACT || contractType == NEXT_WEEK_CONTRACT {
		return "", errors.New("binance only support contract quarter or bi_quarter")
	}

	if contractType == SWAP_CONTRACT {
		return fmt.Sprintf("%s_PERP", pair.AdaptUsdtToUsd().ToSymbol("")), nil
	}

	if bs.exchangeInfo == nil || len(bs.exchangeInfo.Symbols) == 0 {
		bs.GetExchangeInfo()
	}

	for _, info := range bs.exchangeInfo.Symbols {
		if info.ContractType != "PERPETUAL" &&
			info.ContractStatus == "TRADING" &&
			info.DeliveryDate <= time.Now().Unix()*1000 {
			logger.Debugf("pair=%s , contractType=%s, delivery date = %d ,  now= %d", info.Pair, info.ContractType, info.DeliveryDate, time.Now().Unix()*1000)
			bs.GetExchangeInfo()
		}

		if info.Pair == pair.ToSymbol("") {
			if info.ContractStatus != "TRADING" {
				return "", errors.New("contract status " + info.ContractStatus)
			}

			if info.ContractType == "CURRENT_QUARTER" && contractType == QUARTER_CONTRACT {
				return info.Symbol, nil
			}

			if info.ContractType == "NEXT_QUARTER" && contractType == BI_QUARTER_CONTRACT {
				return info.Symbol, nil
			}
		}
	}

	return "", errors.New("binance not support " + pair.ToSymbol("") + " " + contractType)
}
