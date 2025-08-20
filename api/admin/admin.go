package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/config"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/handler"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/custom_validator"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/middleware"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"

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
				errorc.NewHTTPUnauthorized(err.Error()),
			)
		}),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetValidator(custom_validator.New())
	httpx.SetErrorHandler(func(err error) (int, any) {
		switch e := err.(type) {
		case *errorc.GRPCError:
			logx.Errorw("[GRPCError]", logx.Field("detail", e))
			return e.Status, common_res.New(false, fmt.Sprintf("2%05d", e.Code), e.Message)
		case *errorc.HTTPError:
			logx.Errorw("[HTTPError]", logx.Field("detail", e))
			return e.Status, common_res.New(false, fmt.Sprintf("1%05d", e.Code), e.Message)
		default:
			logx.Errorw("[SystemError]", logx.Field("detail", e))
			return http.StatusInternalServerError, common_res.New(false, "000500", msgc.SYSTEM_ERROR)
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
