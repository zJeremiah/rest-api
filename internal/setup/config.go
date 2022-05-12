package setup

import (
	"github.com/go-chi/chi/v5"
	"github.com/rest-api/logger"
	"github.com/rest-api/router"
)

type Config struct {
	Port  int            `toml:"port" json:"port" flag:"port"`
	Debug bool           `toml:"debug" json:"debug" flag:"debug"`
	Log   logger.Options `toml:"log_options" json:"log_options"`

	Router    *chi.Mux
	Endpoints []router.Endpoint
}
