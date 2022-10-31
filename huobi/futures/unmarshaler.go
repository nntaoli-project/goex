package futures

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	. "github.com/nntaoli-project/goex/v2"
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
		kline.Origin = value
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
