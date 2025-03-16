package common

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2/util"
	"net/url"
	"time"
)

func SignParams(params *url.Values, secret string) {
	timestamp := time.Now().UnixMilli()
	params.Set("timestamp", fmt.Sprint(timestamp))
	params.Set("recvWindow", "6000")
	payload := params.Encode()
	sign, _ := util.HmacSHA256Sign(secret, payload)
	params.Set("signature", sign)
}
