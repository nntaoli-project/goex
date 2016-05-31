package okcoin

import
(
	. "rest"
	"strconv"
)

const
(
	url_ticker = "https://www.okcoin.cn/api/v1/ticker.do";
	url_depth = "https://www.okcoin.cn/api/v1/depth.do";
	url_trades = "https://www.okcoin.cn/api/v1/trades.do";
	url_kline = "https://www.okcoin.cn/api/v1/kline.do";

	url_userinfo = "https://www.okcoin.cn/api/v1/userinfo.do";
	url_trade = "https://www.okcoin.cn/api/v1/trade.do";
	url_cancel_order = "https://www.okcoin.cn/api/v1/cancel_order.do";
	url_order_info = "https://www.okcoin.cn/api/v1/order_info.do";	
)

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
	var tickerMap map[string]interface{};
	var ticker Ticker;
	
	url := url_ticker + "?symbol=" + currencyPair2String(currency);
	bodyDataMap, err := HttpGet(url);
	if err != nil{
		return nil, err;
	}

	tickerMap = bodyDataMap["ticker"].(map[string]interface{});
	ticker.Date, _ = strconv.ParseUint(bodyDataMap["date"].(string), 10, 64);
	ticker.Last, _ = strconv.ParseFloat(tickerMap["last"].(string), 64);
	ticker.Buy, _ = strconv.ParseFloat(tickerMap["buy"].(string), 64);
	ticker.Sell, _ = strconv.ParseFloat(tickerMap["sell"].(string), 64);
	ticker.Low, _ = strconv.ParseFloat(tickerMap["low"].(string), 64);
	ticker.High, _ = strconv.ParseFloat(tickerMap["high"].(string), 64);
	ticker.Vol, _ = strconv.ParseFloat(tickerMap["vol"].(string), 64);

	return &ticker, nil;
}

func (ctx * OKCoinCN_API) GetDepth(size int, currency CurrencyPair) (*Depth, error){
	var depth Depth;
	
	url := url_depth + "?symbol=" + currencyPair2String(currency) + "&size=" + strconv.Itoa(size);
	bodyDataMap, err := HttpGet(url);
	if err != nil {
		return nil, err;
	}

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

func (ctx * OKCoinCN_API) GetExchangeName() string{
	return ctx.name;
}