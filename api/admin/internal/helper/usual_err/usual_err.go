package usual_err

import (
	"net/http"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
)

// token 校验失败
func HTTPUnauthorized(detail string) *errorc.HTTPError {
	return errorc.NewHTTPError(
		http.StatusUnauthorized,
		detail,
		common_res.UNAUTHORIZED.Code,
		common_res.UNAUTHORIZED.Msg,
	)
}

// token 刷新失败
func HTTPRefreshFail(detail string) *errorc.HTTPError {
	return errorc.NewHTTPError(
		http.StatusBadRequest,
		detail,
		common_res.REFRESH_FAIL.Code,
		common_res.REFRESH_FAIL.Msg,
	)
}

// 无权限
func HTTPForbidden(detail string) *errorc.HTTPError {
	return errorc.NewHTTPError(
		http.StatusForbidden,
		detail,
		common_res.FORBIDDEN.Code,
		common_res.FORBIDDEN.Msg,
	)
}

func HTTPBadRequest(detail string) *errorc.HTTPError {
	return errorc.NewHTTPError(
		http.StatusBadRequest,
		detail,
		common_res.BADREQUEST.Code,
		common_res.BADREQUEST.Msg,
	)
}

func HTTPBadRequestCustom(biz *common_res.BizCODE, detail string) *errorc.HTTPError {
	return errorc.NewHTTPError(
		http.StatusBadRequest,
		detail,
		biz.Code,
		biz.Msg,
	)
}
