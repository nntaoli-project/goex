package spot

import (
	"encoding/json"
	"errors"
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

	if data[1] != '{' || data[len(data)-2] != '}' {
		logger.Warnf("[UnmarshalTicker] response data not json object ???,data=%s", string(data))
		return tk, errors.New("response data not json object")
	}

	var open float64
	err := jsonparser.ObjectEach(data[1:len(data)-1], func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch string(key) {
		case "last":
			tk.Last = cast.ToFloat64(string(value))
		case "askPx":
			tk.Sell = cast.ToFloat64(string(value))
		case "bidPx":
			tk.Buy = cast.ToFloat64(string(value))
		case "volCcy24h":
			tk.Vol = cast.ToFloat64(string(value))
		case "high24h":
			tk.High = cast.ToFloat64(string(value))
		case "low24h":
			tk.Low = cast.ToFloat64(string(value))
		case "ts":
			tk.Timestamp = cast.ToInt64(string(value))
		case "open24h":
			open = cast.ToFloat64(string(value))
		}
		return nil
	})

	if err != nil {
		logger.Errorf("[UnmarshalTicker] %s", err.Error())
		return nil, err
	}

	tk.Percent = (tk.Last - open) / open * 100

	return tk, nil

}

func (u *RespUnmarshaler) UnmarshalResponse(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
