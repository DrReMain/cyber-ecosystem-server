package common_res

type BizCODE struct {
	Code string
	Msg  string
}

var (
	// "0xxxxx"
	EMPTY   = &BizCODE{Code: "000000", Msg: ""}
	SUCCESS = &BizCODE{Code: "000001", Msg: "SUCCESS"}
	FAIL    = &BizCODE{Code: "000002", Msg: "FAIL"}

	// "100000"
	SYSTEM_ERROR = &BizCODE{Code: "100000", Msg: "http service system error"}

	// "200000"
	UNKNOWN_ERROR = &BizCODE{Code: "200000", Msg: "http service unknown error"}

	// "3xxxxx" grpc error

	// "4xxxxx" http error
	UNAUTHORIZED = &BizCODE{Code: "400001", Msg: "UNAUTHORIZED"}
	REFRESH_FAIL = &BizCODE{Code: "400002", Msg: "REFRESH FAIL"}
	FORBIDDEN    = &BizCODE{Code: "400003", Msg: "FORBIDDEN"}
	BADREQUEST   = &BizCODE{Code: "400004", Msg: "BADREQUEST"}

	// "5xxxxx" custom http error
	USER_BANNED    = &BizCODE{Code: "500001", Msg: "USER BANNED"}
	PASSWORD_ERROR = &BizCODE{Code: "500002", Msg: "PASSWORD ERROR"}
)
