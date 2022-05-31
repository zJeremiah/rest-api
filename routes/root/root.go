package root

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rest-api/internal/setup"
	"github.com/rest-api/internal/version"
)

// adds endpoint routes to the api config
func Setup() {
	setup.AddEndpoints(
		RootEP(),
		PostTestEP(),
	)
}

// build the endpoint and then add the endpoint to the setup endpoints
func RootEP() setup.Endpoint {
	return setup.Endpoint{
		Name:         "Root Path",
		Path:         "/",
		Description:  "The root path returns the api name and version.",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		HandlerFunc:  RootHandler,
		ResponseBody: struct {
			AppName string         `json:"app_name"`
			Version version.Struct `json:"build_info"`
		}{
			AppName: "rest-api",
			Version: version.JSON(),
		},
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

type PostTest struct {
	ID    int                    `json:"id"`
	Name  string                 `json:"name"`
	Float float64                `json:"float"`
	Map   map[string]interface{} `json:"map"`
	Array []interface{}          `json:"array"`
}

// build the endpoint and then add the endpoint to the setup endpoints
func PostTestEP() setup.Endpoint {
	return setup.Endpoint{
		Name:         "Testing Post Route",
		Path:         "/testing",
		Description:  "This is for testing post requests",
		Method:       setup.POST,
		RequestType:  setup.ContentJSON,
		ResponseType: setup.ContentJSON,
		HandlerFunc:  TestHandler,
		RequestBody: PostTest{
			ID:    123,
			Name:  "Alex Doe",
			Float: 123.456,
			Map: map[string]interface{}{
				"test":  "value",
				"test2": 432,
			},
			Array: []interface{}{"string value", 123.12, map[string]string{"key": "value"}, 6575},
		},
		ResponseBody: PostTest{
			ID:    123,
			Name:  "Alex Doe",
			Float: 123.456,
			Map: map[string]interface{}{
				"test":  "value",
				"test2": 432,
			},
			Array: []interface{}{"string value", 123.12, map[string]string{"key": "value"}, 6575},
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
