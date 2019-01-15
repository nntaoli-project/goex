package okex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nntaoli-project/GoEx"
	"internal/log"
	"strings"
	"time"
)

/*
 Get a iso time
  eg: 2018-03-16T18:02:48.284Z
*/
func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

/*
 Get a http request body is a json string and a byte array.
*/
func BuildRequestBody(params interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	data, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
		return "", nil, errors.New("json convert string error")
	}

	jsonBody := string(data)
	binBody := bytes.NewReader(data)

	return jsonBody, binBody, nil
}

func doParamSign(httpMethod, apiSecret, uri, requestBody string) (string, string) {
	timestamp := IsoTime()
	preText := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), uri, requestBody)
	log.Println("preHash", preText)
	sign, _ := goex.GetParamHmacSHA256Base64Sign(apiSecret, preText)
	return sign, timestamp
}
