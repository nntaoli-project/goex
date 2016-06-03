package haobtc

import (
	. "rest"
	//"io/ioutil"
	//"errors"
	//"encoding/json"
	//"strconv"
	//"fmt"
	//"net/url"
	"net/http"
	//"time"
	//"strings"
	"errors"
	"fmt"
)

const
(
	EXCHANGE_NAME = "haobtc";
	API_BASE_URL = "https://haobtc.com/exchange/api/v1/";
	TICKER_URI = "ticker";
	DEPTH_URI = "depth?size=%d";
)

type HaoBtc struct {
	httpClient *http.Client;
	accessKey,
	secretKey  string
}

func New(httpClient *http.Client, accessKey, secretKey string) *HaoBtc {
	return &HaoBtc{httpClient, accessKey, secretKey};
}

func (ctx *HaoBtc) GetTicker(currency CurrencyPair) (*Ticker, error) {
	if currency != BTC_CNY {
		return nil, errors.New("The HaoBtc Unsupport " + CurrencyPairSymbol[currency]);
	}

	var tickerMap map[string]interface{};
	var ticker Ticker;

	url := API_BASE_URL + TICKER_URI;

	bodyDataMap, err := HttpGet(url);
	if err != nil {
		return nil, err;
	}

	tickerMap = bodyDataMap["ticker"].(map[string]interface{});
	ticker.Date = uint64(bodyDataMap["date"].(float64));
	ticker.Last = tickerMap["last"].(float64);
	ticker.Buy = tickerMap["buy"].(float64);
	ticker.Sell = tickerMap["sell"].(float64);
	ticker.Low = tickerMap["low"].(float64);
	ticker.High = tickerMap["high"].(float64);
	ticker.Vol = tickerMap["vol"].(float64);

	return &ticker, nil;
}

func (hb *HaoBtc) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depthUri string;

	switch currency {
	case BTC_CNY:
		depthUri = API_BASE_URL + fmt.Sprintf(DEPTH_URI, size);
	default:
		return nil, errors.New("Unsupport The CurrencyPair " + CurrencyPairSymbol[currency]);
	}

	bodyDataMap, err := HttpGet(depthUri);

	if err != nil {
		return nil, err;
	}

	var depth Depth;

	for _, v := range bodyDataMap["asks"].([]interface{}) {
		var dr DepthRecord;
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64);
			case 1:
				dr.Amount = vv.(float64);
			}
		}
		depth.AskList = append(depth.AskList, dr);
	}

	for _, v := range bodyDataMap["bids"].([]interface{}) {
		var dr DepthRecord;
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64);
			case 1:
				dr.Amount = vv.(float64);
			}
		}
		depth.BidList = append(depth.BidList, dr);
	}

	return &depth, nil;
}
