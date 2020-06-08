package huobi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"strings"
	"sync"
	"time"
)

type WsResponse struct {
	Ch   string
	Ts   int64
	Tick json.RawMessage
}

type TradeResponse struct {
	Id   int64
	Ts   int64
	Data []struct {
		Id        int64
		Amount    float64
		Price     float64
		Direction string
		Ts        int64
	}
}

//"id": 1539842340,
//"mrid": 268041138,
//"open": 6740.47,
//"close": 7800,
//"high": 7800,
//"low": 6726.13,
//"amount": 477.1200312075244664773339914558562673572,
//"vol": 32414,
//"count": 1716
//}
type DetailResponse struct {
	Id     int64
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Amount float64
	Vol    float64
	Count  int64
}

type DepthResponse struct {
	Bids [][]float64
	Asks [][]float64
	Ts   int64 `json:"ts"`
}

type HbdmWs struct {
	*WsBuilder
	sync.Once
	wsConn *WsConn

	tickerCallback func(*FutureTicker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade, string)
}

func NewHbdmWs() *HbdmWs {
	hbdmWs := &HbdmWs{WsBuilder: NewWsBuilder()}
	hbdmWs.WsBuilder = hbdmWs.WsBuilder.
		WsUrl("wss://api.hbdm.com/ws").
		AutoReconnect().
		//Heartbeat([]byte("{\"event\": \"ping\"} "), 30*time.Second).
		//Heartbeat(func() []byte { return []byte("{\"op\":\"ping\"}") }(), 5*time.Second).
		DecompressFunc(GzipDecompress).
		ProtoHandleFunc(hbdmWs.handle)
	go hbdmInit()
	return hbdmWs
}

func (hbdmWs *HbdmWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string)) {
	hbdmWs.tickerCallback = tickerCallback
	hbdmWs.depthCallback = depthCallback
	hbdmWs.tradeCallback = tradeCallback
}

func (hbdmWs *HbdmWs) TickerCallback(call func(ticker *FutureTicker)) {
	hbdmWs.tickerCallback = call
}
func (hbdmWs *HbdmWs) TradeCallback(call func(trade *Trade, contract string)) {
	hbdmWs.tradeCallback = call
}

func (hbdmWs *HbdmWs) DepthCallback(call func(depth *Depth)) {
	hbdmWs.depthCallback = call
}

func (hbdmWs *HbdmWs) SubscribeTicker(pair CurrencyPair, contract string) error {
	if hbdmWs.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	return hbdmWs.subscribe(map[string]interface{}{
		"id":  "ticker_1",
		"sub": fmt.Sprintf("market.%s_%s.detail", pair.CurrencyA.Symbol, hbdmWs.adaptContractSymbol(contract))})
}

func (hbdmWs *HbdmWs) SubscribeDepth(pair CurrencyPair, contract string) error {
	if hbdmWs.depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	return hbdmWs.subscribe(map[string]interface{}{
		"id":  "futures.depth",
		"sub": fmt.Sprintf("market.%s_%s.depth.size_20.high_freq", pair.CurrencyA.Symbol, hbdmWs.adaptContractSymbol(contract))})
}

func (hbdmWs *HbdmWs) SubscribeTrade(pair CurrencyPair, contract string) error {
	if hbdmWs.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	return hbdmWs.subscribe(map[string]interface{}{
		"id":  "trade_3",
		"sub": fmt.Sprintf("market.%s_%s.trade.detail", pair.CurrencyA.Symbol, hbdmWs.adaptContractSymbol(contract))})
}

func (hbdmWs *HbdmWs) subscribe(sub map[string]interface{}) error {
	//	log.Println(sub)
	hbdmWs.connectWs()
	return hbdmWs.wsConn.Subscribe(sub)
}

func (hbdmWs *HbdmWs) connectWs() {
	hbdmWs.Do(func() {
		hbdmWs.wsConn = hbdmWs.WsBuilder.Build()
	})
}

func (hbdmWs *HbdmWs) handle(msg []byte) error {
	//心跳
	if bytes.Contains(msg, []byte("ping")) {
		pong := bytes.ReplaceAll(msg, []byte("ping"), []byte("pong"))
		hbdmWs.wsConn.SendMessage(pong)
		return nil
	}
	
	var resp WsResponse
	err := json.Unmarshal(msg, &resp)
	if err != nil {
		return err
	}

	if resp.Ch == "" {
		logger.Warnf("[%s] ch == \"\" , msg=%s", hbdmWs.wsConn.WsUrl, string(msg))
		return nil
	}

	pair, contract, err := hbdmWs.parseCurrencyAndContract(resp.Ch)
	if err != nil {
		logger.Errorf("[%s] parse currency and contract err=%s", hbdmWs.wsConn.WsUrl, err)
		return err
	}

	if strings.Contains(resp.Ch, ".depth.") {
		var depResp DepthResponse
		err := json.Unmarshal(resp.Tick, &depResp)
		if err != nil {
			return err
		}

		dep := ParseDepthFromResponse(depResp)
		dep.ContractType = contract
		dep.ContractId = hbdmWs.getContractId(contract)
		dep.Pair = pair
		dep.UTime = time.Unix(0, resp.Ts*int64(time.Millisecond))

		hbdmWs.depthCallback(&dep)
		return nil
	}

	if strings.HasSuffix(resp.Ch, "trade.detail") {
		var tradeResp TradeResponse
		err := json.Unmarshal(resp.Tick, &tradeResp)
		if err != nil {
			return err
		}
		trades := hbdmWs.parseTrade(tradeResp)
		for _, v := range trades {
			v.Pair = pair
			hbdmWs.tradeCallback(&v, contract)
		}
		return nil
	}

	if strings.HasSuffix(resp.Ch, ".detail") {
		var detail DetailResponse
		err := json.Unmarshal(resp.Tick, &detail)
		if err != nil {
			return err
		}
		ticker := hbdmWs.parseTicker(detail)
		ticker.ContractType = contract
		ticker.Pair = pair
		hbdmWs.tickerCallback(&ticker)
		return nil
	}

	logger.Errorf("[%s] unknown message, msg=%s", hbdmWs.wsConn.WsUrl, string(msg))

	return nil
}

func (hbdmWs *HbdmWs) parseTicker(r DetailResponse) FutureTicker {
	return FutureTicker{Ticker: &Ticker{High: r.High, Low: r.Low, Vol: r.Amount}}
}

func (hbdmWs *HbdmWs) parseCurrencyAndContract(ch string) (CurrencyPair, string, error) {
	el := strings.Split(ch, ".")
	if len(el) < 2 {
		return UNKNOWN_PAIR, "", errors.New(ch)
	}
	cs := strings.Split(el[1], "_")
	contract := ""
	switch cs[1] {
	case "CQ":
		contract = QUARTER_CONTRACT
	case "NW":
		contract = NEXT_WEEK_CONTRACT
	case "CW":
		contract = THIS_WEEK_CONTRACT
	}
	return NewCurrencyPair(NewCurrency(cs[0], ""), USD), contract, nil
}

func (hbdmWs *HbdmWs) parseTrade(r TradeResponse) []Trade {
	var trades []Trade
	for _, v := range r.Data {
		trades = append(trades, Trade{
			Tid:    v.Id,
			Price:  v.Price,
			Amount: v.Amount,
			Type:   AdaptTradeSide(v.Direction),
			Date:   v.Ts})
	}
	return trades
}

func (hbdmWs *HbdmWs) adaptContractSymbol(contract string) string {
	//log.Println(contract)
	switch contract {
	case QUARTER_CONTRACT:
		return "CQ"
	case NEXT_WEEK_CONTRACT:
		return "NW"
	case THIS_WEEK_CONTRACT:
		return "CW"
	}
	return ""
}

func (hbdmWs *HbdmWs) adaptTime(tm string) int64 {
	format := "2006-01-02 15:04:05"
	day := time.Now().Format("2006-01-02")
	local, _ := time.LoadLocation("Asia/Chongqing")
	t, _ := time.ParseInLocation(format, day+" "+tm, local)
	return t.UnixNano() / 1e6

}

func (hbdmWs *HbdmWs) getContractId(alias string) string {
	for _, info := range FuturesContractInfos {
		if info.ContractType == alias {
			return info.InstrumentID
		}
	}
	return ""
}
