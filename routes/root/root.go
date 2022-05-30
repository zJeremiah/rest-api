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
		Route(),
	)
}

// build the endpoint and then add the endpoint to the setup endpoints
func Route() setup.Endpoint {
	return setup.Endpoint{
		Path:         "/",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		HandlerFunc:  Handler,
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
func Handler(w http.ResponseWriter, r *http.Request) error {
	var err error
	var respBody []byte

	e := Route()

	respBody, err = json.MarshalIndent(e.ResponseBody, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error for response body %w", err)
	}

	w.Header().Set("Content-Type", setup.ContentJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
	return nil
}
