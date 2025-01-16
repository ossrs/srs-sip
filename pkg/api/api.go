package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ossrs/go-oryx-lib/logger"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

func (h *HttpApiServer) Start() {
	router := mux.NewRouter().StrictSlash(true)
	h.RegisterRoutes(router)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	go func() {
		ctx := context.Background()
		addr := fmt.Sprintf(":%v", h.conf.HttpApi.Port)
		logger.Tf(ctx, "http api listen on %s", addr)
		err := http.ListenAndServe(addr, handlers.CORS(headers, methods, origins)(router))
		if err != nil {
			panic(err)
		}
	}()
}
