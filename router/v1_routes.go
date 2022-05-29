package router

import (
	"errors"
	"net/http"

	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
)

type PostResp struct {
	RID    int               `json:"id,omitempty"`
	RName  string            `json:"name,omitempty"`
	RFloat float64           `json:"field_float,omitempty"`
	RMap   map[string]string `json:"field_map,omitempty"`
	RArray []interface{}     `json:"array,omitempty"`
}

type PostReq struct {
	ID    int               `json:"id,omitempty"`
	Name  string            `json:"name,omitempty"`
	Float float64           `json:"field_float,omitempty"`
	Map   map[string]string `json:"field_map,omitempty"`
	Array []interface{}     `json:"array,omitempty"`
}

// an example of an endpoint function
func TestErrRoute() setup.Endpoint {
	return setup.Endpoint{
		Description:  "return error testing",
		Path:         "/error",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		HandlerFunc:  TestError,
	}
}

// TestError is an endpoint example of how an error is returned with the handler
// Example is to test pushing an error to the logger middleware
func TestError(w http.ResponseWriter, r *http.Request) error {
	err := logger.NewError( // defines the error that will be returned
		"internal error message",
		"response body error message",
		400,
		errors.New("testing 400 error"))
	return err
}

// PostRoute defines the endpoint request/response info for documentation
func PostRoute() setup.Endpoint {
	e := setup.Endpoint{
		Description:  "POST testing endpoint",
		Path:         "/post",
		Method:       setup.POST,
		RequestType:  setup.ContentJSON,
		RequestBody:  PostReq{}, // this defines what you expect from the request body
		ResponseType: setup.ContentJSON,
		ResponseBody: PostResp{}, // this defines what will be returned in the response body
		Pretty:       true,
		HandlerFunc:  PostTest,
	}

	return e
}

// PostTest is the request handler function
func PostTest(w http.ResponseWriter, r *http.Request) error {
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
