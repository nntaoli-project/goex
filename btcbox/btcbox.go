package btcbox

import (
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strings"
)

var (
	baseurl = "https://www.btcbox.co.jp/api/v1/"
)

type BtcBox struct {
	client *http.Client
	accessKey,
	secretkey string
}

func New(client *http.Client, apikey, secretkey string) *BtcBox {
	return &BtcBox{client: client, accessKey: apikey, secretkey: secretkey}
}

func (btcbox *BtcBox) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (btcbox *BtcBox) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (btcbox *BtcBox) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (btcbox *BtcBox) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (btcbox *BtcBox) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (btcbox *BtcBox) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (btcbox *BtcBox) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (btcbox *BtcBox) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (btcbox *BtcBox) GetAccount() (*Account, error) {
	panic("not implements")
}

func (btcbox *BtcBox) GetTicker(currency CurrencyPair) (*Ticker, error) {
	respmap, err := HttpGet(btcbox.client, baseurl+"ticker?coin="+strings.ToLower(currency.CurrencyA.Symbol))
	if err != nil {
		return nil, err
	}

	return &Ticker{
		Low:  ToFloat64(respmap["low"]),
		Buy:  ToFloat64(respmap["buy"]),
		Sell: ToFloat64(respmap["sell"]),
		Last: ToFloat64(respmap["last"]),
		Vol:  ToFloat64(respmap["vol"]),
		High: ToFloat64(respmap["high"])}, nil
}

func (btcbox *BtcBox) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	respmap, err := HttpGet(btcbox.client, baseurl+"depth?coin="+strings.ToLower(currency.CurrencyA.Symbol))
	if err != nil {
		return nil, err
	}
	//log.Println(respmap)
	dep := new(Depth)
	asksmap := respmap["asks"].([]interface{})
	bidsmap := respmap["bids"].([]interface{})

	var (
		l     = len(asksmap)
		i int = l - size
		c int = 0
	)

	for ; i < l; i++ {
		ask := asksmap[i]
		var dr DepthRecord
		for i, vv := range ask.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		dep.AskList = append(dep.AskList, dr)
	}

	c = 0
	for _, v := range bidsmap {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		dep.BidList = append(dep.BidList, dr)
		c++
		if c == size {
			break
		}
	}

	return dep, nil
}

func (btcbox *BtcBox) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (btcbox *BtcBox) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (btcbox *BtcBox) GetExchangeName() string {
	return "btcbox.co.jp"
}
