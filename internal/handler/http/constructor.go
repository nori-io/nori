package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	"github.com/nori-io/nori/internal/handler/http/handler/plugin"
	"github.com/nori-io/nori/web"
	"github.com/nori-io/nori/web/static"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env     *env.Env
	PluginHandler plugin.Handler
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
	m.HandleFunc("/system/plugins", params.PluginHandler.Plugins).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/meta", params.PluginHandler.Meta).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/dependencies", params.PluginHandler.Dependencies).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/dependent", params.PluginHandler.Dependent).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/enable", params.PluginHandler.Enable).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/disable", params.PluginHandler.Disable).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/install", params.PluginHandler.Install).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/uninstall", params.PluginHandler.UnInstall).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/start", params.PluginHandler.Start).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/start-all", params.PluginHandler.StartAll).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/stop", params.PluginHandler.Stop).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/stop-all", params.PluginHandler.StopAll).Methods("GET")

	// transfer
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/get", params.PluginHandler.Get).Methods("GET")
	m.HandleFunc("/system/plugin/{id:.+}/@v/{version}/download", params.PluginHandler.Download).Methods("GET")

	m.HandleFunc("/system/interface/{id:.+}/@v/{version}/", func(writer http.ResponseWriter, request *http.Request) {
		//
	}).Methods("GET")


	return &Handler{
		Env: params.Env,
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
