package router

import (
	"errors"
	"net/http"

	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/version"
)

const otherGroup = "other_group"

func OtherRoute() Endpoint {
	e := Endpoint{
		Path:         "/",
		Method:       GET,
		ResponseType: ContentJSON,
		ResponseBody: Home{
			AppName: "rest-api",
			Version: version.JSON(),
		},
	}

	e.HandlerFunc = e.OtherHome
	return e
}

// Welcome is an endpoint example for the root / path
func (e *Endpoint) OtherHome(w http.ResponseWriter, r *http.Request) error {
	e.Respond(w)
	return nil
}

func OtherTestErrRoute() Endpoint {
	e := Endpoint{
		Name:         "v2 test error",
		Path:         "/error",
		Method:       GET,
		ResponseType: ContentJSON,
		ResponseBody: Home{
			AppName: "rest-api",
			Version: version.JSON(),
		},
	}

	e.HandlerFunc = e.OtherTestError
	return e
}

// TestError is an endpoint example of how an error is returned with the handler
// Example is to test pushing an error to the logger middleware
func (e *Endpoint) OtherTestError(w http.ResponseWriter, r *http.Request) error {
	err := logger.NewError(r, // defines the error that will be returned
		"internal error message",
		"response body error message",
		400,
		errors.New("testing 404 error"))
	return err
}

// PostRoute defines the endpoint request/response info for documentation
func OtherPostRoute() Endpoint {
	e := Endpoint{
		Name:         "v2 post test",
		Path:         "/post",
		Method:       POST,
		RequestType:  ContentJSON,
		RequestBody:  PostReq{}, // this defines what you expect from the request body
		ResponseType: ContentJSON,
		ResponseBody: PostResp{}, // this defines what will be returned in the response body
		Pretty:       true,
	}

	e.HandlerFunc = e.OtherPostTest
	return e
}

// PostTest is the request handler function
func (e *Endpoint) OtherPostTest(w http.ResponseWriter, r *http.Request) error {
	req := PostReq{}
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

	resp := PostResp{
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
