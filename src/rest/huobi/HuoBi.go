package huobi

import (
	. "rest"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strconv"
	"fmt"
	"net/http"
)

const
(
	API_BASE_URL = "https://api.huobi.com/";
	TICKER_URI = "staticmarket/ticker_%s_json.js";
	DEPTH_URI = "staticmarket/depth_%s_json.js";
)

type HuoBi struct {
	httpClient *http.Client;
	accessKey,
	secretKey  string
}

func New(httpClient *http.Client, accessKey, secretKey string) *HuoBi {
	return &HuoBi{httpClient, accessKey, secretKey};
}

func httpGet(uri string, client *http.Client) (map[string]interface{}, error) {
	url := API_BASE_URL + uri;
	//println(url);
	resp, err := client.Get(url);
	if err != nil {
		return nil, err;
	}

	defer resp.Body.Close();

	body, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return nil, err;
	}

	var bodyDataMap map[string]interface{};
	json.Unmarshal(body, &bodyDataMap);
	return bodyDataMap, nil;
}

func (hb *HuoBi) GetTicker(currency CurrencyPair) (*Ticker, error) {
	var tickerUri string;

	switch currency {
	case BTC_CNY:
		tickerUri = fmt.Sprintf(TICKER_URI, "btc");
	case LTC_CNY:
		tickerUri = fmt.Sprintf(TICKER_URI, "ltc");
	default:
		return nil, errors.New("Unsupport The CurrencyPair");
	}

	bodyDataMap, err := httpGet(tickerUri, hb.httpClient);

	if err != nil {
		return nil, err;
	}

	var tickerMap map[string]interface{};
	var ticker Ticker;

	tickerMap = bodyDataMap["ticker"].(map[string]interface{});

	ticker.Date, _ = strconv.ParseUint(bodyDataMap["time"].(string), 10, 64);
	ticker.Last = tickerMap["last"].(float64);
	ticker.Buy = tickerMap["buy"].(float64);
	ticker.Sell = tickerMap["sell"].(float64);
	ticker.Low = tickerMap["low"].(float64);
	ticker.High = tickerMap["high"].(float64);
	ticker.Vol = tickerMap["vol"].(float64);

	return &ticker, nil;
}

func (hb *HuoBi) GetDepth(size int32, currency CurrencyPair) (*Depth, error) {
	return nil, nil;
}