package huobi

import "net/http"

type HuobiPro struct {
	*HuoBi_V2
}

func NewHuobiPro(client *http.Client, apikey, secretkey, accountId string) *HuobiPro {
	hbv2 := new(HuoBi_V2)
	hbv2.accountId = accountId
	hbv2.accessKey = apikey
	hbv2.secretKey = secretkey
	hbv2.httpClient = client
	hbv2.baseUrl = "https://api.huobi.pro"
	return &HuobiPro{hbv2}
}

func (hbpro *HuobiPro) GetExchangeName() string {
	return "huobi.pro"
}
