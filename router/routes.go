package router

import (
	"net/http"

	"github.com/rest-api/internal/setup"
	"github.com/rest-api/internal/version"
)

func HomeRoute() setup.Endpoint {
	e := setup.Endpoint{
		Path:         "/",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		HandlerFunc:  Home,
	}

	return e
}

// Home is an setup.Endpoint handler func example for the root / path
func Home(w http.ResponseWriter, r *http.Request) error {
	Respond(w, struct {
		AppName string         `json:"app_name"`
		Version version.Struct `json:"build_info"`
	}{
		AppName: "rest-api",
		Version: version.JSON(),
	}, true)
	return nil
}

func DocsRoute() setup.Endpoint {
	return setup.Endpoint{
		Description:  "api documentation",
		Path:         "/docs",
		Method:       setup.GET,
		ResponseType: "text/html",
		HandlerFunc:  Docs,
		ResponseBody: []byte("<p>docs parsed template goes here</p>"),
	}
}

// Home is an setup.Endpoint handler func example for the root / path
func Docs(w http.ResponseWriter, r *http.Request) error {
	e := DocsRoute()

	w.Header().Set("Content-Type", e.ResponseType)
	w.WriteHeader(http.StatusOK)
	w.Write(e.ResponseBody.([]byte))
	return nil
}
