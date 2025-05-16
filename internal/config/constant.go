package config

// 定义业务响应码常量
const (
	CodeSuccess      = 0     // 成功
	CodeBadRequest   = 40000 // 请求参数错误
	CodeUnauthorized = 40100 // 未授权
	CodeForbidden    = 40300 // 禁止访问
	CodeNotFound     = 40400 // 资源未找到
	CodeConflict     = 40900 // 资源冲突

	CodeInternalError      = 50000 // 服务器内部错误
	CodeServiceUnavailable = 50300 // 服务不可用
	CodeServerError        = 50001 // 服务器错误

	CodeUserExist        = 10001 // 用户已存在
	CodeUserNotFound     = 10002 // 用户不存在
	CodePasswordError    = 10003 // 密码错误
	CodeTokenExpired     = 10004 // Token过期
	CodeTokenInvalid     = 10005 // Token无效
	CodePermissionDenied = 10006 // 权限不足
	CodeUsernameExists   = 10007 // 用户名已存在
	CodeDatabaseError    = 10008 // 数据库错误	//
	CodeEmailExists      = 10009 // 邮箱已存在
)
