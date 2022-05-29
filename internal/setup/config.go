package setup

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/rest-api/internal/logger"
)

const (
	GET Method = iota
	HEAD
	OPTIONS
	PATCH
	POST
	PUT
	DELETE
)

const ContentJSON = "application/json"

type Method uint

// Endpoint object is used to setup an api endpoint in the api
// this is used to make the correct request to the handler
// and also used to generate the slate api docs
type Endpoint struct {
	Version      string      `json:"version"`       // the version pathing v1, v2, etc...
	Group        string      `json:"group"`         // the name of the group path, is part of the path /v1/name/{path}
	Path         string      `json:"path "`         // the URI path for the just the endpoint must start with a forward slash /
	Method       Method      `json:"method"`        // the HTTP method for the endpoint GET, POST etc...
	RequestType  string      `json:"request_type"`  // The request type is the ContentType for the request
	ResponseType string      `json:"response_type"` // This is the ContentType that will be returned in the response (normally application/json)
	RequestBody  interface{} `json:"request_body"`  // This is used to generate the api docs JSON object string for the request
	ResponseBody interface{} `json:"response_body"` // This is used to generate the api docs JSON object string for the response
	HandlerFunc  Handler     `json:"-"`             // the Handler function is called when the router matches the endpoint path
	RespFunc     Response    `json:"-"`             // a Response function is called to create an example of the endpoint response body (as json)
	ReqFunc      Request     `json:"-"`             // a Request function is called to produce an example of the endpoint request body (as json)
	Description  string      `json:"description"`   // The description of the endpoint for api docs
	Pretty       bool        `json:"pretty"`        // output the json string as pretty format when true
}

type Endpoints map[string]Endpoint
type Response func() string
type Request func() string
type Handler func(w http.ResponseWriter, r *http.Request) error

type Key string
type Config struct {
	Port      int             `toml:"port" json:"port" flag:"port" comment:"http port number"`
	Debug     bool            `toml:"debug" json:"debug" flag:"debug" comment:"show debug logging"`
	Log       *logger.Options `toml:"log_options" json:"log_options"`
	PrettyLog bool            `toml:"pretty" json:"pretty" flag:"pretty" comment:"will pretty print request logs"`

	router *chi.Mux
	Routes Endpoints
}

var singelton *Config

func (rm Method) String() string {
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
	req := r.Context().Value(logger.RequestKey).(*logger.Log)
	if err := h(w, r); err != nil {
		a, ok := err.(*logger.APIErr)
		if !ok {
			a = &logger.APIErr{
				Internal: logger.Internal{
					Msg:     "handler error",
					Err:     err,
					ErrText: err.Error(),
				},
				RespBody: logger.RespBody{
					Msg:  "an error has occured",
					Code: http.StatusInternalServerError,
				},
			}
		}
		respBody, _ = json.Marshal(a.RespBody)
		req.APIError = &a.Internal
		a.Status = http.StatusText(a.Code)

		w.Header().Set("Content-Type", ContentJSON)
		w.WriteHeader(a.Code)
		w.Write(respBody)
	}
}

// Validate is to verify that the groups of
// endpoint routes have been setup correctly
func Validate() error {
	for _, r := range singelton.Routes {
		if len(r.Path) == 0 || r.Path[:1] != "/" {
			return fmt.Errorf("the path must start with a forward slash / %s %s (%s)",
				r.Group, r.Method.String(), r.Path)
		}
	}

	return nil
}

// Routes will take each endpoint and build all path routes
//  endpoints, setting the handler func, path and method
// This should not need to be altered
func Routes() {
	if singelton == nil {
		log.Fatal("config singelton is nil")
	}

	for ke, e := range singelton.Routes {
		path := ke
		// drop any trailing forward slashes on the path
		if path[len(path)-1:] == "/" && len(path) > 1 {
			path = path[:len(path)-1]
		}

		if singelton.Debug {
			log.Printf("adding route %4s  path: %s", e.Method, path)
		}

		singelton.router.Method(e.Method.String(), path, e.HandlerFunc)
	}

}

func Debug() bool {
	if singelton == nil {
		log.Fatal("config singelton is nil")
	}
	return singelton.Debug
}

func InitConfig(c *Config) {
	if c == nil {
		log.Fatal("given config is nil")
	}
	singelton = c

	if singelton.Debug {
		log.Println("debug flag enabled")
	}

}

func SetRouter(r *chi.Mux) {
	if singelton == nil {
		log.Fatal("config singelton is nil")
	}
	singelton.router = r
}

func Router() *chi.Mux {
	return singelton.router
}

// Testing the possibility of auto unmarshaling the request body if needed
func (c *Config) ParseReqBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//r = r.WithContext(context.WithValue(r.Context(), Key("request_body"), req))
		next.ServeHTTP(rw, r)
	})
}

// AddEndpoint will add the endpoint /{version}/{group}/{ep.Path}
// an error is returned if the given group name is not found
func AddEndpoints(version, group string, ep ...Endpoint) error {
	for _, e := range ep {
		e.Version = version
		e.Group = group

		path := path.Clean("/" + e.Version + "/" + e.Group + e.Path)
		_, found := singelton.Routes[path]
		if found {
			return fmt.Errorf("endpoint already exists (%s) %s", path, e.Description)
		}
		singelton.Routes[path] = e
	}

	return nil
}
