package account

import (
	"context"
	"strings"
	"time"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/redisc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/jwt"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq, token string) (resp *types.LogoutRes, err error) {
	t := strings.TrimPrefix(token, "Bearer ")
	claims, err := jwt.Parse(t, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errorc.NewHTTPInternal(msgc.SYSTEM_ERROR, err.Error())
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errorc.NewHTTPInternal(msgc.SYSTEM_ERROR, "get exp from claims fail")
	}
	expTime := time.Unix(int64(exp), 0)
	now := time.Now()
	ttl := expTime.Sub(now)
	if ttl > 0 {
		err = l.svcCtx.Redis.Set(l.ctx, redisc.REDIS_TOKEN_PREFIX+t, redisc.REDIS_TOKEN_BANNED, ttl).Err()
		if err != nil {
			return nil, errorc.NewHTTPInternal(msgc.SYSTEM_ERROR, err.Error())
		}
	}

	return &types.LogoutRes{
		CommonRes: common_res.NewYES(""),
	}, nil
}
