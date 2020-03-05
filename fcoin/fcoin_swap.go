package fcoin

import (
	"fmt"
	//"github.com/google/uuid"
	. "github.com/nntaoli-project/goex"
	"github.com/pkg/errors"
	//"strings
	"net/url"
)

const (

	/*Rest Endpoint*/
	Endpoint              = "https://api.testnet.fmex.com"
	GET_ACCOUNTS          = "/v3/contracts/accounts"
	PLACE_ORDER           = "/v3/contracts/orders"
	CANCEL_ORDER          = "/v3/contracts/orders/%s/cancel"
	GET_POSITION          = "/v3/contracts/positions"
	GET_DEPTH             = "/v2/market/depth/L20/%s"
	GET_TICKER            = "/v2/market/ticker/%s"
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
	uri := fmt.Sprintf(GET_TICKER,fm.adaptContractType(currencyPair))
	respmap, err := HttpGet(fm.httpClient, fm.baseUrl+uri)
	if err != nil {
		return nil, err
	}
	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}
	ticker := respmap["data"].(map[string]interface{})["ticker"].([]interface{})
	return &Ticker{Pair:currencyPair,
		       Last:ticker[0].(float64),
	       	       Buy:ticker[2].(float64),
	       	       Sell:ticker[4].(float64),
	       	       High:ticker[7].(float64),
	       	       Low:ticker[8].(float64),
	       	       Vol:ticker[9].(float64)},nil
}

func (fm *FMexSwap) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	var uri string

	uri = fmt.Sprintf(GET_DEPTH, fm.adaptContractType(currencyPair))
	//fmt.Println("get depth uri:",fm.baseUrl+uri)
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
	r, err := fm.doAuthenticatedRequest("GET", GET_ACCOUNTS, url.Values{})
	if err != nil {
		return nil, err
	}
	fmt.Println("get userinfo:",r)
	return nil,nil
}


func (fm *FMexSwap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {

	params := url.Values{}

	params.Set("symbol", fm.adaptContractType(currencyPair))
	if matchPrice == 1{
		params.Set("type", "market")
	}else{
		params.Set("type", "limit")
	}
	if openType == OPEN_BUY || openType == CLOSE_SELL{
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
 */
func (fm *FMexSwap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not support")
}

/**
 */
func (fm *FMexSwap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	panic("not support")
}

func (fm *FMexSwap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	r, err := fm.doAuthenticatedRequest("GET", GET_POSITION, url.Values{})
	if err != nil {
		return nil, err
	}
	data := r.(map[string]interface{})
	var positions []FuturePosition
	for _,info := range data["results"].([]interface{}) {
		cont := info.(map[string]interface{})
		//fmt.Println("position info:",cont["direction"])
		p := FuturePosition{CreateDate:int64(cont["updated_at"].(float64)),
				    LeverRate:int(cont["leverage"].(float64)),
				    Symbol:currencyPair,
				    ContractId:int64(cont["user_id"].(float64)),
				    ForceLiquPrice:cont["liquidation_price"].(float64)}
		if cont["direction"] == "long"{
			p.BuyAmount     = cont["quantity"].(float64)
			p.BuyPriceAvg   = cont["entry_price"].(float64)
			p.BuyPriceCost  = cont["margin"].(float64)
			p.BuyProfitReal = cont["realized_pnl"].(float64)
		}else{
			p.SellAmount     = cont["quantity"].(float64)
			p.SellPriceAvg   = cont["entry_price"].(float64)
			p.SellPriceCost  = cont["margin"].(float64)
			p.SellProfitReal = cont["realized_pnl"].(float64)
		}

		positions = append(positions,p)

	}

	return positions,nil
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
