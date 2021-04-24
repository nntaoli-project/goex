package huobi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/Jameslu041/goex"
	"github.com/Jameslu041/goex/internal/logger"
	"strings"
	"sync"
	"time"
)

type HbdmSwapWs struct {
	*WsBuilder
	sync.Once
	wsConn *WsConn

	tickerCallback func(*FutureTicker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade, string)
}

func NewHbdmSwapWs() *HbdmSwapWs {
	ws := &HbdmSwapWs{WsBuilder: NewWsBuilder()}
	ws.WsBuilder = ws.WsBuilder.
		WsUrl("wss://api.hbdm.com/swap-ws").
		//ProxyUrl("socks5://127.0.0.1:1080").
		AutoReconnect().
		DecompressFunc(GzipDecompress).
		ProtoHandleFunc(ws.handle)
	return ws
}

//构建usdt本位永续合约ws
func NewHbdmLinearSwapWs() *HbdmSwapWs {
	ws := &HbdmSwapWs{WsBuilder: NewWsBuilder()}
	ws.WsBuilder = ws.WsBuilder.
		WsUrl("wss://api.hbdm.com/linear-swap-ws").
		//ProxyUrl("socks5://127.0.0.1:1080").
		AutoReconnect().
		DecompressFunc(GzipDecompress).
		ProtoHandleFunc(ws.handle)
	return ws
}

func (ws *HbdmSwapWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string)) {
	ws.tickerCallback = tickerCallback
	ws.depthCallback = depthCallback
	ws.tradeCallback = tradeCallback
}

func (ws *HbdmSwapWs) TickerCallback(call func(ticker *FutureTicker)) {
	ws.tickerCallback = call
}
func (ws *HbdmSwapWs) TradeCallback(call func(trade *Trade, contract string)) {
	ws.tradeCallback = call
}

func (ws *HbdmSwapWs) DepthCallback(call func(depth *Depth)) {
	ws.depthCallback = call
}

func (ws *HbdmSwapWs) SubscribeTicker(pair CurrencyPair, contract string) error {
	if ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}

	if contract == SWAP_CONTRACT || contract == SWAP_USDT_CONTRACT {
		return ws.subscribe(map[string]interface{}{
			"id":  "ticker_1",
			"sub": fmt.Sprintf("market.%s.detail", pair.ToSymbol("-"))})
	}

	return errors.New("not implement")
}

func (ws *HbdmSwapWs) SubscribeDepth(pair CurrencyPair, contract string) error {
	if ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	if contract == SWAP_CONTRACT || contract == SWAP_USDT_CONTRACT {
		return ws.subscribe(map[string]interface{}{
			"id":  "swap.depth",
			"sub": fmt.Sprintf("market.%s.depth.step6", pair.ToSymbol("-"))})
	}

	return errors.New("not implement")
}

func (ws *HbdmSwapWs) SubscribeTrade(pair CurrencyPair, contract string) error {
	if ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}

	if contract == SWAP_CONTRACT || contract == SWAP_USDT_CONTRACT {
		return ws.subscribe(map[string]interface{}{
			"id":  "swap_trade_3",
			"sub": fmt.Sprintf("market.%s.trade.detail", pair.ToSymbol("-"))})
	}

	return errors.New("not implement")
}

func (ws *HbdmSwapWs) subscribe(sub map[string]interface{}) error {
	//	log.Println(sub)
	ws.connectWs()
	return ws.wsConn.Subscribe(sub)
}

func (ws *HbdmSwapWs) connectWs() {
	ws.Do(func() {
		ws.wsConn = ws.WsBuilder.Build()
	})
}

func (ws *HbdmSwapWs) handle(msg []byte) error {
	logger.Debug("ws message data:", string(msg))
	//心跳
	if bytes.Contains(msg, []byte("ping")) {
		pong := bytes.ReplaceAll(msg, []byte("ping"), []byte("pong"))
		ws.wsConn.SendMessage(pong)
		return nil
	}

	var resp WsResponse
	err := json.Unmarshal(msg, &resp)
	if err != nil {
		return err
	}

	if resp.Ch == "" {
		logger.Warnf("[%s] ch == \"\" , msg=%s", ws.wsConn.WsUrl, string(msg))
		return nil
	}

	ts := time.Now()
	if resp.Ts > 0 {
		ts = time.Unix(0, resp.Ts*int64(time.Millisecond))
	}

	pair, contract, err := ws.parseCurrencyAndContract(resp.Ch)
	if err != nil {
		logger.Errorf("[%s] parse currency and contract err=%s", ws.wsConn.WsUrl, err)
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
		dep.Pair = pair
		dep.UTime = ts

		ws.depthCallback(&dep)

		return nil
	}

	if strings.HasSuffix(resp.Ch, "trade.detail") {
		var tradeResp TradeResponse
		err := json.Unmarshal(resp.Tick, &tradeResp)
		if err != nil {
			return err
		}

		trades := ws.parseTrade(tradeResp)
		for _, v := range trades {
			v.Pair = pair
			ws.tradeCallback(&v, contract)
		}

		return nil
	}

	if strings.HasSuffix(resp.Ch, ".detail") {
		var detail DetailResponse
		err := json.Unmarshal(resp.Tick, &detail)
		if err != nil {
			return err
		}

		ticker := ws.parseTicker(detail)
		ticker.ContractType = contract
		ticker.Pair = pair
		ticker.Date = uint64(detail.Id)

		ws.tickerCallback(&ticker)

		return nil
	}

	logger.Errorf("[%s] unknown message, msg=%s", ws.wsConn.WsUrl, string(msg))

	return nil
}

func (ws *HbdmSwapWs) parseTicker(r DetailResponse) FutureTicker {
	return FutureTicker{Ticker: &Ticker{Last: r.Close, High: r.High, Low: r.Low, Vol: r.Amount}}
}

func (ws *HbdmSwapWs) parseCurrencyAndContract(ch string) (CurrencyPair, string, error) {
	el := strings.Split(ch, ".")

	if len(el) < 2 {
		return UNKNOWN_PAIR, "", errors.New(ch)
	}

	pair := NewCurrencyPair3(el[1], "-")
	if pair.CurrencyB.Eq(USD) {
		return pair, SWAP_CONTRACT, nil
	}

	return pair, SWAP_USDT_CONTRACT, nil
}

func (ws *HbdmSwapWs) parseTrade(r TradeResponse) []Trade {
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

func (ws *HbdmSwapWs) adaptTime(tm string) int64 {
	format := "2006-01-02 15:04:05"
	day := time.Now().Format("2006-01-02")
	local, _ := time.LoadLocation("Asia/Chongqing")
	t, _ := time.ParseInLocation(format, day+" "+tm, local)
	return t.UnixNano() / 1e6

}
