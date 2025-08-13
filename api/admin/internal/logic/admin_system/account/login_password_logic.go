package account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/orm/ent/mixins"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/encrypt"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/jwt"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginPasswordLogic {
	return &LoginPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginPasswordLogic) LoginPassword(req *types.LoginPasswordReq) (resp *types.LoginPasswordRes, err error) {
	user, err := l.svcCtx.RPCAdminSystem.USER.GetUserByEmail(l.ctx, &admin_system.EmailReq{Email: *req.Email})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	if user.Status != nil && *user.Status != uint32(mixins.StatusNormal) {
		return nil, errorc.NewHTTPBadRequest(msgc.USER_BANNED, "user has been banned")
	}

	if !encrypt.EncryptCheck(*req.Password, *user.Password) {
		return nil, errorc.NewHTTPBadRequest(msgc.PASSWORD_ERROR, fmt.Sprintf("password check failed email:%s password:%s", *req.Email, *req.Password))
	}

	userID := ""
	if user.Id != nil {
		userID = *user.Id
	}

	departmentID := ""
	if user.Department != nil && user.Department.Id != nil {
		departmentID = *user.Department.Id
	}

	positionID := make([]string, 0)
	if user.Positions != nil && len(user.Positions) > 0 {
		for _, position := range user.Positions {
			positionID = append(positionID, *position.Id)
		}
	}

	roleCode := make([]string, 0)
	if user.Roles != nil && len(user.Roles) > 0 {
		for _, role := range user.Roles {
			roleCode = append(roleCode, *role.Code)
		}
	}

	now := time.Now()

	accessToken, err := jwt.New(
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
		now.Unix(),
		jwt.WithPayload("userID", userID),
		jwt.WithPayload("departmentID", departmentID),
		jwt.WithPayload("positionID", strings.Join(positionID, ",")),
		jwt.WithPayload("roleCode", strings.Join(roleCode, ",")),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.New(
		l.svcCtx.Config.Auth.RefreshSecret,
		l.svcCtx.Config.Auth.RefreshExpire,
		now.Unix(),
		jwt.WithPayload("userID", userID),
	)
	if err != nil {
		return nil, err
	}

	return &types.LoginPasswordRes{
		CommonRes: common_res.NewYES(""),
		Result: &types.Token{
			AccessToken:   pointc.P(accessToken),
			AccessExpire:  pointc.P(now.Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second).UnixMilli()),
			RefreshToken:  pointc.P(refreshToken),
			RefreshExpire: pointc.P(now.Add(time.Duration(l.svcCtx.Config.Auth.RefreshExpire) * time.Second).UnixMilli()),
		},
	}, nil
}
