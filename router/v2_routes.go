package router

import (
	"errors"
	"net/http"

	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
	"github.com/rest-api/internal/version"
)

func OtherRoute() setup.Endpoint {
	e := setup.Endpoint{
		Path:         "/",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		ResponseBody: "",
		HandlerFunc:  OtherHome,
	}

	return e
}

// Welcome is an endpoint example for the root / path handler function
func OtherHome(w http.ResponseWriter, r *http.Request) error {
	Respond(w, version.JSON(), true)
	return nil
}

func OtherTestErrRoute() setup.Endpoint {
	e := setup.Endpoint{
		Description:  "v2 testing returned error",
		Path:         "/error",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		HandlerFunc:  OtherTestError,
	}

	return e
}

// TestError is an endpoint example of how an error is returned with the handler
// Example is to test pushing an error to the logger middleware
func OtherTestError(w http.ResponseWriter, r *http.Request) error {
	err := errors.New("this is a testing error")
	return err
}

// PostRoute defines the endpoint request/response info for documentation
func OtherPostRoute() setup.Endpoint {
	e := setup.Endpoint{
		Description:  "v2 other post test",
		Path:         "/post",
		Method:       setup.POST,
		RequestType:  setup.ContentJSON,
		RequestBody:  PostReq{}, // this defines what you expect from the request body
		ResponseType: setup.ContentJSON,
		ResponseBody: PostResp{}, // this defines what will be returned in the response body
		Pretty:       true,
		HandlerFunc:  OtherPostTest,
	}

	return e
}

// PostTest is the request handler function
func OtherPostTest(w http.ResponseWriter, r *http.Request) error {
	req := PostReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// defines the error that will be returned, use this to set the request log with the error
		err = logger.NewError(
			"error decoding request body",
			"could not decode request body",
			400, err)
		return err
	}

	resp := PostResp{
		RID:    req.ID,
		RName:  req.Name,
		RFloat: req.Float,
		RMap:   req.Map,
		RArray: req.Array,
	}

	Respond(w, resp, true) // writes the response code and body if there are no errors
	return nil
}
