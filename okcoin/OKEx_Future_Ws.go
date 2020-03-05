package okcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"log"
	"strings"
	"sync"
	"time"
)

type OKExFutureWs struct {
	*WsBuilder
	sync.Once
	wsConn *WsConn

	tickerCallback func(*FutureTicker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade, string)
}

func NewOKExFutureWs() *OKExFutureWs {
	okWs := &OKExFutureWs{WsBuilder: NewWsBuilder()}
	okWs.WsBuilder = okWs.WsBuilder.
		WsUrl("wss://real.okex.com:10440/ws/v1").
		AutoReconnect().
		Heartbeat(func() []byte {
			return []byte("{\"event\": \"ping\"} ")
		}, 30*time.Second).
		UnCompressFunc(FlateUnCompress).
		ProtoHandleFunc(okWs.handle)
	return okWs
}

func (okWs *OKExFutureWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string)) {
	okWs.tickerCallback = tickerCallback
	okWs.depthCallback = depthCallback
	okWs.tradeCallback = tradeCallback
}

func (okWs *OKExFutureWs) SubscribeTicker(pair CurrencyPair, contract string) error {
	if okWs.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	return okWs.subscribe(map[string]interface{}{
		"event":   "addChannel",
		"channel": fmt.Sprintf("ok_sub_futureusd_%s_ticker_%s", strings.ToLower(pair.CurrencyA.Symbol), contract)})
}

func (okWs *OKExFutureWs) SubscribeDepth(pair CurrencyPair, contract string, size int) error {
	if okWs.depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	return okWs.subscribe(map[string]interface{}{
		"event":   "addChannel",
		"channel": fmt.Sprintf("ok_sub_futureusd_%s_depth_%s_%d", strings.ToLower(pair.CurrencyA.Symbol), contract, size)})
}

func (okWs *OKExFutureWs) SubscribeTrade(pair CurrencyPair, contract string) error {
	if okWs.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	return okWs.subscribe(map[string]interface{}{
		"event":   "addChannel",
		"channel": fmt.Sprintf("ok_sub_futureusd_%s_trade_%s", strings.ToLower(pair.CurrencyA.Symbol), contract)})
}

func (okWs *OKExFutureWs) subscribe(sub map[string]interface{}) error {
	okWs.connectWs()
	return okWs.wsConn.Subscribe(sub)
}

func (okWs *OKExFutureWs) connectWs() {
	okWs.Do(func() {
		okWs.wsConn = okWs.WsBuilder.Build()
	})
}

func (okWs *OKExFutureWs) handle(msg []byte) error {
	//log.Println(string(msg))
	if string(msg) == "{\"event\":\"pong\"}" {
		//	log.Println(string(msg))
		return nil
	}

	var resp []WsBaseResp
	err := json.Unmarshal(msg, &resp)
	if err != nil {
		return err
	}

	if len(resp) < 0 {
		return nil
	}

	if resp[0].Channel == "addChannel" {
		log.Println("subscribe:", string(resp[0].Data))
		return nil
	}

	pair, contract, ch := okWs.parseChannel(resp[0].Channel)

	if ch == "ticker" {
		var t FutureTicker
		err := json.Unmarshal(resp[0].Data, &t)
		if err != nil {
			return err
		}
		t.ContractType = contract
		t.Pair = pair
		okWs.tickerCallback(&t)
		return nil
	}

	if ch == "depth" {
		var (
			d    Depth
			data struct {
				Asks      [][]float64 `json:"asks"`
				Bids      [][]float64 `json:"bids"`
				Timestamp int64       `json:"timestamp"`
			}
		)

		err := json.Unmarshal(resp[0].Data, &data)
		if err != nil {
			return err
		}

		for _, a := range data.Asks {
			d.AskList = append(d.AskList, DepthRecord{a[0], a[1]})
		}

		for _, b := range data.Bids {
			d.BidList = append(d.BidList, DepthRecord{b[0], b[1]})
		}

		d.Pair = pair
		d.ContractType = contract
		d.UTime = time.Unix(data.Timestamp/1000, 0)
		okWs.depthCallback(&d)

		return nil
	}

	if ch == "trade" {
		var data TradeData
		err := json.Unmarshal(resp[0].Data, &data)
		if err != nil {
			return err
		}

		for _, td := range data {
			side := TradeSide(SELL)
			if td[4] == "bid" {
				side = BUY
			}
			okWs.tradeCallback(&Trade{Pair: pair, Tid: ToInt64(td[0]), Price: ToFloat64(td[1]),
				Amount: ToFloat64(td[2]), Type: side, Date: okWs.adaptTime(td[3])}, contract)
		}

		return nil
	}

	return errors.New("unknown channel for " + resp[0].Channel)
}

func (okWs *OKExFutureWs) parseChannel(channel string) (pair CurrencyPair, contract string, ch string) {
	metas := strings.Split(channel, "_")
	pair = NewCurrencyPair2(strings.ToUpper(metas[3] + "_USD"))
	contract = metas[5]
	ch = metas[4]
	return pair, contract, ch
}

func (okWs *OKExFutureWs) adaptTime(tm string) int64 {
	format := "2006-01-02 15:04:05"
	day := time.Now().Format("2006-01-02")
	local, _ := time.LoadLocation("Asia/Chongqing")
	t, _ := time.ParseInLocation(format, day+" "+tm, local)
	return t.UnixNano() / 1e6

}
