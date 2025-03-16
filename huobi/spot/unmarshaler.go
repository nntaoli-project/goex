package spot

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/spf13/cast"
)

func UnmarshalResponse(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}

func UnmarshalDepth(data []byte) (*Depth, error) {
	//TODO implement me
	panic("implement me")
}

func UnmarshalTicker(data []byte) (*Ticker, error) {
	var (
		tk   = new(Ticker)
		open float64
	)

	tk.Timestamp, _ = jsonparser.GetInt(data, "ts")
	tickData, _, _, _ := jsonparser.Get(data, "tick")
	err := jsonparser.ObjectEach(tickData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch string(key) {
		case "close":
			tk.Last = cast.ToFloat64(string(value))
		case "high":
			tk.High = cast.ToFloat64(string(value))
		case "low":
			tk.Low = cast.ToFloat64(string(value))
		case "vol":
			tk.Vol = cast.ToFloat64(string(value))
		case "open":
			open = cast.ToFloat64(string(value))
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

	if err != nil {
		return nil, err
	}

	tk.Percent = (tk.Last - open) / open * 100
	return tk, nil
}
