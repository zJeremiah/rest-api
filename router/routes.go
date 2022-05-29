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

func EndPointStruct() setup.Endpoint {
	e := setup.Endpoint{
		Path:         "/endpoint",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		HandlerFunc:  EP,
		ResponseBody: setup.Endpoint{
			Version:      "v1",
			Group:        "group",
			RequestType:  setup.ContentJSON,
			ResponseType: setup.ContentJSON,
			Path:         "/v1/group/endpoint.path?param1=1st&param2=2nd",
			Method:       setup.GET,
			Description:  "endpoint example",
			RequestBody:  "request body struct",
			ResponseBody: "response body struct",
		},
	}
	return e
}

func EP(w http.ResponseWriter, r *http.Request) error {
	e := EndPointStruct()

	err := Respond(w, e.ResponseBody, true)
	if err != nil {
		return err
	}
	return nil
}
