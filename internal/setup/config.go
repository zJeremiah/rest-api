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
	Name         string      // a simple name for the endpoint (key value)
	Version      string      // the version pathing v1, v2, etc...
	Group        string      // the name of the group path, is part of the path /v1/name/{path}
	Path         string      // the URI path for the just the endpoint must start with a forward slash /
	FullPath     string      // (auto built) the full uri path for the endpoint /{version}/{group}/{path}
	Method       Method      // the HTTP method for the endpoint GET, POST etc...
	RequestType  string      // The request type is the ContentType for the request
	ResponseType string      // This is the ContentType that will be returned in the response (normally application/json)
	RequestBody  interface{} // This is used to generate the api docs JSON object string for the request
	ResponseBody interface{} // This is used to generate the api docs JSON object string for the response
	HandlerFunc  Handler     // the Handler function is called when the router matches the endpoint path
	RespFunc     Response    // a Response function is called to create an example of the endpoint response body (as json)
	ReqFunc      Request     // a Request function is called to produce an example of the endpoint request body (as json)
	Description  string      // The description of the endpoint for api docs
	Pretty       bool        // output the json string as pretty format when true
	QueryParams  []Param     // (api docs) listed query params
	URLParams    []Param     // (api docs) listed url's path paramaters i.e., http://mydomain.com/{section}/{id}
}

type Param struct {
	Name        string
	Description string
}

type Endpoints map[string]Endpoint
type Response func() string // should return a json string of the response body
// should return a json string of the request body
type Request func() string

// custom handler func for error handling
type Handler func(w http.ResponseWriter, r *http.Request) error

type Key string
type Config struct {
	Port      int             `toml:"port" json:"port" flag:"port" comment:"http port number"`
	Debug     bool            `toml:"debug" json:"debug" flag:"debug" comment:"show debug logging"`
	Log       *logger.Options `toml:"log_options" json:"log_options"`
	PrettyLog bool            `flag:"pretty" comment:"will pretty print request logs"`
	BuildDocs bool            `flag:"docs" comment:"flag to build the docs md file for slate"`
	mux       *chi.Mux
	Routes    Endpoints
}

var apiConfig *Config

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

// InitConifg sets up the internal configuration for the api
func InitConfig(c *Config) {
	if c == nil {
		log.Fatal("given config is nil")
	}
	apiConfig = c

	if apiConfig.Debug {
		log.Println("debug flag enabled")
	}

	apiConfig.mux = chi.NewRouter()
}

// ServeHTTP is the wrapper method for the http.HandlerFunc
// this is for marshaling the error handling
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
					Msg:  "an error has occured, please see request id: " + req.ID,
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

// validate is to verify that the apiConfig
// endpoints have been setup correctly
func validate() error {
	for _, e := range apiConfig.Routes {
		if len(e.Path) == 0 || e.Path[:1] != "/" {
			return fmt.Errorf("the path must start with a forward slash / %s %s (%s)",
				e.Name, e.Method.String(), e.Path)
		}

		urlParams := 0
		pStart := false
		for _, c := range e.Path {
			if c == '{' {
				pStart = true
			}
			if c == '}' {
				if !pStart {
					return fmt.Errorf("mis-matched url params %s %s (%s)",
						e.Name, e.Method.String(), e.Path)
				}
				pStart = false
				urlParams++
			}
		}
		if len(e.URLParams) != urlParams {
			return fmt.Errorf("params in path do not match url_params %v (%s)",
				e.URLParams, e.FullPath)
		}
	}

	return nil
}

// addRoutes will take each endpoint and build all path routes
//  endpoints, setting the handler func, path and method
// This must happen after all middleware has been initialized
func AddRoutes() {
	if apiConfig == nil {
		log.Fatal("config apiConfig is nil")
	}

	// validate the added endpoints
	err := validate()
	if err != nil {
		log.Fatalln(err)
	}

	// add the endpoints to the chi mux router
	for ke, e := range apiConfig.Routes {
		path := ke
		// drop any trailing forward slashes on the path
		if path[len(path)-1:] == "/" && len(path) > 1 {
			path = path[:len(path)-1]
		}

		if apiConfig.Debug {
			log.Printf("adding route %4s  path: %s", e.Method, path)
		}

		apiConfig.mux.Method(e.Method.String(), path, e.HandlerFunc)
	}
}

// Mux returns the chi.mux (router)
func Mux() *chi.Mux {
	return apiConfig.mux
}

// Testing the possibility of auto unmarshaling the request body if needed
func (c *Config) ParseReqBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//r = r.WithContext(context.WithValue(r.Context(), Key("request_body"), req))
		next.ServeHTTP(rw, r)
	})
}

// AddEndpoint will add the endpoint list by {version}/{group}
// a full route is determined by {domain}/{version}/{group}/{endpoint.path}
// an error is returned if the endpoint has already been created
func AddEndpoints(ep ...Endpoint) error {
	for _, e := range ep {
		e.FullPath = path.Clean("/" + e.Version + "/" + e.Group + e.Path)
		_, found := apiConfig.Routes[e.FullPath]
		if found {
			log.Fatalf("duplicate endpoint (%s) %s", e.FullPath, e.Description)
		}

		apiConfig.Routes[e.FullPath] = e
	}

	return nil
}

func EndpointsList() Endpoints {
	return apiConfig.Routes
}

func GetRoute(name string) (e Endpoint, found bool) {
	for _, v := range apiConfig.Routes {
		if v.Name == name {
			return v, true
		}
	}
	return e, false
}
