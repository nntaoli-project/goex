package common

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2/util"
	"log"
	"strings"
	"time"
)

func SignParams(httpMethod, apiUri, apiSecret, reqBody string) (signStr, timestamp string) {
	timestamp = time.Now().Format("2006-01-02T15:04:05.000Z") //iso time style
	payload := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), apiUri, reqBody)
	log.Println("payload=", payload)
	signStr, _ = util.HmacSHA256Base64Sign(apiSecret, payload)
	return
}
