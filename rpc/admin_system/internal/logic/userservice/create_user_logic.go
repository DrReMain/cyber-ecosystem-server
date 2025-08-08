package userservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/encrypt"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/user"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *admin_system.UserBody) (*admin_system.BaseIDRes, error) {
	if in.Password == nil {
		l.Logger.Errorw("password is nil", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.CREATE_FAILED)
	}
	if in.Email != nil {
		if exist, err := l.svcCtx.DB.User.Query().Where(user.EmailEQ(*in.Email)).Exist(l.ctx); err != nil {
			return nil, ent.DefaultHandleError(l.Logger, err, in)
		} else if exist {
			l.Logger.Errorw("email already exist", logx.Field("detail", in))
			return nil, errorc.GRPCInvalidArgumentError(msgc.CREATE_FAILED)
		}
	}

	create := l.svcCtx.DB.User.Create().
		SetNotNilPassword(pointc.P(encrypt.EncryptGenerate(*in.Password))).
		SetNotNilEmail(in.Email).
		SetNotNilName(in.Name).
		SetNotNilNickname(in.Nickname).
		SetNotNilPhone(in.Phone).
		SetNotNilAvatar(in.Avatar).
		SetNotNilRemark(in.Remark)

	if in.Department != nil {
		create = create.SetNotNilDepartmentID(in.Department.Id)
	}
	if in.Positions != nil && len(in.Positions) > 0 {
		ids := make([]string, len(in.Positions))
		for i, p := range in.Positions {
			ids[i] = *p.Id
		}
		create = create.AddPositionIDs(ids...)
	}
	if in.Roles != nil && len(in.Roles) > 0 {
		ids := make([]string, len(in.Roles))
		for i, p := range in.Roles {
			ids[i] = *p.Id
		}
		create = create.AddRoleIDs(ids...)
	}

	item, err := create.Save(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseIDRes{Id: item.ID, Msg: msgc.CREATE_SUCCESS}, nil
}
