package coinapi

//http request 工具函数
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"log"
)

func _httpRequest(client *http.Client, reqType string, reqUrl string, postData url.Values, requstHeaders map[string]string) ([]byte, error) {
	req, _ := http.NewRequest(reqType, reqUrl, strings.NewReader(postData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	if requstHeaders != nil {
		for k, v := range requstHeaders {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//var bodyDataMap map[string]interface{};
	//err = json.Unmarshal(bodyData, &bodyDataMap);
	//if err != nil {
	//	println(string(bodyData));
	//	return nil, err;
	//}

	return bodyData, nil
}

func HttpGet(client *http.Client, reqUrl string) (map[string]interface{}, error) {
	respData, err := _httpRequest(client, "GET", reqUrl, url.Values{}, nil)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	//fmt.Printf("\n%s\n", respData);
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Println(string(respData))
		return nil, err
	}
	return bodyDataMap, nil
}

func HttpPostForm(client *http.Client, reqUrl string, postData url.Values) ([]byte, error) {
	return _httpRequest(client, "POST", reqUrl, postData, nil)
}

func HttpPostForm2(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	return _httpRequest(client, "POST", reqUrl, postData, headers)
}
