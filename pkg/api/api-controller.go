package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	apiV1Router.HandleFunc("/query-record", h.ApiQueryRecord).Methods(http.MethodPost)

	// 媒体服务器相关接口，查询，新增，删除，用restful风格
	apiV1Router.HandleFunc("/media-servers", h.ApiListMediaServers).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/media-servers", h.ApiAddMediaServer).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/media-servers/{id}", h.ApiDeleteMediaServer).Methods(http.MethodDelete)
	apiV1Router.HandleFunc("/media-servers/default/{id}", h.ApiSetDefaultMediaServer).Methods(http.MethodPost)

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

// request: {"media_server_addr": "192.168.1.1:1935", "device_id": "1", "channel_id": "1", "sub_stream": 0}
// response: {"code": 0, "data": {"channel_id": "1", "url": "webrtc://"}}
func (h *HttpApiServer) ApiInvite(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req InviteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := 0
	url := ""

	defer func() {
		response := InviteResponse{
			ChannelID: req.ChannelID,
			URL:       url,
		}
		h.RespondWithJSON(w, code, response)
	}()

	if err := h.sipSvr.Uas.Invite(req.MediaServerAddr, req.DeviceID, req.ChannelID); err != nil {
		code = http.StatusInternalServerError
		return
	}
	c, ok := h.sipSvr.Uas.GetVideoChannelStatue(req.ChannelID)
	if !ok {
		code = http.StatusInternalServerError
		return
	}
	url = "webrtc://" + req.MediaServerAddr + "/live/" + c.Ssrc
}

func (h *HttpApiServer) ApiBye(w http.ResponseWriter, r *http.Request) {
	h.RespondWithJSONSimple(w, `{"msg":"Not implemented"}`)
}

// request: {"device_id": "1", "channel_id": "1", "ptz": "up", "speed": "1}
func (h *HttpApiServer) ApiPTZControl(w http.ResponseWriter, r *http.Request) {
	var req PTZControlRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := 0
	msg := ""
	defer func() {
		h.RespondWithJSON(w, code, map[string]string{"msg": msg})
	}()
	if err := h.sipSvr.Uas.ControlPTZ(req.DeviceID, req.ChannelID, req.PTZ, req.Speed); err != nil {
		code = http.StatusInternalServerError
		msg = err.Error()
		return
	}
	msg = "success"
}

func (h *HttpApiServer) ApiQueryRecord(w http.ResponseWriter, r *http.Request) {
	var req QueryRecordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	records, err := h.sipSvr.Uas.QueryRecord(req.DeviceID, req.ChannelID, req.StartTime, req.EndTime)
	if err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}
	h.RespondWithJSON(w, 0, records)
}

func (h *HttpApiServer) ApiListMediaServers(w http.ResponseWriter, r *http.Request) {
	servers, err := h.mediaDB.ListMediaServers()
	if err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, servers)
}

// request: {"name": "srs1", "ip": "192.168.1.100", "port": 1935, "type": "SRS", "username": "admin", "password": "123456"}
func (h *HttpApiServer) ApiAddMediaServer(w http.ResponseWriter, r *http.Request) {
	var req MediaServerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": "invalid request"})
		return
	}

	// 验证必填字段
	if req.Name == "" || req.IP == "" || req.Port == 0 || req.Type == "" {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": "name, ip, port and type are required"})
		return
	}

	// 添加到数据库
	if err := h.mediaDB.AddMediaServer(req.Name, req.IP, req.Port, req.Username, req.Password, req.Type, req.IsDefault); err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, map[string]string{"msg": "success"})
}

func (h *HttpApiServer) ApiDeleteMediaServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": "invalid id"})
		return
	}

	if err := h.mediaDB.DeleteMediaServer(id); err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, map[string]string{"msg": "success"})
}

func (h *HttpApiServer) ApiSetDefaultMediaServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": "invalid id"})
		return
	}

	if err := h.mediaDB.SetDefaultMediaServer(id); err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, map[string]string{"msg": "success"})
}
