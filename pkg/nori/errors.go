package nori

import (
	"fmt"

	"github.com/nori-io/common/v5/pkg/domain/meta"
)

type ErrorPluginAlreadyInited struct {
	ID meta.ID
}

func (e ErrorPluginAlreadyInited) Error() string {
	return fmt.Sprintf("plugin %s already initialized", e.ID.String())
}

type ErrorPluginAlreadyStarted struct {
	ID meta.ID
}

func (e ErrorPluginAlreadyStarted) Error() string {
	return fmt.Sprintf("plugin %s already started", e.ID.String())
}

type ErrorPluginAlreadyStopped struct {
	ID meta.ID
}

func (e ErrorPluginAlreadyStopped) Error() string {
	return fmt.Sprintf("plugin %s already stopped", e.ID.String())
}
