package router

import (
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
		DocsRoute())
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
func Respond(w http.ResponseWriter, response interface{}, pretty bool) {
	var respBody []byte
	respBody, _ = json.Marshal(response)
	if pretty {
		respBody, _ = json.MarshalIndent(response, "", "  ")
	}
	w.Header().Set("Content-Type", setup.ContentJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
