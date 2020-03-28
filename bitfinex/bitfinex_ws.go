package bitfinex

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	. "github.com/nntaoli-project/goex"
)

const ticker = "ticker"
const subscribe = "subscribe"
const subscribed = "subscribed"

type BitfinexWs struct {
	*WsBuilder
	sync.Once
	wsConn   *WsConn
	eventMap map[int64]SubscribeEvent

	tickerCallback func(*Ticker)
}

type SubscribeEvent struct {
	Event     string `json:"event"`
	SubID     string `json:"subId"`
	Channel   string `json:"channel"`
	ChanID    int64  `json:"chanId"`
	Symbol    string `json:"symbol"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair"`
}

type EventMap map[int64]SubscribeEvent

func NewWs() *BitfinexWs {
	bws := &BitfinexWs{WsBuilder: NewWsBuilder(), eventMap: make(map[int64]SubscribeEvent)}
	bws.WsBuilder = bws.WsBuilder.
		WsUrl("wss://api-pub.bitfinex.com/ws/2").
		AutoReconnect().
		ProtoHandleFunc(bws.handle)
	return bws
}

func (bws *BitfinexWs) SetCallbacks(tickerCallback func(*Ticker)) {
	bws.tickerCallback = tickerCallback
}

func (bws *BitfinexWs) SubscribeTicker(pair CurrencyPair) error {
	if bws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	return bws.subscribe(map[string]interface{}{
		"event":   subscribe,
		"channel": ticker,
		"symbol":  convertPairToBitfinexSymbol(pair)})
}

func (bws *BitfinexWs) subscribe(sub map[string]interface{}) error {
	bws.connectWs()
	return bws.wsConn.Subscribe(sub)
}

func (bws *BitfinexWs) connectWs() {
	bws.Do(func() {
		bws.wsConn = bws.WsBuilder.Build()
	})
}

func (bws *BitfinexWs) handle(msg []byte) error {
	var event SubscribeEvent
	if err := json.Unmarshal(msg, &event); err == nil {
		if event.Event == subscribed && event.Channel == ticker {
			bws.eventMap[event.ChanID] = event
			return nil
		}
	}

	var resp []interface{}
	if err := json.Unmarshal(msg, &resp); err == nil {
		if rawTicker, ok := resp[1].([]interface{}); ok {
			channelID := ToInt64(resp[0])
			t := bws.tickerFromRaw(channelID, rawTicker)
			bws.tickerCallback(t)
			return nil
		}
	}

	return nil
}

func (bws *BitfinexWs) resolveCurrencyPair(channelID int64) CurrencyPair {
	ev, ok := bws.eventMap[channelID]
	if ok {
		return symbolToCurrencyPair(ev.Pair)
	}
	return UNKNOWN_PAIR
}

func (bws *BitfinexWs) tickerFromRaw(channelID int64, rawTicker []interface{}) *Ticker {
	pair := bws.resolveCurrencyPair(channelID)
	return &Ticker{
		Pair: pair,
		Buy:  ToFloat64(rawTicker[0]),
		Sell: ToFloat64(rawTicker[2]),
		Last: ToFloat64(rawTicker[6]),
		Vol:  ToFloat64(rawTicker[7]),
		High: ToFloat64(rawTicker[8]),
		Low:  ToFloat64(rawTicker[9]),
		Date: uint64(time.Now().UnixNano() / int64(time.Millisecond)),
	}

}

func convertPairToBitfinexSymbol(pair CurrencyPair) string {
	symbol := pair.ToSymbol("")
	return "t" + symbol
}
