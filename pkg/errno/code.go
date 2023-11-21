package errno

const (
	// For api-gateway
	StatusSuccessCode = 0
	StatusSuccessMsg  = "ok"

	// For microservices
	SuccessCode = 10000
	SuccessMsg  = "ok"

	// Error
	ServiceErrorCode           = 10001 // 未知微服务错误
	ParamErrorCode             = 10002 // 参数错误
	AuthorizationFailedErrCode = 10003 // 鉴权失败
	UnexpectedTypeErrorCode    = 10004 // 未知类型
	NotImplementErrorCode      = 10005 // 未实装
	SensitiveWordsErrorCode    = 10006 // 敏感词
	PwdErrorCode               = 10007 // 密码错误
	UserNameAuthErrorCode      = 10008 //用户名错误

	//
	TimeErrorCode     = 23001
	CreateErrorCode   = 23002
	GetErrorCode      = 23003
	DelErrorCode      = 23004
	UpdateErrorCode   = 23005
	NotExistErrorCode = 23006
	ExistErrorCode    = 23007
)
