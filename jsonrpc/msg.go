package jsonrpc

type Msg struct {
	Id      uint32      `json:"id"`
	Version string      `json:"version"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func NewMsg(id uint32, method string, params interface{}) *Msg {
	return &Msg{
		Id:      id,
		Version: "2.0",
		Method:  method,
		Params:  params,
	}
}
