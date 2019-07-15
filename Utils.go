package goex

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
)

func ToFloat64(v interface{}) float64 {
	if v == nil {
		return 0.0
	}

	switch v.(type) {
	case float64:
		return v.(float64)
	case string:
		vStr := v.(string)
		vF, _ := strconv.ParseFloat(vStr, 64)
		return vF
	default:
		panic("to float64 error.")
	}
}

func ToInt(v interface{}) int {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case string:
		vStr := v.(string)
		vInt, _ := strconv.Atoi(vStr)
		return vInt
	case int:
		return v.(int)
	case float64:
		vF := v.(float64)
		return int(vF)
	default:
		panic("to int error.")
	}
}

func ToUint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case int:
		return uint64(v.(int))
	case float64:
		return uint64((v.(float64)))
	case string:
		uV, _ := strconv.ParseUint(v.(string), 10, 64)
		return uV
	default:
		panic("to uint64 error.")
	}
}

func ToInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case float64:
		return int64(v.(float64))
	default:
		vv := fmt.Sprint(v)

		if vv == "" {
			return 0
		}

		vvv, err := strconv.ParseInt(vv, 0, 64)
		if err != nil {
			return 0
		}

		return vvv
	}
}

//n :保留的小数点位数
func FloatToString(v float64, n int) string {
	return fmt.Sprintf(fmt.Sprintf("%%.%df", n), v)
}

func ValuesToJson(v url.Values) ([]byte, error) {
	parammap := make(map[string]interface{})
	for k, vv := range v {
		if len(vv) == 1 {
			parammap[k] = vv[0]
		} else {
			parammap[k] = vv
		}
	}
	return json.Marshal(parammap)
}

func GzipUnCompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func FlateUnCompress(data []byte) ([]byte, error) {
	return ioutil.ReadAll(flate.NewReader(bytes.NewReader(data)))
}
