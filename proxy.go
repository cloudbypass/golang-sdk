package cloudbypass

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type ProxyIterator struct {
	proxy *CloudbypassProxy
	index int
}

type LoopIterator struct {
	proxy  *CloudbypassProxy
	length int
	index  int
	pool   []string
}

type CloudbypassProxy struct {
	username  string
	password  string
	region    string
	expire    int
	sessionId string
	gateway   string
}

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var (
	seededRand *rand.Rand = rand.New(
		rand.NewSource(rand.Int63()),
	)
	checkAuthCompile = regexp.MustCompile(`^(\w+-(res|dat)):(\w+)$`)
)

func checkAuth(auth string) (string, string, error) {
	// auth: ^\w+-(res|dat):\w+$
	if auth == "" {
		return "", "", fmt.Errorf("auth is empty")
	}
	// 正则判断auth格式并提取username, password
	contents := checkAuthCompile.FindStringSubmatch(auth)

	if len(contents) != 4 {
		return "", "", fmt.Errorf("auth format error")
	}

	return contents[1], contents[3], nil
}

func NewProxy(auth string) (*CloudbypassProxy, error) {
	username, password, err := checkAuth(auth)
	if err != nil {
		return nil, err
	}

	return &CloudbypassProxy{
		username: username,
		password: password,
		gateway:  "gw.cloudbypass.com:1288",
	}, nil
}

func (proxy *CloudbypassProxy) SetExpire(expire int) *CloudbypassProxy {
	proxy.expire = expire
	proxy.sessionId = ""
	return proxy
}

func (proxy *CloudbypassProxy) SetDynamic() *CloudbypassProxy {
	return proxy.SetExpire(0)
}

func (proxy *CloudbypassProxy) SetGateway(gateway string) *CloudbypassProxy {
	proxy.gateway = gateway
	proxy.sessionId = ""
	return proxy
}

func (proxy *CloudbypassProxy) SetRegion(region string) *CloudbypassProxy {
	proxy.region = region
	proxy.sessionId = ""
	return proxy
}

func (proxy *CloudbypassProxy) ClearRegion() *CloudbypassProxy {
	proxy.region = ""
	proxy.sessionId = ""
	return proxy
}

func (proxy *CloudbypassProxy) GetUsername() string {
	return proxy.username
}

func (proxy *CloudbypassProxy) GetPassword() string {
	return proxy.password
}

func (proxy *CloudbypassProxy) GetRegion() string {
	return proxy.region
}

func (proxy *CloudbypassProxy) GetExpire() int {
	return proxy.expire
}

func (proxy *CloudbypassProxy) GetSessionId() string {
	if proxy.sessionId == "" {
		b := make([]byte, 11)
		for i := range b {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
		proxy.sessionId = string(b)
	}
	return proxy.sessionId
}

func (proxy *CloudbypassProxy) parseOptions() string {
	options := []string{
		proxy.username,
	}
	if proxy.region != "" {
		options = append(options, strings.ReplaceAll(proxy.region, " ", "+"))
	}
	expire := proxy.expire
	if expire > 0 && expire < 5184000 {
		for _, unit := range [][2]interface{}{
			{60, "s"}, {60, "m"}, {24, "h"}, {999, "d"},
		} {
			if expire < unit[0].(int) || expire%unit[0].(int) != 0 {
				options = append(options, fmt.Sprintf("%s-%d%s", proxy.GetSessionId(), expire, unit[1]))
				break
			}
			expire = expire / unit[0].(int)
		}
	}

	return strings.Join(options, "_")
}

func (proxy *CloudbypassProxy) String() string {
	return fmt.Sprintf("%s:%s@%s", proxy.parseOptions(), proxy.password, proxy.gateway)
}

func (proxy *CloudbypassProxy) StringFormat(format string) string {
	if format == "" {
		return proxy.String()
	}

	format = strings.ReplaceAll(format, "username", proxy.parseOptions())
	format = strings.ReplaceAll(format, "password", proxy.password)
	format = strings.ReplaceAll(format, "gateway", proxy.gateway)

	return format
}

// Copy a new proxy
func (proxy *CloudbypassProxy) Copy() *CloudbypassProxy {
	return &CloudbypassProxy{
		username: proxy.username,
		password: proxy.password,
		region:   proxy.region,
		expire:   proxy.expire,
		gateway:  proxy.gateway,
	}
}

func (proxy *CloudbypassProxy) Iterate(count int) *ProxyIterator {
	return &ProxyIterator{
		proxy: proxy,
		index: count,
	}
}

func (iterator *ProxyIterator) HasNext() bool {
	return iterator.index > 0
}

func (iterator *ProxyIterator) Next() string {
	iterator.index--
	iterator.proxy.sessionId = ""
	return iterator.proxy.String()
}

func (proxy *CloudbypassProxy) Loop(count int) *LoopIterator {
	return &LoopIterator{
		proxy:  proxy,
		length: count,
		pool:   make([]string, count),
	}
}

func (iterator *LoopIterator) HasNext() bool {
	return true
}

func (iterator *LoopIterator) Next() string {
	if iterator.index == iterator.length {
		iterator.index = 0
	}

	if iterator.pool[iterator.index] == "" {
		iterator.proxy.sessionId = ""
		iterator.pool[iterator.index] = iterator.proxy.String()
	}

	iterator.index++
	return iterator.pool[iterator.index-1]
}
