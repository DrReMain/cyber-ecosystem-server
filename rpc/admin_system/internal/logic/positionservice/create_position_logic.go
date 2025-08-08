package positionservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePositionLogic {
	return &CreatePositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePositionLogic) CreatePosition(in *admin_system.PositionBody) (*admin_system.BaseIDRes, error) {
	item, err := l.svcCtx.DB.Position.Create().
		SetNotNilSort(in.Sort).
		SetNotNilPositionName(in.PositionName).
		SetNotNilCode(in.Code).
		SetNotNilRemark(in.Remark).
		Save(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseIDRes{Id: item.ID, Msg: msgc.CREATE_SUCCESS}, nil
}
