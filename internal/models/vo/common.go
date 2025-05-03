package vo

import "cengkeHelperBackGo/internal/config"

type RespData struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewBadResp(msg string) RespData {
	return RespData{
		Code: config.CodeBadRequest,
		Data: nil,
		Msg:  msg,
	}
}

// NewSuccessResp 成功响应
func NewSuccessResp(data interface{}, msg string) RespData {
	return RespData{
		Code: config.CodeSuccess,
		Msg:  msg,
		Data: data,
	}
}
