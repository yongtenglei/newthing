package e

const (
	OK     = "OK"
	Failed = "Failed"

	UserAlreadyExists = "用户已存在"
	UserDoesNotFound  = "用户不存在"
	PasswordErr       = "密码错误"
	InternalBusy      = "内部错误"

	InvalidTokenErr = "无效TOKEN"
	ExpiredTokenErr = "过期的Token"
	TokenMethodErr  = "非HMAC方法"
)
