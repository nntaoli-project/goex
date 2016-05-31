package okcoin

import
(
	. "rest"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

const url_ticker = "https://www.okcoin.cn/api/v1/ticker.do";
const url_depth = "https://www.okcoin.cn/api/v1/depth.do";
const url_trades = "https://www.okcoin.cn/api/v1/trades.do";
const url_kline = "https://www.okcoin.cn/api/v1/kline.do";

const url_userinfo = "https://www.okcoin.cn/api/v1/userinfo.do";
const url_trade = "https://www.okcoin.cn/api/v1/trade.do";
const url_cancel_order = "https://www.okcoin.cn/api/v1/cancel_order.do";
const url_order_info = "https://www.okcoin.cn/api/v1/order_info.do";

type OKCoinCN_API struct{
	name string
	api_key string
	secret_key string
}

func currencyPair2String(currency CurrencyPair) string{
	switch currency{
		case BTC_CNY:
			return "btc_cny";
		case LTC_CNY:
			return "ltc_cny";
		default:
			return "";
	}
}

func New(name, api_key, secret_key string) * OKCoinCN_API{
	return &OKCoinCN_API{name, api_key, secret_key};
}

func (ctx * OKCoinCN_API) LimitBuy(amount, price string, currency CurrencyPair) (string, error){
	return "", nil;
}

func (ctx * OKCoinCN_API) LimitSell(amount, price string, currency CurrencyPair) (string, error){
	return "", nil;
}

func (ctx * OKCoinCN_API) CancelOrder(orderId string, currency CurrencyPair) (string, error){
	return "", nil;
}

func (ctx * OKCoinCN_API) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error){
	return nil, nil;
}

func (ctx * OKCoinCN_API) GetUnfinishOrders(currency CurrencyPair) ([]Order, error){
	return nil, nil;
}

func (ctx * OKCoinCN_API) GetAccount() (*Account, error){
	return nil, nil;
}

func (ctx * OKCoinCN_API) GetTicker(currency CurrencyPair) (*Ticker, error){
	type ticker_data struct{
		Last string `json:"last"`
		Buy  string `json:"buy"`
		Sell string `json:"sell"`
		High string `json:"high"`
		Low  string `json:"low"`
		Vol  string `json:"vol"`
	}
	type ticker struct{
		Date string `json:"date"`
		Data ticker_data `json:"ticker"`
	}
	var tk ticker;
	url := url_ticker + "?symbol=" + currencyPair2String(currency);
	resp, err := http.Get(url);
	if err != nil{
		return nil, err;
	}
	defer resp.Body.Close();
	body, err := ioutil.ReadAll(resp.Body);
	if err != nil{
		return nil, err;
	}
	err = json.Unmarshal(body, &tk);
	if err != nil{
		return nil, err;
	}
	last, err := strconv.ParseFloat(tk.Data.Last, 64);
	if err != nil{
		return nil, err;
	}
	buy, err := strconv.ParseFloat(tk.Data.Buy, 64);
	if err != nil{
		return nil, err;
	}
	sell, err := strconv.ParseFloat(tk.Data.Sell, 64);
	if err != nil{
		return nil, err;
	}
	high, err := strconv.ParseFloat(tk.Data.High, 64);
	if err != nil{
		return nil, err;
	}
	low, err := strconv.ParseFloat(tk.Data.Low, 64);
	if err != nil{
		return nil, err;
	}
	vol, err := strconv.ParseFloat(tk.Data.Vol, 64);
	if err != nil{
		return nil, err;
	}
	date, err := strconv.ParseUint(tk.Date, 10, 64);
	if err != nil{
		return nil, err;
	}	
	return &Ticker{last, buy, sell, high, low, vol, date}, nil;
}

func (ctx * OKCoinCN_API) GetDepth(size int32, currency CurrencyPair) (*Depth, error){
	return nil, nil;
}

func (ctx * OKCoinCN_API) GetExchangeName() string{
	return ctx.name;
}