package cloudbypass

import "fmt"

type APIError struct {
	error
	Message string `json:"message"`
}

// 实现 error 接口
type BypassException struct {
	error
	Id      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 实现 error 接口
func (e BypassException) Error() string {
	return fmt.Sprintf("BypassException %s: %s - %s", e.Id, e.Code, e.Message)
}
