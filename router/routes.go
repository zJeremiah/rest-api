package router

import (
	"errors"
	"net/http"

	"github.com/rest-api/internal/apierr"
	"github.com/rest-api/internal/version"
	"github.com/rest-api/logger"
	"github.com/rest-api/request"
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
	req.APIError = apierr.NewError( // defines the error that will be returned
		"internal error message",
		"response body error message",
		400,
		errors.New("testing 404 error"))
	return req.APIError
}

// PostRoute defines the endpoint request/response info for documentation
func PostRoute() Endpoint {
	e := Endpoint{
		Path:         "/post",
		Method:       POST,
		RequestType:  ContentJSON,
		RequestBody:  request.PostReq{}, // this defines what you expect from the request body
		ResponseType: ContentJSON,
		ResponseBody: response.PostResp{}, // this defines what will be returned in the response body
		Pretty:       true,
	}

	e.HandlerFunc = e.PostTest
	return e
}

// PostTest is the request handler function
func (e *Endpoint) PostTest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Req(r)
	req := request.PostReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		e1 := apierr.NewError( // defines the error that will be returned
			"error decoding request body",
			"could not decode request body",
			400, err)

	}

	e.RequestBody = req

	resp := response.PostResp{
		RID:    req.ID,
		RName:  req.Name,
		RFloat: req.Float,
		RMap:   req.Map,
		RArray: req.Array,
	}
	e.ResponseBody = resp
	e.Respond(w)
	return nil
}
