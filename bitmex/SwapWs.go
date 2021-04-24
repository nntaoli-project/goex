package bitmex

import (
	"encoding/json"
	"fmt"
	. "github.com/Jameslu041/goex"
	"github.com/Jameslu041/goex/internal/logger"
	"sort"
	"sync"
	"time"
)

type SubscribeOp struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

type wsMessage struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   json.RawMessage
}

type tickerData struct {
	Symbol          string  `json:"symbol"`
	MakerFee        float64 `json:"makerFee"`
	TakerFee        float64 `json:"takerFee"`
	LastPrice       float64 `json:"lastPrice"`
	HighPrice       float64 `json:"highPrice"`
	LowPrice        float64 `json:"lowPrice"`
	AskPrice        float64 `json:"askPrice"`
	BidPrice        float64 `json:"bidPrice"`
	HomeNotional24h float64 `json:"homeNotional24h"`
	Turnover24h     float64 `json:"turnover24h"`
	Timestamp       string  `json:"timestamp"`
}

type depthData struct {
	Symbol    string          `json:"symbol"`
	Bids      [][]interface{} `json:"bids"`
	Asks      [][]interface{} `json:"asks"`
	Timestamp string          `json:"timestamp"`
}

type SwapWs struct {
	c         *WsConn
	once      sync.Once
	wsBuilder *WsBuilder

	depthCall  func(depth *Depth)
	tickerCall func(ticker *FutureTicker)

	tickerCacheMap map[string]FutureTicker
}

func NewSwapWs() *SwapWs {
	s := new(SwapWs)
	s.wsBuilder = NewWsBuilder().DisableEnableCompression().WsUrl("wss://www.bitmex.com/realtime")
	s.wsBuilder = s.wsBuilder.Heartbeat(func() []byte { return []byte("ping") }, 5*time.Second)
	s.wsBuilder = s.wsBuilder.ProtoHandleFunc(s.handle).AutoReconnect()
	//s.c = wsBuilder.Build()
	s.tickerCacheMap = make(map[string]FutureTicker, 10)
	return s
}

func (s *SwapWs) connect() {
	s.once.Do(func() {
		s.c = s.wsBuilder.Build()
	})
}

func (s *SwapWs) DepthCallback(f func(depth *Depth)) {
	s.depthCall = f
}

func (s *SwapWs) TickerCallback(f func(ticker *FutureTicker)) {
	s.tickerCall = f
}

func (s *SwapWs) TradeCallback(f func(trade *Trade, contract string)) {
	panic("implement me")
}

func (s *SwapWs) SubscribeDepth(pair CurrencyPair, contractType string) error {
	//{"op": "subscribe", "args": ["orderBook10:XBTUSD"]}
	s.connect()

	op := SubscribeOp{
		Op: "subscribe",
		Args: []string{
			fmt.Sprintf("orderBook10:%s", AdaptCurrencyPairToSymbol(pair, contractType)),
		},
	}
	return s.c.Subscribe(op)
}

func (s *SwapWs) SubscribeTicker(pair CurrencyPair, contractType string) error {
	s.connect()

	return s.c.Subscribe(SubscribeOp{
		Op: "subscribe",
		Args: []string{
			"instrument:" + AdaptCurrencyPairToSymbol(pair, contractType),
		},
	})
}

func (s *SwapWs) SubscribeTrade(pair CurrencyPair, contractType string) error {
	panic("implement me")
}

func (s *SwapWs) handle(data []byte) error {
	if string(data) == "pong" {
		return nil
	}

	var msg wsMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Errorf("unmarshal error , message: %s", string(data))
		return err
	}

	switch msg.Table {
	case "orderBook10":
		if msg.Action != "update" {
			return nil
		}

		var (
			depthData []depthData
			dep       Depth
		)

		err = json.Unmarshal(msg.Data, &depthData)
		if err != nil {
			logger.Errorf("unmarshal depth data error , data: %s", string(msg.Data))
			return nil
		}

		if len(depthData) == 0 {
			logger.Warn("depth data len==0 ??")
			return nil
		}

		dep.UTime, _ = time.Parse(time.RFC3339, depthData[0].Timestamp)
		dep.Pair, dep.ContractType = AdaptWsSymbol(depthData[0].Symbol)

		for _, item := range depthData[0].Bids {
			dep.BidList = append(dep.BidList, DepthRecord{
				Price:  ToFloat64(item[0]),
				Amount: ToFloat64(item[1]),
			})
		}

		for _, item := range depthData[0].Asks {
			dep.AskList = append(dep.AskList, DepthRecord{
				Price:  ToFloat64(item[0]),
				Amount: ToFloat64(item[1]),
			})
		}

		sort.Sort(sort.Reverse(dep.AskList))

		s.depthCall(&dep)
	case "instrument":
		var tickerData []tickerData

		err = json.Unmarshal(msg.Data, &tickerData)
		if err != nil {
			logger.Errorf("ticker data unmarshal error , data: %s", string(msg.Data))
			return err
		}

		if msg.Action == "partial" {
			ticker := s.tickerCacheMap[tickerData[0].Symbol]
			ticker.Ticker = new(Ticker)
			ticker.Pair, ticker.ContractType = AdaptWsSymbol(tickerData[0].Symbol)
			ticker.Vol = tickerData[0].HomeNotional24h
			ticker.Last = tickerData[0].LastPrice
			ticker.Sell = tickerData[0].AskPrice
			ticker.Buy = tickerData[0].BidPrice
			ticker.High = tickerData[0].HighPrice
			ticker.Low = tickerData[0].LowPrice

			tickerTime, _ := time.Parse(time.RFC3339, tickerData[0].Timestamp)
			ticker.Date = uint64(tickerTime.Unix())

			s.tickerCacheMap[tickerData[0].Symbol] = ticker
			s.tickerCall(&ticker)
		}

		if msg.Action == "update" {
			ticker := s.tickerCacheMap[tickerData[0].Symbol]
			tickerTime, _ := time.Parse(time.RFC3339, tickerData[0].Timestamp)
			ticker.Date = uint64(tickerTime.Unix())

			if tickerData[0].LastPrice > 0 {
				ticker.Last = tickerData[0].LastPrice
			}

			if tickerData[0].AskPrice > 0 {
				ticker.Sell = tickerData[0].AskPrice
			}

			if tickerData[0].BidPrice > 0 {
				ticker.Buy = tickerData[0].BidPrice
			}

			if tickerData[0].HighPrice > 0 {
				ticker.High = tickerData[0].HighPrice
			}

			if tickerData[0].LowPrice > 0 {
				ticker.Low = tickerData[0].LowPrice
			}

			if tickerData[0].HomeNotional24h > 0 {
				ticker.Vol = tickerData[0].HomeNotional24h
			}

			s.tickerCacheMap[tickerData[0].Symbol] = ticker
			s.tickerCall(&ticker)
		}
	default:
		logger.Warnf("unknown ws message: %s", string(data))
	}

	return nil
}
