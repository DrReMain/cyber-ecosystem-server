package base

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitLogic {
	return &InitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitLogic) Init(req *types.InitReq) (resp *types.CommonRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.BASE.InitDB(l.ctx, &admin_system.Empty{})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	if err := l.svcCtx.Casbin.LoadPolicy(); err != nil {
		return nil, errorc.NewHTTPInternal(msgc.SYSTEM_ERROR, err.Error())
	}

	return common_res.NewYES(data.Msg), nil
}
