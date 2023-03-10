package jsonrpc

import (
	"fmt"
)

type Request struct {
	Id      uint32      `json:"id"`
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func NewRequest(id uint32, method string, params interface{}) *Request {
	return &Request{
		Id:      id,
		Version: "2.0",
		Method:  method,
		Params:  params,
	}
}

type Response struct {
	Id        uint32         `json:"id"`
	Version   string         `json:"jsonrpc"`
	Result    interface{}    `json:"result,omitempty"`
	ErrorInfo *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (re *ResponseError) Error() string {
	return fmt.Sprintf("{\n\t\"code\": %d,\n\t\"message\": \"%s\"\n}", re.Code, re.Message)
}
