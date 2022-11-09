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

func (u *RespUnmarshaler) UnmarshalGetDepthResponse(data []byte) (*Depth, error) {
	//TODO implement me
	panic("implement me")
}

func (u *RespUnmarshaler) UnmarshalGetTickerResponse(data []byte) (*Ticker, error) {
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

func (u *RespUnmarshaler) UnmarshalGetKlineResponse(data []byte) ([]Kline, error) {
	var (
		err    error
		klines []Kline
	)

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var (
			i = 0
			k Kline
		)
		_, err = jsonparser.ArrayEach(value, func(val []byte, dataType jsonparser.ValueType, offset int, err error) {
			switch i {
			case 0:
				k.Timestamp, _ = jsonparser.ParseInt(val)
			case 1:
				k.Open = cast.ToFloat64(string(val))
			case 2:
				k.High = cast.ToFloat64(string(val))
			case 3:
				k.Low = cast.ToFloat64(string(val))
			case 4:
				k.Close = cast.ToFloat64(string(val))
			case 5:
				//ignore
			case 6:
				k.Vol = cast.ToFloat64(string(val))
			}
			i += 1
		})
		k.Origin = value
		klines = append(klines, k)
	})

	return klines, err
}

func (u *RespUnmarshaler) UnmarshalResponse(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
