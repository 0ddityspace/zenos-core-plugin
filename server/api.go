package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/0ddityspace/zenos-core-plugin/pluginapi"
)

// API exposes a Controller's actions via a standard HTTP interface.
type API struct {
	ctrl             pluginapi.Controller
	logger           *slog.Logger
	additionalRoutes map[string]http.HandlerFunc
}

// NewAPI initializes a new API instance with the provided controller and logger.
func NewAPI(ctrl pluginapi.Controller, logger *slog.Logger) *API {
	return &API{
		ctrl:             ctrl,
		logger:           logger,
		additionalRoutes: make(map[string]http.HandlerFunc),
	}
}

// AddRoute adds a custom route before registering standard routes.
func (api *API) AddRoute(pattern string, handler http.HandlerFunc) {
	api.additionalRoutes[pattern] = handler
}

// RegisterRoutes attaches standard endpoints to the provided HTTP Mux.
func (api *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/manifest", api.handleManifest)
	mux.HandleFunc("GET /v1/status", api.handleStatus)
	mux.HandleFunc("GET /v1/now-playing", api.handleNowPlaying)

	mux.HandleFunc("POST /v1/player/play", api.handleAction(api.ctrl.Play))
	mux.HandleFunc("POST /v1/player/pause", api.handleAction(api.ctrl.Pause))
	mux.HandleFunc("POST /v1/player/stop", api.handleAction(api.ctrl.Stop))
	mux.HandleFunc("POST /v1/player/next", api.handleAction(api.ctrl.Next))
	mux.HandleFunc("POST /v1/player/previous", api.handleAction(api.ctrl.Previous))

	mux.HandleFunc("POST /v1/player/seek", api.handleSeek)
	mux.HandleFunc("POST /v1/player/volume", api.handleVolume)

	for pattern, h := range api.additionalRoutes {
		mux.HandleFunc(pattern, h)
	}
}

func (api *API) handleManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.ctrl.GetManifest())
}

func (api *API) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	up := api.ctrl.IsUp()
	json.NewEncoder(w).Encode(map[string]any{
		"ok":        true,
		"engine_up": up,
	})
}

func (api *API) handleNowPlaying(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.ctrl.GetNowPlaying())
}

func (api *API) handleAction(action func() error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := action(); err != nil {
			api.logger.Error("action failed", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func (api *API) handleSeek(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		http.Error(w, "missing offset", http.StatusBadRequest)
		return
	}
	offsetMS, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := api.ctrl.Seek(offsetMS); err != nil {
		api.logger.Error("seek failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (api *API) handleVolume(w http.ResponseWriter, r *http.Request) {
	levelStr := r.URL.Query().Get("level")
	if levelStr == "" {
		http.Error(w, "missing level", http.StatusBadRequest)
		return
	}
	level, err := strconv.ParseFloat(levelStr, 64)
	if err != nil {
		http.Error(w, "invalid level", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := api.ctrl.SetVolume(level); err != nil {
		api.logger.Error("volume failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
