package common_res

import (
	"fmt"
	"time"

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

func NewYES(msg string) *types.CommonRes {
	if msg == "" {
		msg = "OK"
	}
	return New(true, "000000", msg)
}
