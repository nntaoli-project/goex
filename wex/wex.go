package wex

import (
	"github.com/btcsuite/goleveldb/leveldb/errors"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"strings"
)

//https://wex.nz

type Wex struct {
	client *http.Client
	accesskey,
	secretkey string
}

var (
	baseurl = "https://wex.nz/api/3"
)

func New(client *http.Client, accesskey, secretkey string) *Wex {
	return &Wex{client, accesskey, secretkey}
}

func (wex *Wex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (wex *Wex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (wex *Wex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (wex *Wex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (wex *Wex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (wex *Wex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (wex *Wex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (wex *Wex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (wex *Wex) GetAccount() (*Account, error) {
	panic("not implements")
}

func (wex *Wex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	respmap, err := HttpGet(wex.client, baseurl+"/ticker/"+strings.ToLower(currency.ToSymbol("_")))
	if err != nil {
		return nil, err
	}

	if errmsg, isok := respmap["error"].(string); isok {
		log.Println(errmsg)
		return nil, errors.New(errmsg)
	}

	for _, v := range respmap {
		tickermap := v.(map[string]interface{})
		return &Ticker{
			Low:  ToFloat64(tickermap["low"]),
			Buy:  ToFloat64(tickermap["buy"]),
			Sell: ToFloat64(tickermap["sell"]),
			Last: ToFloat64(tickermap["last"]),
			Vol:  ToFloat64(tickermap["vol_cur"]),
			High: ToFloat64(tickermap["high"])}, nil
	}

	return nil, nil
}

func (wex *Wex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	return nil, nil
}

func (wex *Wex) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (wex *Wex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (wex *Wex) GetExchangeName() string {
	return "wex.nz"
}
