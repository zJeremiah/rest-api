package root

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
	"github.com/rest-api/internal/version"
)

// adds endpoint routes to the api config
func Setup() {
	setup.AddEndpoints(
		RootEP(),
		PostTestEP(),
		ErrorEP(),
	)
}

// build the endpoint and then add the endpoint to the setup endpoints
func RootEP() setup.Endpoint {
	return setup.Endpoint{
		Name:         "Root Path",
		Path:         "/",
		Description:  "The root path returns the api name and version.",
		Methods:      setup.Methods{setup.ANY},
		ResponseType: setup.ContentJSON,
		HandlerFunc:  RootHandler,
		ResponseBody: version.JSON(),
	}
}

// Handler func for the root / path
func RootHandler(w http.ResponseWriter, r *http.Request) error {
	var err error
	var respBody []byte

	e := RootEP()

	respBody, err = json.MarshalIndent(e.ResponseBody, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error for response body %w", err)
	}

	w.Header().Set("Content-Type", setup.ContentJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
	return nil
}

func ErrorEP() setup.Endpoint {
	return setup.Endpoint{
		Name:         "Error Path",
		Path:         "/error",
		Description:  "This route returns an error depending on given id (none, 1, or 2)",
		Methods:      setup.Methods{setup.GET},
		ResponseType: setup.ContentJSON,
		HandlerFunc:  ErrorHandler,
		QueryParams: []setup.Param{
			{Name: "id", Description: "use none, 1 or 2 for different errors"},
		},
		ResponseBody: logger.RespBody{
			Msg:    "you have made a bad request",
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
		},
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")

	if id == "1" {
		return &logger.APIErr{
			Internal: logger.Internal{
				Msg: "internal message number 1",
				Err: fmt.Errorf("internal error 1"),
			},
			RespBody: logger.RespBody{
				Msg:  "you have made a bad request",
				Code: http.StatusBadRequest,
			},
		}
	}

	if id == "2" {
		return &logger.APIErr{
			Internal: logger.Internal{
				Msg: "internal message number 2",
				Err: fmt.Errorf("internal error 2"),
			},
			RespBody: logger.RespBody{
				Msg:  "this is not acceptable",
				Code: http.StatusNotAcceptable,
			},
		}
	}
	return fmt.Errorf("an error has occured ")
}

type PostTest struct {
	ID    int            `json:"id"`
	Name  string         `json:"name"`
	Float float64        `json:"float"`
	Map   map[string]any `json:"map"`
	Array []any          `json:"array"`
}

// build the endpoint and then add the endpoint to the setup endpoints
func PostTestEP() setup.Endpoint {
	return setup.Endpoint{
		Name:         "Testing Post Route",
		Path:         "/testing",
		Description:  "This is for testing post requests",
		Methods:      setup.Methods{setup.POST},
		RequestType:  setup.ContentJSON,
		ResponseType: setup.ContentJSON,
		HandlerFunc:  TestHandler,
		RequestBody: PostTest{
			ID:    123,
			Name:  "Alex Doe",
			Float: 123.456,
			Map: map[string]any{
				"test":  "value",
				"test2": 432,
			},
			Array: []any{"string value", 123.12, map[string]string{"key": "value"}, 6575},
		},
		ResponseBody: PostTest{
			ID:    123,
			Name:  "Alex Doe",
			Float: 123.456,
			Map: map[string]any{
				"test":  "value",
				"test2": 432,
			},
			Array: []any{"string value", 123.12, map[string]string{"key": "value"}, 6575},
		},
		JSONFields: []setup.Param{
			{Name: "id"}, {Name: "name"}, {Name: "float"}, {Name: "map"}, {Name: "array"},
		},
	}
}

// Handler func
func TestHandler(w http.ResponseWriter, r *http.Request) error {
	var err error
	var respBody []byte

	e := RootEP()

	respBody, err = json.MarshalIndent(e.ResponseBody, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error for response body %w", err)
	}

	w.Header().Set("Content-Type", setup.ContentJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
	return nil
}
