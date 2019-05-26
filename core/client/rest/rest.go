package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sync"

	"github.com/nori-io/nori-common/logger"

	"github.com/cheebo/gorest"

	"github.com/gobuffalo/packr"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori/core/plugins"

	"github.com/gorilla/mux"
)

func New(addr, base string,
	manager plugins.Manager,
	wg *sync.WaitGroup,
	shutdownCh <-chan struct{},
	logger logger.Logger,
) {
	r := mux.NewRouter()

	RegisterRoutes(base, r, manager)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	wg.Add(1)

	go func() {
		<-shutdownCh
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Errorf("Nori REST API server error: %v", err)
		}
		logger.Info("Stopped Nori Core REST API Service")
		wg.Done()
	}()

	go func() {
		logger.Infof("Starting Nori REST API server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Nori REST API server error: %v", err)
			wg.Done()
		}
	}()
}

func RegisterRoutes(base string, r *mux.Router, m plugins.Manager) {
	box := packr.NewBox("../../../html")
	fs := http.FileServer(box)

	r.Handle(base, fs)
	r.Handle(path.Join(base, "/plugins"), restPlugins(m))
	r.Handle(path.Join(base, "/plugins/installable"), restPluginsInstallable(m))
	r.Handle(path.Join(base, "/plugins/install"), restPluginsInstall(m))
	r.Handle(path.Join(base, "/plugins/running"), restPluginsRunning(m))
	r.Handle(path.Join(base, "/plugins/stop"), restPluginsStop(m))
	r.Handle(path.Join(base, "/plugins/start"), restPluginsStart(m))
}

func restPlugins(m plugins.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var output []byte
		output, err := json.MarshalIndent(m.Metas(plugins.FilterRunnable), "", "\t")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(output)
		}
	}
}

func restPluginsInstallable(m plugins.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var output []byte
		output, err := json.MarshalIndent(m.Metas(plugins.FilterInstallable), "", "\t")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(output)
		}
	}
}

func restPluginsInstall(m plugins.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ID := meta.ID{
			ID:      meta.PluginID(r.URL.Query().Get("id")),
			Version: r.URL.Query().Get("ver"),
		}

		var output []byte

		err := m.Install(ID, context.Background())
		var resp interface{}
		if err != nil {
			resp = rest.ErrResp{
				Meta: rest.ErrMeta{
					ErrCode:    500,
					ErrMessage: err.Error(),
				},
			}
		} else {
			resp = rest.ListResp{
				Items: []interface{}{
					struct {
						Msg string
					}{
						Msg: fmt.Sprintf("Plugin %s successfully installed", ID.String()),
					},
				},
			}
		}

		output, err = json.MarshalIndent(resp, "", "\t")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(output)
		}
	}
}

func restPluginsRunning(m plugins.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var output []byte
		output, err := json.MarshalIndent(m.Metas(plugins.FilterRunning), "", "\t")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(output)
		}
	}
}

func restPluginsStop(m plugins.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html")

		id := meta.ID{
			ID:      meta.PluginID(r.URL.Query().Get("id")),
			Version: r.URL.Query().Get("ver"),
		}

		_, err := m.Meta(id)
		if err == nil {
			err = m.Stop(id, context.Background())
		}

		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("stopped successfully"))
		}
	}
}

func restPluginsStart(m plugins.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html")

		id := meta.ID{
			ID:      meta.PluginID(r.URL.Query().Get("id")),
			Version: r.URL.Query().Get("ver"),
		}

		_, err := m.Meta(id)
		if err == nil {
			err = m.Start(id, context.Background())
		}

		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("started successfully"))
		}
	}
}
