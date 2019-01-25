package huobi

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"internal/log"
	"sort"
	"time"
)

type Hbdm struct {
	config *APIConfig
}

type BaseResponse struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int64  `json:"ts"`
	ErrMsg string `json:"err-msg"`
}

var (
	apiUrl = "https://api.hbdm.com"
)

func NewHbdm(conf *APIConfig) *Hbdm {
	return &Hbdm{conf}
}

func (dm *Hbdm) GetExchangeName() string {
	return HBDM
}

func (dm *Hbdm) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	ret, err := HttpGet(dm.config.HttpClient, apiUrl+"/api/v1//contract_delivery_price?symbol="+currencyPair.CurrencyA.Symbol)
	if err != nil {
		return -1, err
	}

	if ret["status"].(string) != "ok" {
		return -1, errors.New(fmt.Sprintf("%+v", ret))
	}

	return ToFloat64(ret["data"].(map[string]interface{})["delivery_price"]), nil
}

func (dm *Hbdm) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	symbol := dm.adaptSymbol(currencyPair, contractType)
	ret, err := HttpGet(dm.config.HttpClient, apiUrl+"/market/detail/merged?symbol="+symbol)
	if err != nil {
		return nil, err
	}
	//log.Println(ret)
	s := ret["status"].(string)
	if s == "error" {
		return nil, errors.New(ret["err-msg"].(string))
	}

	tick := ret["tick"].(map[string]interface{})
	return &Ticker{
		Pair: currencyPair,
		Last: ToFloat64(tick["close"]),
		Vol:  ToFloat64(tick["amount"]),
		Low:  ToFloat64(tick["low"]),
		High: ToFloat64(tick["high"]),
		Sell: ToFloat64(tick["ask"].([]interface{})[0]),
		Buy:  ToFloat64(tick["bid"].([]interface{})[0]),
		Date: ToUint64(ret["ts"])}, nil
}

func (dm *Hbdm) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	symbol := dm.adaptSymbol(currencyPair, contractType)
	url := apiUrl + "/market/depth?type=step0&symbol=" + symbol
	ret, err := HttpGet(dm.config.HttpClient, url)
	if err != nil {
		return nil, err
	}

	s := ret["status"].(string)
	if s == "error" {
		return nil, errors.New(ret["err-msg"].(string))
	}
	log.Println(ret)
	dep := new(Depth)
	dep.Pair = currencyPair
	dep.ContractType = symbol

	mills := ToUint64(ret["ts"])
	dep.UTime = time.Unix(int64(mills/1000), int64(mills%1000)*int64(time.Millisecond))

	tick := ret["tick"].(map[string]interface{})
	asks := tick["asks"].([]interface{})
	bids := tick["bids"].([]interface{})

	for _, item := range asks {
		askItem := item.([]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(askItem[0]), ToFloat64(askItem[1])})
	}

	for _, item := range bids {
		bidItem := item.([]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(bidItem[0]), ToFloat64(bidItem[1])})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, nil
}

func (dm *Hbdm) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	ret, err := HttpGet(dm.config.HttpClient, apiUrl+"/api/v1/contract_index?symbol="+currencyPair.CurrencyA.Symbol)
	if err != nil {
		return -1, err
	}

	if ret["status"].(string) != "ok" {
		return -1, errors.New(fmt.Sprintf("%+v", ret))
	}

	datamap := ret["data"].([]interface{})
	index := datamap[0].(map[string]interface{})["index_price"]
	return ToFloat64(index), nil
}

func (dm *Hbdm) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	symbol := dm.adaptSymbol(currency, contract_type)
	periodS := dm.adaptKLinePeriod(period)
	url := fmt.Sprintf("%s/market/history/kline?symbol=%s&period=%s&size=%d", apiUrl, symbol, periodS, size)

	var ret struct {
		BaseResponse
		Data []struct {
			Id     int64   `json:"id"`
			Amount float64 `json:"amount"`
			Close  float64 `json:"close"`
			High   float64 `json:"high"`
			Low    float64 `json:"low"`
			Open   float64 `json:"open"`
			Vol    float64 `json:"vol"`
		} `json:"data"`
	}

	err := HttpGet4(dm.config.HttpClient, url, nil, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Status != "ok" {
		return nil, errors.New(ret.ErrMsg)
	}

	var klines []FutureKline
	for _, d := range ret.Data {
		klines = append(klines, FutureKline{
			Kline: &Kline{
				Pair:      currency,
				Vol:       d.Vol,
				Open:      d.Open,
				Close:     d.Close,
				High:      d.High,
				Low:       d.Low,
				Timestamp: d.Id},
			Vol2: d.Vol})
	}

	return klines, nil
}

func (dm *Hbdm) adaptSymbol(pair CurrencyPair, contractType string) string {
	symbol := pair.CurrencyA.Symbol + "_"
	switch contractType {
	case THIS_WEEK_CONTRACT:
		symbol += "CW"
	case NEXT_WEEK_CONTRACT:
		symbol += "NW"
	case QUARTER_CONTRACT:
		symbol += "CQ"
	}
	return symbol
}

func (dm *Hbdm) adaptKLinePeriod(period int) string {
	switch period {
	case KLINE_PERIOD_1MIN:
		return "1min"
	case KLINE_PERIOD_5MIN:
		return "5min"
	case KLINE_PERIOD_15MIN:
		return "15min"
	case KLINE_PERIOD_30MIN:
		return "30min"
	case KLINE_PERIOD_60MIN:
		return "60min"
	case KLINE_PERIOD_1H:
		return "1h"
	case KLINE_PERIOD_4H:
		return "4h"
	case KLINE_PERIOD_1DAY:
		return "1day"
	case KLINE_PERIOD_1WEEK:
		return "1week"
	case KLINE_PERIOD_1MONTH:
		return "1mon"
	default:
		return "1day"
	}
}
