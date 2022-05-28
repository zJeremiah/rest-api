package setup

import (
	"github.com/go-chi/chi/v5"
	"github.com/rest-api/internal/logger"
	"github.com/rest-api/router"
)

type Config struct {
	Host  string         `toml:"host" comment:"environment hostname"`
	Port  int            `toml:"port" json:"port" flag:"port" comment:"http port number"`
	Debug bool           `toml:"debug" json:"debug" flag:"debug" comment:"show debug logging"`
	Log   logger.Options `toml:"log_options" json:"log_options"`

	Router      *chi.Mux
	RouteGroups router.Groups
}
