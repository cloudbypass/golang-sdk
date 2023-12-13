<p align="center">
  <a href="https://cloudbypass.com/" target="_blank" rel="noopener noreferrer" >
    <div align="center">
        <img src="https://github.com/cloudbypass/example/blob/main/assets/img.png?raw=true" alt="Cloudbypass" height="50">
    </div>
  </a>
</p>

## Cloudbypass SDK for Go

### 开始使用

> 继承 [go-resty/resty#supported-go-versions](https://github.com/go-resty/resty#supported-go-versions) v2支持的Go版本


[![GoDoc](https://godoc.org/github.com/cloudbypass/golang-sdk?status.svg)](https://godoc.org/github.com/cloudbypass/golang-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudbypass/golang-sdk)](https://goreportcard.com/report/github.com/cloudbypass/golang-sdk)

在`go-resty/resty.v2`基础上封装的穿云SDK。

### 安装

```bash
# Go Modules
require github.com/cloudbypass/golang-sdk V0.0.1
```

### 用法

```go
import "github.com/cloudbypass/golang-sdk"
```

### 发起请求

使用 `cloudbypass.New()` 创建一个新的 `resty.Client` 实例。

增加初始化参数`apikey`和`proxy`，分别用于设置穿云API服务密钥和代理IP。

定制用户可以通过设置`api_host`参数来指定服务地址。

> 以上参数可使用环境变量`CB_APIKEY`、`CB_PROXY`和`CB_APIHOST`进行配置。

```go
package main

import (
	"fmt"
	"github.com/cloudbypass/golang-sdk/cloudbypass"
)

func main() {
	client := cloudbypass.New(cloudbypass.BypassConfig{
		Apikey: "/* APIKEY */",
	})

	resp, err := client.R().
		EnableTrace().
		Get("https://opensea.io/category/memberships")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.StatusCode(), resp.Header().Get("X-Cb-Status"))
	fmt.Println(resp.String())
}
```

### 使用V2

穿云API V2适用于需要通过JS质询验证的网站。例如访问https://etherscan.io/accounts/label/lido，请求示例：

```go
package main

import (
	"fmt"
	"github.com/cloudbypass/golang-sdk/cloudbypass"
)

func main() {
	client := cloudbypass.New(cloudbypass.BypassConfig{
		Apikey: "/* APIKEY */",
		Part:   "0",
		Proxy:  "/* PROXY */",
	})

	resp, err := client.R().
		EnableTrace().
		Get("https://etherscan.io/accounts/label/lido")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.StatusCode(), resp.Header().Get("X-Cb-Status"))
	fmt.Println(resp.String())
}

```

### 查询余额

```go
package main

import (
	"fmt"
	"github.com/cloudbypass/golang-sdk/cloudbypass"
)

func main() {
	balance, err := cloudbypass.GetBalance( /* APIKEY */)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Balance:", balance)
}

```


### 关于重定向问题

使用SDK发起请求时，重定向操作会自动处理，无需手动处理。且重定向响应也会消耗积分。

### 关于服务密钥

请访问[穿云控制台](https://console.cloudbypass.com/#/api/account)获取服务密钥。