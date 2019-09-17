package fcoin

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"github.com/pkg/errors"
	"net/url"
)

const (

	/*Rest Endpoint*/
	Endpoint              = "https://api.testnet.fmex.com"
	// GET_ACCOUNTS          = "/api/swap/v3/accounts"
	PLACE_ORDER           = "/v3/contracts/orders"
	CANCEL_ORDER          = "/v3/contracts/orders/%s/cancel"
	// GET_ORDER             = "/api/swap/v3/orders/%s/%s"
	// GET_POSITION          = "/api/swap/v3/%s/position"
	GET_DEPTH             = "/v2/market/depth/L20/%s"
	// GET_TICKER            = "/api/swap/v3/instruments/%s/ticker"
	GET_UNFINISHED_ORDERS = "/v3/contracts/orders/open"
)



type FMexSwap struct {
	*FCoin
	config *APIConfig
}

func NewFMexSwap(config *APIConfig) *FMexSwap {
	fcoin :=&FCoin{baseUrl: "https://api.testnet.fmex.com", accessKey: config.ApiKey, secretKey: config.ApiSecretKey, httpClient: config.HttpClient}
	return &FMexSwap{FCoin: fcoin, config: config}
}

func (fm *FMexSwap) GetExchangeName() string {
	return "fmex.com"
}

func (fm *FMexSwap) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	panic("not support")
}

func (fm *FMexSwap) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	var uri string

	uri = fmt.Sprintf(GET_DEPTH, fm.adaptContractType(currencyPair))
	fmt.Println("get depth uri:",fm.baseUrl+uri)
	respmap, err := HttpGet(fm.httpClient, fm.baseUrl+uri)
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})

	bids, ok1 := datamap["bids"].([]interface{})
	asks, ok2 := datamap["asks"].([]interface{})

	if !ok1 || !ok2 {
		return nil, errors.New("depth error")
	}

	depth := new(Depth)
	depth.Pair = currencyPair

	n := 0
	for i := 0; i < len(bids); {
		depth.BidList = append(depth.BidList, DepthRecord{ToFloat64(bids[i]), ToFloat64(bids[i+1])})
		i += 2
		n++
		if n == size {
			break
		}
	}

	n = 0
	for i := 0; i < len(asks); {
		depth.AskList = append(depth.AskList, DepthRecord{ToFloat64(asks[i]), ToFloat64(asks[i+1])})
		i += 2
		n++
		if n == size {
			break
		}
	}

	return depth, nil
}

func (fm *FMexSwap) GetFutureUserinfo() (*FutureAccount, error) {
	panic("not support")
}


func (fm *FMexSwap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {

	params := url.Values{}

	params.Set("symbol", fm.adaptContractType(currencyPair))
	params.Set("type", "limit")
	if openType == BUY{
		params.Set("direction", "long")
	}else{
		params.Set("direction", "short")
	}
	params.Set("price", price)
	params.Set("quantity", amount)

	r, err := fm.doAuthenticatedRequest("POST", PLACE_ORDER, params)
	if err != nil {
		return "", err
	}

	data := r.(map[string]interface{})

	return fmt.Sprintf("%d",int64(data["id"].(float64))),nil

}

func (fm *FMexSwap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	uri := fmt.Sprintf(CANCEL_ORDER, orderId)
	_, err := fm.doAuthenticatedRequest("POST", uri, url.Values{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fm *FMexSwap) parseOrder(ord interface{}) FutureOrder {
	order := ord.(map[string]interface{})
	return FutureOrder{
		OrderID2:   fmt.Sprintf("%d",int64(order["id"].(float64))),
		Amount:     order["quantity"].(float64),
		OrderTime:  int64(order["created_at"].(float64))}
}

func (fm *FMexSwap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	r, err := fm.doAuthenticatedRequest("GET", GET_UNFINISHED_ORDERS, url.Values{})
	if err != nil {
		return nil, err
	}
	data := r.(map[string]interface{})
	var orders []FutureOrder
	for _, info := range data["results"].([]interface{}) {
		ord := fm.parseOrder(info)
		ord.Currency = currencyPair
		ord.ContractName = contractType
		orders = append(orders, ord)
	}

	return orders, nil

}

/**
 *获取订单信息
 */
func (fm *FMexSwap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not support")
}

/**
 *获取单个订单信息
 */
func (fm *FMexSwap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	panic("not support")
}

func (fm *FMexSwap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	panic("not support")
}

/**
 */
func (fm *FMexSwap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not support")
}

func (fm *FMexSwap) GetFee() (float64, error) {
	panic("not support")
}

func (fm *FMexSwap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not support")
}

func (fm *FMexSwap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("not support")
}

func (fm *FMexSwap) GetDeliveryTime() (int, int, int, int) {
	panic("not support")
}

func (fm *FMexSwap) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	panic("not support")
}

func (fm *FMexSwap) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not support")
}

func (fm *FMexSwap) GetExchangeRate() (float64, error) {
	panic("not support")
}

func (fm *FMexSwap) GetHistoricalFunding(contractType string, currencyPair CurrencyPair, page int) ([]HistoricalFunding, error) {
	panic("not support")
}

func (fm *FMexSwap) AdaptTradeStatus(status int) TradeStatus {
	panic("not support")
}

func (fm *FMexSwap) adaptContractType(currencyPair CurrencyPair) string {
	return fmt.Sprintf("%s_P", currencyPair.ToSymbol(""))
}
