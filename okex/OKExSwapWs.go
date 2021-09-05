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

type OKExV3SwapWs struct {
	base             *OKEx
	v3Ws             *OKExV3Ws
	tickerCallback   func(*FutureTicker)
	depthCallback    func(*Depth)
	tradeCallback    func(*Trade, string)
	klineCallback    func(*FutureKline, int)
	orderCallback    func(order *FutureOrder)
	accountCallback  func(acct *FutureAccount)
	positionCallback func(pos *FuturePosition)
}

func NewOKExV3SwapWs(base *OKEx) *OKExV3SwapWs {
	okV3Ws := &OKExV3SwapWs{
		base: base,
	}
	okV3Ws.v3Ws = NewOKExV3Ws(base, okV3Ws.handle)
	return okV3Ws
}

func NewOKExV3SwapWsWithAuth(base *OKEx) *OKExV3SwapWs {
	okV3Ws := &OKExV3SwapWs{
		base: base,
	}
	okV3Ws.v3Ws = NewOKExV3WsWithAuth(base, okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3SwapWs) TickerCallback(tickerCallback func(*FutureTicker)) {
	okV3Ws.tickerCallback = tickerCallback
}

func (okV3Ws *OKExV3SwapWs) DepthCallback(depthCallback func(*Depth)) {
	okV3Ws.depthCallback = depthCallback
}

func (okV3Ws *OKExV3SwapWs) TradeCallback(tradeCallback func(*Trade, string)) {
	okV3Ws.tradeCallback = tradeCallback
}

func (okV3Ws *OKExV3SwapWs) KlineCallback(klineCallback func(*FutureKline, int)) {
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3SwapWs) OrderCallback(orderCallback func(*FutureOrder)) {
	okV3Ws.orderCallback = orderCallback
}

func (okV3Ws *OKExV3SwapWs) AccountCallback(acctCallback func(*FutureAccount)) {
	okV3Ws.accountCallback = acctCallback
}

func (okV3Ws *OKExV3SwapWs) PositionCallback(positionCallback func(position *FuturePosition)) {
	okV3Ws.positionCallback = positionCallback
}

func (okV3Ws *OKExV3SwapWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string),
	klineCallback func(*FutureKline, int)) {
	okV3Ws.tickerCallback = tickerCallback
	okV3Ws.depthCallback = depthCallback
	okV3Ws.tradeCallback = tradeCallback
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3SwapWs) getChannelName(currencyPair CurrencyPair, contractType string) string {
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
		//	logger.Info("contractid=", contractId)
	}

	if contractId == "" {
		return ""
	}

	channelName = prefix + "/%s:" + contractId

	return channelName
}

func (okV3Ws *OKExV3SwapWs) SubscribeDepth(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "depth5")}})
}

func (okV3Ws *OKExV3SwapWs) SubscribeTicker(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "ticker")}})
}

func (okV3Ws *OKExV3SwapWs) SubscribeTrade(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "trade")}})
}

func (okV3Ws *OKExV3SwapWs) SubscribeOrder(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.orderCallback == nil {
		return errors.New("please set order callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "order")}})
}

//isUSDT: 是否币本位, true: 不是, false: 是
func (okV3Ws *OKExV3SwapWs) SubscribeAccount(currency Currency, isUSDT bool) error {
	if okV3Ws.accountCallback == nil {
		return errors.New("please set account callback func")
	}

	chName := "swap/account:" + currency.String()
	if isUSDT {
		chName += "-USDT"
	}
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3SwapWs) SubscribePosition(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.positionCallback == nil {
		return errors.New("please set position callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "position")}})
}

func (okV3Ws *OKExV3SwapWs) SubscribeKline(currencyPair CurrencyPair, contractType string, period int) error {
	if okV3Ws.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	seconds := adaptKLinePeriod(KlinePeriod(period))
	if seconds == -1 {
		return fmt.Errorf("unsupported kline period %d in okex", period)
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, fmt.Sprintf("candle%ds", seconds))}})
}

func (okV3Ws *OKExV3SwapWs) getContractAliasAndCurrencyPairFromInstrumentId(instrumentId string) (alias string, pair CurrencyPair) {
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

func (okV3Ws *OKExV3SwapWs) handle(channel string, data json.RawMessage) error {
	var (
		err       error
		ch        string
		tickers   []tickerResponse
		depthResp []depthResponse
		dep       Depth
		orders    []futureOrderResponse
		holding   []struct {
			InstrumentId         string    `json:"instrument_id"`
			LongQty              float64   `json:"long_qty,string"` //多
			LongAvailQty         float64   `json:"long_avail_qty,string"`
			LongAvgCost          float64   `json:"long_avg_cost,string"`
			LongSettlementPrice  float64   `json:"long_settlement_price,string"`
			LongMargin           float64   `json:"long_margin,string"`
			LongPnl              float64   `json:"long_pnl,string"`
			LongPnlRatio         float64   `json:"long_pnl_ratio,string"`
			LongUnrealisedPnl    float64   `json:"long_unrealised_pnl,string"`
			RealisedPnl          float64   `json:"realised_pnl,string"`
			Leverage             float64   `json:"leverage,string"`
			ShortQty             float64   `json:"short_qty,string"`
			ShortAvailQty        float64   `json:"short_avail_qty,string"`
			ShortAvgCost         float64   `json:"short_avg_cost,string"`
			ShortSettlementPrice float64   `json:"short_settlement_price,string"`
			ShortMargin          float64   `json:"short_margin,string"`
			ShortPnl             float64   `json:"short_pnl,string"`
			ShortPnlRatio        float64   `json:"short_pnl_ratio,string"`
			ShortUnrealisedPnl   float64   `json:"short_unrealised_pnl,string"`
			LiquidationPrice     float64   `json:"liquidation_price,string"`
			CreatedAt            time.Time `json:"created_at,string"`
			UpdatedAt            time.Time `json:"updated_at"`
		}
		tradeResponse []struct {
			Side         string  `json:"side"`
			TradeId      int64   `json:"trade_id,string"`
			Price        float64 `json:"price,string"`
			Qty          float64 `json:"qty,string"`
			InstrumentId string  `json:"instrument_id"`
			Timestamp    string  `json:"timestamp"`
		}
		klineResponse []struct {
			Candle       []string `json:"candle"`
			InstrumentId string   `json:"instrument_id"`
		}
	)

	if strings.Contains(channel, "futures/candle") ||
		strings.Contains(channel, "swap/candle") {
		ch = "candle"
	} else {
		ch, err = okV3Ws.v3Ws.parseChannel(channel)
		if err != nil {
			logger.Errorf("[%s] parse channel err=%s ,  originChannel=%s", okV3Ws.base.GetExchangeName(), err, ch)
			return nil
		}
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
				ContractId:   t.InstrumentId,
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
		dep.ContractId = depthResp[0].InstrumentId
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
		err := json.Unmarshal(data, &orders)
		if err != nil {
			logger.Error("unmarshal error :", err)
			return err
		}
		for _, r := range orders {
			order := okV3Ws.base.OKExFuture.adaptOrder(r)
			okV3Ws.orderCallback(&order)
		}
		return nil
	case "account":
		var resp []map[string]CrossedAccountInfo
		err := json.Unmarshal(data, &resp)
		if err != nil {
			logger.Error("unmarshal error :", err)
			return err
		}
		fa := FutureAccount{
			FutureSubAccounts: make(map[Currency]FutureSubAccount, 1),
		}
		for _, m := range resp {
			for k, v := range m {
				a := FutureSubAccount{
					Currency:      NewCurrency(k, ""),
					AccountRights: v.Equity,
					ProfitReal:    v.RealizedPnl,
					ProfitUnreal:  v.UnrealizedPnl,
					KeepDeposit:   v.MarginFrozen,
					RiskRate:      v.MarginRatio,
				}
				fa.FutureSubAccounts[a.Currency] = a

			}
		}
		okV3Ws.accountCallback(&fa)
		return nil
	case "position":
		err := json.Unmarshal(data, &holding)
		if err != nil {
			logger.Error("unmarshal error :", err)
			return err
		}
		for _, pos := range holding {
			p := FuturePosition{
				ContractId:     ToInt64(pos.InstrumentId[8:]),
				BuyAmount:      pos.LongQty,
				BuyAvailable:   pos.LongAvailQty,
				BuyPriceAvg:    pos.LongAvgCost,
				BuyPriceCost:   pos.LongAvgCost,
				BuyProfitReal:  pos.LongPnl,
				SellAmount:     pos.ShortQty,
				SellAvailable:  pos.ShortAvailQty,
				SellPriceAvg:   pos.ShortAvgCost,
				SellPriceCost:  pos.ShortAvgCost,
				SellProfitReal: pos.ShortPnl,
				ForceLiquPrice: pos.LiquidationPrice,
				LeverRate:      pos.Leverage,
				CreateDate:     pos.CreatedAt.Unix(),
				ShortPnlRatio:  pos.ShortPnlRatio,
				LongPnlRatio:   pos.LongPnlRatio,
			}
			okV3Ws.positionCallback(&p)
		}
	}

	return fmt.Errorf("[%s] unknown websocket message: %s", ch, string(data))
}

func (okV3Ws *OKExV3SwapWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}
