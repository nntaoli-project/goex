package binance

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
)

type FuturesWs struct {
	base  *BinanceFutures
	fOnce sync.Once
	dOnce sync.Once

	wsBuilder *goex.WsBuilder
	f         *goex.WsConn
	d         *goex.WsConn

	depthCallFn  func(depth *goex.Depth)
	tickerCallFn func(ticker *goex.FutureTicker)
	tradeCalFn   func(trade *goex.Trade, contract string)
}

func NewFuturesWs() *FuturesWs {
	futuresWs := new(FuturesWs)

	futuresWs.wsBuilder = goex.NewWsBuilder().
		ProxyUrl(os.Getenv("HTTPS_PROXY")).
		ProtoHandleFunc(futuresWs.handle).AutoReconnect()

	httpCli := &http.Client{
		Timeout: 10 * time.Second,
	}

	if os.Getenv("HTTPS_PROXY") != "" {
		httpCli = &http.Client{
			Transport: &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) {
					return url.Parse(os.Getenv("HTTPS_PROXY"))
				},
			},
			Timeout: 10 * time.Second,
		}
	}

	futuresWs.base = NewBinanceFutures(&goex.APIConfig{
		HttpClient: httpCli,
	})

	return futuresWs
}

func (s *FuturesWs) connectUsdtFutures() {
	s.fOnce.Do(func() {
		s.f = s.wsBuilder.WsUrl(TESTNET_FUTURE_USD_WS_BASE_URL).Build()
	})
}

func (s *FuturesWs) connectFutures() {
	s.dOnce.Do(func() {
		s.d = s.wsBuilder.WsUrl(TESTNET_FUTURE_COIN_WS_BASE_URL).Build()
	})
}

func (s *FuturesWs) DepthCallback(f func(depth *goex.Depth)) {
	s.depthCallFn = f
}

func (s *FuturesWs) TickerCallback(f func(ticker *goex.FutureTicker)) {
	s.tickerCallFn = f
}

func (s *FuturesWs) TradeCallback(f func(trade *goex.Trade, contract string)) {
	s.tradeCalFn = f
}

func (s *FuturesWs) SubscribeDepth(pair goex.CurrencyPair, contractType string) error {
	switch contractType {
	case goex.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@depth10@100ms"},
			Id:     1,
		})
	default:
		s.connectFutures()
		sym, _ := s.base.adaptToSymbol(pair.AdaptUsdtToUsd(), contractType)
		return s.d.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{strings.ToLower(sym) + "@depth10@100ms"},
			Id:     2,
		})
	}
}

func (s *FuturesWs) SubscribeTicker(pair goex.CurrencyPair, contractType string) error {
	switch contractType {
	case goex.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@ticker"},
			Id:     1,
		})
	default:
		s.connectFutures()
		sym, _ := s.base.adaptToSymbol(pair.AdaptUsdtToUsd(), contractType)
		return s.d.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{strings.ToLower(sym) + "@ticker"},
			Id:     2,
		})
	}
}

func (s *FuturesWs) SubscribeTrade(pair goex.CurrencyPair, contractType string) error {
	switch contractType {
	case goex.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@aggTrade"},
			Id:     1,
		})
	default:
		s.connectFutures()
		sym, _ := s.base.adaptToSymbol(pair.AdaptUsdtToUsd(), contractType)
		return s.d.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{strings.ToLower(sym) + "@aggTrade"},
			Id:     1,
		})
	}
}

func (s *FuturesWs) handle(data []byte) error {
	var m = make(map[string]interface{}, 4)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	if e, ok := m["e"].(string); ok && e == "depthUpdate" {
		dep := s.depthHandle(m["b"].([]interface{}), m["a"].([]interface{}))
		dep.ContractType = m["s"].(string)
		symbol, ok := m["ps"].(string)

		if ok {
			dep.Pair = adaptSymbolToCurrencyPair(symbol)
		} else {
			dep.Pair = adaptSymbolToCurrencyPair(dep.ContractType) //usdt swap
		}

		dep.UTime = time.Unix(0, goex.ToInt64(m["T"])*int64(time.Millisecond))
		s.depthCallFn(dep)

		return nil
	}

	if e, ok := m["e"].(string); ok && e == "24hrTicker" {
		s.tickerCallFn(s.tickerHandle(m))
		return nil
	}

	if e, ok := m["e"].(string); ok && e == "aggTrade" {

		contractType := m["s"].(string)
		s.tradeCalFn(s.tradeHandle(m), contractType)
		return nil
	}

	logger.Warn("unknown ws response:", string(data))

	return nil
}

func (s *FuturesWs) depthHandle(bids []interface{}, asks []interface{}) *goex.Depth {
	var dep goex.Depth

	for _, item := range bids {
		bid := item.([]interface{})
		dep.BidList = append(dep.BidList,
			goex.DepthRecord{
				Price:  goex.ToFloat64(bid[0]),
				Amount: goex.ToFloat64(bid[1]),
			})
	}

	for _, item := range asks {
		ask := item.([]interface{})
		dep.AskList = append(dep.AskList, goex.DepthRecord{
			Price:  goex.ToFloat64(ask[0]),
			Amount: goex.ToFloat64(ask[1]),
		})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return &dep
}

func (s *FuturesWs) tickerHandle(m map[string]interface{}) *goex.FutureTicker {
	var ticker goex.FutureTicker
	ticker.Ticker = new(goex.Ticker)

	symbol, ok := m["ps"].(string)
	if ok {
		ticker.Pair = adaptSymbolToCurrencyPair(symbol)
	} else {
		ticker.Pair = adaptSymbolToCurrencyPair(m["s"].(string)) //usdt swap
	}

	ticker.ContractType = m["s"].(string)
	ticker.Date = goex.ToUint64(m["E"])
	ticker.High = goex.ToFloat64(m["h"])
	ticker.Low = goex.ToFloat64(m["l"])
	ticker.Last = goex.ToFloat64(m["c"])
	ticker.Vol = goex.ToFloat64(m["v"])

	return &ticker
}

func (s *FuturesWs) tradeHandle(m map[string]interface{}) *goex.Trade {
	var trade goex.Trade

	symbol, ok := m["s"].(string) // Symbol
	if ok {
		trade.Pair = adaptSymbolToCurrencyPair(symbol) //usdt swap
	}

	trade.Tid = goex.ToInt64(m["a"])      // Aggregate trade ID
	trade.Date = goex.ToInt64(m["E"])     // Event time
	trade.Amount = goex.ToFloat64(m["q"]) // Quantity
	trade.Price = goex.ToFloat64(m["p"])  // Price

	if m["m"].(bool) {
		trade.Type = goex.BUY_MARKET
	} else {
		trade.Type = goex.SELL_MARKET
	}

	return &trade
}
