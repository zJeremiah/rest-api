package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/logger"
)

const (
	GET = iota
	HEAD
	OPTIONS
	PATCH
	POST
	PUT
	DELETE
)
const ContentJSON = "application/json"

type RequestMethod uint

type Group struct {
	Name        string    // the name of the group path, is part of the path /v1/name/{path}
	Version     string    // the version pathing v1, v2
	Description string    // a simple name to describe this group
	Routes      Endpoints // the routes in this group
}

type Groups []Group

// Endpoint object is used to setup an api endpoint in the api
// this is used to make the correct request to the handler
// and also used to generate the slate api docs
type Endpoint struct {
	Name         string        // this is the reference name for the endpoint
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

type Endpoints []Endpoint

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

// InitGroups will initialize the endpoints organized by group
// add new endpoint rout methods here to include them in the api request list
// these are auto documented based on the endpoint values
func InitGroups() Groups {
	var g Groups
	g = append(g, Group{
		Name:        "", // base path has no name
		Version:     "", // base path has no route
		Description: "base route",
		Routes:      Endpoints{HomeRoute()},
	})

	g = append(g, Group{
		Name:        "test",
		Version:     "v1",
		Description: "v1 test routes",
		Routes: Endpoints{
			TestErrRoute(),
			PostRoute(),
		},
	},
		Group{
			Name:        "test",
			Version:     "v2",
			Description: "v2 test routes",
			Routes: Endpoints{
				OtherRoute(),
				OtherTestErrRoute(),
				OtherPostRoute(),
			},
		})

	return g
}

// SetupRoutes will take each endpoint and build
//  all the rest api endpoints, setting the handler func, path and method
func (gs Groups) SetupRoutes(r chi.Router, debug bool) {
	for _, g := range gs {
		for _, e := range g.Routes {
			path := e.Path
			if g.Name != "" {
				path = "/" + g.Name + path
			}
			if g.Version != "" {
				path = "/" + g.Version + path
			}

			if debug {
				log.Println("adding route", path)
			}

			r.Method(e.Method.String(), path, e.HandlerFunc)
		}
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

// Validate is to verify that the groups of
// endpoint routes have been setup correctly
func Validate(groups Groups) error {
	for _, g := range groups {
		for _, e := range g.Routes {
			if e.Path[:1] != "/" {
				return fmt.Errorf("the path must start with a forward slash / %s %s %s (%s)",
					e.Name, e.Method.String(), g.Name, e.Path)
			}
		}
	}

	return nil
}
