package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	common_meta "github.com/nori-io/common/v5/pkg/meta"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	"github.com/nori-io/nori/web"
	"github.com/nori-io/nori/web/static"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env     *env.Env
	FileService      service.FileService
	InstalledService service.PluginOptionService
	PluginManager    service.PluginManager
	PluginService    service.PluginService
}

type Handler struct {
	Env     *env.Env
	FileService      service.FileService
	InstalledService service.PluginOptionService
	PluginManager    service.PluginManager
	PluginService    service.PluginService

	Server *http.Server
}

func New(params Params) *Handler {
	m := mux.NewRouter()

	// static
	m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(web.IndexHtml)
		if err != nil {
			params.Env.Logger.Error(err.Error())
		}
	}).Methods("GET")
	m.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Content))))


	// api
	// GET /system/plugins
	m.HandleFunc("/system/plugins", func(writer http.ResponseWriter, request *http.Request) {
		plugins := params.PluginService.GetAll()

		pd := []meta.Meta{}
		for _, p := range plugins {
			pd = append(pd, p.Meta())
		}

		data, err := json.Marshal(pd)
		if err != nil {
			writer.Write([]byte("err: " + err.Error()))
			return
		}

		writer.Write(data)
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/meta", func(writer http.ResponseWriter, request *http.Request) {
		//
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/dependencies", func(writer http.ResponseWriter, request *http.Request) {
		//
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/dependent", func(writer http.ResponseWriter, request *http.Request) {
		//
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/enable", func(writer http.ResponseWriter, request *http.Request) {
		id := getID(request)
		if id == nil {
			writer.Write([]byte("incorrect id format"))
		}
		if err := params.PluginManager.Enable(request.Context(), id); err != nil {
			writer.Write([]byte("err: " + err.Error()))
		} else {
			writer.Write([]byte("enabled"))
		}
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/disable", func(writer http.ResponseWriter, request *http.Request) {
		id := getID(request)
		if id == nil {
			writer.Write([]byte("incorrect id format"))
		}
		if err := params.PluginManager.Disable(request.Context(), id); err != nil {
			writer.Write([]byte("err: " + err.Error()))
		} else {
			writer.Write([]byte("disabled"))
		}
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/install", func(writer http.ResponseWriter, request *http.Request) {
		id := getID(request)
		if id == nil {
			writer.Write([]byte("incorrect id format"))
		}
		if err := params.PluginManager.Install(request.Context(), id); err != nil {
			writer.Write([]byte("err: " + err.Error()))
		} else {
			writer.Write([]byte("installed"))
		}
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/uninstall", func(writer http.ResponseWriter, request *http.Request) {
		id := getID(request)
		if id == nil {
			writer.Write([]byte("incorrect id format"))
		}
		if err := params.PluginManager.UnInstall(request.Context(), id); err != nil {
			writer.Write([]byte("err: " + err.Error()))
		} else {
			writer.Write([]byte("uninstalled"))
		}
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/start", func(writer http.ResponseWriter, request *http.Request) {
		id := getID(request)
		if id == nil {
			writer.Write([]byte("incorrect id format"))
		}
		if err := params.PluginManager.Start(request.Context(), id); err != nil {
			writer.Write([]byte("err: " + err.Error()))
		} else {
			writer.Write([]byte("started"))
		}
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/stop", func(writer http.ResponseWriter, request *http.Request) {
		id := getID(request)
		if id == nil {
			writer.Write([]byte("incorrect id format"))
		}
		if err := params.PluginManager.Stop(request.Context(), id); err != nil {
			writer.Write([]byte("err: " + err.Error()))
		} else {
			writer.Write([]byte("stopped"))
		}
	}).Methods("GET")

	// transfer
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/get", func(writer http.ResponseWriter, request *http.Request) {
		//
	}).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/download", func(writer http.ResponseWriter, request *http.Request) {
		//
	}).Methods("GET")

	m.HandleFunc("/system/interface/{id:.+}/@v/{version}/", func(writer http.ResponseWriter, request *http.Request) {
		//
	}).Methods("GET")


	return &Handler{
		Env: params.Env,
		FileService:      params.FileService,
		InstalledService: params.InstalledService,
		PluginManager:    params.PluginManager,
		PluginService:    params.PluginService,
		Server: &http.Server{
			Addr:           ":8080",
			Handler:        m,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 2 << 20,
		},
	}
}

func (h *Handler) Start() error {
	h.Env.Logger.Info("Nori Admin HTTP server starting")
	return h.Server.ListenAndServe()
}

func (h *Handler) Stop(ctx context.Context) error {
	h.Env.Logger.Info("Nori Admin HTTP server stopping")
	return h.Server.Shutdown(ctx)
}

func getID(r *http.Request) meta.ID {
	vars := mux.Vars(r)
	return common_meta.ID{
		ID:      meta.PluginID(vars["id"]),
		Version: vars["version"],
	}
}