package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/apierr"
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

type Endpoint struct {
	Path         string
	Method       RequestMethod
	RequestType  string
	ResponseType string
	RequestBody  interface{}
	ResponseBody interface{}
	HandlerFunc  Handler
	Description  string
	Pretty       bool
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
		a, ok := err.(*apierr.APIError)
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

func (ep Endpoints) SetupRoutes(r chi.Router) {
	for _, e := range ep {
		r.Method(e.Method.String(), e.Path, e.HandlerFunc)
	}
}

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
