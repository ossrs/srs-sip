package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ossrs/srs-sip/pkg/models"
	"github.com/ossrs/srs-sip/pkg/service"
)

func (h *HttpApiServer) RegisterRoutes(router *mux.Router) {
	// Add Auth middleware
	//apiV1Router.Use(authMiddleware)

	router.HandleFunc("/devices", h.ApiListDevices).Methods(http.MethodGet)
	router.HandleFunc("/devices/{id}/channels", h.ApiGetChannelByDeviceId).Methods(http.MethodGet)
	router.HandleFunc("/channels", h.ApiGetAllChannels).Methods(http.MethodGet)

	router.HandleFunc("/invite", h.ApiInvite).Methods(http.MethodPost)
	router.HandleFunc("/bye", h.ApiBye).Methods(http.MethodPost)
	router.HandleFunc("/ptz", h.ApiPTZControl).Methods(http.MethodPost)
	router.HandleFunc("/pause", h.ApiPause).Methods(http.MethodPost)
	router.HandleFunc("/resume", h.ApiResume).Methods(http.MethodPost)
	router.HandleFunc("/speed", h.ApiSpeed).Methods(http.MethodPost)

	router.HandleFunc("/query-record", h.ApiQueryRecord).Methods(http.MethodPost)

	// 媒体服务器相关接口，查询，新增，删除，用restful风格
	router.HandleFunc("/media-servers", h.ApiListMediaServers).Methods(http.MethodGet)
	router.HandleFunc("/media-servers", h.ApiAddMediaServer).Methods(http.MethodPost)
	router.HandleFunc("/media-servers/{id}", h.ApiDeleteMediaServer).Methods(http.MethodDelete)
	router.HandleFunc("/media-servers/default/{id}", h.ApiSetDefaultMediaServer).Methods(http.MethodPost)

	router.HandleFunc("", h.GetAPIRoutes(router)).Methods(http.MethodGet)
}

func (h *HttpApiServer) RespondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	wrapper := models.CommonResponse{
		Code: code,
		Data: data,
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

func (h *HttpApiServer) ApiInvite(w http.ResponseWriter, r *http.Request) {
	var req models.InviteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	session, err := h.sipSvr.Uas.Invite(req)
	if err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	response := models.InviteResponse{
		ChannelID: req.ChannelID,
		URL:       session.URL,
	}
	h.RespondWithJSON(w, 0, response)
}

func (h *HttpApiServer) ApiBye(w http.ResponseWriter, r *http.Request) {
	var req models.ByeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	if err := h.sipSvr.Uas.Bye(req); err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, map[string]string{"msg": "success"})
}

func (h *HttpApiServer) ApiPause(w http.ResponseWriter, r *http.Request) {
	var req models.PauseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	if err := h.sipSvr.Uas.Pause(req); err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, map[string]string{"msg": "success"})
}

func (h *HttpApiServer) ApiResume(w http.ResponseWriter, r *http.Request) {
	var req models.ResumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	if err := h.sipSvr.Uas.Resume(req); err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, map[string]string{"msg": "success"})
}

func (h *HttpApiServer) ApiSpeed(w http.ResponseWriter, r *http.Request) {
	var req models.SpeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	if err := h.sipSvr.Uas.Speed(req); err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, map[string]string{"msg": "success"})
}

// request: {"device_id": "1", "channel_id": "1", "ptz": "up", "speed": "1}
func (h *HttpApiServer) ApiPTZControl(w http.ResponseWriter, r *http.Request) {
	var req models.PTZControlRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": err.Error()})
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
	var req models.QueryRecordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": err.Error()})
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
	servers, err := service.MediaDB.ListMediaServers()
	if err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, servers)
}

// request: {"name": "srs1", "ip": "192.168.1.100", "port": 1935, "type": "SRS", "username": "admin", "password": "123456"}
func (h *HttpApiServer) ApiAddMediaServer(w http.ResponseWriter, r *http.Request) {
	var req models.MediaServerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	// 验证必填字段
	if req.Name == "" || req.IP == "" || req.Port == 0 || req.Type == "" {
		h.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"msg": "name, ip, port and type are required"})
		return
	}

	// 添加到数据库
	if err := service.MediaDB.AddMediaServer(req.Name, req.Type, req.IP, req.Port, req.Username, req.Password, req.Secret, req.IsDefault); err != nil {
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

	if err := service.MediaDB.DeleteMediaServer(id); err != nil {
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

	if err := service.MediaDB.SetDefaultMediaServer(id); err != nil {
		h.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	h.RespondWithJSON(w, 0, map[string]string{"msg": "success"})
}
