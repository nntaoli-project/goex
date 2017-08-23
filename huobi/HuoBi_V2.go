package huobi

import "net/http"

type HuoBi_V2 struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func NewV2(httpClient *http.Client, accessKey, secretKey string) *HuoBi_V2 {
	return &HuoBi_V2{httpClient, accessKey, secretKey}
}


