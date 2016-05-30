package huobi

import (
	. "rest"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strconv"
	"fmt"
)

type HuoBi struct {
	cfg       APIConfig;
	tickerUri string;
}

func New(cfg APIConfig) *HuoBi {
	_hb := new(HuoBi);
	_hb.cfg = cfg;
	_hb.tickerUri = "staticmarket/ticker_%s_json.js";
	return _hb;
}

func (hb HuoBi) GetTicker(currency CurrencyPair) (Ticker, error) {
	var tickerUri string;

	switch currency {
	case BTC_CNY:
		tickerUri = fmt.Sprintf(hb.tickerUri, "btc");
	case LTC_CNY:
		tickerUri = fmt.Sprintf(hb.tickerUri, "ltc");
	default:
		return Ticker{}, errors.New("Unsupport The CurrencyPair");
	}

	url := hb.cfg.ApiUrl + tickerUri;

	//println(url);
	resp, err := hb.cfg.HttpClient.Get(url);
	if err != nil {
		return Ticker{}, errors.New("Get Ticker Error ?");
	}

	defer resp.Body.Close();

	body, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return Ticker{}, errors.New("Read Body Error ?");
	}

	//println(string(body))

	var bodyDataMap map[string]interface{};
	var tickerMap map[string]interface{};
	var ticker Ticker;

	json.Unmarshal(body, &bodyDataMap);

	tickerMap = bodyDataMap["ticker"].(map[string]interface{});

	ticker.Date, _ = strconv.Atoi(bodyDataMap["time"].(string));
	ticker.Last = tickerMap["last"].(float64);
	ticker.Buy = tickerMap["buy"].(float64);
	ticker.Sell = tickerMap["sell"].(float64);
	ticker.Low = tickerMap["low"].(float64);
	ticker.High = tickerMap["high"].(float64);
	ticker.Vol = fmt.Sprintf("%.4f" , tickerMap["vol"].(float64));

	return ticker, nil;
}