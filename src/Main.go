package main

import
(
    . "rest"
    "fmt"
    "utils"
    "rest/okcoin"
)

func main(){
    cfg, err := utils.NewConfig("D:\\config.json");
    if err != nil{
        return;
    }
    apiKey := cfg.Get("okcoin_cn.api_key", "null");
    secretKey := cfg.Get("okcoin_cn.secret_key", "null");
    fmt.Printf("apiKey:%s\nsecretKey:%s\n", apiKey, secretKey);
    
    var api API;
    api = okcoin.New(apiKey, secretKey);
    tk, err := api.GetTicker(BTC_CNY);
    if err != nil{
        fmt.Printf("%s", err.Error());
    }
    fmt.Printf("last:%f buy:%f sell:%f high:%f low:%f vol:%f date:%d",
        tk.Last, tk.Buy, tk.Sell, tk.High, tk.Low, tk.Vol, tk.Date);
    _, err = api.GetAccount();
    _, err = api.LimitBuy("0.01", "10", LTC_CNY);
    return;
}