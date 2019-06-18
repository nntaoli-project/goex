package okcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deckarep/golang-set"
	. "github.com/nntaoli-project/GoEx"
)

const (
	V3_FUTURE_HOST_URL        = "https://www.okex.com/"
	V3_FUTURE_API_BASE_URL    = "api/futures/v3/"
	V3_SWAP_API_BASE_URL      = "api/swap/v3/"
	V3_FUTRUR_INSTRUMENTS_URL = "instruments"
	V3_FUTURE_TICKER_URI      = "instruments/%s/ticker"
	V3_SWAP_DEPTH_URI         = "instruments/%s/depth" //wtf api
	V3_FUTURE_DEPTH_URI       = "instruments/%s/book"  //wtf api
	V3_FUTURE_ESTIMATED_PRICE = "instruments/%s/estimated_price"
	V3_FUTURE_INDEX_PRICE     = "instruments/%s/index"
	V3_FUTURE_USERINFOS_URI   = "accounts/%s"
	V3_FUTURE_CANCEL_URI      = "cancel_order/%s/%s"
	V3_FUTURE_ORDER_INFO_URI  = "orders/%s/%s"
	V3_FUTURE_ORDERS_INFO_URI = "orders/%s"
	V3_FUTURE_POSITION_URI    = "%s/position"
	V3_FUTURE_ORDER_URI       = "order"
	V3_FUTURE_TRADES_URI      = "instruments/%s/trades"
	V3_FUTURE_FILLS_URI       = "fills"
	V3_EXCHANGE_RATE_URI      = "instruments/%s/rate"
	V3_GET_KLINE_URI          = "instruments/%s/candles"
)

// common utils, maybe should be extracted in future
func timeStringToInt64(t string) (int64, error) {
	timestamp, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return 0, err
	}
	return timestamp.UnixNano() / int64(time.Millisecond), nil
}

func int64ToTime(ti int64) time.Time {
	return time.Unix(0, ti*int64(time.Millisecond)).UTC()
}

func int64ToTimeString(ti int64) string {
	t := int64ToTime(ti)
	return t.Format(time.RFC3339)
}

// contract information
type futureContract struct {
	InstrumentID    string `json:"instrument_id"`
	UnderlyingIndex string `json:"underlying_index"`
	QuoteCurrency   string `json:"quote_currency"`
	Coin            string `json:"coin"`
	TickSize        string `json:"tick_size"`
	ContractVal     string `json:"contract_val"`
	Listing         string `json:"listing"`
	Delivery        string `json:"delivery"`
	SizeIncrement   string `json:"size_increment"`
	TradeIncrement  string `json:"trade_increment"`
	Alias           string `json:"alias"`
}

func (fc futureContract) normalizePrice(price float64) (string, error) {
	tickSize := fc.TickSize
	if len(tickSize) == 0 {
		return "", fmt.Errorf("no tick size info in contract %v", fc)
	}

	precision := 0
	i := strings.Index(tickSize, ".")
	if i > -1 {
		decimal := tickSize[i+1:]
		precision = len(decimal) - len(strings.TrimPrefix(decimal, "0")) + 1
	}
	return fmt.Sprintf("%."+fmt.Sprintf("%df", precision), price), nil
}

func (fc futureContract) normalizePriceString(price string) (string, error) {
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", err
	}
	return fc.normalizePrice(p)
}

func (fc futureContract) getSizeIncrement() string {
	if len(fc.SizeIncrement) > 0 {
		return fc.SizeIncrement
	} else if len(fc.TradeIncrement) > 0 {
		return fc.TradeIncrement
	}
	return ""
}

func (fc futureContract) normalizeAmount(amount float64) (string, error) {
	increment := fc.getSizeIncrement()
	if len(increment) == 0 {
		return "", fmt.Errorf("no trade incrument info in contract %v", fc)
	}

	precision := 0
	i := strings.Index(increment, ".")
	if i > -1 {
		decimal := increment[i+1:]
		precision = len(decimal) - len(strings.TrimPrefix(decimal, "0")) + 1
	}
	return fmt.Sprintf("%."+fmt.Sprintf("%df", precision), amount), nil
}

func (fc futureContract) normalizeAmountString(amount string) (string, error) {
	a, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", err
	}
	return fc.normalizeAmount(a)
}

type futureContracts []futureContract

type futureContractsMapKey struct {
	UnderlyingIndex string
	QuoteCurrency   string
	Alias           string
}

type futureContractsMap map[futureContractsMapKey]*futureContract
type futureContractsIDMap map[string]*futureContract

func newFutureContractsMap(contracts futureContracts) futureContractsMap {
	contractsMap := futureContractsMap{}
	for _, v := range contracts {
		func(v futureContract) {
			key := futureContractsMapKey{
				UnderlyingIndex: v.UnderlyingIndex,
				QuoteCurrency:   v.QuoteCurrency,
				Alias:           v.Alias,
			}
			contractsMap[key] = &v
		}(v)
	}
	return contractsMap
}

func newFutureContractsIDMap(contracts futureContracts) futureContractsIDMap {
	contractsIDMap := futureContractsIDMap{}
	for _, v := range contracts {
		func(v futureContract) {
			contractsIDMap[v.InstrumentID] = &v
		}(v)
	}
	return contractsIDMap
}

// NOTE:
// contracts 默认五分钟更新一次。
// 由于V3没有自动合约日期的映射，到了周五交割的时候，最好还是手动平仓，关闭策略，交割完后重启。
type OKExV3 struct {
	apiKey,
	apiSecretKey,
	passphrase,
	endpoint string
	dataParser     *OKExV3DataParser
	client         *http.Client
	contractsMap   futureContractsMap
	contractsIDMap map[string]*futureContract
	contractsRW    *sync.RWMutex
}

func NewOKExV3(client *http.Client, api_key, secret_key, passphrase, endpoint string) *OKExV3 {
	okv3 := new(OKExV3)
	okv3.apiKey = api_key
	okv3.apiSecretKey = secret_key
	okv3.client = client
	okv3.passphrase = passphrase
	okv3.contractsRW = &sync.RWMutex{}
	okv3.dataParser = NewOKExV3DataParser(okv3)
	okv3.endpoint = endpoint
	contracts, err := okv3.getAllContracts()
	if err != nil {
		panic(err)
	}
	okv3.setContracts(contracts)
	return okv3
}

func (okv3 *OKExV3) GetUrlRoot() (url string) {
	if okv3.endpoint != "" {
		url = okv3.endpoint
	} else {
		url = V3_FUTURE_HOST_URL
	}
	if url[len(url)-1] != '/' {
		url = url + "/"
	}
	return
}

func (okv3 *OKExV3) setContracts(contracts futureContracts) {
	contractsMap := newFutureContractsMap(contracts)
	contractsIDMap := newFutureContractsIDMap(contracts)
	okv3.contractsRW.Lock()
	defer okv3.contractsRW.Unlock()
	okv3.contractsMap = contractsMap
	okv3.contractsIDMap = contractsIDMap
}

func (okv3 *OKExV3) getAllContracts() (futureContracts, error) {
	var err error
	var futureContracts, swapContracts futureContracts
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		var err2 error
		futureContracts, err2 = okv3.getFutureContracts()
		if err2 != nil {
			err = err2
		}
	}()
	go func() {
		defer wg.Done()
		var err2 error
		swapContracts, err2 = okv3.getSwapContracts()
		if err2 != nil {
			err = err2
		}
	}()
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return append(futureContracts, swapContracts...), nil
}

func (okv3 *OKExV3) getContractByKey(key futureContractsMapKey) (*futureContract, error) {
	okv3.contractsRW.RLock()
	defer okv3.contractsRW.RUnlock()
	data, ok := okv3.contractsMap[key]
	if !ok {
		msg := fmt.Sprintf("no contract in okex contracts map for %v", key)
		return nil, errors.New(msg)
	}
	return data, nil
}

func (okv3 *OKExV3) GetContract(currencyPair CurrencyPair, contractType string) (*futureContract, error) {
	key := futureContractsMapKey{
		UnderlyingIndex: currencyPair.CurrencyA.Symbol,
		QuoteCurrency:   currencyPair.CurrencyB.Symbol,
		Alias:           contractType,
	}
	return okv3.getContractByKey(key)
}

func (okv3 *OKExV3) GetContractID(currencyPair CurrencyPair, contractType string) (string, error) {
	fallback := ""
	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return fallback, err
	}
	return contract.InstrumentID, nil
}

func (okv3 *OKExV3) ParseContractID(contractID string) (CurrencyPair, string, error) {
	contract, err := okv3.getContractByID(contractID)
	if err != nil {
		return UNKNOWN_PAIR, "", err
	}
	currencyA := NewCurrency(contract.UnderlyingIndex, "")
	currencyB := NewCurrency(contract.QuoteCurrency, "")
	return NewCurrencyPair(currencyA, currencyB), contract.Alias, nil
}

func (okv3 *OKExV3) getContractByID(instrumentID string) (*futureContract, error) {
	okv3.contractsRW.RLock()
	defer okv3.contractsRW.RUnlock()
	data, ok := okv3.contractsIDMap[instrumentID]
	if !ok {
		msg := fmt.Sprintf("no contract in okex contracts with id %s", instrumentID)
		return nil, errors.New(msg)
	}
	return data, nil
}

func (okv3 *OKExV3) startUpdateContractsLoop() {
	interval := 5 * time.Minute
	go func() {
		for {
			time.Sleep(interval)
			contracts, err := okv3.getAllContracts()
			if err == nil {
				okv3.setContracts(contracts)
			}
		}
	}()
}

func (okv3 *OKExV3) GetExchangeName() string {
	return OKEX_FUTURE
}

func (okv3 *OKExV3) getTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func (okv3 *OKExV3) getSign(timestamp, method, url, body string) (string, error) {
	relURL := "/" + strings.TrimPrefix(url, okv3.GetUrlRoot())
	data := timestamp + method + relURL + body
	return GetParamHmacSHA256Base64Sign(okv3.apiSecretKey, data)
}

func (okv3 *OKExV3) getSignedHTTPHeader(method, url string) (map[string]string, error) {
	timestamp := okv3.getTimestamp()
	sign, err := okv3.getSign(timestamp, method, url, "")
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Content-Type":         "application/json",
		"OK-ACCESS-KEY":        okv3.apiKey,
		"OK-ACCESS-SIGN":       sign,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": okv3.passphrase,
	}, nil
}

func (okv3 *OKExV3) getSignedHTTPHeader2(method, url, body string) (map[string]string, error) {
	timestamp := okv3.getTimestamp()
	sign, err := okv3.getSign(timestamp, method, url, body)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Content-Type":         "application/json",
		"OK-ACCESS-KEY":        okv3.apiKey,
		"OK-ACCESS-SIGN":       sign,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": okv3.passphrase,
	}, nil
}

func (okv3 *OKExV3) getSignedHTTPHeader3(method, url string, postData map[string]string) (map[string]string, error) {
	body, _ := json.Marshal(postData)
	return okv3.getSignedHTTPHeader2(method, url, string(body))
}

func (okv3 *OKExV3) getSwapContracts() (futureContracts, error) {
	url := okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTRUR_INSTRUMENTS_URL
	headers, err := okv3.getSignedHTTPHeader("GET", url)
	if err != nil {
		return nil, err
	}

	body, err := HttpGet5(okv3.client, url, headers)
	if err != nil {
		return nil, err
	}

	contracts := futureContracts{}
	err = json.Unmarshal(body, &contracts)
	if err != nil {
		return nil, err
	}
	for i := range contracts {
		contracts[i].Alias = SWAP_CONTRACT
	}
	return contracts, nil
}

func (okv3 *OKExV3) getFutureContracts() (futureContracts, error) {
	url := okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTRUR_INSTRUMENTS_URL
	headers, err := okv3.getSignedHTTPHeader("GET", url)
	if err != nil {
		return nil, err
	}

	body, err := HttpGet5(okv3.client, url, headers)
	if err != nil {
		return nil, err
	}

	contracts := futureContracts{}
	err = json.Unmarshal(body, &contracts)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (okv3 *OKExV3) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	var url string
	if contractType == SWAP_CONTRACT {
		url = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_SWAP_DEPTH_URI
	} else {
		url = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_DEPTH_URI
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	symbol := contract.InstrumentID

	url = fmt.Sprintf(url, symbol)
	if size > 0 {
		url = fmt.Sprintf(url+"?size=%d", size)
	}

	headers, err := okv3.getSignedHTTPHeader("GET", url)
	if err != nil {
		return nil, err
	}

	body, err := HttpGet5(okv3.client, url, headers)
	if err != nil {
		return nil, err
	}

	// println(string(body))

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	if bodyMap["code"] != nil {
		log.Println(bodyMap)
		return nil, errors.New(string(body))
	} else if bodyMap["error_code"] != nil {
		log.Println(bodyMap)
		return nil, errors.New(string(body))
	}

	depth := new(Depth)
	depth.Pair = currencyPair
	depth.ContractType = contractType
	return okv3.dataParser.ParseDepth(depth, bodyMap, size)
}

func (okv3 *OKExV3) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	var url string
	if contractType == SWAP_CONTRACT {
		url = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_TICKER_URI
	} else {
		url = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_TICKER_URI
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	symbol := contract.InstrumentID

	url = fmt.Sprintf(url, symbol)

	headers, err := okv3.getSignedHTTPHeader("GET", url)
	if err != nil {
		return nil, err
	}

	body, err := HttpGet5(okv3.client, url, headers)
	if err != nil {
		return nil, err
	}

	//println(string(body))
	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	if bodyMap["code"] != nil {
		log.Println(bodyMap)
		return nil, errors.New(string(body))
	} else if bodyMap["error_code"] != nil {
		log.Println(bodyMap)
		return nil, errors.New(string(body))
	}

	ticker := new(Ticker)
	ticker.Pair = currencyPair
	timestamp, _ := timeStringToInt64(bodyMap["timestamp"].(string))
	ticker.Date = uint64(timestamp)
	ticker.Buy, _ = strconv.ParseFloat(bodyMap["best_ask"].(string), 64)
	ticker.Sell, _ = strconv.ParseFloat(bodyMap["best_bid"].(string), 64)
	ticker.Last, _ = strconv.ParseFloat(bodyMap["last"].(string), 64)
	ticker.High, _ = strconv.ParseFloat(bodyMap["high_24h"].(string), 64)
	ticker.Low, _ = strconv.ParseFloat(bodyMap["low_24h"].(string), 64)
	ticker.Vol, _ = strconv.ParseFloat(bodyMap["volume_24h"].(string), 64)

	//fmt.Println(bodyMap)
	return ticker, nil
}

func (okv3 *OKExV3) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_ORDER_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_ORDER_URI
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return "", err
	}
	symbol := contract.InstrumentID

	price, err = contract.normalizePriceString(price)
	if err != nil {
		return "", err
	}

	amount, err = contract.normalizeAmountString(amount)
	if err != nil {
		return "", err
	}

	postData := make(map[string]string)
	postData["instrument_id"] = symbol
	postData["price"] = price
	postData["size"] = amount
	postData["type"] = strconv.Itoa(openType)
	// postData["order_type"] = strconv.Itoa(2)
	if contractType != SWAP_CONTRACT {
		postData["leverage"] = strconv.Itoa(leverRate)
	}
	postData["match_price"] = strconv.Itoa(matchPrice)

	headers, err := okv3.getSignedHTTPHeader3("POST", requestURL, postData)
	if err != nil {
		return "", err
	}

	body, err := HttpPostForm4(okv3.client, requestURL, postData, headers)

	if err != nil {
		return "", err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return "", err
	}

	//println(string(body));
	if respMap["error_code"] != nil {
		errorCode, err := strconv.Atoi(respMap["error_code"].(string))
		if err != nil || errorCode != 0 {
			return "", errors.New(string(body))
		}
	}

	return respMap["order_id"].(string), nil
}

func (okv3 *OKExV3) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderID string) (bool, error) {
	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_CANCEL_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_CANCEL_URI
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return false, err
	}
	symbol := contract.InstrumentID

	requestURL = fmt.Sprintf(requestURL, symbol, orderID)

	headers, err := okv3.getSignedHTTPHeader("POST", requestURL)
	if err != nil {
		return false, err
	}

	body, err := HttpPostForm3(okv3.client, requestURL, "", headers)

	if err != nil {
		return false, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return false, err
	}

	//println(string(body));
	if respMap["error_code"] != nil {
		errorCode, err := strconv.Atoi(respMap["error_code"].(string))
		if err != nil || errorCode != 0 {
			return false, errors.New(string(body))
		}
	}

	// wtf api
	switch v := respMap["result"].(type) {
	case bool:
		return v, nil
	case string:
		b, err := strconv.ParseBool(v)
		return b, err
	default:
		return false, errors.New(string(body))
	}
}

func (okv3 *OKExV3) parseOrder(respMap map[string]interface{}, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	order, _, err := okv3.dataParser.ParseFutureOrder(respMap)
	return order, err
}

func (okv3 *OKExV3) parseOrders(body []byte, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	respMap := make(map[string]interface{})
	err := json.Unmarshal(body, &respMap)
	if err != nil {
		return nil, err
	}

	var result bool
	if respMap["result"] != nil {
		switch v := respMap["result"].(type) {
		case bool:
			result = v
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				b = false
			}
			result = b
		default:
			result = false
		}
	} else {
		result = true // no result field, should look at order_info field
	}

	//{"error_message": "You have not uncompleted order at the moment",
	// "result": false, "error_code": "32004", "order_id": "2924351754742784"}
	if !result && respMap["error_code"] != nil {
		log.Println(respMap["error_code"].(string))
		if respMap["error_code"].(string) == "32004" {
			return []FutureOrder{}, nil
		}
	}

	if result && respMap["order_info"] != nil {
		orderInfos := respMap["order_info"].([]interface{})
		orders := make([]FutureOrder, 0, len(orderInfos))
		for _, o := range orderInfos {
			orderInfo := o.(map[string]interface{})
			order, err := okv3.parseOrder(orderInfo, currencyPair, contractType)
			if err != nil {
				return nil, err
			}
			orders = append(orders, *order)
		}
		return orders, nil
	}

	return nil, errors.New(string(body))
}

func (okv3 *OKExV3) GetFutureOrder(orderID string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_ORDER_INFO_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_ORDER_INFO_URI
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	symbol := contract.InstrumentID

	requestURL = fmt.Sprintf(requestURL, symbol, orderID)

	headers, err := okv3.getSignedHTTPHeader("GET", requestURL)
	if err != nil {
		return nil, err
	}

	body, err := HttpGet5(okv3.client, requestURL, headers)

	if err != nil {
		return nil, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return nil, err
	}

	//println(string(body));
	if respMap["error_code"] != nil {
		errorCode, err := strconv.Atoi(respMap["error_code"].(string))
		if err != nil || errorCode != 0 {
			return nil, errors.New(string(body))
		}
	}

	return okv3.parseOrder(respMap, currencyPair, contractType)
}

var orderStateOrder = map[TradeStatus]int{
	ORDER_UNFINISH:    0,
	ORDER_PART_FINISH: 1,
	ORDER_REJECT:      2,
	ORDER_CANCEL_ING:  3,
	ORDER_CANCEL:      4,
	ORDER_FINISH:      5,
}

func (okv3 *OKExV3) mergeOrders(ordersList ...[]FutureOrder) []FutureOrder {
	orderMap := make(map[string]FutureOrder)
	for _, orders := range ordersList {
		if orders != nil {
			for _, order := range orders {
				if o, ok := orderMap[order.OrderID2]; ok {
					if orderStateOrder[o.Status] < orderStateOrder[order.Status] {
						orderMap[order.OrderID2] = order
					}
				} else {
					orderMap[order.OrderID2] = order
				}
			}
		}
	}
	orders := make([]FutureOrder, 0, len(orderMap))
	for _, value := range orderMap {
		orders = append(orders, value)
	}
	return orders
}

func (okv3 *OKExV3) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	var err error
	var orders1, orders2, orders3 []FutureOrder

	// TODO: when cancelling orders via api which query orders, order maybe not valiable in all the state,
	// api which query single order is more stable, so when orderIDs's length is less than 3~5,
	// just using query sigle order api.
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		var _err error
		orders2, _err = okv3.GetUnfinishFutureOrdersByIDs(orderIds, currencyPair, contractType)
		if _err != nil {
			err = _err
		}
	}()

	// cancelling orders
	go func() {
		defer wg.Done()
		var _err error
		orders3, _err = okv3.GetFutureOrdersByIDsAndState(orderIds, "4", currencyPair, contractType)
		if _err != nil {
			err = _err
		}
	}()

	go func() {
		defer wg.Done()
		var _err error
		orders1, _err = okv3.GetFinishedFutureOrdersByIDs(orderIds, currencyPair, contractType)
		if _err != nil {
			err = _err
		}
	}()

	wg.Wait()
	if err != nil {
		return nil, err
	}

	return okv3.mergeOrders(orders1, orders2, orders3), nil
}

func (okv3 *OKExV3) filterOrdersByIDs(orderIDs []string, orders []FutureOrder) []FutureOrder {
	//if no orderIDs specific, return all the orders
	if len(orderIDs) == 0 {
		return orders
	}

	orderSet := mapset.NewSet()
	for _, o := range orderIDs {
		orderSet.Add(o)
	}
	newOrders := make([]FutureOrder, 0, len(orderIDs))
	for _, order := range orders {
		if orderSet.Contains(order.OrderID2) {
			newOrders = append(newOrders, order)
		}
	}
	return newOrders
}

func (okv3 *OKExV3) GetFutureOrdersByIDsAndState(orderIDs []string, state string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_ORDERS_INFO_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_ORDERS_INFO_URI
	}

	postData := url.Values{}
	if len(orderIDs) > 0 {
		sort.Strings(orderIDs)
		// 设计api像cxk, from 参数对应的订单竟然不包含在查询结果中
		// postData.Set("from", orderIDs[len(orderIDs) - 1])
		postData.Set("to", orderIDs[0])
	}
	postData.Set("state", state)

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	symbol := contract.InstrumentID

	requestURL = fmt.Sprintf(requestURL, symbol) + "?" + postData.Encode()

	headers, err := okv3.getSignedHTTPHeader("GET", requestURL)
	if err != nil {
		return nil, err
	}

	body, err := HttpGet5(okv3.client, requestURL, headers)
	if err != nil {
		return nil, err
	}

	orders, err := okv3.parseOrders(body, currencyPair, contractType)
	if err != nil {
		return nil, err
	}

	return okv3.filterOrdersByIDs(orderIDs, orders), nil
}

func (okv3 *OKExV3) GetUnfinishFutureOrdersByIDs(orderIDs []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return okv3.GetFutureOrdersByIDsAndState(orderIDs, "6", currencyPair, contractType)
}

func (okv3 *OKExV3) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return okv3.GetUnfinishFutureOrdersByIDs([]string{}, currencyPair, contractType)
}

func (okv3 *OKExV3) GetFutureOrdersByState(state string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return okv3.GetFutureOrdersByIDsAndState([]string{}, state, currencyPair, contractType)
}

func (okv3 *OKExV3) GetFinishedFutureOrdersByIDs(orderIDs []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return okv3.GetFutureOrdersByIDsAndState(orderIDs, "7", currencyPair, contractType)
}

func (okv3 *OKExV3) GetFinishedFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return okv3.GetFinishedFutureOrdersByIDs([]string{}, currencyPair, contractType)
}

var OKexOrderTypeMap = map[int]int{
	ORDER_TYPE_LIMIT:     0,
	ORDER_TYPE_POST_ONLY: 1,
	ORDER_TYPE_FAK:       2,
	ORDER_TYPE_IOC:       3,
}

func (okv3 *OKExV3) PlaceFutureOrder2(currencyPair CurrencyPair, contractType, price, amount string, orderType, openType, matchPrice, leverRate int) (string, error) {
	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_ORDER_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_ORDER_URI
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return "", err
	}
	symbol := contract.InstrumentID

	price, err = contract.normalizePriceString(price)
	if err != nil {
		return "", err
	}

	amount, err = contract.normalizeAmountString(amount)
	if err != nil {
		return "", err
	}

	postData := make(map[string]string)
	postData["instrument_id"] = symbol
	postData["price"] = price
	postData["size"] = amount
	postData["type"] = strconv.Itoa(openType)
	if ot, ok := OKexOrderTypeMap[orderType]; ok {
		postData["order_type"] = strconv.Itoa(ot)
	} else {
		return "", fmt.Errorf("unsupport order type %s in %s", OrderType(ot).String(), okv3.GetExchangeName())
	}
	if contractType != SWAP_CONTRACT {
		postData["leverage"] = strconv.Itoa(leverRate)
	}
	postData["match_price"] = strconv.Itoa(matchPrice)

	headers, err := okv3.getSignedHTTPHeader3("POST", requestURL, postData)
	if err != nil {
		return "", err
	}

	body, err := HttpPostForm4(okv3.client, requestURL, postData, headers)

	if err != nil {
		return "", err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return "", err
	}

	//println(string(body));
	if respMap["error_code"] != nil {
		errorCode, err := strconv.Atoi(respMap["error_code"].(string))
		if err != nil || errorCode != 0 {
			return "", errors.New(string(body))
		}
	}

	return respMap["order_id"].(string), nil
}

func (okv3 *OKExV3) GetDeliveryTime() (int, int, int, int) {
	return 4, 16, 0, 0
}

func (okv3 *OKExV3) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	// seems contractType should be one of the function parameters.
	contractType := QUARTER_CONTRACT
	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return -1, err
	}
	f, err := strconv.ParseFloat(contract.ContractVal, 64)
	if err != nil {
		return -1, err
	}
	return f, nil
}

func (okv3 *OKExV3) GetFee() (float64, error) {
	return 0.03, nil //期货固定0.03%手续费
}

func (okv3 *OKExV3) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	fallback := -1.0

	contractType := THIS_WEEK_CONTRACT
	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return fallback, err
	}
	requestURL := okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_ESTIMATED_PRICE
	requestURL = fmt.Sprintf(requestURL, contract.InstrumentID)
	headers, err := okv3.getSignedHTTPHeader("GET", requestURL)
	if err != nil {
		return fallback, err
	}

	body, err := HttpGet5(okv3.client, requestURL, headers)

	if err != nil {
		return fallback, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return fallback, err
	}

	//println(string(body));
	if respMap["error_code"] != nil {
		errorCode, err := strconv.Atoi(respMap["error_code"].(string))
		if err != nil || errorCode != 0 {
			return fallback, errors.New(string(body))
		}
	}

	f, err := strconv.ParseFloat(respMap["settlement_price"].(string), 64)
	if err != nil {
		return fallback, err
	}
	return f, nil
}

// long and short position should be separated,
// current FuturePosition struct should be redesigned.
var (
	OKEX_MARGIN_MODE_CROSSED = "crossed"
	OKEX_MARGIN_MODE_FIXED   = "fixed"
)

func (okv3 *OKExV3) parsePositionsFuture(respMap map[string]interface{}, currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	// wtf okex
	var posAtr []FuturePosition
	var pos FuturePosition

	holdings := respMap["holding"].([]interface{})

	for _, v := range holdings {
		holdingMap := v.(map[string]interface{})
		marginMode := holdingMap["margin_mode"].(string)

		pos = FuturePosition{}
		pos.ContractType = contractType
		pos.Symbol = currencyPair
		if marginMode == OKEX_MARGIN_MODE_FIXED {
			pos.ForceLiquPrice, _ = strconv.ParseFloat(holdingMap["long_liqui_price"].(string), 64)
			i, _ := strconv.ParseInt(holdingMap["long_leverage"].(string), 10, 64)
			pos.LeverRate = int(i)
		} else {
			pos.ForceLiquPrice, _ = strconv.ParseFloat(holdingMap["liquidation_price"].(string), 64)
			i, _ := strconv.ParseInt(holdingMap["leverage"].(string), 10, 64)
			pos.LeverRate = int(i)
		}
		pos.BuyAmount, _ = strconv.ParseFloat(holdingMap["long_qty"].(string), 64)
		pos.BuyAvailable, _ = strconv.ParseFloat(holdingMap["long_avail_qty"].(string), 64)
		pos.BuyPriceAvg, _ = strconv.ParseFloat(holdingMap["long_avg_cost"].(string), 64)
		// TODO: what's mean of BuyPriceCost
		pos.BuyPriceCost, _ = strconv.ParseFloat(holdingMap["long_avg_cost"].(string), 64)
		pos.BuyProfitReal, _ = strconv.ParseFloat(holdingMap["long_pnl"].(string), 64)
		pos.CreateDate, _ = timeStringToInt64(holdingMap["created_at"].(string))
		if pos.BuyAmount > 0 {
			posAtr = append(posAtr, pos)
		}

		pos = FuturePosition{}
		pos.ContractType = contractType
		pos.Symbol = currencyPair
		if marginMode == OKEX_MARGIN_MODE_FIXED {
			pos.ForceLiquPrice, _ = strconv.ParseFloat(holdingMap["short_liqui_price"].(string), 64)
			i, _ := strconv.ParseInt(holdingMap["short_leverage"].(string), 10, 64)
			pos.LeverRate = int(i)
		} else {
			pos.ForceLiquPrice, _ = strconv.ParseFloat(holdingMap["liquidation_price"].(string), 64)
			i, _ := strconv.ParseInt(holdingMap["leverage"].(string), 10, 64)
			pos.LeverRate = int(i)
		}
		pos.SellAmount, _ = strconv.ParseFloat(holdingMap["short_qty"].(string), 64)
		pos.SellAvailable, _ = strconv.ParseFloat(holdingMap["short_avail_qty"].(string), 64)
		pos.SellPriceAvg, _ = strconv.ParseFloat(holdingMap["short_avg_cost"].(string), 64)
		pos.SellPriceCost, _ = strconv.ParseFloat(holdingMap["short_avg_cost"].(string), 64)
		pos.SellProfitReal, _ = strconv.ParseFloat(holdingMap["short_pnl"].(string), 64)
		pos.CreateDate, _ = timeStringToInt64(holdingMap["created_at"].(string))

		if pos.SellAmount > 0 {
			posAtr = append(posAtr, pos)
		}
	}
	return posAtr, nil
}

func (okv3 *OKExV3) parsePositionsSwap(respMap map[string]interface{}, currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	// wtf okex
	var pos FuturePosition
	var posAr []FuturePosition
	holdings := respMap["holding"].([]interface{})

	for _, v := range holdings {
		holdingMap := v.(map[string]interface{})

		pos = FuturePosition{}
		pos.ContractType = contractType
		pos.Symbol = currencyPair
		pos.ForceLiquPrice, _ = strconv.ParseFloat(holdingMap["liquidation_price"].(string), 64)
		i, _ := strconv.ParseInt(holdingMap["leverage"].(string), 10, 64)
		pos.LeverRate = int(i)
		pos.CreateDate, _ = timeStringToInt64(holdingMap["timestamp"].(string))

		side := holdingMap["side"].(string)
		if side == "short" {
			pos.SellAmount, _ = strconv.ParseFloat(holdingMap["position"].(string), 64)
			pos.SellAvailable, _ = strconv.ParseFloat(holdingMap["avail_position"].(string), 64)
			pos.SellPriceAvg, _ = strconv.ParseFloat(holdingMap["avg_cost"].(string), 64)
			// TODO: what's mean of BuyPriceCost
			pos.SellPriceCost, _ = strconv.ParseFloat(holdingMap["avg_cost"].(string), 64)
			pos.SellProfitReal, _ = strconv.ParseFloat(holdingMap["realized_pnl"].(string), 64)
		} else if side == "long" {
			pos.BuyAmount, _ = strconv.ParseFloat(holdingMap["position"].(string), 64)
			pos.BuyAvailable, _ = strconv.ParseFloat(holdingMap["avail_position"].(string), 64)
			pos.BuyPriceAvg, _ = strconv.ParseFloat(holdingMap["avg_cost"].(string), 64)
			pos.BuyPriceCost, _ = strconv.ParseFloat(holdingMap["avg_cost"].(string), 64)
			pos.BuyProfitReal, _ = strconv.ParseFloat(holdingMap["realized_pnl"].(string), 64)
		} else {
			return nil, fmt.Errorf("unknown swap position side, respMap is: %v", respMap)
		}
		posAr = append(posAr, pos)
	}
	return posAr, nil
}

func (okv3 *OKExV3) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	var fallback []FuturePosition

	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_POSITION_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_POSITION_URI
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return fallback, err
	}
	requestURL = fmt.Sprintf(requestURL, contract.InstrumentID)

	headers, err := okv3.getSignedHTTPHeader("GET", requestURL)
	if err != nil {
		return fallback, err
	}

	body, err := HttpGet5(okv3.client, requestURL, headers)
	if err != nil {
		return fallback, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return fallback, err
	}

	//println(string(body));
	if respMap["error_code"] != nil {
		errorCode, err := strconv.Atoi(respMap["error_code"].(string))
		if err != nil || errorCode != 0 {
			return fallback, errors.New(string(body))
		}
	}

	if contractType == SWAP_CONTRACT {
		return okv3.parsePositionsSwap(respMap, currencyPair, contractType)
	}
	return okv3.parsePositionsFuture(respMap, currencyPair, contractType)
}

var KlineTypeSecondsMap = map[int]int{
	KLINE_PERIOD_1MIN:  60,
	KLINE_PERIOD_3MIN:  180,
	KLINE_PERIOD_5MIN:  300,
	KLINE_PERIOD_15MIN: 900,
	KLINE_PERIOD_30MIN: 1800,
	KLINE_PERIOD_60MIN: 3600,
	KLINE_PERIOD_1H:    3600,
	KLINE_PERIOD_2H:    7200,
	KLINE_PERIOD_4H:    14400,
	KLINE_PERIOD_6H:    21600,
	KLINE_PERIOD_12H:   43200,
	KLINE_PERIOD_1DAY:  86400,
	KLINE_PERIOD_1WEEK: 604800,
}

func (okv3 *OKExV3) mergeKlineRecords(klines [][]*FutureKline) []FutureKline {
	ret := make([]FutureKline, 0)
	for _, kline := range klines {
		for _, k := range kline {
			ret = append(ret, *k)
		}
	}
	return ret
}

func (okv3 *OKExV3) getKlineRecords(contractType string, currencyPair CurrencyPair, seconds int, start, end *time.Time) ([]*FutureKline, error) {
	var fallback []*FutureKline
	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_GET_KLINE_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_GET_KLINE_URI
	}
	params := url.Values{}
	params.Set("granularity", strconv.Itoa(seconds))

	if start != nil {
		params.Set("start", start.Format(time.RFC3339))
	}
	if end != nil {
		params.Set("end", end.Format(time.RFC3339))
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return fallback, err
	}
	requestURL = fmt.Sprintf(requestURL, contract.InstrumentID)
	requestURL = requestURL + "?" + params.Encode()

	headers, err := okv3.getSignedHTTPHeader("GET", requestURL)
	if err != nil {
		return fallback, err
	}

	body, err := HttpGet5(okv3.client, requestURL, headers)

	if err != nil {
		return fallback, err
	}

	var klines [][]interface{}
	err = json.Unmarshal(body, &klines)
	if err != nil {
		return fallback, err
	}

	var klineRecords []*FutureKline
	for _, record := range klines {
		r := FutureKline{}
		r.Kline = new(Kline)
		for i, e := range record {
			switch i {
			case 0:
				r.Timestamp, _ = timeStringToInt64(e.(string)) //to unix timestramp
			case 1:
				r.Open, _ = strconv.ParseFloat(e.(string), 64)
			case 2:
				r.High, _ = strconv.ParseFloat(e.(string), 64)
			case 3:
				r.Low, _ = strconv.ParseFloat(e.(string), 64)
			case 4:
				r.Close, _ = strconv.ParseFloat(e.(string), 64)
			case 5:
				r.Vol, _ = strconv.ParseFloat(e.(string), 64)
			case 6:
				r.Vol2, _ = strconv.ParseFloat(e.(string), 64)
			}
		}
		klineRecords = append(klineRecords, &r)
	}
	reverse(klineRecords)
	return klineRecords, nil
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func (okv3 *OKExV3) GetKlineRecords(contractType string, currencyPair CurrencyPair, period, size, since int) ([]FutureKline, error) {
	var fallback []FutureKline

	var maxSize int
	if contractType == SWAP_CONTRACT {
		maxSize = 200
	} else {
		maxSize = 300
	}

	seconds, ok := KlineTypeSecondsMap[period]
	if !ok {
		return nil, fmt.Errorf("invalid kline period for okex %d", period)
	}

	starts := make([]*time.Time, 0)
	ends := make([]*time.Time, 0)

	if since > 0 {
		startTime := int64ToTime(int64(since))
		for start, left := startTime, size; true; {
			if left > maxSize {
				s := start
				e := s.Add(time.Duration(seconds*maxSize) * time.Second)
				starts = append(starts, &s)
				ends = append(ends, &e)
				start = e.Add(1) // a little trick
				left -= maxSize
			} else {
				s := start
				starts = append(starts, &s)
				ends = append(ends, nil)
				break
			}
		}
		lock := &sync.Mutex{}
		klinesSlice := make([][]*FutureKline, len(starts))
		var err error
		wg := &sync.WaitGroup{}
		wg.Add(len(starts))
		for i := 0; i < len(starts); i++ {
			go func(i int) {
				defer wg.Done()
				klines, err2 := okv3.getKlineRecords(contractType, currencyPair, seconds, starts[i], ends[i])
				lock.Lock()
				defer lock.Unlock()
				if err2 != nil {
					if err == nil {
						err = err2
					}
					log.Println("error when get kline: ", err2)
				}
				klinesSlice[i] = klines
			}(i)
		}
		wg.Wait()
		if err != nil {
			return fallback, err
		}
		klines := okv3.mergeKlineRecords(klinesSlice)
		l := len(klines)
		if l > size {
			klines = klines[0:size]
		}
		return klines, nil
	} else {
		endTime := time.Now().UTC()
		for end, left := endTime, size; true; {
			if left > maxSize {
				e := end
				s := e.Add(-time.Duration(seconds*maxSize) * time.Second)
				starts = append(starts, &s)
				ends = append(ends, &e)
				end = s.Add(-1) // a little trick
				left -= maxSize
			} else {
				e := end
				starts = append(starts, nil)
				ends = append(ends, &e)
				break
			}
		}
		reverse(starts)
		reverse(ends)

		lock := &sync.Mutex{}
		klinesSlice := make([][]*FutureKline, len(starts))
		var err error
		wg := &sync.WaitGroup{}
		wg.Add(len(starts))
		for i := 0; i < len(starts); i++ {
			go func(i int) {
				defer wg.Done()
				klines, err2 := okv3.getKlineRecords(contractType, currencyPair, seconds, starts[i], ends[i])
				lock.Lock()
				defer lock.Unlock()
				if err2 != nil {
					if err == nil {
						err = err2
					}
					log.Println("error when get kline: ", err2)
				}
				klinesSlice[i] = klines
			}(i)
		}
		wg.Wait()
		if err != nil {
			return fallback, err
		}
		klines := okv3.mergeKlineRecords(klinesSlice)
		l := len(klines)
		if l > size {
			klines = klines[l-size : l]
		}
		return klines, nil
	}
}

func (okv3 *OKExV3) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	var fallback []Trade

	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_TRADES_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_TRADES_URI
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return fallback, err
	}
	requestURL = fmt.Sprintf(requestURL, contract.InstrumentID)

	headers, err := okv3.getSignedHTTPHeader("GET", requestURL)
	if err != nil {
		return fallback, err
	}

	body, err := HttpGet5(okv3.client, requestURL, headers)
	if err != nil {
		return fallback, err
	}

	var trades []Trade
	var resp []interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, errors.New(string(body))
	}

	for _, v := range resp {
		item := v.(map[string]interface{})
		trade := new(Trade)
		trade.Pair = currencyPair
		trade, _, err := okv3.dataParser.ParseTrade(trade, contractType, item)
		if err != nil {
			return nil, err
		}
		trades = append(trades, *trade)
	}

	return trades, nil
}

func (okv3 *OKExV3) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	fallback := -1.0

	contractType := SWAP_CONTRACT

	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_INDEX_PRICE
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_INDEX_PRICE
	}

	contract, err := okv3.GetContract(currencyPair, contractType)
	if err != nil {
		return fallback, err
	}
	requestURL = fmt.Sprintf(requestURL, contract.InstrumentID)

	headers, err := okv3.getSignedHTTPHeader("GET", requestURL)
	if err != nil {
		return fallback, err
	}

	body, err := HttpGet5(okv3.client, requestURL, headers)
	if err != nil {
		return fallback, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return fallback, err
	}

	f, err := strconv.ParseFloat(respMap["index"].(string), 64)
	if err != nil {
		return fallback, err
	}
	return f, nil
}

func (okv3 *OKExV3) getAllCurrencies() []string {
	okv3.contractsRW.RLock()
	defer okv3.contractsRW.RUnlock()

	set := make(map[string]string)
	for _, v := range okv3.contractsMap {
		currency := strings.ToUpper(v.UnderlyingIndex)
		set[currency] = currency
	}
	currencies := make([]string, 0, len(set))
	for _, v := range set {
		currencies = append(currencies, v)
	}
	return currencies
}

// GetFutureUserinfo 还只能查交割合约接口
func (okv3 *OKExV3) GetFutureUserinfo() (*FutureAccount, error) {
	var fallback *FutureAccount

	contractType := QUARTER_CONTRACT // 目前只查交割合约

	var requestURL string
	if contractType == SWAP_CONTRACT {
		requestURL = okv3.GetUrlRoot() + V3_SWAP_API_BASE_URL + V3_FUTURE_USERINFOS_URI
	} else {
		requestURL = okv3.GetUrlRoot() + V3_FUTURE_API_BASE_URL + V3_FUTURE_USERINFOS_URI
	}

	currencies := okv3.getAllCurrencies()

	wg := &sync.WaitGroup{}
	wg.Add(len(currencies))

	account := new(FutureAccount)
	account.FutureSubAccounts = make(map[Currency]FutureSubAccount)
	var err error
	lock := &sync.Mutex{}

	for _, v := range currencies {
		go func(currency string) {
			defer wg.Done()
			requestURL := fmt.Sprintf(requestURL, strings.ToLower(currency))

			headers, err2 := okv3.getSignedHTTPHeader("GET", requestURL)

			lock.Lock()
			defer lock.Unlock()

			if err2 != nil {
				err = err2
				return
			}

			body, err2 := HttpGet5(okv3.client, requestURL, headers)
			if err2 != nil {
				err = err2
				return
			}

			respMap := make(map[string]interface{})
			err = json.Unmarshal(body, &respMap)
			if err2 != nil {
				err = err2
				return
			}
			marginMode := respMap["margin_mode"].(string)
			symbol := NewCurrency(currency, "")
			subAccount := FutureSubAccount{}
			subAccount.Currency = symbol
			subAccount.AccountRights, _ = strconv.ParseFloat(respMap["equity"].(string), 64)
			if marginMode == OKEX_MARGIN_MODE_CROSSED {
				subAccount.KeepDeposit, _ = strconv.ParseFloat(respMap["margin"].(string), 64)
				subAccount.ProfitReal, _ = strconv.ParseFloat(respMap["realized_pnl"].(string), 64)
				subAccount.ProfitUnreal, _ = strconv.ParseFloat(respMap["unrealized_pnl"].(string), 64)
				subAccount.RiskRate, _ = strconv.ParseFloat(respMap["margin_ratio"].(string), 64)
			} else {
				totalMargin := 0.0
				totalRealizedPnl := 0.0
				totalUnrealizedPnl := 0.0
				if respMap["contracts"] != nil {
					contracts := respMap["contracts"].([]interface{})
					for _, v := range contracts {
						cMap := v.(map[string]interface{})
						marginFrozen, _ := strconv.ParseFloat(cMap["margin_frozen"].(string), 64)
						marginUnfilled, _ := strconv.ParseFloat(cMap["margin_for_unfilled"].(string), 64)
						realizedPnl, _ := strconv.ParseFloat(cMap["realized_pnl"].(string), 64)
						unrealizedPnl, _ := strconv.ParseFloat(cMap["unrealized_pnl"].(string), 64)
						totalMargin = totalMargin + marginFrozen + marginUnfilled
						totalRealizedPnl = totalRealizedPnl + realizedPnl
						totalUnrealizedPnl = totalUnrealizedPnl + unrealizedPnl
					}
					subAccount.KeepDeposit = totalMargin
					subAccount.ProfitReal = totalRealizedPnl
					subAccount.ProfitUnreal = totalUnrealizedPnl
				}
			}
			account.FutureSubAccounts[symbol] = subAccount
		}(v)
	}

	wg.Wait()

	if err != nil {
		return fallback, err
	}

	return account, nil
}
