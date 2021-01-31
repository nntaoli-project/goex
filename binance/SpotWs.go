package binance

import (
	json2 "encoding/json"
	"fmt"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"os"
	"sort"
	"strings"
	"time"
)

type req struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

type resp struct {
	Stream string           `json:"stream"`
	Data   json2.RawMessage `json:"data"`
}

type depthResp struct {
	LastUpdateId int             `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
}

type SpotWs struct {
	c *goex.WsConn

	reqId int

	depthCallFn  func(depth *goex.Depth)
	tickerCallFn func(ticker *goex.Ticker)
	tradeCallFn  func(trade *goex.Trade)
}

func NewSpotWs() *SpotWs {
	spotWs := &SpotWs{}
	logger.Debugf("proxy url: %s", os.Getenv("HTTPS_PROXY"))

	wsBuilder := goex.NewWsBuilder().
		WsUrl("wss://stream.binance.com:9443/stream?streams=depth/miniTicker/ticker/trade").
		ProxyUrl(os.Getenv("HTTPS_PROXY")).
		ProtoHandleFunc(spotWs.handle).AutoReconnect()

	spotWs.c = wsBuilder.Build()
	spotWs.reqId = 1

	return spotWs
}

func (s *SpotWs) DepthCallback(f func(depth *goex.Depth)) {
	s.depthCallFn = f
}

func (s *SpotWs) TickerCallback(f func(ticker *goex.Ticker)) {
	s.tickerCallFn = f
}

func (s *SpotWs) TradeCallback(f func(trade *goex.Trade)) {
	s.tradeCallFn = f
}

func (s *SpotWs) SubscribeDepth(pair goex.CurrencyPair) error {
	return s.SubscribeDepthArgs(pair, 10, 100)
}

func (s *SpotWs) SubscribeDepthArgs(pair goex.CurrencyPair, nums int, speed int) error {
	defer func() {
		s.reqId++
	}()

	return s.c.Subscribe(req{
		Method: "SUBSCRIBE",
		Params: []string{
			fmt.Sprintf("%s@depth%d@%dms", pair.ToLower().ToSymbol(""), nums, speed),
		},
		Id: s.reqId,
	})
}

func (s *SpotWs) SubscribeTicker(pair goex.CurrencyPair) error {
	defer func() {
		s.reqId++
	}()

	return s.c.Subscribe(req{
		Method: "SUBSCRIBE",
		Params: []string{pair.ToLower().ToSymbol("") + "@ticker"},
		Id:     s.reqId,
	})
}

func (s *SpotWs) SubscribeTrade(pair goex.CurrencyPair) error {
	panic("implement me")
}

func (s *SpotWs) handle(data []byte) error {
	var r resp
	err := json2.Unmarshal(data, &r)
	if err != nil {
		logger.Errorf("json unmarshal ws response error [%s] , response data = %s", err, string(data))
		return err
	}

	if strings.HasSuffix(r.Stream, "@depth10@100ms") {
		return s.depthHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	}

	if strings.HasSuffix(r.Stream, "@ticker") {
		return s.tickerHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	}

	logger.Warn("unknown ws response:", string(data))

	return nil
}

func (s *SpotWs) depthHandle(data json2.RawMessage, pair goex.CurrencyPair) error {
	var (
		depthR depthResp
		dep    goex.Depth
		err    error
	)

	err = json2.Unmarshal(data, &depthR)
	if err != nil {
		logger.Errorf("unmarshal depth response error %s[] , response data = %s", err, string(data))
		return err
	}

	dep.UTime = time.Now()
	dep.Pair = pair

	for _, bid := range depthR.Bids {
		dep.BidList = append(dep.BidList, goex.DepthRecord{
			Price:  goex.ToFloat64(bid[0]),
			Amount: goex.ToFloat64(bid[1]),
		})
	}

	for _, ask := range depthR.Asks {
		dep.AskList = append(dep.AskList, goex.DepthRecord{
			Price:  goex.ToFloat64(ask[0]),
			Amount: goex.ToFloat64(ask[1]),
		})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	s.depthCallFn(&dep)

	return nil
}

func (s *SpotWs) tickerHandle(data json2.RawMessage, pair goex.CurrencyPair) error {
	var (
		tickerData = make(map[string]interface{}, 4)
		ticker     goex.Ticker
	)

	err := json2.Unmarshal(data, &tickerData)
	if err != nil {
		logger.Errorf("unmarshal ticker response data error [%s] , data = %s", err, string(data))
		return err
	}

	ticker.Pair = pair
	ticker.Vol = goex.ToFloat64(tickerData["v"])
	ticker.Last = goex.ToFloat64(tickerData["c"])
	ticker.Sell = goex.ToFloat64(tickerData["a"])
	ticker.Buy = goex.ToFloat64(tickerData["b"])
	ticker.High = goex.ToFloat64(tickerData["h"])
	ticker.Low = goex.ToFloat64(tickerData["l"])
	ticker.Date = goex.ToUint64(tickerData["E"])

	s.tickerCallFn(&ticker)

	return nil
}
