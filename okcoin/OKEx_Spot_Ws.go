package okcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"strings"
	"sync"
	"time"
)

type WsBaseResp struct {
	Channel string
	Data    json.RawMessage
}

type AddChannelData struct {
	Result  bool
	Channel string
}

type TickerData struct {
	Last      float64 `json:"last,string"`
	Sell      float64 `json:"sell,string"`
	Buy       float64 `json:"buy,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Vol       float64 `json:"vol,string"`
	Timestamp int64
}

type DepthData struct {
	Asks      [][]string `json:"asks""`
	Bids      [][]string `json:"bids"`
	Timestamp int64
}

type TradeData [][]string

type KlineData [][]string

type OKExSpotWs struct {
	*WsBuilder
	sync.Once
	wsConn *WsConn

	tickerCallback func(*Ticker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade)
	klineCallback  func(*Kline, int)
}

func NewOKExSpotWs() *OKExSpotWs {
	okWs := &OKExSpotWs{}

	okWs.WsBuilder = NewWsBuilder().
		WsUrl("wss://real.okex.com:10440/ws/v1").
		Heartbeat([]byte("{\"event\": \"ping\"} "), 30*time.Second).
		ReconnectIntervalTime(24 * time.Hour).
		UnCompressFunc(FlateUnCompress).
		ProtoHandleFunc(okWs.handle)

	return okWs
}

func (okWs *OKExSpotWs) SetCallbacks(tickerCallback func(*Ticker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade),
	klineCallback func(*Kline, int)) {
	okWs.tickerCallback = tickerCallback
	okWs.depthCallback = depthCallback
	okWs.tradeCallback = tradeCallback
	okWs.klineCallback = klineCallback
}

func (okWs *OKExSpotWs) subscribe(sub map[string]interface{}) error {
	okWs.connectWs()
	return okWs.wsConn.Subscribe(sub)
}

func (okWs *OKExSpotWs) SubscribeDepth(pair CurrencyPair, size int) error {
	if okWs.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	return okWs.subscribe(map[string]interface{}{
		"event":   "addChannel",
		"channel": fmt.Sprintf("ok_sub_spot_%s_depth_%d", strings.ToLower(pair.ToSymbol("_")), size)})
}

func (okWs *OKExSpotWs) SubscribeTicker(pair CurrencyPair) error {
	if okWs.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}

	return okWs.subscribe(map[string]interface{}{
		"event":   "addChannel",
		"channel": fmt.Sprintf("ok_sub_spot_%s_ticker", strings.ToLower(pair.ToSymbol("_")))})
}

func (okWs *OKExSpotWs) SubscribeTrade(pair CurrencyPair) error {
	if okWs.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}

	return okWs.subscribe(map[string]interface{}{
		"event":   "addChannel",
		"channel": fmt.Sprintf("ok_sub_spot_%s_deals", strings.ToLower(pair.ToSymbol("_")))})
}

func (okWs *OKExSpotWs) SubscribeKline(pair CurrencyPair, period int) error {
	if okWs.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	return okWs.subscribe(map[string]interface{}{
		"event": "addChannel",
		"channel": fmt.Sprintf("ok_sub_spot_%s_kline_%s",
			strings.ToLower(pair.ToSymbol("_")), AdaptKlinePeriodForOKEx(period))})
}

func (okWs *OKExSpotWs) connectWs() {
	okWs.Do(func() {
		okWs.wsConn = okWs.WsBuilder.Build()
		okWs.wsConn.ReceiveMessage()
	})
}

func (okWs *OKExSpotWs) handle(msg []byte) error {
	if string(msg) == "{\"event\":\"pong\"}" {
		okWs.wsConn.UpdateActiveTime()
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

	pair, err := okWs.getPairFormChannel(resp[0].Channel)
	if err != nil {
		log.Println(err, string(msg))
		return nil
	}

	if strings.Contains(resp[0].Channel, "depth") {
		var (
			depthData DepthData
			dep       Depth
		)

		err := json.Unmarshal(resp[0].Data, &depthData)
		if err != nil {
			return err
		}

		for _, ask := range depthData.Asks {
			dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(ask[0]), ToFloat64(ask[1])})
		}

		for _, bid := range depthData.Bids {
			dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(bid[0]), ToFloat64(bid[1])})
		}

		dep.Pair = pair
		dep.UTime = time.Unix(depthData.Timestamp/1000, 0)

		okWs.depthCallback(&dep)
		return nil
	}

	if strings.Contains(resp[0].Channel, "ticker") {
		var tickerData TickerData
		err := json.Unmarshal(resp[0].Data, &tickerData)
		if err != nil {
			return err
		}
		okWs.tickerCallback(&Ticker{
			Pair: pair,
			Last: tickerData.Last,
			Low:  tickerData.Low,
			High: tickerData.High,
			Sell: tickerData.Sell,
			Buy:  tickerData.Buy,
			Vol:  tickerData.Vol,
			Date: uint64(tickerData.Timestamp)})
		return nil
	}

	if strings.Contains(resp[0].Channel, "deals") {
		var (
			tradeData TradeData
		)

		err := json.Unmarshal(resp[0].Data, &tradeData)
		if err != nil {
			return err
		}

		for _, td := range tradeData {
			side := TradeSide(SELL)
			if td[4] == "bid" {
				side = BUY
			}
			okWs.tradeCallback(&Trade{Pair: pair, Tid: ToInt64(td[0]),
				Price: ToFloat64(td[1]), Amount: ToFloat64(td[2]), Type: side, Date: okWs.adaptTime(td[3])})
		}

		return nil
	}

	if strings.Contains(resp[0].Channel, "kline") {
		var k KlineData
		period := okWs.getKlinePeriodFormChannel(resp[0].Channel)
		err := json.Unmarshal(resp[0].Data, &k)
		if err != nil {
			return err
		}
		okWs.klineCallback(&Kline{
			Pair:      pair,
			Timestamp: ToInt64(k[0][0]),
			Open:      ToFloat64(k[0][1]),
			Close:     ToFloat64(k[0][4]),
			High:      ToFloat64(k[0][2]),
			Low:       ToFloat64(k[0][3]),
			Vol:       ToFloat64(k[0][5])}, period)
		return nil
	}

	return errors.New("unknown message " + resp[0].Channel)
}

func (okWs *OKExSpotWs) getPairFormChannel(channel string) (CurrencyPair, error) {
	metas := strings.Split(channel, "_")
	if len(metas) < 5 {
		return UNKNOWN_PAIR, errors.New("channel format error")
	}
	return NewCurrencyPair2(metas[3] + "_" + metas[4]), nil
}
func (okWs *OKExSpotWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, "_")
	if len(metas) < 7 {
		return 0
	}

	switch metas[6] {
	case "1hour":
		return KLINE_PERIOD_1H
	case "2hour":
		return KLINE_PERIOD_2H
	case "4hour":
		return KLINE_PERIOD_4H
	case "1min":
		return KLINE_PERIOD_1MIN
	case "5min":
		return KLINE_PERIOD_5MIN
	case "15min":
		return KLINE_PERIOD_15MIN
	case "30min":
		return KLINE_PERIOD_30MIN
	case "day":
		return KLINE_PERIOD_1DAY
	case "week":
		return KLINE_PERIOD_1WEEK
	default:
		return 0
	}
}

func (okWs *OKExSpotWs) adaptTime(tm string) int64 {
	format := "2006-01-02 15:04:05"
	day := time.Now().Format("2006-01-02")
	local, _ := time.LoadLocation("Asia/Chongqing")
	t, _ := time.ParseInLocation(format, day+" "+tm, local)
	return t.UnixNano() / 1e6

}
