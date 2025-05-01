package vo

type RespData struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewBadResp(msg string) RespData {
	return RespData{
		Code: 400,
		Data: nil,
		Msg:  msg,
	}
}
