package cloudbypass

import (
	"encoding/json"
	"fmt"
	resty "github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"strings"
)

const Version = "0.0.1"

type BypassConfig struct {
	Apikey  string
	Proxy   string
	ApiHost string
	UseV2   bool
	Part    string
	Options []string
}

func New(config BypassConfig) *resty.Client {
	apikey := getEnv("CB_APIKEY", config.Apikey)
	proxy := getEnv("CB_PROXY", config.Proxy)
	if config.ApiHost == "" {
		config.ApiHost = "https://api.cloudbypass.com"
	}
	ApiHost, err := url.Parse(getEnv("CB_APIHOST", config.ApiHost))

	if err != nil || ApiHost.Host == "" {
		panic(fmt.Sprintf("ApiHost [%s] is wrong or incomplete, requires a complete URL such as https://api.cloudbypass.com", ApiHost))
	}

	if ApiHost.Scheme == "" {
		ApiHost.Scheme = "https"
	}

	client := resty.New()
	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		Url, _ := url.Parse(r.URL)
		r.SetHeader("X-Cb-Host", Url.Host)
		r.SetHeader("X-Cb-Apikey", apikey)
		if config.Proxy == "" {
			r.SetHeader("X-Cb-Proxy", proxy)
		} else {
			r.SetHeader("X-Cb-Proxy", config.Proxy)
		}
		Url.Scheme = ApiHost.Scheme
		Url.Host = ApiHost.Host

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
		if config.UseV2 {
			r.SetHeader("X-Cb-Version", "2")
		}
		r.URL = Url.String()
		return nil
	})
	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.Header().Get("X-Cb-Status") != "ok" {
			var bypassException BypassException
			err := json.Unmarshal(r.Body(), &bypassException)
			if err != nil {
				return err
			}
			return bypassException
		}
		return nil
	})
	client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		req.Header.Set("X-Cb-Host", req.URL.Host)

		req.URL.Scheme = ApiHost.Scheme
		req.URL.Host = ApiHost.Host
		return nil
	}))
	return client
}

const balanceAPIURL = "https://console.cloudbypass.com/api/v1/balance"

// BalanceResult is the JSON response from the balance API.
// Total is present for traffic types (BalanceTypeRes / BalanceTypeDat).
type BalanceResult struct {
	Total   *float64 `json:"total,omitempty"`
	Balance float64  `json:"balance"`
}

// GetBalance calls POST /api/v1/balance with JSON body. typ is BalanceTypePoints (default if empty), BalanceTypeRes, or BalanceTypeDat.
func GetBalance(apikey string, email string, typ string) (*BalanceResult, error) {
	if typ == "" {
		typ = BalanceTypePoints
	}
	key := getEnv("CB_APIKEY", apikey)
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"apikey": key,
			"email":  email,
			"type":   typ,
		}).
		Post(balanceAPIURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}
	var out BalanceResult
	if err := json.Unmarshal(resp.Body(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
