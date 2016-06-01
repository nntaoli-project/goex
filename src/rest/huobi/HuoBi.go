package huobi

import (
	. "rest"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strconv"
	"fmt"
	"net/url"
	"net/http"
	"time"
)

const
(
	EXCHANGE_NAME = "huobi";
	API_BASE_URL = "https://api.huobi.com/";
	TRADE_API_V3 = "apiv3";
	TICKER_URI = "staticmarket/ticker_%s_json.js";
	DEPTH_URI = "staticmarket/depth_%s_%d.js";
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

func (hb *HuoBi) GetExchangeName() string {
	return EXCHANGE_NAME;
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

func (hb *HuoBi) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depthUri string;

	switch currency {
	case BTC_CNY:
		depthUri = fmt.Sprintf(DEPTH_URI, "btc", size);
	case LTC_CNY:
		depthUri = fmt.Sprintf(DEPTH_URI, "ltc", size);
	default:
		return nil, errors.New("Unsupport The CurrencyPair");
	}

	bodyDataMap, err := httpGet(depthUri, hb.httpClient);

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

func (hb *HuoBi) GetAccount() (*Account, error) {
	accountUrl := API_BASE_URL + TRADE_API_V3;

	postData := url.Values{};
	postData.Set("method", "get_account_info");
	postData.Set("created", fmt.Sprintf("%d", time.Now().Unix()));
	postData.Set("access_key", hb.accessKey);
	postData.Set("secret_key", hb.secretKey);

	sign, _ := GetParamMD5Sign(hb.secretKey, postData.Encode());
	postData.Set("sign", sign);
	postData.Del("secret_key");

	bodyDataMap, err := HttpPostForm(hb.httpClient, accountUrl, postData);
	if err != nil {
		return nil, err;
	}

	//fmt.Println(bodyDataMap);

	account := new(Account);
	account.Exchange = hb.GetExchangeName();
	account.Asset, _ = strconv.ParseFloat(bodyDataMap["total"].(string), 64);
	account.NetAsset, _ = strconv.ParseFloat(bodyDataMap["net_asset"].(string), 64);

	var btcSubAccount SubAccount;
	var ltcSubAccount SubAccount;
	var cnySubAccount SubAccount;

	btcSubAccount.Currency = BTC;
	btcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_btc_display"].(string), 64);
	btcSubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_btc_display"].(string), 64);
	btcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_btc_display"].(string), 64);

	ltcSubAccount.Currency = LTC;
	ltcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_ltc_display"].(string), 64);
	ltcSubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_ltc_display"].(string), 64);
	ltcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_ltc_display"].(string), 64);

	cnySubAccount.Currency = CNY;
	cnySubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_cny_display"].(string), 64);
	cnySubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_cny_display"].(string), 64);
	cnySubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_cny_display"].(string), 64);

	account.SubAccounts = make(map[Currency]SubAccount, 3);
	account.SubAccounts[BTC] = btcSubAccount;
	account.SubAccounts[LTC] = ltcSubAccount;
	account.SubAccounts[CNY] = cnySubAccount;

	return account, nil;
}