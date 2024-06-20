package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/ossrs/srs-sip/pkg/service"
)

type HttpServer struct {
	conf   *config.MainConfig
	sipSvr *service.Service
}

func NewHttpServer(r0 interface{}, svr *service.Service) (*HttpServer, error) {
	return &HttpServer{
		conf:   r0.(*config.MainConfig),
		sipSvr: svr,
	}, nil
}

func (h *HttpServer) Start() {
	router := mux.NewRouter().StrictSlash(true)
	h.RegisterRoutes(router)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	go func() {
		err := http.ListenAndServe(":2020", handlers.CORS(headers, methods, origins)(router))
		if err != nil {
			panic(err)
		}
	}()
}
