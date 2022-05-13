package router

import (
	"net/http"

	"github.com/rest-api/internal/version"
	"github.com/rest-api/response"
)

func HomeRoute() Endpoint {
	e := Endpoint{
		Path:         "/",
		Method:       GET,
		ResponseType: ContentJSON,
		ResponseBody: response.Home{
			AppName: "rest-api",
			Version: version.JSON(),
		},
	}

	e.HandlerFunc = e.Home
	return e
}

// Welcome is an endpoint example for the root / path
func (e *Endpoint) Home(w http.ResponseWriter, r *http.Request) error {
	e.Respond(w)
	return nil
}

func TestRoute() Endpoint {
	e := Endpoint{
		Path:         "/",
		Method:       OPTIONS,
		ResponseType: ContentJSON,
		ResponseBody: response.Home{
			AppName: "rest-api",
			Version: version.JSON(),
		},
	}

	e.HandlerFunc = e.Test
	return e
}

// Welcome is an endpoint example for the root / path
func (e *Endpoint) Test(w http.ResponseWriter, r *http.Request) error {
	e.Respond(w)
	return nil
}
