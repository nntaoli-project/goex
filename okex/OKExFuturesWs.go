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

type OKExV3FuturesWs struct {
	base           *OKEx
	v3Ws           *OKExV3Ws
	tickerCallback func(*FutureTicker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade, string)
	klineCallback  func(*FutureKline, int)
	orderCallback  func(*FutureOrder, string)
}

func NewOKExV3FuturesWs(base *OKEx) *OKExV3FuturesWs {
	okV3Ws := &OKExV3FuturesWs{
		base: base,
	}
	okV3Ws.v3Ws = NewOKExV3Ws(base, okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3FuturesWs) TickerCallback(tickerCallback func(*FutureTicker)) *OKExV3FuturesWs {
	okV3Ws.tickerCallback = tickerCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FuturesWs) DepthCallback(depthCallback func(*Depth)) *OKExV3FuturesWs {
	okV3Ws.depthCallback = depthCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FuturesWs) TradeCallback(tradeCallback func(*Trade, string)) *OKExV3FuturesWs {
	okV3Ws.tradeCallback = tradeCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FuturesWs) OrderCallback(orderCallback func(*FutureOrder, string)) *OKExV3FuturesWs {
	okV3Ws.orderCallback = orderCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FuturesWs) KlineCallback(klineCallback func(*FutureKline, int)) *OKExV3FuturesWs {
	okV3Ws.klineCallback = klineCallback
	return okV3Ws
}

func (okV3Ws *OKExV3FuturesWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string),
	klineCallback func(*FutureKline, int),
	orderCallback func(*FutureOrder, string)) {
	okV3Ws.tickerCallback = tickerCallback
	okV3Ws.depthCallback = depthCallback
	okV3Ws.tradeCallback = tradeCallback
	okV3Ws.klineCallback = klineCallback
	okV3Ws.orderCallback = orderCallback
}

func (okV3Ws *OKExV3FuturesWs) getChannelName(currencyPair CurrencyPair, contractType string) string {
	var (
		prefix      string
		contractId  string
		channelName string
	)

	if contractType == SWAP_CONTRACT {
		prefix = "swap"
		contractId = fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-"))
	} else {
		prefix = "futures"
		contractId = okV3Ws.base.OKExFuture.GetFutureContractId(currencyPair, contractType)
	}

	channelName = prefix + "/%s:" + contractId

	return channelName
}

func (okV3Ws *OKExV3FuturesWs) SubscribeDepth(currencyPair CurrencyPair, contractType string, size int) error {
	if (size > 0) && (size != 5) {
		return errors.New("only support depth 5")
	}

	if okV3Ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "depth5")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeTicker(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	chName := okV3Ws.getChannelName(currencyPair, contractType)
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "ticker")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeTrade(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	chName := okV3Ws.getChannelName(currencyPair, contractType)
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "trade")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeKline(currencyPair CurrencyPair, contractType string, period int) error {
	if okV3Ws.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	seconds := adaptKLinePeriod(KlinePeriod(period))
	if seconds == -1 {
		return fmt.Errorf("unsupported kline period %d in okex", period)
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, fmt.Sprintf("candle%ds", seconds))}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeOrder(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.orderCallback == nil {
		return errors.New("place set order callback func")
	}
	okV3Ws.v3Ws.Login()
	chName := okV3Ws.getChannelName(currencyPair, contractType)
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "order")}})
}

func (okV3Ws *OKExV3FuturesWs) getContractAliasAndCurrencyPairFromInstrumentId(instrumentId string) (alias string, pair CurrencyPair) {
	if strings.HasSuffix(instrumentId, "SWAP") {
		ar := strings.Split(instrumentId, "-")
		return instrumentId, NewCurrencyPair2(fmt.Sprintf("%s_%s", ar[0], ar[1]))
	} else {
		contractInfo, err := okV3Ws.base.OKExFuture.GetContractInfo(instrumentId)
		if err != nil {
			logger.Error("instrument id invalid:", err)
			return "", UNKNOWN_PAIR
		}
		alias = contractInfo.Alias
		pair = NewCurrencyPair2(fmt.Sprintf("%s_%s", contractInfo.UnderlyingIndex, contractInfo.QuoteCurrency))
		return alias, pair
	}
}

func (okV3Ws *OKExV3FuturesWs) handle(ch string, data json.RawMessage) error {
	var (
		err           error
		tickers       []tickerResponse
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
		orderResp     []futureOrderResponse
		klineResponse []struct {
			Candle       []string `json:"candle"`
			InstrumentId string   `json:"instrument_id"`
		}
	)
	if strings.Contains(ch, "candle") {
		ch = "candle"
	}
	switch ch {
	case "ticker":
		err = json.Unmarshal(data, &tickers)
		if err != nil {
			return err
		}

		for _, t := range tickers {
			alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(t.InstrumentId)
			date, _ := time.Parse(time.RFC3339, t.Timestamp)
			okV3Ws.tickerCallback(&FutureTicker{
				Ticker: &Ticker{
					Pair: pair,
					Last: t.Last,
					Buy:  t.BestBid,
					Sell: t.BestAsk,
					High: t.High24h,
					Low:  t.Low24h,
					Vol:  t.Volume24h,
					Date: uint64(date.UnixNano() / int64(time.Millisecond)),
				},
				ContractType: alias,
			})
		}
		return nil
	case "candle":
		err = json.Unmarshal(data, &klineResponse)
		if err != nil {
			return err
		}

		for _, t := range klineResponse {
			_, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(t.InstrumentId)
			ts, _ := time.Parse(time.RFC3339, t.Candle[0])
			//granularity := adaptKLinePeriod(KlinePeriod(period))
			okV3Ws.klineCallback(&FutureKline{
				Kline: &Kline{
					Pair:      pair,
					High:      ToFloat64(t.Candle[2]),
					Low:       ToFloat64(t.Candle[3]),
					Timestamp: ts.Unix(),
					Open:      ToFloat64(t.Candle[1]),
					Close:     ToFloat64(t.Candle[4]),
					Vol:       ToFloat64(t.Candle[5]),
				},
				Vol2: ToFloat64(t.Candle[6]),
			}, 1)
		}
		return nil
	case "depth5":
		err := json.Unmarshal(data, &depthResp)
		if err != nil {
			logger.Error(err)
			return err
		}
		if len(depthResp) == 0 {
			return nil
		}
		alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(depthResp[0].InstrumentId)
		dep.Pair = pair
		dep.ContractType = alias
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
	case "trade":
		err := json.Unmarshal(data, &tradeResponse)
		if err != nil {
			logger.Error("unmarshal error :", err)
			return err
		}

		for _, resp := range tradeResponse {
			alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(resp.InstrumentId)

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
				Pair:   pair,
			}, alias)
		}
		return nil
	case "order":
		//2020/03/18 18:05:00 OKExFuturesWs.go:257: [D] [ws] [response] {"table":"futures/order","data":[{"leverage":"20","last_fill_time":"2020-03-18T10:05:00.790Z","filled_qty":"4","fee":"-0.00010655","price_avg":"112.62","type":"1","client_oid":"ce1661e5cb614fd690d0463de7a2eeb0","last_fill_qty":"4","instrument_id":"BSV-USD-200327","last_fill_px":"112.62","pnl":"0","size":"4","price":"112.73","last_fill_id":"15229749","error_code":"0","state":"2","contract_val":"10","order_id":"4573750935784449","order_type":"0","timestamp":"2020-03-18T10:05:00.790Z","status":"2"}]}
		err := json.Unmarshal(data, &orderResp)
		if err != nil {
			return err
		}
		for _, o := range orderResp {
			alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(o.InstrumentId)
			okV3Ws.orderCallback(&FutureOrder{
				ClientOid:    o.ClientOid,
				OrderID2:     o.OrderId,
				Price:        o.Price,
				Amount:       o.Size,
				AvgPrice:     o.PriceAvg,
				DealAmount:   o.FilledQty,
				Status:       okV3Ws.base.adaptOrderState(o.State),
				Currency:     pair,
				OrderType:    o.OrderType,
				OType:        o.Type,
				LeverRate:    o.Leverage,
				Fee:          o.Fee,
				ContractName: o.InstrumentId,
				OrderTime:    o.Timestamp.UnixNano() / int64(time.Millisecond),
			}, alias)
		}
		return nil
	}

	return fmt.Errorf("[%s] unknown websocket message: %s", ch, string(data))
}

func (okV3Ws *OKExV3FuturesWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}
