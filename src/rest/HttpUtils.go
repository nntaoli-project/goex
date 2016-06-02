package rest
//http request 工具函数
import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strings"
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

func HttpPost(url string) (map[string]interface{}, error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", 
		strings.NewReader("name=cjb"));
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

func HttpPostForm(client *http.Client, reqUrl string, postData url.Values) (map[string]interface{}, error) {
	req, _ := http.NewRequest("POST", reqUrl, strings.NewReader(postData.Encode()));
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36");

	resp, err := client.Do(req);
	if err != nil {
		return nil, err;
	}

	defer resp.Body.Close();

	bodyData, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return nil, err;
	}

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		return nil, err;
	}

	return bodyDataMap, nil;
}