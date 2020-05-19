package okex

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OKExV3SpotWs struct {
	base           *OKEx
	v3Ws           *OKExV3Ws
	tickerCallback func(*Ticker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade)
	klineCallback  func(*Kline, KlinePeriod)
}

func NewOKExSpotV3Ws(base *OKEx) *OKExV3SpotWs {
	okV3Ws := &OKExV3SpotWs{
		base: base,
	}
	okV3Ws.v3Ws = NewOKExV3Ws(base, okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3SpotWs) TickerCallback(tickerCallback func(*Ticker)) {
	okV3Ws.tickerCallback = tickerCallback
}

func (okV3Ws *OKExV3SpotWs) DepthCallback(depthCallback func(*Depth)) {
	okV3Ws.depthCallback = depthCallback
}

func (okV3Ws *OKExV3SpotWs) TradeCallback(tradeCallback func(*Trade)) {
	okV3Ws.tradeCallback = tradeCallback
}

func (okV3Ws *OKExV3SpotWs) KLineCallback(klineCallback func(kline *Kline, period KlinePeriod)) {
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3SpotWs) SetCallbacks(tickerCallback func(*Ticker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade),
	klineCallback func(*Kline, KlinePeriod)) {
	okV3Ws.tickerCallback = tickerCallback
	okV3Ws.depthCallback = depthCallback
	okV3Ws.tradeCallback = tradeCallback
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3SpotWs) SubscribeDepth(currencyPair CurrencyPair) error {
	if okV3Ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf("spot/depth5:%s", currencyPair.ToSymbol("-"))}})
}

func (okV3Ws *OKExV3SpotWs) SubscribeTicker(currencyPair CurrencyPair) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf("spot/ticker:%s", currencyPair.ToSymbol("-"))}})
}

func (okV3Ws *OKExV3SpotWs) SubscribeTrade(currencyPair CurrencyPair) error {
	if okV3Ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf("spot/trade:%s", currencyPair.ToSymbol("-"))}})
}

func (okV3Ws *OKExV3SpotWs) SubscribeKline(currencyPair CurrencyPair, period int) error {
	if okV3Ws.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	seconds := adaptKLinePeriod(KlinePeriod(period))
	if seconds == -1 {
		return fmt.Errorf("unsupported kline period %d in okex", period)
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf("spot/candle%ds:%s", seconds, currencyPair.ToSymbol("-"))}})
}

func (okV3Ws *OKExV3SpotWs) getCurrencyPair(instrumentId string) CurrencyPair {
	return NewCurrencyPair3(instrumentId, "-")
}

func (okV3Ws *OKExV3SpotWs) handle(ch string, data json.RawMessage) error {
	var (
		err           error
		tickers       []spotTickerResponse
		depthResp     []depthResponse
		dep           Depth
		tradeResponse []struct {
			Side         string  `json:"side"`
			TradeId      int64   `json:"trade_id,string"`
			Price        float64 `json:"price,string"`
			Qty          float64 `json:"qty,string"`
			InstrumentId string  `json:"instrument_id"`
			Timestamp    string  `json:"timestamp"`
		}
		candleResponse []struct {
			Candle       []string `json:"candle"`
			InstrumentId string   `json:"instrument_id"`
		}
	)

	switch ch {
	case "spot/ticker":
		err = json.Unmarshal(data, &tickers)
		if err != nil {
			return err
		}

		for _, t := range tickers {
			date, _ := time.Parse(time.RFC3339, t.Timestamp)
			okV3Ws.tickerCallback(&Ticker{
				Pair: okV3Ws.getCurrencyPair(t.InstrumentId),
				Last: t.Last,
				Buy:  t.BestBid,
				Sell: t.BestAsk,
				High: t.High24h,
				Low:  t.Low24h,
				Vol:  t.BaseVolume24h,
				Date: uint64(date.UnixNano() / int64(time.Millisecond)),
			})
		}
		return nil
	case "spot/depth5":
		err := json.Unmarshal(data, &depthResp)
		if err != nil {
			logger.Error(err)
			return err
		}
		if len(depthResp) == 0 {
			return nil
		}

		dep.Pair = okV3Ws.getCurrencyPair(depthResp[0].InstrumentId)
		dep.UTime, _ = time.Parse(time.RFC3339, depthResp[0].Timestamp)
		for _, itm := range depthResp[0].Asks {
			dep.AskList = append(dep.AskList, DepthRecord{
				Price:  ToFloat64(itm[0]),
				Amount: ToFloat64(itm[1])})
		}
		for _, itm := range depthResp[0].Bids {
			dep.BidList = append(dep.BidList, DepthRecord{
				Price:  ToFloat64(itm[0]),
				Amount: ToFloat64(itm[1])})
		}
		sort.Sort(sort.Reverse(dep.AskList))
		//call back func
		okV3Ws.depthCallback(&dep)
		return nil
	case "spot/trade":
		err := json.Unmarshal(data, &tradeResponse)
		if err != nil {
			logger.Error("unmarshal error :", err)
			return err
		}

		for _, resp := range tradeResponse {
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
				Pair:   okV3Ws.getCurrencyPair(resp.InstrumentId),
			})
		}
		return nil
	default:
		if strings.HasPrefix(ch, "spot/candle") {
			err := json.Unmarshal(data, &candleResponse)
			if err != nil {
				return err
			}
			periodMs := strings.TrimPrefix(ch, "spot/candle")
			periodMs = strings.TrimSuffix(periodMs, "s")
			for _, k := range candleResponse {
				pair := okV3Ws.getCurrencyPair(k.InstrumentId)
				tm, _ := time.Parse(time.RFC3339, k.Candle[0])
				okV3Ws.klineCallback(&Kline{
					Pair:      pair,
					Timestamp: tm.Unix(),
					Open:      ToFloat64(k.Candle[1]),
					Close:     ToFloat64(k.Candle[4]),
					High:      ToFloat64(k.Candle[2]),
					Low:       ToFloat64(k.Candle[3]),
					Vol:       ToFloat64(k.Candle[5]),
				}, adaptSecondsToKlinePeriod(ToInt(periodMs)))
			}
			return nil
		}
	}

	return fmt.Errorf("unknown websocket message: %s", string(data))
}

func (okV3Ws *OKExV3SpotWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}
