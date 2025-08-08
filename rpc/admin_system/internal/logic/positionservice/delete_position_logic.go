package positionservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/position"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/user"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePositionLogic {
	return &DeletePositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePositionLogic) DeletePosition(in *admin_system.IDsReq) (*admin_system.BaseRes, error) {
	if exist, err := l.svcCtx.DB.User.Query().Where(user.HasPositionsWith(position.IDIn(in.Ids...))).Exist(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	} else if exist {
		l.Logger.Errorw("position used by a user", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.DELETE_FAILED)
	}

	if _, err := l.svcCtx.DB.Position.Delete().Where(position.IDIn(in.Ids...)).Exec(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.DELETE_SUCCESS}, nil
}
