package common

import (
	"fmt"
	"github.com/nntaoli-project/goex/v2/options"
	"github.com/nntaoli-project/goex/v2/util"
	"net/url"
	"time"
)

func DoSignParam(httpMethod, reqUrl string, apiOpt options.ApiOptions) *url.Values {
	///////////////////// 参数签名 ////////////////////////
	signParams := url.Values{}
	signParams.Set("AccessKeyId", apiOpt.Key)
	signParams.Set("SignatureMethod", "HmacSHA256")
	signParams.Set("SignatureVersion", "2")
	signParams.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))

	reqURL, _ := url.Parse(reqUrl)
	path := reqURL.RequestURI()
	domain := reqURL.Hostname()

	payload := fmt.Sprintf("%s\n%s\n%s\n%s", httpMethod, domain, path, signParams.Encode())
	sign, _ := util.HmacSHA256Base64Sign(apiOpt.Secret, payload)

	signParams.Set("Signature", sign)
	///////////////////签名结束////////////////////

	return &signParams
}
