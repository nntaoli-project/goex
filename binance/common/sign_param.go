package common

import (
	"github.com/nntaoli-project/goex/v2"
	"net/url"
	"strconv"
	"time"
)

func SignParams(params *url.Values, secret string) {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
	params.Set("timestamp", timestamp)
	//params.Set("recvWindow", "60000")
	payload := params.Encode()
	sign, _ := goex.HmacSHA256Sign(secret, payload)
	params.Set("signature", sign)
}
