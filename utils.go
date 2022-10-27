package goex

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

//FloatToString 保留的小数点位数,去除末尾多余的0(StripTrailingZeros)
func FloatToString(v float64, n int) string {
	ret := strconv.FormatFloat(v, 'f', n, 64)
	return strconv.FormatFloat(cast.ToFloat64(ret), 'f', -1, 64) //StripTrailingZeros
}

func ValuesToJson(v url.Values) ([]byte, error) {
	paramMap := make(map[string]interface{})
	for k, vv := range v {
		if len(vv) == 1 {
			paramMap[k] = vv[0]
		} else {
			paramMap[k] = vv
		}
	}
	return json.Marshal(paramMap)
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

func GenerateOrderClientId(size int) string {
	uuidStr := strings.Replace(uuid.New().String(), "-", "", 32)
	return "goex-" + uuidStr[0:size-5]
}

func MergeOptionParams(params *url.Values, opts ...OptionParameter) {
	for _, opt := range opts {
		params.Set(opt.Key, opt.Value)
	}
}
