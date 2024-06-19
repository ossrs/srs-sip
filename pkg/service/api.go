package service

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (s *UAS) startHttpServer() {
	router := mux.NewRouter().StrictSlash(true)
	s.RegisterRoutes(router)

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
