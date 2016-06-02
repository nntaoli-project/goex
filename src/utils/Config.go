package utils

import
(
    "os"
    "io/ioutil"
    "encoding/json"
    "strings"
)

/*
Json格式配置文件解析
{
    "okcoin_cn":
    {
        "api_key":"",
        "secret_key":""
        ...
    }
    "huobi":
    {
         "api_key":"",
        "secret_key":""
        ...       
    }
    ...
}
*/

type Config struct{
    FileName string
    ConfigMap map[string]interface{};
}

func NewConfig(fn string)(*Config, error){
    var cfgMap map[string]interface{};
    fi, err := os.Open(fn);
    if err != nil{
        return nil, err;
    }
    defer fi.Close();
    content, err := ioutil.ReadAll(fi);
    if err != nil{
        return nil, err;
    }    
    err = json.Unmarshal(content, &cfgMap);
	if err != nil {
		return nil, err;
	}
    return &Config{fn, cfgMap}, nil;
}

func (ctx * Config) Get(key, def string) string {
    var value interface{};
    value = nil;
    keys := strings.Split(key, ".");
    for _, it := range keys{
        if value == nil{
            value = ctx.ConfigMap[it];
        }else{
            var m map[string]interface{};
            m = value.(map[string]interface{});
            value = m[it];
        }
        if value == nil{
            return def;
        }
    }
    return value.(string);
}