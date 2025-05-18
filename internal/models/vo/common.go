package vo

import (
	"cengkeHelperBackGo/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
func NewSuccessResp(msg string, data ...interface{}) RespData {
	var respData interface{}
	if len(data) > 0 {
		respData = data[0]
	}

	return RespData{
		Code: config.CodeSuccess,
		Msg:  msg,
		Data: respData,
	}
}

// --- Gin Handler 辅助函数 ---

// RespondSuccess 使用 NewSuccessResp 发送成功响应
func RespondSuccess(c *gin.Context, msg string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}
	c.JSON(http.StatusOK, NewSuccessResp(msg, responseData))
}

// NewCustomCodeResp 创建一个自定义错误码的 RespData
func NewCustomCodeResp(code int, msg string, data ...interface{}) RespData {
	var respData interface{}
	if len(data) > 0 {
		respData = data[0]
	}
	return RespData{
		Code: code,
		Msg:  msg,
		Data: respData,
	}
}

// RespondError 使用 NewBadResp 或 NewCustomCodeResp 发送错误响应
// statusCode: HTTP 状态码
// errCode: 业务错误码 (例如 config.CodeBadRequest)
// msg: 用户可读的错误信息
// err: (可选) 内部错误详情，用于日志记录
func RespondError(c *gin.Context, statusCode int, errCode int, msg string, err ...error) {
	// log.Println("Error occurred:", err) // 考虑记录实际的 error
	c.JSON(statusCode, NewCustomCodeResp(errCode, msg, nil)) // data 通常为 nil
}

// RespondBadRequest 快速发送一个 BadRequest 错误
func RespondBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, NewBadResp(msg))
}
