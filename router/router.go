package router

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/setup"
)

var json = jsoniter.ConfigFastest

// InitGroups will initialize the endpoints organized by group
// add new endpoint rout methods here to include them in the api request routes
// these are auto documented based on the endpoint values
func InitGroups() {

	setup.AddEndpoints("", "",
		HomeRoute(),
		DocsRoute(),
		EndPointStruct(),
	)
	setup.AddEndpoints("v1", "test",
		TestErrRoute(),
		PostRoute())
	setup.AddEndpoints("v2", "test",
		OtherRoute(),
		OtherTestErrRoute(),
		OtherPostRoute())

}

// Respond is a helper method to automate the data for the response
// sets the response code, content type, and body
// pretty will pretty print the json output for the response
func Respond(w http.ResponseWriter, response interface{}, pretty bool) error {
	var err error
	var respBody []byte
	if pretty {
		respBody, err = json.MarshalIndent(response, "", "  ")
	} else {
		respBody, err = json.Marshal(response)
	}
	if err != nil {
		return fmt.Errorf("marshal error for resposne body %w", err)
	}
	w.Header().Set("Content-Type", setup.ContentJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
	return nil
}
