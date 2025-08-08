package userservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/user"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *admin_system.IDReq) (*admin_system.UserBody, error) {
	item, err := l.svcCtx.DB.User.Query().
		Where(user.IDEQ(in.Id)).
		WithDepartment().
		WithPositions().
		WithRoles().
		First(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.UserBody{
		Id:         pointc.P(item.ID),
		CreatedAt:  pointc.P(item.CreatedAt.UnixMilli()),
		UpdatedAt:  pointc.P(item.UpdatedAt.UnixMilli()),
		Status:     pointc.PStatus8t32(&item.Status),
		Password:   pointc.P(item.Password),
		Email:      pointc.P(item.Email),
		Name:       pointc.P(item.Name),
		Nickname:   pointc.P(item.Nickname),
		Phone:      pointc.P(item.Phone),
		Avatar:     pointc.P(item.Avatar),
		Remark:     pointc.P(item.Remark),
		Department: buildBDepartment(item.Edges.Department),
		Positions:  buildBPositions(item.Edges.Positions),
		Roles:      buildBRoles(item.Edges.Roles),
	}, nil
}
