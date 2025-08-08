package common_res

import (
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
)

func New(success bool, code, msg string) *types.CommonRes {
	t := fmt.Sprintf("%d", time.Now().UnixMilli())
	return &types.CommonRes{
		T:       t,
		Success: success,
		Code:    code,
		Msg:     msg,
	}
}

func NewSystemRes() *types.CommonRes {
	return New(false, SYSTEM_ERROR.Code, SYSTEM_ERROR.Msg)
}

func NewUnknownRes() *types.CommonRes {
	return New(false, UNKNOWN_ERROR.Code, UNKNOWN_ERROR.Msg)
}

func NewGRPCRes(code codes.Code, msg string) *types.CommonRes {
	return New(false, fmt.Sprintf("3%05d", code), msg)
}

func NewYES(msg string) *types.CommonRes {
	if msg == "" {
		msg = SUCCESS.Msg
	}
	return New(true, SUCCESS.Code, msg)
}

func NewNO(msg string) *types.CommonRes {
	if msg == "" {
		msg = FAIL.Msg
	}
	return New(false, FAIL.Code, msg)
}
