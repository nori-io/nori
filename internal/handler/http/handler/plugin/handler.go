package plugin

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	common_meta "github.com/nori-io/common/v5/pkg/meta"
)

func (h Handler) Plugins(w http.ResponseWriter, r *http.Request) {
	plugins := h.PluginService.GetAll()

	pd := []meta.Meta{}
	for _, p := range plugins {
		pd = append(pd, p.Meta())
	}

	data, err := json.Marshal(pd)
	if err != nil {
		w.Write([]byte("err: " + err.Error()))
		return
	}

	w.Write(data)
}

func (h Handler) Meta(w http.ResponseWriter, r *http.Request) {}

func (h Handler) Dependencies(w http.ResponseWriter, r *http.Request) {}

func (h Handler) Dependent(w http.ResponseWriter, r *http.Request) {}

func (h Handler) Enable(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == nil {
		w.Write([]byte("incorrect id format"))
		return
	}
	if err := h.PluginManager.Enable(r.Context(), id); err != nil {
		w.Write([]byte("err: " + err.Error()))
	} else {
		w.Write([]byte("enabled"))
	}
}

func (h Handler) Disable(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == nil {
		w.Write([]byte("incorrect id format"))
	}
	if err := h.PluginManager.Disable(r.Context(), id); err != nil {
		w.Write([]byte("err: " + err.Error()))
	} else {
		w.Write([]byte("disabled"))
	}
}

func (h Handler) Install(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == nil {
		w.Write([]byte("incorrect id format"))
		return
	}
	if err := h.PluginManager.Install(r.Context(), id); err != nil {
		w.Write([]byte("err: " + err.Error()))
	} else {
		w.Write([]byte("installed"))
	}
}

func (h Handler) UnInstall(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == nil {
		w.Write([]byte("incorrect id format"))
		return
	}
	if err := h.PluginManager.UnInstall(r.Context(), id); err != nil {
		w.Write([]byte("err: " + err.Error()))
	} else {
		w.Write([]byte("uninstalled"))
	}
}

func (h Handler) Start(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == nil {
		w.Write([]byte("incorrect id format"))
	}
	if err := h.PluginManager.Start(r.Context(), id); err != nil {
		w.Write([]byte("err: " + err.Error()))
	} else {
		w.Write([]byte("started"))
	}
}

func (h Handler) StartAll(w http.ResponseWriter, r *http.Request) {
	if err := h.PluginManager.StartAll(r.Context()); err != nil {
		w.Write([]byte("err: " + err.Error()))
	} else {
		w.Write([]byte("started"))
	}
}

func (h Handler) Stop(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == nil {
		w.Write([]byte("incorrect id format"))
	}
	if err := h.PluginManager.Stop(r.Context(), id); err != nil {
		w.Write([]byte("err: " + err.Error()))
	} else {
		w.Write([]byte("stopped"))
	}
}

func (h Handler) StopAll(w http.ResponseWriter, r *http.Request) {
	if err := h.PluginManager.StopAll(r.Context()); err != nil {
		w.Write([]byte("err: " + err.Error()))
	} else {
		w.Write([]byte("stopped"))
	}
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {}

func (h Handler) Download(w http.ResponseWriter, r *http.Request) {}

func (h Handler) Upload(w http.ResponseWriter, r *http.Request) {}


func getID(r *http.Request) meta.ID {
	vars := mux.Vars(r)
	return common_meta.ID{
		ID:      meta.PluginID(vars["id"]),
		Version: vars["version"],
	}
}