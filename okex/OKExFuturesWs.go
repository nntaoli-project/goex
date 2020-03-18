package okex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nntaoli-project/goex/internal/logger"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/nntaoli-project/goex"
)

type wsResp struct {
	Event     string `json:"event"`
	Channel   string `json:"channel"`
	Table     string `json:"table"`
	Data      json.RawMessage
	Success   bool        `json:"success"`
	ErrorCode interface{} `json:"errorCode"`
}

type OKExV3FutureWs struct {
	*OKEx
	*WsBuilder
	sync.Once
	wsConn        *WsConn
	loginCh       chan wsResp
	logined       bool
	loginLock     *sync.Mutex
	authoriedSubs []map[string]interface{}

	tickerCallback func(*FutureTicker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade, string)
	klineCallback  func(*FutureKline, int)
	orderCallback  func(*FutureOrder, string)
}

func NewOKExV3FuturesWs() *OKExV3FutureWs {
	okV3Ws := &OKExV3FutureWs{}
	okV3Ws.loginCh = make(chan wsResp)
	okV3Ws.logined = false
	okV3Ws.loginLock = &sync.Mutex{}
	okV3Ws.authoriedSubs = make([]map[string]interface{}, 0)
	okV3Ws.WsBuilder = NewWsBuilder().
		WsUrl("wss://real.okex.com:8443/ws/v3").
		ReconnectInterval(2*time.Second).
		AutoReconnect().
		Heartbeat(func() []byte {
			return []byte("ping")
		}, 28*time.Second).
		UnCompressFunc(FlateUnCompress).
		ProtoHandleFunc(okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3FutureWs) getSign(timestamp, method, url, body string) (string, error) {
	data := timestamp + method + url + body
	return GetParamHmacSHA256Base64Sign(okV3Ws.config.ApiSecretKey, data)
}

func (okV3Ws *OKExV3FutureWs) TickerCallback(tickerCallback func(*FutureTicker)) *OKExV3FutureWs {
	okV3Ws.tickerCallback = tickerCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FutureWs) DepthCallback(depthCallback func(*Depth)) *OKExV3FutureWs {
	okV3Ws.depthCallback = depthCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FutureWs) TradeCallback(tradeCallback func(*Trade, string)) *OKExV3FutureWs {
	okV3Ws.tradeCallback = tradeCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FutureWs) OrderCallback(orderCallback func(*FutureOrder, string)) *OKExV3FutureWs {
	okV3Ws.orderCallback = orderCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FutureWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string),
	klineCallback func(*FutureKline, int),
	orderCallback func(*FutureOrder, string)) {
	okV3Ws.tickerCallback = tickerCallback
	okV3Ws.depthCallback = depthCallback
	okV3Ws.tradeCallback = tradeCallback
	okV3Ws.klineCallback = klineCallback
	okV3Ws.orderCallback = orderCallback
}

func (okV3Ws *OKExV3FutureWs) Login() error {
	// already logined
	if okV3Ws.logined {
		return nil
	}
	okV3Ws.connectWs()
	err := okV3Ws.login()
	if err == nil {
		okV3Ws.logined = true
	}
	return err
}

func (okV3Ws *OKExV3FutureWs) getTimestamp() string {
	seconds := float64(time.Now().UTC().UnixNano()) / float64(time.Second)
	return fmt.Sprintf("%.3f", seconds)
}

func (okV3Ws *OKExV3FutureWs) clearChan(c chan wsResp) {
	for {
		if len(c) > 0 {
			<-c
		} else {
			break
		}
	}
}

func (okV3Ws *OKExV3FutureWs) login() error {
	okV3Ws.loginLock.Lock()
	defer okV3Ws.loginLock.Unlock()

	okV3Ws.clearChan(okV3Ws.loginCh)
	timestamp := okV3Ws.getTimestamp()
	method := "GET"
	url := "/users/self/verify"
	sign, _ := okV3Ws.getSign(timestamp, method, url, "")
	op := map[string]interface{}{
		"op":   "login",
		"args": []string{okV3Ws.config.ApiKey, okV3Ws.config.ApiPassphrase, timestamp, sign}}
	err := okV3Ws.wsConn.SendJsonMessage(op)
	if err != nil {
		logger.Error(err)
		return err
	}

	re := <-okV3Ws.loginCh

	if !re.Success {
		return fmt.Errorf("login failed: %v", re)
	}

	logger.Info("ws login success")
	return nil
}

func (okV3Ws *OKExV3FutureWs) subscribe(sub map[string]interface{}) error {
	okV3Ws.connectWs()
	return okV3Ws.wsConn.Subscribe(sub)
}

func (okV3Ws *OKExV3FutureWs) getTablePrefix(currencyPair CurrencyPair, contractType string) string {
	if contractType == SWAP_CONTRACT {
		return "swap"
	}
	return "futures"
}

func (okV3Ws *OKExV3FutureWs) SubscribeDepth(currencyPair CurrencyPair, contractType string, size int) error {
	if (size > 0) && (size != 5) {
		return errors.New("only support depth 5")
	}
	if okV3Ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	contractId := okV3Ws.OKExFuture.GetFutureContractId(currencyPair, contractType)
	chName := fmt.Sprintf("%s/depth5:%s", okV3Ws.getTablePrefix(currencyPair, contractType), contractId)

	return okV3Ws.subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) SubscribeTicker(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}

	contractId := okV3Ws.OKExFuture.GetFutureContractId(currencyPair, contractType)
	chName := fmt.Sprintf("%s/ticker:%s", okV3Ws.getTablePrefix(currencyPair, contractType), contractId)

	return okV3Ws.subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) SubscribeTrade(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}

	contractId := okV3Ws.OKExFuture.GetFutureContractId(currencyPair, contractType)
	chName := fmt.Sprintf("%s/trade:%s", okV3Ws.getTablePrefix(currencyPair, contractType), contractId)

	return okV3Ws.subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) SubscribeKline(currencyPair CurrencyPair, contractType string, period int) error {
	if okV3Ws.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	contractId := okV3Ws.OKExFuture.GetFutureContractId(currencyPair, contractType)
	seconds := adaptKLinePeriod(KlinePeriod(period))
	if seconds == -1 {
		return fmt.Errorf("unsupported kline period %d in okex", period)
	}

	chName := fmt.Sprintf("%s/candle%ds:%s", okV3Ws.getTablePrefix(currencyPair, contractType), seconds, contractId)

	return okV3Ws.subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) SubscribeOrder(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.orderCallback == nil {
		return errors.New("place set order callback func")
	}

	contractId := okV3Ws.OKExFuture.GetFutureContractId(currencyPair, contractType)
	chName := fmt.Sprintf("%s/order:%s", okV3Ws.getTablePrefix(currencyPair, contractType), contractId)

	return okV3Ws.authoriedSubscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) authoriedSubscribe(data map[string]interface{}) error {
	okV3Ws.authoriedSubs = append(okV3Ws.authoriedSubs, data)
	return okV3Ws.subscribe(data)
}

func (okV3Ws *OKExV3FutureWs) reSubscribeAuthoriedChannel() {
	for _, d := range okV3Ws.authoriedSubs {
		okV3Ws.wsConn.SendJsonMessage(d)
	}
}

func (okV3Ws *OKExV3FutureWs) connectWs() {
	okV3Ws.Do(func() {
		okV3Ws.wsConn = okV3Ws.WsBuilder.Build()
	})
}

func (okV3Ws *OKExV3FutureWs) handle(msg []byte) error {
	logger.Debug("[ws] [response] ", string(msg))
	if string(msg) == "pong" {
		return nil
	}

	var wsResp wsResp
	err := json.Unmarshal(msg, &wsResp)
	if err != nil {
		return err
	}

	if wsResp.ErrorCode != nil {
		logger.Error(string(msg))
		return fmt.Errorf("%s", string(msg))
	}

	if wsResp.Event != "" {
		switch wsResp.Event {
		case "subscribe":
			logger.Info("subscribed:", wsResp.Channel)
			return nil
		case "login":
			select {
			case okV3Ws.loginCh <- wsResp:
				return nil
			default:
				return nil
			}
		case "error":
			var errorCode int
			switch v := wsResp.ErrorCode.(type) {
			case int:
				errorCode = v
			case float64:
				errorCode = int(v) // float64 okex牛逼嗷
			case string:
				i, _ := strconv.ParseInt(v, 10, 64)
				errorCode = int(i)
			}

			switch errorCode {
			// event:error message:Already logged in errorCode:30042
			case 30041:
				if okV3Ws.logined { // have logined successfully
					go func() {
						okV3Ws.login()
						okV3Ws.reSubscribeAuthoriedChannel()
					}()
				} // else skip, or better hanle?
				return nil
			case 30042:
				return nil
			}

			// TODO: clearfy which errors should be treated as login result.
			select {
			case okV3Ws.loginCh <- wsResp:
				return nil
			default:
				return fmt.Errorf("error in websocket: %v", wsResp)
			}
		}
		return fmt.Errorf("unknown websocet message: %v", wsResp)
	}

	if wsResp.Table != "" {
		ch, err := okV3Ws.parseChannel(wsResp.Table)
		if err != nil {
			return err
		}

		switch ch {
		case "ticker":
			var tickers []tickerResponse
			err = json.Unmarshal(wsResp.Data, &tickers)
			if err != nil {
				return err
			}

			for _, t := range tickers {
				contractInfo, err := okV3Ws.OKExFuture.GetContractInfo(t.InstrumentId)
				if err != nil {
					logger.Warn(t.InstrumentId, " contract id error ,  ", err)
					continue
				}
				date, _ := time.Parse(time.RFC3339, t.Timestamp)
				okV3Ws.tickerCallback(&FutureTicker{
					Ticker: &Ticker{
						Pair: NewCurrencyPair2(fmt.Sprintf("%s_%s", contractInfo.UnderlyingIndex, contractInfo.QuoteCurrency)),
						Last: t.Last,
						Buy:  t.BestBid,
						Sell: t.BestAsk,
						High: t.High24h,
						Low:  t.Low24h,
						Vol:  t.Volume24h,
						Date: uint64(date.UnixNano() / int64(time.Millisecond)),
					},
					ContractType: contractInfo.Alias,
				})
			}
		case "depth5":
			var (
				depthResp depthResponse
				dep       Depth
			)
			err := json.Unmarshal(wsResp.Data, &depthResp)
			if err != nil {
				return err
			}
			contractInfo, err := okV3Ws.OKExFuture.GetContractInfo(depthResp.InstrumentId)
			if err != nil {
				logger.Warn("")
				return err
			}
			dep.Pair = NewCurrencyPair2(fmt.Sprintf("%s_%s", contractInfo.UnderlyingIndex, contractInfo.QuoteCurrency))
			dep.ContractType = contractInfo.Alias
			dep.UTime, _ = time.Parse(time.RFC3339, depthResp.Timestamp)
			for _, itm := range depthResp.Asks {
				dep.AskList = append(dep.AskList, DepthRecord{
					Price:  ToFloat64(itm[0]),
					Amount: ToFloat64(itm[1])})
			}
			for _, itm := range depthResp.Bids {
				dep.BidList = append(dep.BidList, DepthRecord{
					Price:  ToFloat64(itm[0]),
					Amount: ToFloat64(itm[1])})
			}
			sort.Sort(sort.Reverse(dep.AskList))
			//call back func
			okV3Ws.depthCallback(&dep)
		case "trade":
			var (
				tradeResponse []struct {
					Side         string  `json:"side"`
					TradeId      int64   `json:"trade_id,string"`
					Price        float64 `json:"price,string"`
					Qty          float64 `json:"qty,string"`
					InstrumentId string  `json:"instrument_id"`
					Timestamp    string  `json:"timestamp"`
				}
			)
			err := json.Unmarshal(wsResp.Data, &tradeResponse)
			if err != nil {
				logger.Error("unmarshal error :", err)
				return err
			}

			for _, resp := range tradeResponse {
				contractInfo, err := okV3Ws.OKExFuture.GetContractInfo(resp.InstrumentId)
				if err != nil {
					return err
				}

				tradeSide := SELL
				switch resp.Side {
				case "buy":
					tradeSide = BUY
				}

				t, err := time.Parse(time.RFC3339, resp.Timestamp)
				if err != nil {
					logger.Warn("parse timestamp error:", err)
				}

				okV3Ws.tradeCallback(&Trade{
					Tid:    resp.TradeId,
					Type:   tradeSide,
					Amount: resp.Qty,
					Price:  resp.Price,
					Date:   t.Unix(),
					Pair:   NewCurrencyPair2(fmt.Sprintf("%s_%s", contractInfo.UnderlyingIndex, contractInfo.QuoteCurrency)),
				}, contractInfo.Alias)
			}
		case "order":
			//2020/03/18 18:05:00 OKExFuturesWs.go:257: [D] [ws] [response] {"table":"futures/order","data":[{"leverage":"20","last_fill_time":"2020-03-18T10:05:00.790Z","filled_qty":"4","fee":"-0.00010655","price_avg":"112.62","type":"1","client_oid":"ce1661e5cb614fd690d0463de7a2eeb0","last_fill_qty":"4","instrument_id":"BSV-USD-200327","last_fill_px":"112.62","pnl":"0","size":"4","price":"112.73","last_fill_id":"15229749","error_code":"0","state":"2","contract_val":"10","order_id":"4573750935784449","order_type":"0","timestamp":"2020-03-18T10:05:00.790Z","status":"2"}]}
			var orderResp []futureOrderResponse
			err := json.Unmarshal(wsResp.Data, &orderResp)
			if err != nil {
				return err
			}
			for _, o := range orderResp {
				contractInfo, err := okV3Ws.OKExFuture.GetContractInfo(o.InstrumentId)
				if err != nil {
					logger.Warn("get contract info error , instrument id:", o.InstrumentId)
					continue
				}
				okV3Ws.orderCallback(&FutureOrder{
					ClientOid:    o.ClientOid,
					OrderID2:     o.OrderId,
					Price:        o.Price,
					Amount:       o.Size,
					AvgPrice:     o.PriceAvg,
					DealAmount:   o.FilledQty,
					Status:       okV3Ws.adaptOrderState(o.State),
					Currency:     CurrencyPair{},
					OrderType:    o.OrderType,
					OType:        o.Type,
					LeverRate:    o.Leverage,
					Fee:          o.Fee,
					ContractName: o.InstrumentId,
					OrderTime:    o.Timestamp.UnixNano() / int64(time.Millisecond),
				}, contractInfo.Alias)
			}
		}
	}

	return fmt.Errorf("unknown websocet message: %v", wsResp)
}

func (okV3Ws *OKExV3FutureWs) parseChannel(channel string) (string, error) {
	metas := strings.Split(channel, "/")
	if len(metas) != 2 {
		return "", fmt.Errorf("unknown channel: %s", channel)
	}
	return metas[1], nil
}

func (okV3Ws *OKExV3FutureWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}
