package positionservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePositionLogic {
	return &UpdatePositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePositionLogic) UpdatePosition(in *admin_system.PositionBody) (*admin_system.BaseRes, error) {
	if in.Id == nil || *in.Id == "" {
		l.Logger.Errorw("id is empty", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.UPDATE_FAILED)
	}

	if err := l.svcCtx.DB.Position.UpdateOneID(*in.Id).
		SetNotNilSort(in.Sort).
		SetNotNilPositionName(in.PositionName).
		SetNotNilCode(in.Code).
		SetNotNilRemark(in.Remark).
		Exec(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.UPDATE_SUCCESS}, nil
}
