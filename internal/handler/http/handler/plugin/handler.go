package plugin

import (
	"net/http"

	"github.com/nori-io/nori/internal/domain/service"
)

func (h *Handler) Plugins(w http.ResponseWriter, r *http.Request) {
	plugins, err := h.PluginManager.GetByFilter(service.GetByFilterData{})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	for _, plugin := range plugins {
		w.Write([]byte(plugin.Plugin.Meta().GetID().String() + " " ))
	}
}
