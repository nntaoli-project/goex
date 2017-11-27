package cexio

import (
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strconv"
)

const (
	EXCHANGE_NAME = "cex.io"
)

var (
	balanceCurrencies = map[string]Currency{
		"USD":  USD,
		"EUR":  EUR,
		"GBP":  NewCurrency("GBP", ""),
		"RUB":  NewCurrency("RUB", ""),
		"GHS":  NewCurrency("GHS", ""),
		"BTC":  BTC,
		"BCH":  BCH,
		"ETH":  ETH,
		"DASH": NewCurrency("DASH", ""),
		"ZEC":  ZEC,
	}
)

type Cex struct {
	httpClient *http.Client
	userId     string
	Api_key    string
	Api_secret string
}

func New(client *http.Client, accessKey, secretKey string, clientId string) *Cex {
	return &Cex{client, clientId, accessKey, secretKey}
}

func (cex *Cex) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (cex *Cex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	pair := currency.ToSymbol("/")
	resp, err := cex.ticker(pair)
	if err != nil {
		return nil, err
	}
	ticker := new(Ticker)
	intTime, _ := strconv.Atoi(resp["timestamp"].(string))
	ticker.Last = ToFloat64(resp["last"])
	ticker.Vol = ToFloat64(resp["volume"])
	ticker.High = ToFloat64(resp["high"])
	ticker.Low = ToFloat64(resp["low"])
	ticker.Sell = ToFloat64(resp["ask"])
	ticker.Buy = ToFloat64(resp["bid"])
	ticker.Date = uint64(intTime)
	return ticker, nil
}

func (cex *Cex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}

func (cex *Cex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}

func (cex *Cex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (cex *Cex) GetAccount() (*Account, error) {
	resp, err := cex.balance()
	if err != nil {
		return nil, err
	}
	//fmt.Println(resp)

	acc := &Account{}
	acc.Exchange = cex.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	for name, currency := range balanceCurrencies {
		data := resp[name].(map[string]interface{})
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(data["available"]),
			ForzenAmount: ToFloat64(data["orders"]),
			LoanAmount:   0,
		}
	}

	return acc, nil
}

func (cex *Cex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	pair := currency.ToSymbol("/")
	resp, err := cex.orderBook(size, pair)
	if err != nil {
		return nil, err
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	depth := new(Depth)

	for _, bid := range bids {
		_bid := bid.([]interface{})
		amount := ToFloat64(_bid[1].(float64))
		price := ToFloat64(_bid[0].(float64))
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range asks {
		_ask := ask.([]interface{})
		amount := ToFloat64(_ask[1].(float64))
		price := ToFloat64(_ask[0].(float64))
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (cex *Cex) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (cex *Cex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
