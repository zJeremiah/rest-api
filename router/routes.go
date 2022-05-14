package router

import (
	"errors"
	"net/http"

	"github.com/rest-api/internal/apierr"
	"github.com/rest-api/internal/version"
	"github.com/rest-api/logger"
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

func TestErrRoute() Endpoint {
	e := Endpoint{
		Path:         "/test_error",
		Method:       GET,
		ResponseType: ContentJSON,
		ResponseBody: response.Home{
			AppName: "rest-api",
			Version: version.JSON(),
		},
	}

	e.HandlerFunc = e.TestError
	return e
}

// Welcome is an endpoint example for the root / path
// Example is to test pushing an error to the logger middleware
func (e *Endpoint) TestError(w http.ResponseWriter, r *http.Request) error {
	req := r.Context().Value("request").(*logger.Request)
	req.APIError = apierr.NewError("internal error message", "response body error message", 400, errors.New("testing 404 error"))
	return req.APIError
}
