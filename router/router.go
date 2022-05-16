package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/logger"
)

type Endpoints []Endpoint

type RequestMethod uint

const (
	GET = iota
	HEAD
	OPTIONS
	PATCH
	POST
	PUT
	DELETE
)

// Endpoint object is used to setup an api endpoint in the api
// this is used to make the correct request to the handler
// and also used to generate the slate api docs
type Endpoint struct {
	Group        string        // this is the group that the endpoint belongs to
	Path         string        // the URI path for the endpoint
	Method       RequestMethod // the HTTP method for the endpoint GET, POST etc...
	RequestType  string        // The request type is the ContentType for the request
	ResponseType string        // This is the ContentType that will be returned in the response (normally application/json)
	RequestBody  interface{}   // This is used to generate the api docs JSON object string for the request
	ResponseBody interface{}   // This is used to generate the api docs JSON object string for the response
	HandlerFunc  Handler       // the Handler function is called when the router matches the endpoint path
	Description  string        // The description of the endpoint for api docs
	Pretty       bool          // output the json string as pretty format when true
}

const ContentJSON = "application/json"

type Handler func(w http.ResponseWriter, r *http.Request) error

var json = jsoniter.ConfigFastest

func (rm RequestMethod) String() string {
	switch rm {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case HEAD:
		return "HEAD"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case PATCH:
		return "PATCH"
	case OPTIONS:
		return "OPTIONS"
	default:
		return ""
	}
}

// ServeHTTP is the wrapper method for the http.HandlerFunc
// this is for error handling
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var respBody []byte
	if err := h(w, r); err != nil {
		a, ok := err.(*logger.APIErr)
		if ok {
			respBody, _ = json.Marshal(a.RespBody)
		} else {
			respBody = []byte(`{"msg":"an error has occured"}`)
			a.Code = http.StatusInternalServerError
			a.Status = http.StatusText(a.Code)
		}
		w.Header().Set("Content-Type", ContentJSON)
		w.WriteHeader(a.Code)
		w.Write(respBody)
	}
}

func InitRoutes() Endpoints {
	var all Endpoints
	all = append(all,
		HomeRoute(),
		TestErrRoute(),
		PostRoute(),
	)

	return all
}

// SetupRoutes will take each endpoint and build
//  all the rest api endpoints, setting the handler func, path and method
func (ep Endpoints) SetupRoutes(r chi.Router) {
	for _, e := range ep {
		r.Method(e.Method.String(), e.Path, e.HandlerFunc)
	}
}

// Respond is a helper method to automate the data for the response
// sets the response code, content type, and body
func (e *Endpoint) Respond(w http.ResponseWriter) {
	var respBody []byte
	respBody, _ = json.Marshal(e.ResponseBody)
	if e.Pretty {
		respBody, _ = json.MarshalIndent(e.ResponseBody, "", "  ")
	}
	w.Header().Set("Content-Type", e.ResponseType)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
