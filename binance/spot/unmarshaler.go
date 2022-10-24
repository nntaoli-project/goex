package spot

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"github.com/spf13/cast"
)

type RespUnmarshaler struct {
}

func (u *RespUnmarshaler) UnmarshalDepth(data []byte) (*Depth, error) {
	//TODO implement me
	panic("implement me")
}

func (u *RespUnmarshaler) UnmarshalTicker(data []byte) (*Ticker, error) {
	var tk = &Ticker{}

	if data[0] != '{' || data[len(data)-1] != '}' {
		logger.Warnf("[UnmarshalTicker] response data not json object ???")
		return tk, nil
	}

	err := jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch string(key) {
		case "lastPrice":
			tk.Last = cast.ToFloat64(string(value))
		case "askPrice":
			tk.Sell = cast.ToFloat64(string(value))
		case "bidPrice":
			tk.Buy = cast.ToFloat64(string(value))
		case "volume":
			tk.Vol = cast.ToFloat64(string(value))
		case "highPrice":
			tk.High = cast.ToFloat64(string(value))
		case "lowPrice":
			tk.Low = cast.ToFloat64(string(value))
		case "closeTime":
			tk.Timestamp = cast.ToInt64(string(value))
		case "priceChangePercent":
			tk.Percent = cast.ToFloat64(string(value))
		}
		return nil
	})
	if err != nil {
		logger.Errorf("[UnmarshalTicker] %s", err.Error())
		return nil, err
	}

	return tk, nil

}

func (u *RespUnmarshaler) UnmarshalResponse(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
