package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/version"
	"github.com/rest-api/response"
)

type Endpoints []Endpoint

type RequestMethod int

const (
	GET RequestMethod = iota
	POST
	HEAD
	PUT
	DELETE
)

type Endpoint struct {
	Path         string
	Method       RequestMethod
	ContentType  string
	ResponseType string
	RequestBody  interface{}
	ResponseBody interface{}
	HandlerFunc  http.HandlerFunc
	Description  string
}

var json = jsoniter.ConfigFastest

func InitRoutes() Endpoints {
	home := Endpoint{
		Path:         "/",
		Method:       GET,
		ContentType:  "application/json",
		ResponseType: "application/json",
		ResponseBody: response.Home{
			AppName: "rest-api",
			Version: version.JSON(),
		},
	}

	home.HandlerFunc = home.Welcome

	return Endpoints{home}
}

func (ep Endpoints) SetupRoutes(r chi.Router) {
	for _, e := range ep {
		switch e.Method {
		case GET:
			r.Get(e.Path, e.HandlerFunc)
		case POST:
			r.Post(e.Path, e.HandlerFunc)
		case PUT:
			r.Put(e.Path, e.HandlerFunc)
		case DELETE:
			r.Delete(e.Path, e.HandlerFunc)
		case HEAD:
			r.Head(e.Path, e.HandlerFunc)
		}
	}
}

// Welcome is the root endpoint example for the / path
func (e *Endpoint) Welcome(w http.ResponseWriter, r *http.Request) {
	respBody, err := json.MarshalIndent(e.ResponseBody, "", "  ")
	if err != nil {
		body, _ := json.Marshal(struct{ Msg string }{Msg: "you should never see this error"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}

	w.Header().Set("Content-Type", e.ResponseType)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
