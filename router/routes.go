package router

import (
	"net/http"

	"github.com/rest-api/internal/version"
)

type Home struct {
	AppName string `json:"app_name"`
	Version version.Struct
}

func HomeRoute() Endpoint {
	e := Endpoint{
		Name:         "Base Path",
		Path:         "/",
		Method:       GET,
		ResponseType: ContentJSON,
		ResponseBody: Home{
			AppName: "rest-api",
			Version: version.JSON(),
		},
	}

	e.HandlerFunc = e.Home
	return e
}

// Home is an endpoint handler func example for the root / path
func (e *Endpoint) Home(w http.ResponseWriter, r *http.Request) error {
	e.Respond(w)
	return nil
}
