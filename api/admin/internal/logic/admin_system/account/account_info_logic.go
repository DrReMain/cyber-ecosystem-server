package account

import (
	"context"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccountInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountInfoLogic {
	return &AccountInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountInfoLogic) AccountInfo(req *types.AccountInfoReq) (resp *types.AccountInfoRes, err error) {
	userID, ok := l.ctx.Value("userID").(string)
	if !ok {
		return nil, errorc.NewHTTPInternal(msgc.SYSTEM_ERROR, "get userID from token fail")
	}
	roleCode, ok := l.ctx.Value("roleCode").(string)
	if !ok {
		return nil, errorc.NewHTTPInternal(msgc.SYSTEM_ERROR, "get roleCode from token fail")
	}
	role := strings.Split(roleCode, ",")

	user, err := l.svcCtx.RPCAdminSystem.USER.GetUser(l.ctx, &admin_system.IDReq{Id: userID})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	menus, err := l.svcCtx.RPCAdminSystem.MENU.QueryMenuByRoleCode(l.ctx, &admin_system.MenuListByRoleCodeReq{RoleCode: role})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.AccountInfoRes{
		CommonRes: common_res.NewYES(""),
		Result: &types.InfoBody{
			Email:    user.Email,
			Name:     user.Name,
			Nickname: user.Nickname,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Menus:    buildTMenus(menus.List),
		},
	}, nil
}

func buildTMenus(b []*admin_system.MenuBody) (result []*types.Menu) {
	result = make([]*types.Menu, len(b))
	for i, v := range b {
		result[i] = &types.Menu{
			CodePath: v.CodePath,
			Icon:     v.Icon,
			MenuType: v.MenuType,
			Children: buildTMenus(v.Children),
		}
	}
	return
}
