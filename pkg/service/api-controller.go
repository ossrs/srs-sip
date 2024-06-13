package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *GB28181Server) RegisterRoutes(router *mux.Router) {

	apiV1Router := router.PathPrefix("/gb/v1").Subrouter()

	// Add Auth middleware
	//apiV1Router.Use(authMiddleware)

	apiV1Router.HandleFunc("/version", s.ApiGetVersion).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/devices", s.ApiListDevices).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/devices/{id}/channels", s.ApiGetChannelByDeviceId).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/channels", s.ApiGetAllChannels).Methods(http.MethodGet)

	apiV1Router.HandleFunc("/invite", s.ApiInvite).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/bye", s.ApiBye).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/ptz", s.ApiPTZControl).Methods(http.MethodPost)

	apiV1Router.HandleFunc("", s.GetAPIRoutes(apiV1Router)).Methods(http.MethodGet)

	router.HandleFunc("/gb", s.ApiGetAPIVersion).Methods(http.MethodGet)
}

func (s *GB28181Server) RespondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	wrapper := map[string]interface{}{
		"code": code,
		"data": data,
	}
	json.NewEncoder(w).Encode(wrapper)
}

func (s *GB28181Server) RespondWithJSONSimple(w http.ResponseWriter, jsonStr string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonStr))
}

func (s *GB28181Server) ApiGetAPIVersion(w http.ResponseWriter, r *http.Request) {
	versionInfo := map[string]string{
		"version": "v1",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionInfo)
}

func (s *GB28181Server) GetAPIRoutes(router *mux.Router) http.HandlerFunc {
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

func (s *GB28181Server) ApiGetVersion(w http.ResponseWriter, r *http.Request) {
	s.RespondWithJSONSimple(w, `{"version":"1.0.0"}`)
}

func (s *GB28181Server) ApiListDevices(w http.ResponseWriter, r *http.Request) {
	list := dm.GetDevices()
	s.RespondWithJSON(w, 0, list)
}

func (s *GB28181Server) ApiGetChannelByDeviceId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	channels := dm.ApiGetChannelByDeviceId(id)
	s.RespondWithJSON(w, 0, channels)
}

func (s *GB28181Server) ApiGetAllChannels(w http.ResponseWriter, r *http.Request) {
	channels := dm.GetAllVideoChannels()
	s.RespondWithJSON(w, 0, channels)
}

// request: {"device_id": "1", "channel_id": "1", "sub_stream": 0}
// response: {"code": 0, "data": {"channel_id": "1", "url": "webrtc://"}}
func (s *GB28181Server) ApiInvite(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req map[string]string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get device and channel
	deviceID := req["device_id"]
	channelID := req["channel_id"]
	//subStream := req["sub_stream"]

	code := 0
	url := ""

	defer func() {
		data := map[string]string{
			"channel_id": channelID,
			"url":        url,
		}
		s.RespondWithJSON(w, code, data)
	}()

	c, ok := s.GetVideoChannelStatue(channelID)
	if ok {
		code = 0
		url = "webrtc://" + s.conf.MediaAddr + "/live/" + c.Ssrc
		return
	}

	if err := s.Invite(deviceID, channelID); err != nil {
		code = http.StatusInternalServerError
		return
	}
	c, ok = s.GetVideoChannelStatue(channelID)
	if !ok {
		code = http.StatusInternalServerError
		return
	}
	url = "webrtc://" + s.conf.MediaAddr + "/live/" + c.Ssrc
}

func (s *GB28181Server) ApiBye(w http.ResponseWriter, r *http.Request) {
	s.RespondWithJSONSimple(w, `{"msg":"Not implemented"}`)
}

func (s *GB28181Server) ApiPTZControl(w http.ResponseWriter, r *http.Request) {
	s.RespondWithJSONSimple(w, `{"msg":"Not implemented"}`)
}
