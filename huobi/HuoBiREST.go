package huobi

import (
    . "github.com/nntaoli/crypto_coin_api"
    "errors"
    "fmt"
    "log"
)

const
(
    REST_BASE_URL = "https://be.huobi.com/";
    REST_MARKET_URL = REST_BASE_URL + "market/"
    REST_TRADE_URL = REST_BASE_URL + "v1/"
    REST_DEPTH_URL = REST_MARKET_URL + "depth"
)

func currencyPair2String(currency CurrencyPair) string {
    switch currency {
    case ETH_CNY:
        return "ethcny"
    default:
        return ""
    }
}

func (hb *HuoBi) GetDepthREST(depthType string, currency CurrencyPair) (*Depth, error) {
    url := REST_DEPTH_URL + "?symbol=" + currencyPair2String(currency) + "&type=" + depthType

    bodyDataMap, err := HttpGet(hb.httpClient , url);

    if err != nil{
        return nil, err;
    }

    if bodyDataMap["status"] != nil && bodyDataMap["status"] != "ok"{
        log.Println(bodyDataMap);
        return nil, errors.New(fmt.Sprintf("%s" ,bodyDataMap));
    }

    var depth Depth;

    asks, isOK := bodyDataMap["tick"].(map[string]interface{})["asks"].([]interface{})
    if !isOK {
        return nil, errors.New("asks assert error")
    }

    for _, v := range asks {
        var dr DepthRecord;
        for i, vv := range v.([]interface{}) {
            switch i {
            case 0:
                dr.Price = vv.(float64);
            case 1:
                dr.Amount = vv.(float64);
            }
        }
        depth.AskList = append(depth.AskList, dr);
    }

    for _, v := range bodyDataMap["tick"].(map[string]interface{})["bids"].([]interface{}) {
        var dr DepthRecord;
        for i, vv := range v.([]interface{}) {
            switch i {
            case 0:
                dr.Price = vv.(float64);
            case 1:
                dr.Amount = vv.(float64);
            }
        }
        depth.BidList = append(depth.BidList, dr);
    }

    return &depth, nil;
}