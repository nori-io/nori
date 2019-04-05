package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/nori-io/nori-common/meta"

	"github.com/nori-io/nori/core/plugins"

	"github.com/cheebo/go-config"
	"github.com/gorilla/mux"
)

func RegisterRoutes(c go_config.Config, r *mux.Router, m plugins.Manager) {
	base := c.String("nori.rest.base")
	//r.Handle(path.Join(base, "/"), noriHome(m))
	r.Handle(path.Join(base, "/plugins"), restPlugins(m))
	r.Handle(path.Join(base, "/plugins/installable"), restPluginsInstallable(m))
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
			ID:      meta.PluginID(r.URL.Query().Get("pid")),
			Version: r.URL.Query().Get("ver"),
		}

		_, err := m.Meta(id)
		if err == nil {
			//println("stop: ", id.String())
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
			ID:      meta.PluginID(r.URL.Query().Get("pid")),
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
