package futures

import (
	"encoding/json"
	"errors"
	"github.com/buger/jsonparser"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/spf13/cast"
)

func UnmarshalResponse(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}

func UnmarshalKline(data []byte) ([]Kline, error) {
	var lines []Kline
	klineData, _, _, err := jsonparser.Get(data, "data")
	if err != nil {
		return nil, err
	}
	_, err = jsonparser.ArrayEach(klineData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var kline Kline
		err = jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			switch string(key) {
			case "id":
				kline.Timestamp = cast.ToInt64(string(value))
			case "open":
				kline.Open = cast.ToFloat64(string(value))
			case "close":
				kline.Close = cast.ToFloat64(string(value))
			case "low":
				kline.Low = cast.ToFloat64(string(value))
			case "high":
				kline.High = cast.ToFloat64(string(value))
			case "vol":
				kline.Vol = cast.ToFloat64(string(value))
			}
			return nil
		})
		lines = append(lines, kline)
	})
	return lines, err
}

func UnmarshalTicker(data []byte) (*Ticker, error) {
	tk := &Ticker{}

	tkData, _, _, err := jsonparser.Get(data, "tick")
	if err != nil {
		return nil, err
	}

	jsonparser.ObjectEach(tkData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch string(key) {
		case "vol":
			tk.Vol = cast.ToFloat64(string(value))
		case "high":
			tk.High = cast.ToFloat64(string(value))
		case "low":
			tk.Low = cast.ToFloat64(string(value))
		case "close":
			tk.Last = cast.ToFloat64(string(value))
		case "ts":
			tk.Timestamp = cast.ToInt64(string(value))
		case "bid":
			var bids []float64
			err := UnmarshalResponse(value, &bids)
			if err != nil {
				return err
			}
			tk.Buy = bids[0]
		case "ask":
			var asks []float64
			err := UnmarshalResponse(value, &asks)
			if err != nil {
				return err
			}
			tk.Sell = asks[0]
		}
		return nil
	})

	return tk, nil
}

func UnmarshalCreateOrderResponse(data []byte) (*Order, error) {
	var order = new(Order)
	err := jsonparser.ObjectEach(data, func(key []byte, val []byte, dataType jsonparser.ValueType, offset int) error {
		switch string(key) {
		case "order_id_str":
			order.Id = string(val)
		case "client_order_id":
			order.CId = string(val)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func UnmarshalCancelOrderResponse(data []byte) error {
	val, _, _, _ := jsonparser.Get(data, "errors")
	if val != nil && len(val) > 0 {
		return errors.New(string(val))
	}
	return nil
}

func UnmarshalGetOrderInfoResponse(data []byte) (*Order, error) {
	var (
		order *Order
		err   error
	)
	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		order, err = unmarshalOrderResponse(value)
	})

	return order, err
}

func UnmarshalGetPendingOrdersResponse(data []byte) ([]Order, error) {
	var pendingOrders []Order

	ordersData, _, _, err := jsonparser.Get(data, "orders")
	if err != nil {
		return nil, err
	}

	_, err = jsonparser.ArrayEach(ordersData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ord, err := unmarshalOrderResponse(value)
		if err != nil {
			return
		}
		pendingOrders = append(pendingOrders, *ord)
	})

	return pendingOrders, err
}

func unmarshalOrderResponse(data []byte) (*Order, error) {
	var (
		order                  = new(Order)
		orderOffset, direction string
	)

	err := jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch string(key) {
		case "order_id_str":
			order.Id = string(value)
		case "client_order_id":
			order.CId = string(value)
		case "volume":
			order.Qty = cast.ToFloat64(string(value))
		case "price":
			order.Price = cast.ToFloat64(string(value))
		case "trade_volume":
			order.ExecutedQty = cast.ToFloat64(string(value))
		case "trade_avg_price":
			order.PriceAvg = cast.ToFloat64(string(value))
		case "fee":
			order.Fee = cast.ToFloat64(string(value))
		case "status":
			order.Status = AdaptStatus(cast.ToInt(string(value)))
		case "created_at", "create_date":
			order.CreatedAt = cast.ToInt64(string(value))
		case "canceled_at":
			order.CanceledAt = cast.ToInt64(string(value))
		case "direction":
			direction = string(value)
		case "offset":
			orderOffset = string(value)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	order.Side = AdaptOffsetDirectionToOrderSide(orderOffset, direction)

	return order, nil
}

func UnmarshalGetHistoryOrdersResponse(data []byte) ([]Order, error) {
	var (
		err    error
		orders []Order
	)
	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ord, err := unmarshalOrderResponse(value)
		if err != nil {
			err = err
			return
		}
		orders = append(orders, *ord)
	})
	return orders, err
}
