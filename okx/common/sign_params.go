package common

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2"
	"log"
	"strings"
)

func SignParams(httpMethod, apiUri, apiSecret, reqBody string) (signStr, timestamp string) {
	timestamp = goex.IsoTime()
	payload := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), apiUri, reqBody)
	log.Println("payload=", payload)
	signStr, _ = goex.HmacSHA256Base64Sign(apiSecret, payload)
	return
}
