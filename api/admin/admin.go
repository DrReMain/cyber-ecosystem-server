package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/config"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/handler"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/custom_validator"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/usual_err"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/middleware"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/chain"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/admin.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithCors(c.CORS.Address),
		rest.WithCorsHeaders(c.Project.AppNameHeader),
		rest.WithChain(chain.New(
			middleware.NewAppNameMiddleware(c.Project.AppNameHeader, c.Project.AppNameValue).Handle,
		)),
		rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
			httpx.Error(
				w,
				usual_err.HTTPUnauthorized(err.Error()),
			)
		}),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetValidator(custom_validator.New())
	httpx.SetErrorHandler(func(err error) (int, any) {
		switch e := err.(type) {
		case *errorc.HTTPError:
			logx.Errorw(
				"[HTTPError]",
				logx.Field("status", e.Status),
				logx.Field("detail", e.Detail),
			)
			return http.StatusOK, common_res.New(false, e.Code, e.Message)
		case *errorc.GRPCError:
			logx.Errorw(
				"[GRPCError]",
				logx.Field("status", e.Status),
				logx.Field("detail", e.Detail),
			)
			return http.StatusOK, common_res.NewGRPCRes(e.Code, e.Message)
		case *errorc.UnknownError:
			logx.Errorw(
				"[UnknownError]",
				logx.Field("detail", err.Error()),
			)
			return http.StatusOK, common_res.NewUnknownRes()
		default:
			logx.Errorw(
				"[SystemError]",
				logx.Field("detail", err.Error()),
			)
			return http.StatusInternalServerError, common_res.NewSystemRes()
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
