package common

import (
	"encoding/json"
	"errors"
	"github.com/buger/jsonparser"
	"github.com/nntaoli-project/goex/v2/logger"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/spf13/cast"
	"time"
)

type RespUnmarshaler struct {
}

func (un *RespUnmarshaler) UnmarshalDepth(data []byte) (*Depth, error) {
	var (
		dep Depth
		err error
	)

	err = jsonparser.ObjectEach(data[1:len(data)-1],
		func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			switch string(key) {
			case "ts":
				dep.UTime = time.UnixMilli(cast.ToInt64(string(value)))
			case "asks":
				items, _ := un.unmarshalDepthItem(value)
				dep.Asks = items
			case "bids":
				items, _ := un.unmarshalDepthItem(value)
				dep.Bids = items
			}
			return nil
		})

	return &dep, err
}

func (un *RespUnmarshaler) unmarshalDepthItem(data []byte) (DepthItems, error) {
	var items DepthItems
	_, err := jsonparser.ArrayEach(data, func(asksItemData []byte, dataType jsonparser.ValueType, offset int, err error) {
		item := DepthItem{}
		i := 0
		_, err = jsonparser.ArrayEach(asksItemData, func(itemVal []byte, dataType jsonparser.ValueType, offset int, err error) {
			valStr := string(itemVal)
			switch i {
			case 0:
				item.Price = cast.ToFloat64(valStr)
			case 1:
				item.Amount = cast.ToFloat64(valStr)
			}
			i += 1
		})
		items = append(items, item)
	})
	return items, err
}

func (un *RespUnmarshaler) UnmarshalTicker(data []byte) (*Ticker, error) {
	var tk = &Ticker{}

	var open float64
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, val []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(val)
			switch string(key) {
			case "last":
				tk.Last = cast.ToFloat64(valStr)
			case "askPx":
				tk.Sell = cast.ToFloat64(valStr)
			case "bidPx":
				tk.Buy = cast.ToFloat64(valStr)
			case "vol24h":
				tk.Vol = cast.ToFloat64(valStr)
			case "high24h":
				tk.High = cast.ToFloat64(valStr)
			case "low24h":
				tk.Low = cast.ToFloat64(valStr)
			case "ts":
				tk.Timestamp = cast.ToInt64(valStr)
			case "open24h":
				open = cast.ToFloat64(valStr)
			}
			return nil
		})
	})

	if err != nil {
		logger.Errorf("[UnmarshalTicker] %s", err.Error())
		return nil, err
	}

	tk.Percent = (tk.Last - open) / open * 100

	return tk, nil
}

func (un *RespUnmarshaler) UnmarshalGetKlineResponse(data []byte) ([]Kline, error) {
	var klines []Kline
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var (
			k Kline
			i int
		)
		_, err = jsonparser.ArrayEach(value, func(val []byte, dataType jsonparser.ValueType, offset int, err error) {
			valStr := string(val)
			switch i {
			case 0:
				k.Timestamp = cast.ToInt64(valStr)
			case 1:
				k.Open = cast.ToFloat64(valStr)
			case 2:
				k.High = cast.ToFloat64(valStr)
			case 3:
				k.Low = cast.ToFloat64(valStr)
			case 4:
				k.Close = cast.ToFloat64(valStr)
			case 5:
				k.Vol = cast.ToFloat64(valStr)
			}
			i += 1
		})
		klines = append(klines, k)
	})

	return klines, err
}

func (un *RespUnmarshaler) UnmarshalCreateOrderResponse(data []byte) (*Order, error) {
	var ord = new(Order)
	err := jsonparser.ObjectEach(data[1:len(data)-1], func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valStr := string(value)
		switch string(key) {
		case "ordId":
			ord.Id = valStr
		case "clOrdId":
			ord.CId = valStr
		}
		return nil
	})
	return ord, err
}

func (un *RespUnmarshaler) UnmarshalGetPendingOrdersResponse(data []byte) ([]Order, error) {
	var (
		orders []Order
		err    error
	)

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ord, err := un.UnmarshalGetOrderInfoResponse(value)
		if err != nil {
			return
		}
		orders = append(orders, *ord)
	})

	return orders, err
}

func (un *RespUnmarshaler) UnmarshalGetOrderInfoResponse(data []byte) (ord *Order, err error) {
	var side, posSide string
	var utime int64
	ord = new(Order)

	err = jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valStr := string(value)
		switch string(key) {
		case "ordId":
			ord.Id = valStr
		case "px":
			ord.Price = cast.ToFloat64(valStr)
		case "sz":
			ord.Qty = cast.ToFloat64(valStr)
		case "cTime":
			ord.CreatedAt = cast.ToInt64(valStr)
		case "avgPx":
			ord.PriceAvg = cast.ToFloat64(valStr)
		case "accFillSz":
			ord.ExecutedQty = cast.ToFloat64(valStr)
		case "fee":
			ord.Fee = cast.ToFloat64(valStr)
		case "clOrdId":
			ord.CId = valStr
		case "side":
			side = valStr
		case "posSide":
			posSide = valStr
		case "ordType":
			ord.OrderTy = adaptSymToOrderTy(valStr)
		case "state":
			ord.Status = adaptSymToOrderStatus(valStr)
		case "uTime":
			utime = cast.ToInt64(valStr)
		}
		return nil
	})

	ord.Origin = data
	ord.Side = adaptSymToOrderSide(side, posSide)
	if ord.Status == OrderStatus_Canceled {
		ord.CanceledAt = utime
	}

	return
}

func (un *RespUnmarshaler) UnmarshalGetAccountResponse(data []byte) (map[string]Account, error) {
	var accMap = make(map[string]Account, 2)

	_, err := jsonparser.ArrayEach(data[1:len(data)-1], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var acc Account
		err = jsonparser.ObjectEach(value, func(key []byte, accData []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(accData)
			switch string(key) {
			case "ccy":
				acc.Coin = valStr
			case "availEq":
				acc.AvailableBalance = cast.ToFloat64(valStr)
			case "eq":
				acc.Balance = cast.ToFloat64(valStr)
			case "frozenBal":
				acc.FrozenBalance = cast.ToFloat64(valStr)
			}
			return err
		})

		if err != nil {
			return
		}

		accMap[acc.Coin] = acc
	}, "details")

	return accMap, err
}

func (un *RespUnmarshaler) UnmarshalGetFuturesAccountResponse(data []byte) (map[string]FuturesAccount, error) {
	var accMap = make(map[string]FuturesAccount, 2)

	_, err := jsonparser.ArrayEach(data[1:len(data)-1], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var acc FuturesAccount
		err = jsonparser.ObjectEach(value, func(key []byte, accData []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(accData)
			switch string(key) {
			case "ccy":
				acc.Coin = valStr
			case "availEq":
				acc.AvailEq = cast.ToFloat64(valStr)
			case "eq":
				acc.Eq = cast.ToFloat64(valStr)
			case "frozenBal":
				acc.FrozenBal = cast.ToFloat64(valStr)
			case "upl":
				acc.Upl = cast.ToFloat64(valStr)
			case "mgnRatio":
				acc.MgnRatio = cast.ToFloat64(valStr)
			}
			return err
		})

		if err != nil {
			return
		}

		accMap[acc.Coin] = acc
	}, "details")

	return accMap, err
}

func (un *RespUnmarshaler) UnmarshalCancelOrderResponse(data []byte) error {
	sCodeData, _, _, err := jsonparser.Get(data[1:len(data)-1], "sCode")
	if err != nil {
		return err
	}

	if cast.ToInt64(string(sCodeData)) == 0 {
		return nil
	}

	return errors.New(string(data))
}

func (un *RespUnmarshaler) UnmarshalGetPositionsResponse(data []byte) ([]FuturesPosition, error) {
	var (
		positions []FuturesPosition
		err       error
	)

	_, err = jsonparser.ArrayEach(data, func(posData []byte, dataType jsonparser.ValueType, offset int, err error) {
		var pos FuturesPosition
		err = jsonparser.ObjectEach(posData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(value)
			switch string(key) {
			case "availPos":
				pos.AvailQty = cast.ToFloat64(valStr)
			case "avgPx":
				pos.AvgPx = cast.ToFloat64(valStr)
			case "pos":
				pos.Qty = cast.ToFloat64(valStr)
			case "posSide":
				if valStr == "long" {
					pos.PosSide = Futures_OpenBuy
				}
				if valStr == "short" {
					pos.PosSide = Futures_OpenSell
				}
			case "upl":
				pos.Upl = cast.ToFloat64(valStr)
			case "uplRatio":
				pos.UplRatio = cast.ToFloat64(valStr)
			case "lever":
				pos.Lever = cast.ToFloat64(valStr)
			}
			return nil
		})
		positions = append(positions, pos)
	})

	return positions, err
}

func (un *RespUnmarshaler) UnmarshalResponse(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
