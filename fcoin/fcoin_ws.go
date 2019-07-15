package fcoin

import (
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	. "github.com/nntaoli-project/GoEx"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	FCoinWSTicker        = "ticker.%s"
	FCoinWSOrderBook     = "depth.L%d.%s"
	FCoinWSOrderBookL20  = "depth.L20.%s"
	FCoinWSOrderBookL150 = "depth.L150.%s"
	FCoinWSOrderBookFull = "depth.full.%s"
	FCoinWSTrades        = "trade.%s"
	FCoinWSKLines        = "candle.%s.%s"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type FCoinWs struct {
	*WsBuilder
	sync.Once
	wsConn *WsConn

	tickerCallback func(*Ticker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade)
	klineCallback  func(*Kline, int)

	clientId      string
	subcribeTypes []string
	timeoffset    int64
	tradeSymbols  []TradeSymbol
}

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:   "M1",
	KLINE_PERIOD_3MIN:   "M3",
	KLINE_PERIOD_5MIN:   "M5",
	KLINE_PERIOD_15MIN:  "M15",
	KLINE_PERIOD_30MIN:  "M30",
	KLINE_PERIOD_60MIN:  "H1",
	KLINE_PERIOD_4H:     "H4",
	KLINE_PERIOD_6H:     "H6",
	KLINE_PERIOD_1DAY:   "D1",
	KLINE_PERIOD_1WEEK:  "W1",
	KLINE_PERIOD_1MONTH: "MN",
}
var _INERNAL_KLINE_PERIOD_REVERTER = map[string]int{
	"M1":  KLINE_PERIOD_1MIN,
	"M3":  KLINE_PERIOD_3MIN,
	"M5":  KLINE_PERIOD_5MIN,
	"M15": KLINE_PERIOD_15MIN,
	"M30": KLINE_PERIOD_30MIN,
	"H1":  KLINE_PERIOD_60MIN,
	"H4":  KLINE_PERIOD_4H,
	"H6":  KLINE_PERIOD_6H,
	"D1":  KLINE_PERIOD_1DAY,
	"W1":  KLINE_PERIOD_1WEEK,
	"MN":  KLINE_PERIOD_1MONTH,
}

func NewFCoinWs(client *http.Client) *FCoinWs {
	fcWs := &FCoinWs{}
	fcWs.clientId = getRandomString(8)
	fcWs.WsBuilder = NewWsBuilder().
		WsUrl("wss://api.fcoin.com/v2/ws").
		Heartbeat2(func() interface{} {
			ts := time.Now().Unix()*1000 + fcWs.timeoffset*1000
			args := make([]interface{}, 0)
			args = append(args, ts)
			return map[string]interface{}{
				"cmd":  "ping",
				"id":   fcWs.clientId,
				"args": args}

		}, 25*time.Second).
		ReconnectIntervalTime(24 * time.Hour).
		UnCompressFunc(FlateUnCompress).
		ProtoHandleFunc(fcWs.handle)
	fc := NewFCoin(client, "", "")
	fcWs.tradeSymbols = fc.tradeSymbols
	if len(fcWs.tradeSymbols) == 0 {
		panic("trade symbol is empty, pls check connection...")
	}
	return fcWs
}

//生成随机字符串
func getRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (fcWs *FCoinWs) SetCallbacks(
	tickerCallback func(*Ticker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade),
	klineCallback func(*Kline, int),
) {
	fcWs.tickerCallback = tickerCallback
	fcWs.depthCallback = depthCallback
	fcWs.tradeCallback = tradeCallback
	fcWs.klineCallback = klineCallback
}

func (fcWs *FCoinWs) subscribe(sub map[string]interface{}) error {
	fcWs.connectWs()
	return fcWs.wsConn.Subscribe(sub)
}

func (fcWs *FCoinWs) SubscribeDepth(pair CurrencyPair, size int) error {
	if fcWs.depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	arg := fmt.Sprintf(FCoinWSOrderBook, size, strings.ToLower(pair.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fcWs.subscribe(map[string]interface{}{
		"id":   fcWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fcWs *FCoinWs) SubscribeTicker(pair CurrencyPair) error {
	if fcWs.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	arg := fmt.Sprintf(FCoinWSTicker, strings.ToLower(pair.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fcWs.subscribe(map[string]interface{}{
		"id":   fcWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fcWs *FCoinWs) SubscribeTrade(pair CurrencyPair) error {
	if fcWs.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	arg := fmt.Sprintf(FCoinWSTrades, strings.ToLower(pair.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fcWs.subscribe(map[string]interface{}{
		"id":   fcWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fcWs *FCoinWs) SubscribeKline(pair CurrencyPair, period int) error {
	if fcWs.klineCallback == nil {
		return errors.New("place set kline callback func")
	}
	periodS, isOk := _INERNAL_KLINE_PERIOD_CONVERTER[period]
	if isOk != true {
		periodS = "M1"
	}

	arg := fmt.Sprintf(FCoinWSKLines, periodS, strings.ToLower(pair.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fcWs.subscribe(map[string]interface{}{
		"id":   fcWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fcWs *FCoinWs) connectWs() {
	fcWs.Do(func() {
		fcWs.wsConn = fcWs.WsBuilder.Build()
		fcWs.wsConn.ReceiveMessage()
	})
}

func (fcWs *FCoinWs) parseTickerData(tickmap []interface{}) *Ticker {
	t := new(Ticker)
	t.Date = uint64(time.Now().UnixNano() / 1000000)
	t.Last = ToFloat64(tickmap[0])
	t.Vol = ToFloat64(tickmap[9])
	t.Low = ToFloat64(tickmap[8])
	t.High = ToFloat64(tickmap[7])
	t.Buy = ToFloat64(tickmap[2])
	t.Sell = ToFloat64(tickmap[4])

	return t
}

func (fcWs *FCoinWs) parseDepthData(bids, asks []interface{}) *Depth {
	depth := new(Depth)
	n := 0
	for i := 0; i < len(bids); {
		depth.BidList = append(depth.BidList, DepthRecord{ToFloat64(bids[i]), ToFloat64(bids[i+1])})
		i += 2
		n++
	}

	n = 0
	for i := 0; i < len(asks); {
		depth.AskList = append(depth.AskList, DepthRecord{ToFloat64(asks[i]), ToFloat64(asks[i+1])})
		i += 2
		n++
	}

	return depth
}

func (fcWs *FCoinWs) parseKlineData(tickmap []interface{}) *Ticker {
	t := new(Ticker)
	t.Date = uint64(time.Now().UnixNano() / 1000000)
	t.Last = ToFloat64(tickmap[0])
	t.Vol = ToFloat64(tickmap[9])
	t.Low = ToFloat64(tickmap[8])
	t.High = ToFloat64(tickmap[7])
	t.Buy = ToFloat64(tickmap[2])
	t.Sell = ToFloat64(tickmap[4])

	return t
}

func (fcWs *FCoinWs) handle(msg []byte) error {
	//fmt.Println("ws msg:", string(msg))
	datamap := make(map[string]interface{})
	err := json.Unmarshal(msg, &datamap)
	if err != nil {
		fmt.Println("json unmarshal error for ", string(msg))
		return err
	}

	msgType, isOk := datamap["type"].(string)
	if isOk {
		resp := strings.Split(msgType, ".")
		switch resp[0] {
		case "hello", "ping":
			fcWs.wsConn.UpdateActiveTime()
			stime := int64(ToInt(datamap["ts"]))
			st := time.Unix(0, stime*1000*1000)
			lt := time.Now()
			offset := st.Sub(lt).Seconds()
			fcWs.timeoffset = int64(offset)
		case "ticker":
			tick := fcWs.parseTickerData(datamap["ticker"].([]interface{}))
			pair, err := fcWs.getPairFromType(resp[1])
			if err != nil {
				panic(err)
			}
			tick.Pair = pair
			fcWs.tickerCallback(tick)
			return nil
		case "depth":
			dep := fcWs.parseDepthData(datamap["bids"].([]interface{}), datamap["asks"].([]interface{}))
			stime := int64(ToInt(datamap["ts"]))
			dep.UTime = time.Unix(stime/1000, 0)
			pair, err := fcWs.getPairFromType(resp[2])
			if err != nil {
				panic(err)
			}
			dep.Pair = pair

			fcWs.depthCallback(dep)
			return nil
		case "candle":
			period := _INERNAL_KLINE_PERIOD_REVERTER[resp[1]]
			kline := &Kline{
				Timestamp: int64(ToInt(datamap["id"])),
				Open:      ToFloat64(datamap["open"]),
				Close:     ToFloat64(datamap["close"]),
				High:      ToFloat64(datamap["high"]),
				Low:       ToFloat64(datamap["low"]),
				Vol:       ToFloat64(datamap["quote_vol"]),
			}
			pair, err := fcWs.getPairFromType(resp[2])
			if err != nil {
				panic(err)
			}
			kline.Pair = pair
			fcWs.klineCallback(kline, period)
			return nil
		case "trade":
			side := BUY
			if datamap["side"] == "sell" {
				side = SELL
			}
			trade := &Trade{
				Tid:    int64(ToUint64(datamap["id"])),
				Type:   TradeSide(side),
				Amount: ToFloat64(datamap["amount"]),
				Price:  ToFloat64(datamap["price"]),
				Date:   int64(ToUint64(datamap["ts"])),
			}
			pair, err := fcWs.getPairFromType(resp[1])
			if err != nil {
				panic(err)
			}
			trade.Pair = pair
			fcWs.tradeCallback(trade)
			return nil
		default:
			return errors.New("unknown message " + msgType)

		}
	}
	return nil
}

func (fcWs *FCoinWs) getPairFromType(pair string) (CurrencyPair, error) {
	for _, v := range fcWs.tradeSymbols {
		if v.Name == pair {
			return NewCurrencyPair2(v.BaseCurrency + "_" + v.QuoteCurrency), nil
		}
	}
	return NewCurrencyPair2("" + "_" + ""), errors.New("pair not support :" + pair)
}
