package okcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/nntaoli-project/GoEx"
)

type OKExV3FutureWs struct {
	*WsBuilder
	sync.Once
	wsConn             *WsConn
	loginCh            chan interface{}
	logined            bool
	loginLock          *sync.Mutex
	dataParser         *OKExV3DataParser
	contractIDProvider IContractIDProvider
	apiKey             string
	apiSecretKey       string
	passphrase         string
	authoriedSubs      []map[string]interface{}

	tickerCallback func(*FutureTicker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade, string)
	klineCallback  func(*FutureKline, int)
	orderCallback  func(*FutureOrder, string)
}

func getSign(apiSecretKey, timestamp, method, url, body string) (string, error) {
	data := timestamp + method + url + body
	return GetParamHmacSHA256Base64Sign(apiSecretKey, data)
}

func NewOKExV3FutureWs(contractIDProvider IContractIDProvider) *OKExV3FutureWs {
	okV3Ws := &OKExV3FutureWs{}
	okV3Ws.contractIDProvider = contractIDProvider
	okV3Ws.dataParser = NewOKExV3DataParser(contractIDProvider)
	okV3Ws.loginCh = make(chan interface{})
	okV3Ws.logined = false
	okV3Ws.loginLock = &sync.Mutex{}
	okV3Ws.authoriedSubs = make([]map[string]interface{}, 0)
	okV3Ws.WsBuilder = NewWsBuilder().
		WsUrl("wss://real.okex.com:10440/ws/v3").
		Heartbeat([]byte("ping"), 30*time.Second).
		ReconnectIntervalTime(24 * time.Hour).
		UnCompressFunc(FlateUnCompress).
		ProtoHandleFunc(okV3Ws.handle)

	return okV3Ws
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

func (okV3Ws *OKExV3FutureWs) Login(apiKey string, apiSecretKey string, passphrase string) error {
	// already logined
	if okV3Ws.logined {
		return nil
	}
	okV3Ws.connectWs()
	okV3Ws.apiKey = apiKey
	okV3Ws.apiSecretKey = apiSecretKey
	okV3Ws.passphrase = passphrase
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

func (okV3Ws *OKExV3FutureWs) clearChan(c chan interface{}) {
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
	apiKey := okV3Ws.apiKey
	apiSecretKey := okV3Ws.apiSecretKey
	passphrase := okV3Ws.passphrase
	//clear last login result
	timestamp := okV3Ws.getTimestamp()
	method := "GET"
	url := "/users/self/verify"
	sign, _ := getSign(apiSecretKey, timestamp, method, url, "")
	op := map[string]interface{}{
		"op":   "login",
		"args": []string{apiKey, passphrase, timestamp, sign}}

	err := okV3Ws.wsConn.SendJsonMessage(op)
	if err != nil {
		return err
	}
	event := <-okV3Ws.loginCh
	if v, ok := event.(map[string]interface{}); ok {
		var success = false
		switch s := v["success"].(type) {
		case bool:
			success = s
		case string:
			success = s == "true"
		}
		if success {
			log.Println("login success:", event)
			return nil
		}
	}
	log.Println("login failed:", event)
	return fmt.Errorf("login failed: %v", event)
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

	symbol, err := okV3Ws.contractIDProvider.GetContractID(currencyPair, contractType)
	if err != nil {
		return err
	}

	chName := fmt.Sprintf("%s/depth5:%s", okV3Ws.getTablePrefix(currencyPair, contractType), symbol)

	return okV3Ws.subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) SubscribeTicker(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}

	symbol, err := okV3Ws.contractIDProvider.GetContractID(currencyPair, contractType)
	if err != nil {
		return err
	}

	chName := fmt.Sprintf("%s/ticker:%s", okV3Ws.getTablePrefix(currencyPair, contractType), symbol)

	return okV3Ws.subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) SubscribeTrade(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}

	symbol, err := okV3Ws.contractIDProvider.GetContractID(currencyPair, contractType)
	if err != nil {
		return err
	}

	chName := fmt.Sprintf("%s/trade:%s", okV3Ws.getTablePrefix(currencyPair, contractType), symbol)

	return okV3Ws.subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) SubscribeKline(currencyPair CurrencyPair, contractType string, period int) error {
	if okV3Ws.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	symbol, err := okV3Ws.contractIDProvider.GetContractID(currencyPair, contractType)
	if err != nil {
		return err
	}

	seconds, ok := KlineTypeSecondsMap[period]
	if !ok {
		return fmt.Errorf("unsupported kline period %d in okex", period)
	}

	chName := fmt.Sprintf("%s/candle%ds:%s", okV3Ws.getTablePrefix(currencyPair, contractType), seconds, symbol)

	return okV3Ws.subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FutureWs) SubscribeOrder(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.orderCallback == nil {
		return errors.New("place set order callback func")
	}

	symbol, err := okV3Ws.contractIDProvider.GetContractID(currencyPair, contractType)
	if err != nil {
		return err
	}

	chName := fmt.Sprintf("%s/order:%s", okV3Ws.getTablePrefix(currencyPair, contractType), symbol)

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
		okV3Ws.wsConn.ReceiveMessage()
	})
}

func (okV3Ws *OKExV3FutureWs) handle(msg []byte) error {
	if string(msg) == "pong" {
		log.Println(string(msg))
		okV3Ws.wsConn.UpdateActiveTime()
		return nil
	}

	var resp map[string]interface{}

	err := json.Unmarshal(msg, &resp)
	if err != nil {
		return err
	}

	if resp["event"] != nil {
		switch resp["event"].(string) {
		case "subscribe":
			log.Println("subscribed:", resp["channel"].(string))
			return nil
		case "login":
			select {
			case okV3Ws.loginCh <- resp:
				return nil
			default:
				return nil
			}
		case "error":
			var errorCode int
			switch v := resp["errorCode"].(type) {
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
			case okV3Ws.loginCh <- resp:
				return nil
			default:
				return fmt.Errorf("error in websocket: %v", resp)
			}
		}
		return fmt.Errorf("unknown websocet message: %v", resp)
	}

	if resp["table"] != nil {
		ch, err := okV3Ws.parseChannel(resp["table"].(string))
		if err != nil {
			return err
		}

		switch ch {
		case "ticker":
			data, ok := resp["data"].([]interface{})
			if ok {
				valid := true
				for _, v := range data {
					ticker, err := okV3Ws.dataParser.ParseFutureTicker(v)
					if err != nil {
						valid = false
						break
					}
					okV3Ws.tickerCallback(ticker)
				}
				if valid {
					return nil
				}
			}
		case "depth5":
			data, ok := resp["data"].([]interface{})
			if ok {
				valid := true
				for _, v := range data {
					depth, err := okV3Ws.dataParser.ParseDepth(nil, v, 5)
					if err != nil {
						valid = false
						break
					}
					okV3Ws.depthCallback(depth)
				}
				if valid {
					return nil
				}
			}
		case "trade":
			data, ok := resp["data"].([]interface{})
			if ok {
				valid := true
				for _, v := range data {
					trade, contractType, err := okV3Ws.dataParser.ParseTrade(nil, "", v)
					if err != nil {
						valid = false
						break
					}
					okV3Ws.tradeCallback(trade, contractType)
				}
				if valid {
					return nil
				}
			}
		case "order":
			data, ok := resp["data"].([]interface{})
			if ok {
				valid := true
				for _, v := range data {
					order, contractType, err := okV3Ws.dataParser.ParseFutureOrder(v)
					if err != nil {
						valid = false
						break
					}
					okV3Ws.orderCallback(order, contractType)
				}
				if valid {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("unknown websocet message: %v", resp)
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

func (okV3Ws *OKExV3FutureWs) adaptTime(tm string) int64 {
	format := "2006-01-02 15:04:05"
	day := time.Now().Format("2006-01-02")
	local, _ := time.LoadLocation("Asia/Chongqing")
	t, _ := time.ParseInLocation(format, day+" "+tm, local)
	return t.UnixNano() / 1e6
}
