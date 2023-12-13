package cloudbypass

import (
	"encoding/json"
	"fmt"
	resty "github.com/go-resty/resty/v2"
	"net/url"
	"strings"
)

const Version = "0.0.1"

type BypassConfig struct {
	Apikey  string
	Proxy   string
	ApiHost string
	Part    string
	Options []string
}

func New(config BypassConfig) *resty.Client {
	apikey := getEnv("CB_APIKEY", config.Apikey)
	Proxy := getEnv("CB_PROXY", config.Proxy)
	ApiHost := getEnv("CB_APIHOST", config.ApiHost)

	client := resty.New()
	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		fmt.Println("Before Request", r.Method, r.URL, r.Header)
		Url, _ := url.Parse(r.URL)
		r.SetHeader("X-Cb-Host", Url.Host)
		r.SetHeader("X-Cb-Apikey", apikey)
		if config.Proxy != "" {
			r.SetHeader("X-Cb-Proxy", Proxy)
		}
		if ApiHost != "" {
			Url.Host = ApiHost
		} else {
			Url.Host = "api.cloudbypass.com"
		}
		optionSet := make(map[string]bool)
		for _, option := range config.Options {
			optionSet[option] = true
		}
		optionSet["disable-redirect"] = true
		optionSet["full-cookie"] = true
		options := make([]string, 0)
		for option := range optionSet {
			options = append(options, option)
		}
		r.SetHeader("X-Cb-Options", strings.Join(options, ","))

		if config.Part != "" {
			r.SetHeader("X-Cb-Version", "2")
			r.SetHeader("X-Cb-Part", config.Part)
		}
		r.URL = Url.String()
		return nil
	})
	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		fmt.Println("After Response", r.StatusCode(), r.Header())
		if r.Header().Get("X-Cb-Status") != "ok" {
			// 解析响应体错误为BypassException
			var bypassException BypassException
			err := json.Unmarshal(r.Body(), &bypassException)
			if err != nil {
				return err
			}
			return bypassException
		}
		return nil
	})
	return client
}

type BypassInfo struct {
	Balance int `json:"balance"`
}

func GetBalance(apikey ...string) (int, error) {
	resp, err := resty.New().R().Get("https://console.cloudbypass.com/api/v1/balance?apikey=" + getEnv("CB_APIKEY", strings.Join(apikey, "")))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("status code %d", resp.StatusCode())
	}
	// 解析Json获取余额balance
	var bypassInfo BypassInfo
	err = json.Unmarshal(resp.Body(), &bypassInfo)
	if err != nil {
		return 0, err
	}
	return bypassInfo.Balance, nil
}
