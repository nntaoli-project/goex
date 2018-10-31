package bigone

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log"
	"time"

	"github.com/nntaoli-project/GoEx"
	"github.com/nubo/jwt"

	"github.com/google/uuid"
)

const (
	API_BASE_URL = "https://big.one/api/v2"
	TICKER_URI   = API_BASE_URL + "/markets/%s/ticker"
	DEPTH_URI    = API_BASE_URL + "/markets/%s/depth"
	ACCOUNT_URI  = API_BASE_URL + "/viewer/accounts"
	ORDERS_URI   = API_BASE_URL + "/viewer/orders"
	//TRADE_URI    = "orders"
)

type Bigone struct {
	accessKey,
	secretKey string
	httpClient *http.Client
	uid        string
}

func New(client *http.Client, api_key, secret_key string) *Bigone {
	return &Bigone{api_key, secret_key, client, uuid.New().String()}
}

func (bo *Bigone) GetExchangeName() string {
	return goex.BIGONE
}

type TickerResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		Ask struct {
			Amount string `json:"amount"`
			Price  string `json:"price"`
		} `json:"ask"`
		Bid struct {
			Amount string `json:"amount"`
			Price  string `json:"price"`
		} `json:"bid"`
		Close           string `json:"close"`
		DailyChange     string `json:"daily_change"`
		DailyChangePerc string `json:"daily_change_perc"`
		High            string `json:"high"`
		Low             string `json:"low"`
		MarketID        string `json:"market_id"`
		MarketUUID      string `json:"market_uuid"`
		Open            string `json:"open"`
		Volume          string `json:"volume"`
	} `json:"data"`
}

func (bo *Bigone) GetTicker(currency goex.CurrencyPair) (*goex.Ticker, error) {
	tickerURI := fmt.Sprintf(TICKER_URI, currency.ToSymbol("-"))

	var resp TickerResp
	log.Printf("GetTicker -> %s", tickerURI)
	err := goex.HttpGet4(bo.httpClient, tickerURI, nil, &resp)

	if err != nil {
		log.Printf("GetTicker - HttpGet4 failed : %v", err)
		return nil, err
	}

	var ticker goex.Ticker
	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = goex.ToFloat64(resp.Data.Close)
	ticker.Buy = goex.ToFloat64(resp.Data.Bid.Price)
	ticker.Sell = goex.ToFloat64(resp.Data.Ask.Price)
	ticker.Low = goex.ToFloat64(resp.Data.Low)
	ticker.High = goex.ToFloat64(resp.Data.High)
	ticker.Vol = goex.ToFloat64(resp.Data.Volume)
	return &ticker, nil
}

type PlaceOrderResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		Amount       string `json:"amount"`
		AvgDealPrice string `json:"avg_deal_price"`
		FilledAmount string `json:"filled_amount"`
		ID           string `json:"id"`
		InsertedAt   string `json:"inserted_at"`
		MarketID     string `json:"market_id"`
		MarketUUID   string `json:"market_uuid"`
		Price        string `json:"price"`
		Side         string `json:"side"`
		State        string `json:"state"`
		UpdatedAt    string `json:"updated_at"`
	} `json:"data"`
}

func (bo *Bigone) placeOrder(amount, price string, pair goex.CurrencyPair, orderType, orderSide string) (*goex.Order, error) {
	path := ORDERS_URI
	params := make(map[string]string)
	params["market_id"] = pair.ToSymbol("-")
	params["side"] = orderSide
	params["amount"] = amount
	params["price"] = price

	var resp PlaceOrderResp
	buf, err := goex.HttpPostForm4(bo.httpClient, path, params, bo.privateHeader())

	if err != nil {
		log.Printf("placeOrder - HttpPostForm4 failed : %v", err)
		return nil, err
	}

	if err = json.Unmarshal(buf, &resp); nil != err {
		log.Printf("buf : %s", string(buf))
		log.Printf("placeOrder - json.Unmarshal failed : %v", err)
		return nil, err
	}

	if len(resp.Errors) > 0 {
		log.Printf("placeOrder - failed : %v", resp.Errors)
		return nil, fmt.Errorf(resp.Errors[0].Message)
	}

	side := goex.BUY
	if orderSide == "ASK" {
		side = goex.SELL
	}

	return &goex.Order{
		Currency:   pair,
		OrderID2:   fmt.Sprintf("%d", resp.Data.ID),
		Price:      goex.ToFloat64(resp.Data.Price),
		Amount:     goex.ToFloat64(resp.Data.Amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       goex.TradeSide(side),
		Status:     goex.ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

func (bo *Bigone) LimitBuy(amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	return bo.placeOrder(amount, price, currency, "LIMIT", "BID")
}

func (bo *Bigone) LimitSell(amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	return bo.placeOrder(amount, price, currency, "LIMIT", "ASK")
}

func (bo *Bigone) MarketBuy(amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	panic("not implements")
}

func (bo *Bigone) MarketSell(amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	panic("not implements")
}

func (bo *Bigone) privateHeader() map[string]string {
	claims := jwt.ClaimSet{
		"type":  "OpenAPI",
		"sub":   bo.accessKey,
		"nonce": time.Now().UnixNano(),
	}
	token, err := claims.Sign(bo.secretKey)
	if nil != err {
		log.Printf("privateHeader - cliam.Sign failed : %v", err)
		return nil
	}
	return map[string]string{"Authorization": "Bearer " + token}
}

type OrderListResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		Edges []struct {
			Cursor string `json:"cursor"`
			Node   struct {
				Amount       string `json:"amount"`
				AvgDealPrice string `json:"avg_deal_price"`
				FilledAmount string `json:"filled_amount"`
				ID           string `json:"id"`
				InsertedAt   string `json:"inserted_at"`
				MarketID     string `json:"market_id"`
				MarketUUID   string `json:"market_uuid"`
				Price        string `json:"price"`
				Side         string `json:"side"`
				State        string `json:"state"`
				UpdatedAt    string `json:"updated_at"`
			} `json:"node"`
		} `json:"edges"`
		PageInfo struct {
			EndCursor       string `json:"end_cursor"`
			HasNextPage     bool   `json:"has_next_page"`
			HasPreviousPage bool   `json:"has_previous_page"`
			StartCursor     string `json:"start_cursor"`
		} `json:"page_info"`
	} `json:"data"`
}

func (bo *Bigone) getOrdersList(currencyPair goex.CurrencyPair, size int, sts goex.TradeStatus) ([]goex.Order, error) {
	apiURL := ""
	apiURL = fmt.Sprintf("%s?market_id=%s",
		ORDERS_URI, currencyPair.ToSymbol("-"))

	if sts == goex.ORDER_FINISH {
		apiURL += "&state=FILLED"
	} else {
		apiURL += "&state=PENDING"
	}
	var resp OrderListResp
	err := goex.HttpGet4(bo.httpClient, apiURL, bo.privateHeader(), &resp)
	if err != nil {
		log.Printf("getOrdersList - HttpGet4 failed : %v", err)
		return nil, err
	}

	orders := make([]goex.Order, 0)
	for _, edge := range resp.Data.Edges {
		order := edge.Node
		status := order.State
		side := order.Side

		ord := goex.Order{}

		switch status {
		case "PENDING":
			ord.Status = goex.ORDER_UNFINISH
		case "FILLED":
			ord.Status = goex.ORDER_FINISH
		case "CANCELED":
			ord.Status = goex.ORDER_CANCEL
		}
		if ord.Status != sts {
			continue // discard
		}

		ord.Currency = currencyPair
		ord.OrderID2 = order.ID

		if side == "ASK" {
			ord.Side = goex.SELL
		} else {
			ord.Side = goex.BUY
		}

		ord.Amount = goex.ToFloat64(order.Amount)
		ord.Price = goex.ToFloat64(order.Price)
		ord.DealAmount = goex.ToFloat64(order.FilledAmount)
		ord.AvgPrice = goex.ToFloat64(order.Price)
		orders = append(orders, ord)
	}

	return orders, nil
}

type CancelOrderResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		ID           string `json:"id"`
		MarketUUID   string `json:"market_uuid"`
		Price        string `json:"price"`
		Amount       string `json:"amount"`
		FilledAmount string `json:"filled_amount"`
		AvgDealPrice string `json:"avg_deal_price"`
		Side         string `json:"side"`
		State        string `json:"state"`
	}
}

func (bo *Bigone) CancelOrder(orderId string, currency goex.CurrencyPair) (bool, error) {
	path := ORDERS_URI + "/" + orderId + "/cancel"
	params := make(map[string]string)
	params["order_id"] = orderId

	buf, err := goex.HttpPostForm4(bo.httpClient, path, params, bo.privateHeader())

	if err != nil {
		log.Printf("CancelOrder - faield : %v", err)
		return false, err
	}
	var resp CancelOrderResp
	if err = json.Unmarshal(buf, &resp); nil != err {
		log.Printf("CancelOrder - json.Unmarshal failed : %v", err)
		return false, err
	}
	if len(resp.Errors) > 0 {
		log.Printf("getOrdersList - response error : %v", resp.Errors)
		return false, fmt.Errorf("%s", resp.Errors[0].Message)
	}
	return true, nil
}

func (bo *Bigone) GetOneOrder(orderId string, currencyPair goex.CurrencyPair) (*goex.Order, error) {
	return nil, fmt.Errorf("GetOneOrder - not support yet")

}
func (bo *Bigone) GetUnfinishOrders(currencyPair goex.CurrencyPair) ([]goex.Order, error) {
	return bo.getOrdersList(currencyPair, -1, goex.ORDER_UNFINISH)
}
func (bo *Bigone) GetOrderHistorys(currencyPair goex.CurrencyPair, currentPage, pageSize int) ([]goex.Order, error) {
	return bo.getOrdersList(currencyPair, -1, goex.ORDER_FINISH)
}

type AccountResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data []struct {
		AssetID       string `json:"asset_id"`
		AssetUUID     string `json:"asset_uuid"`
		Balance       string `json:"balance"`
		LockedBalance string `json:"locked_balance"`
	} `json:"data"`
}

func (bo *Bigone) GetAccount() (*goex.Account, error) {
	var resp AccountResp
	apiUrl := ACCOUNT_URI

	err := goex.HttpGet4(bo.httpClient, apiUrl, bo.privateHeader(), &resp)
	if err != nil {
		log.Println("GetAccount error:", err)
		return nil, err
	}

	acc := goex.Account{}
	acc.Exchange = bo.GetExchangeName()
	acc.SubAccounts = make(map[goex.Currency]goex.SubAccount)

	for _, v := range resp.Data {
		//log.Println(v)
		currency := goex.NewCurrency(v.AssetID, "")

		acc.SubAccounts[currency] = goex.SubAccount{
			Currency:     currency,
			Amount:       goex.ToFloat64(v.Balance),
			ForzenAmount: goex.ToFloat64(v.LockedBalance),
		}
	}

	return &acc, nil
}

type DepthResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		MarketID string `json:"market_id"`
		Bids     []struct {
			Price      string `json:"price"`
			OrderCount int    `json:"order_count"`
			Amount     string `json:"amount"`
		} `json:"bids"`
		Asks []struct {
			Price      string `json:"price"`
			OrderCount int    `json:"order_count"`
			Amount     string `json:"amount"`
		} `json:"asks"`
	}
}

func (bo *Bigone) GetDepth(size int, currencyPair goex.CurrencyPair) (*goex.Depth, error) {
	var resp DepthResp
	apiURL := fmt.Sprintf(DEPTH_URI, currencyPair.ToSymbol("-"))
	err := goex.HttpGet4(bo.httpClient, apiURL, nil, &resp)
	if err != nil {
		log.Println("GetDepth error:", err)
		return nil, err
	}

	depth := new(goex.Depth)

	for _, bid := range resp.Data.Bids {
		amount := goex.ToFloat64(bid.Amount)
		price := goex.ToFloat64(bid.Price)
		dr := goex.DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range resp.Data.Asks {
		amount := goex.ToFloat64(ask.Amount)
		price := goex.ToFloat64(ask.Price)
		dr := goex.DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (bo *Bigone) GetKlineRecords(currency goex.CurrencyPair, period, size, since int) ([]goex.Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (bo *Bigone) GetTrades(currencyPair goex.CurrencyPair, since int64) ([]goex.Trade, error) {
	panic("not implements")
}
