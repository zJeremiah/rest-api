package router

import (
	"errors"
	"net/http"

	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/version"
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

// TestError is an endpoint example of how an error is returned with the handler
// Example is to test pushing an error to the logger middleware
func (e *Endpoint) TestError(w http.ResponseWriter, r *http.Request) error {
	err := logger.NewError(r, // defines the error that will be returned
		"internal error message",
		"response body error message",
		400,
		errors.New("testing 404 error"))
	return err
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
	req := request.PostReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// defines the error that will be returned, use this to set the request log with the error
		err = logger.NewError(r,
			"error decoding request body",
			"could not decode request body",
			400, err)
		return err
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
	e.Respond(w) // writes the response code and body if there are no errors
	return nil
}
