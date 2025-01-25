package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/ossrs/srs-sip/pkg/service"
)

type HttpApiServer struct {
	conf   *config.MainConfig
	sipSvr *service.Service
}

func NewHttpApiServer(r0 interface{}, svr *service.Service) (*HttpApiServer, error) {
	return &HttpApiServer{
		conf:   r0.(*config.MainConfig),
		sipSvr: svr,
	}, nil
}

func (h *HttpApiServer) Start(router *mux.Router) {
	// 添加版本检查路由到主路由器
	router.HandleFunc("/srs-sip", h.ApiGetAPIVersion).Methods(http.MethodGet)

	// 创建一个子路由，所有API都以/srs-sip/v1为前缀
	apiRouter := router.PathPrefix("/srs-sip/v1").Subrouter()

	logger.Tf(context.Background(), "Registering API routes under /srs-sip/v1")
	h.RegisterRoutes(apiRouter)

	// 打印所有注册的路由，包含更详细的信息
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		pathRegexp, _ := route.GetPathRegexp()
		methods, _ := route.GetMethods()
		queries, _ := route.GetQueriesTemplates()
		logger.Tf(context.Background(), "Route Details: Path=%v, Regexp=%v, Methods=%v, Queries=%v",
			pathTemplate, pathRegexp, methods, queries)
		return nil
	})
}
