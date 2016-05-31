package rest

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "fmt"
)

func HttpGet(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url);
	if err != nil {
		return nil, err;
	}
	defer resp.Body.Close();
	body, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return nil, err;
	}
	var bodyDataMap map[string]interface{};
    fmt.Printf("\n%s\n", body);
	err = json.Unmarshal(body, &bodyDataMap);
	if err != nil {
		return nil, err;
	}
	return bodyDataMap, nil;
}