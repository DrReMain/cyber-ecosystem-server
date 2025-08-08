package baseservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/redisc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/bsm/redislock"
	"github.com/zeromicro/go-zero/core/logx"
)

type InitDBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitDBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDBLogic {
	return &InitDBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitDBLogic) InitDB(in *admin_system.Empty) (*admin_system.BaseRes, error) {
	// If your db speed is high, comment the code below.
	//l.ctx = context.Background()

	locker := redislock.New(l.svcCtx.Redis)
	lock, err := locker.Obtain(l.ctx, redisc.INIT_DB_LOCK, 10*time.Minute, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		logx.Error("last initialization is running")
		return nil, errorc.GRPCInternalError(msgc.DB_INIT_ERROR)
	} else if err != nil {
		logx.Errorw("redis error", logx.Field("detail", err.Error()))
		return nil, errorc.GRPCInternalError(msgc.DB_INIT_ERROR)
	}
	defer func() {
		_ = lock.Release(l.ctx)
	}()

	if err := l.svcCtx.DB.Schema.Create(
		l.ctx,
		schema.WithForeignKeys(false),
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	); err != nil {
		return handleError(l, err)
	}

	if err = l.initBaseData(); err != nil {
		return handleError(l, err)
	}

	if err = l.initCasbin(); err != nil {
		return handleError(l, err)
	}

	_ = l.svcCtx.Redis.Set(l.ctx, redisc.INIT_DB_STATE, redisc.INIT_DB_STATE_YES, 24*time.Hour)
	return &admin_system.BaseRes{Msg: msgc.SUCCESS}, nil
}

func handleError(l *InitDBLogic, err error) (*admin_system.BaseRes, error) {
	logx.Errorw("database init error", logx.Field("detail", err.Error()))
	_ = l.svcCtx.Redis.Set(l.ctx, redisc.INIT_DB_ERROR, err.Error(), 300*time.Second)
	return nil, errorc.GRPCInternalError(msgc.DB_INIT_ERROR)
}
