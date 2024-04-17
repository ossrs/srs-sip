package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ossrs/srs-sip/pkg/service"
)

func (h *HttpApiServer) RegisterRoutes(router *mux.Router) {

	apiV1Router := router.PathPrefix("/srs-sip/v1").Subrouter()

	// Add Auth middleware
	//apiV1Router.Use(authMiddleware)

	apiV1Router.HandleFunc("/devices", h.ApiListDevices).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/devices/{id}/channels", h.ApiGetChannelByDeviceId).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/channels", h.ApiGetAllChannels).Methods(http.MethodGet)

	apiV1Router.HandleFunc("/invite", h.ApiInvite).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/bye", h.ApiBye).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/ptz", h.ApiPTZControl).Methods(http.MethodPost)

	apiV1Router.HandleFunc("", h.GetAPIRoutes(apiV1Router)).Methods(http.MethodGet)

	router.HandleFunc("/srs-sip", h.ApiGetAPIVersion).Methods(http.MethodGet)
}

func (h *HttpApiServer) RespondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	wrapper := map[string]interface{}{
		"code": code,
		"data": data,
	}
	json.NewEncoder(w).Encode(wrapper)
}

func (h *HttpApiServer) RespondWithJSONSimple(w http.ResponseWriter, jsonStr string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonStr))
}

func (h *HttpApiServer) GetAPIRoutes(router *mux.Router) http.HandlerFunc {
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

		h.RespondWithJSON(w, 0, routes)
	}
}

func (h *HttpApiServer) ApiGetAPIVersion(w http.ResponseWriter, r *http.Request) {
	h.RespondWithJSONSimple(w, `{"version": "v1"}`)
}

func (h *HttpApiServer) ApiListDevices(w http.ResponseWriter, r *http.Request) {
	list := service.DM.GetDevices()
	h.RespondWithJSON(w, 0, list)
}

func (h *HttpApiServer) ApiGetChannelByDeviceId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	channels := service.DM.ApiGetChannelByDeviceId(id)
	h.RespondWithJSON(w, 0, channels)
}

func (h *HttpApiServer) ApiGetAllChannels(w http.ResponseWriter, r *http.Request) {
	channels := service.DM.GetAllVideoChannels()
	h.RespondWithJSON(w, 0, channels)
}

// request: {"device_id": "1", "channel_id": "1", "sub_stream": 0}
// response: {"code": 0, "data": {"channel_id": "1", "url": "webrtc://"}}
func (h *HttpApiServer) ApiInvite(w http.ResponseWriter, r *http.Request) {
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
		h.RespondWithJSON(w, code, data)
	}()

	if err := h.sipSvr.Uas.Invite(deviceID, channelID); err != nil {
		code = http.StatusInternalServerError
		return
	}
	c, ok := h.sipSvr.Uas.GetVideoChannelStatue(channelID)
	if !ok {
		code = http.StatusInternalServerError
		return
	}
	url = "webrtc://" + h.conf.MediaAddr + "/live/" + c.Ssrc
}

func (h *HttpApiServer) ApiBye(w http.ResponseWriter, r *http.Request) {
	h.RespondWithJSONSimple(w, `{"msg":"Not implemented"}`)
}

func (h *HttpApiServer) ApiPTZControl(w http.ResponseWriter, r *http.Request) {
	h.RespondWithJSONSimple(w, `{"msg":"Not implemented"}`)
}
