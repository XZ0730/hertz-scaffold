package errno

var (
	// Success
	Success = NewErrNo(SuccessCode, "Success")

	ServiceError             = NewErrNo(ServiceErrorCode, "service is unable to start successfully")
	ServiceInternalError     = NewErrNo(ServiceErrorCode, "service internal error")
	ParamError               = NewErrNo(ParamErrorCode, "parameter error")
	AuthorizationFailedError = NewErrNo(AuthorizationFailedErrCode, "authorization failed")

	// User
	UserExistedError = NewErrNo(ParamErrorCode, "user existed")
	UserNameError    = NewErrNo(UserNameAuthErrorCode, "user name is not exist")
	PWDError         = NewErrNo(PwdErrorCode, "pwd not match")

	//
	TimeError   = NewErrNo(TimeErrorCode, "time set error")
	CreateError = NewErrNo(CreateErrorCode, "  create error")
	GetError    = NewErrNo(GetErrorCode, "  get error")
	DelError    = NewErrNo(DelErrorCode, "  del error")
	UpdateError = NewErrNo(UpdateErrorCode, "  update error")

	//
	NotExistError = NewErrNo(NotExistErrorCode, "record not exist")
	ExistError    = NewErrNo(ExistErrorCode, "record have existed")
)
