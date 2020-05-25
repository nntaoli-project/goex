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

type SpotWs struct {
	*WsBuilder
	sync.Once
	wsConn *WsConn

	tickerCallback func(*Ticker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade)
}

func NewSpotWs() *SpotWs {
	ws := &SpotWs{
		WsBuilder: NewWsBuilder(),
	}
	ws.WsBuilder = ws.WsBuilder.
		WsUrl("wss://api.huobi.pro/ws").
		AutoReconnect().
		DecompressFunc(GzipDecompress).
		ProtoHandleFunc(ws.handle)
	return ws
}

func (ws *SpotWs) DepthCallback(call func(depth *Depth)) {
	ws.depthCallback = call
}

func (ws *SpotWs) TickerCallback(call func(ticker *Ticker)) {
	ws.tickerCallback = call
}
func (ws *SpotWs) TradeCallback(call func(trade *Trade)) {
	ws.tradeCallback = call
}

func (ws *SpotWs) connectWs() {
	ws.Do(func() {
		ws.wsConn = ws.WsBuilder.Build()
	})
}

func (ws *SpotWs) subscribe(sub map[string]interface{}) error {
	ws.connectWs()
	return ws.wsConn.Subscribe(sub)
}

func (ws *SpotWs) SubscribeDepth(pair CurrencyPair) error {
	if ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	return ws.subscribe(map[string]interface{}{
		"id":  "spot.depth",
		"sub": fmt.Sprintf("market.%s.mbp.refresh.20", pair.ToLower().ToSymbol(""))})
}

func (ws *SpotWs) SubscribeTicker(pair CurrencyPair) error {
	if ws.tickerCallback == nil {
		return errors.New("please set ticker call back func")
	}
	return ws.subscribe(map[string]interface{}{
		"id":  "spot.ticker",
		"sub": fmt.Sprintf("market.%s.detail", pair.ToLower().ToSymbol("")),
	})
	return nil
}

func (ws *SpotWs) SubscribeTrade(pair CurrencyPair) error {
	return nil
}

func (ws *SpotWs) handle(msg []byte) error {
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

	currencyPair := ParseCurrencyPairFromSpotWsCh(resp.Ch)
	if strings.Contains(resp.Ch, "mbp.refresh") {
		var (
			depthResp DepthResponse
		)

		err := json.Unmarshal(resp.Tick, &depthResp)
		if err != nil {
			return err
		}

		dep := ParseDepthFromResponse(depthResp)
		dep.Pair = currencyPair
		dep.UTime = time.Unix(0, resp.Ts*int64(time.Millisecond))
		ws.depthCallback(&dep)

		return nil
	}

	if strings.Contains(resp.Ch, ".detail") {
		var tickerResp DetailResponse
		err := json.Unmarshal(resp.Tick, &tickerResp)
		if err != nil {
			return err
		}
		ws.tickerCallback(&Ticker{
			Pair: currencyPair,
			Last: tickerResp.Close,
			High: tickerResp.High,
			Low:  tickerResp.Low,
			Vol:  tickerResp.Amount,
			Date: uint64(resp.Ts),
		})
		return nil
	}

	logger.Errorf("[%s] unknown message ch , msg=%s", ws.wsConn.WsUrl, string(msg))

	return nil
}
