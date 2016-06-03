package haobtc

import (
	. "rest"
	"net/http"
	"errors"
	"fmt"
	"net/url"
	"encoding/json"
	"strings"
)

const
(
	EXCHANGE_NAME = "haobtc";
	API_BASE_URL = "https://haobtc.com/exchange/api/v1/";
	TICKER_URI = "ticker";
	DEPTH_URI = "depth?size=%d";
	ACCOUNT_URI = "account_info";
)

type HaoBtc struct {
	httpClient *http.Client;
	accessKey,
	secretKey  string
}

func New(httpClient *http.Client, accessKey, secretKey string) *HaoBtc {
	return &HaoBtc{httpClient, accessKey, secretKey};
}

func (ctx *HaoBtc) buildPostForm(postForm *url.Values) error {
	postForm.Set("api_key", ctx.accessKey);
	//postForm.Set("secret_key", ctx.secret_key);

	payload := postForm.Encode();
	payload = payload + "&secret_key=" + ctx.secretKey;

	sign, err := GetParamMD5Sign(ctx.secretKey, payload);
	if err != nil {
		return err;
	}

	postForm.Set("sign", strings.ToUpper(sign));
	//postForm.Del("secret_key")
	return nil;
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

func (ctx *HaoBtc) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
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

func (ctx *HaoBtc) GetAccount() (*Account, error) {
	postData := url.Values{};
	ctx.buildPostForm(&postData);

	bodyData, err := HttpPostForm(ctx.httpClient, API_BASE_URL + ACCOUNT_URI, postData);
	if err != nil {
		return nil, err;
	}

	fmt.Println(string(bodyData));

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		println(string(bodyData));
		return nil, err;
	}

	if bodyDataMap["code"] != nil {
		return nil, errors.New(string(bodyData));
	}

	account := new(Account);
	account.Exchange = ctx.GetExchangeName();

	var btcSubAccount SubAccount;
	var cnySubAccount SubAccount;

	btcSubAccount.Currency = BTC;
	btcSubAccount.Amount = bodyDataMap["exchange_btc"].(float64);
	btcSubAccount.ForzenAmount = bodyDataMap["exchange_frozen_btc"].(float64);

	cnySubAccount.Currency = CNY;
	cnySubAccount.Amount = bodyDataMap["exchange_cny"].(float64);
	cnySubAccount.ForzenAmount = bodyDataMap["exchange_frozen_cny"].(float64);

	account.SubAccounts = make(map[Currency]SubAccount, 2);
	account.SubAccounts[BTC] = btcSubAccount;
	account.SubAccounts[CNY] = cnySubAccount;

	return account, nil;
}

func (ctx *HaoBtc) GetExchangeName() string {
	return EXCHANGE_NAME;
}