package okcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	. "rest"
	"strconv"
	"net/url"
	"strings"
	"errors"
)

const (
	FUTURE_API_BASE_URL    = "https://www.okcoin.com/api/v1/"
	FUTURE_TICKER_URI      = "future_ticker.do?symbol=%s&contract_type=%s"
	FUTURE_DEPTH_URI       = "future_depth.do?symbol=%s&contract_type=%s"
	FUTURE_USERINFO_URI    = "future_userinfo.do"
	FUTURE_CANCEL_URI      = "future_cancel.do"
	FUTURE_ORDER_INFO_URI  = "future_order_info.do"
	FUTURE_ORDERS_INFO_URI = "future_orders_info.do"
	FUTURE_POSITION_URI    = "future_position.do"
	FUTURE_TRADE_URI       = "future_trade.do"
)

type OKCoinFuture struct {
	apiKey,
	apiSecretKey string
	client *http.Client
}

func NewFuture(client *http.Client, api_key, secret_key string) *OKCoinFuture {
	ok := new(OKCoinFuture)
	ok.apiKey = api_key
	ok.apiSecretKey = secret_key
	ok.client = client
	return ok
}

func (ok *OKCoinFuture) buildPostForm(postForm *url.Values) error {
	postForm.Set("api_key", ok.apiKey);
	//postForm.Set("secret_key", ctx.secret_key);

	payload := postForm.Encode();
	payload = payload + "&secret_key=" + ok.apiSecretKey;

	sign, err := GetParamMD5Sign(ok.apiSecretKey, payload);
	if err != nil {
		return err;
	}

	postForm.Set("sign", strings.ToUpper(sign));
	//postForm.Del("secret_key")
	fmt.Println(postForm)
	return nil;
}

func (ok *OKCoinFuture) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	url := FUTURE_API_BASE_URL + FUTURE_TICKER_URI
	//fmt.Println(fmt.Sprintf(url, CurrencyPairSymbol[currencyPair], contractType));
	resp, err := ok.client.Get(fmt.Sprintf(url, CurrencyPairSymbol[currencyPair], contractType))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	//println(string(body))

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	tickerMap := bodyMap["ticker"].(map[string]interface{})

	ticker := new(Ticker)
	ticker.Date, _ = strconv.ParseUint(bodyMap["date"].(string), 10, 64)
	ticker.Buy = tickerMap["buy"].(float64)
	ticker.Sell = tickerMap["sell"].(float64)
	ticker.Last = tickerMap["last"].(float64)
	ticker.High = tickerMap["high"].(float64)
	ticker.Low = tickerMap["low"].(float64)
	ticker.Vol = tickerMap["vol"].(float64)

	//fmt.Println(bodyMap)
	return ticker, nil
}

func (ok *OKCoinFuture) GetFutureDepth(currencyPair CurrencyPair, contractType string) (*Depth, error) {
	url := FUTURE_API_BASE_URL + FUTURE_DEPTH_URI
	//fmt.Println(fmt.Sprintf(url, CurrencyPairSymbol[currencyPair], contractType));
	resp, err := ok.client.Get(fmt.Sprintf(url, CurrencyPairSymbol[currencyPair], contractType))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	//println(string(body))

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	depth := new(Depth)

	for _, v := range bodyMap["asks"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}

	for _, v := range bodyMap["bids"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.BidList = append(depth.BidList, dr)
	}

	//fmt.Println(bodyMap)
	return depth, nil
}

func (ok *OKCoinFuture) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	return 0, nil
}

type futureUserInfoResponse struct {
	Info   struct {
		       Btc map[string]float64 `json:btc`
		       Ltc map[string]float64 `json:ltc`
	       } `json:info`
	Result bool `json:"result,bool"`
}

func (ok *OKCoinFuture) GetFutureUserinfo() (*FutureAccount, error) {
	userInfoUrl := FUTURE_API_BASE_URL + FUTURE_USERINFO_URI;

	postData := url.Values{};
	ok.buildPostForm(&postData);

	body , err := HttpPostForm(ok.client , userInfoUrl , postData);

	if err != nil {
		return nil , err;
	}

	//println(string(body));
	resp := futureUserInfoResponse{};
	json.Unmarshal(body , &resp)
	if !resp.Result {
		return nil , errors.New(string(body));
	}

	account := new(FutureAccount);
	account.FutureSubAccounts = make(map[Currency]FutureSubAccount , 2);

	btcMap := resp.Info.Btc;
	ltcMap := resp.Info.Ltc;

	account.FutureSubAccounts[BTC] = FutureSubAccount{BTC, btcMap["account_rights"], btcMap["keep_deposit"], btcMap["profit_real"], btcMap["profit_unreal"], btcMap["risk_rate"]};
	account.FutureSubAccounts[LTC] = FutureSubAccount{LTC, ltcMap["account_rights"], ltcMap["keep_deposit"], ltcMap["profit_real"], ltcMap["profit_unreal"], ltcMap["risk_rate"]};

	return account, nil
}

func (ok *OKCoinFuture) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice , leverRate int) (string, error) {
	postData := url.Values{};
	postData.Set("symbol" , CurrencyPairSymbol[currencyPair]);
	postData.Set("price" , price);
	postData.Set("contract_type" , contractType);
	postData.Set("amount" , amount);
	postData.Set("type" , strconv.Itoa(openType));
	postData.Set("lever_rate" , strconv.Itoa(leverRate))
	postData.Set("match_price" , strconv.Itoa(matchPrice));

	ok.buildPostForm(&postData);

	placeOrderUrl := FUTURE_API_BASE_URL + FUTURE_TRADE_URI;
	body , err := HttpPostForm(ok.client , placeOrderUrl , postData);

	if err != nil {
		return "" , err;
	}

	respMap := make(map[string]interface{});
	json.Unmarshal(body , &respMap);

	println(string(body));

	if !respMap["result"].(bool) {
		return "" , errors.New(string(body));
	}

	return fmt.Sprintf("%.0f" , respMap["order_id"].(float64)) , nil
}

func (ok *OKCoinFuture) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	postData := url.Values{};
	postData.Set("symbol" , CurrencyPairSymbol[currencyPair]);
	postData.Set("order_id" , orderId);
	postData.Set("contract_type" , contractType);

	ok.buildPostForm(&postData);

	cancelUrl := FUTURE_API_BASE_URL + FUTURE_CANCEL_URI;

	body , err := HttpPostForm(ok.client , cancelUrl , postData);
	if err != nil {
		return false , err;
	}

	respMap := make(map[string]interface{});
	json.Unmarshal(body , &respMap);

	if respMap["result"] != nil && !respMap["result"].(bool) {
		return false , errors.New(string(body));
	}

	return true, nil
}

func (ok *OKCoinFuture) GetFuturePosition(currencyPair CurrencyPair, contractType string) (*FuturePosition, error) {
	positionUrl := FUTURE_API_BASE_URL + FUTURE_POSITION_URI;

	postData := url.Values{};
	postData.Set("contract_type" , contractType);
	postData.Set("symbol" , CurrencyPairSymbol[currencyPair]);

	ok.buildPostForm(&postData);

	body , err := HttpPostForm(ok.client , positionUrl , postData);

	if err != nil {
		return nil , err;
	}

	respMap := make(map[string]interface{});

	json.Unmarshal(body , &respMap);

	if !respMap["result"].(bool) {
		return nil , errors.New(string(body));
	}
	println(string(body))

	pos := new(FuturePosition);
	pos.ForceLiquPrice , _ = strconv.ParseFloat(respMap["force_liqu_price"].(string) , 64);

	/*
	holding := respMap["holding"].([]interface{});
	pos.LeverRate = holding["lever_rate"].(int);
	pos.ContractType = holding["contract_type"].(string);
	pos.ContractId = holding["contract_id"].(int64);
	pos.BuyAmount = holding["buy_amount"].(float64);
	pos.BuyAvailable = holding["buy_available"].(float64);
	pos.BuyPriceAvg = holding["buy_price_avg"].(float64);
	pos.BuyPriceCost = holding["buy_price_cost"].(float64);
	pos.BuyProfitReal = holding["buy_profit_real"].(float64);
	pos.SellAmount = holding["sell_amount"].(float64);
	pos.SellAvailable = holding["sell_available"].(float64);
	pos.SellPriceAvg = holding["sell_price_avg"].(float64);
	pos.SellPriceCost = holding["sell_price_cost"].(float64);
	pos.SellProfitReal = holding["sell_profit_real"].(float64);
	pos.CreateDate = holding["create_date"].(int64);
	pos.Symbol = currencyPair;
        */
	return pos, nil
}

func (ok *OKCoinFuture) GetFutureOrders(orderId int64, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return nil, nil
}

func (ok *OKCoinFuture) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return nil, nil
}
