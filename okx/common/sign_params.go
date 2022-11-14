package common

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2"
	"net/url"
	"strings"
)

func SignParams(httpMethod, apiUri, apiSecret, reqBody string) (signStr, timestamp string) {
	timestamp = goex.IsoTime()
	payload := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), apiUri, reqBody)
	signStr, _ = goex.HmacSHA256Base64Sign(apiSecret, payload)
	return
}

func DoAuthRequest(httpMethod, reqUrl string, params *url.Values, apiOpts goex.ApiOptions) ([]byte, error) {
	reqBody, _ := goex.ValuesToJson(*params)
	reqBodyStr := string(reqBody)
	signStr, timestamp := SignParams(httpMethod, "", "", string(reqBodyStr))
	headers := map[string]string{
		"Content-Type":         "application/json; charset=UTF-8",
		"Accept":               "application/json",
		"OK-ACCESS-KEY":        apiOpts.Key,
		"OK-ACCESS-PASSPHRASE": apiOpts.Passphrase,
		"OK-ACCESS-SIGN":       signStr,
		"OK-ACCESS-TIMESTAMP":  fmt.Sprint(timestamp)}
	return goex.GetHttpCli().DoRequest(httpMethod, reqUrl, reqBodyStr, headers)
}
