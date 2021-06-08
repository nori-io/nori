package helper

import (
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/helper/github"
	"go.uber.org/dig"
)

func Provide(container *dig.Container) {
	if err := container.Provide(
		github.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}
}
