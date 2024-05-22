package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {

	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	apiV1Router.HandleFunc("/version", GetVersion).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/devices", ListDevices).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/devices/{id}/channels", GetChannels).Methods(http.MethodGet)

	apiV1Router.HandleFunc("/channels/invite", InviteChannel).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/channels/bye", ByeChannel).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/channels/ptz", PTZControl).Methods(http.MethodPost)

	apiV1Router.HandleFunc("", GetAPIRoutes(apiV1Router)).Methods(http.MethodGet)

	router.HandleFunc("/api", GetAPIVersion).Methods(http.MethodGet)
}

func GetAPIVersion(w http.ResponseWriter, r *http.Request) {
	versionInfo := map[string]string{
		"version": "v1",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionInfo)
}

func GetAPIRoutes(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var routes []map[string]string

		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, err := route.GetPathTemplate()
			if err != nil {
				return err
			}
			methods, err := route.GetMethods()
			if err != nil {
				return err
			}
			for _, method := range methods {
				routes = append(routes, map[string]string{
					"method": method,
					"path":   path,
				})
			}
			return nil
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes)
	}
}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"version":"1.0.0"}`))
}

func ListDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list := dm.GetDevices()
	json.NewEncoder(w).Encode(list)
}

func GetChannels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	channels := dm.GetChannels(id)
	json.NewEncoder(w).Encode(channels)
}

func InviteChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result":"Not implemented"}`))
}

func ByeChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result":"Not implemented"}`))
}

func PTZControl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result":"Not implemented"}`))
}
