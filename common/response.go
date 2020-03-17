package common

// 基础的返回
type Response struct{
	Status int `json:"status"`
	Msg string `json:"msg"`
}

// 成功信息
func Success(msg string) *Response {
	return &Response{Status:0,Msg:msg}
}

// 失败信息
func Error(msg string) *Response {
	return &Response{Status:1,Msg:msg}
}