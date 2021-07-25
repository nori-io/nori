package http

import (
	"context"
	"net/http"
	"time"

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

func NewHandler(params Params)*Handler {
	mux := http.NewServeMux()

	// static
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(web.IndexHtml)
		if err != nil {
			params.Env.Logger.Error(err.Error())
		}
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Content))))

	// api
	mux.HandleFunc("/system/plugins", func(writer http.ResponseWriter, request *http.Request) {

	})

	return &Handler{
		Env: params.Env,
		FileService:      params.FileService,
		InstalledService: params.InstalledService,
		PluginManager:    params.PluginManager,
		PluginService:    params.PluginService,
		Server: &http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			MaxHeaderBytes:    2 << 20,
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