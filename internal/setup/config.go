package setup

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/logger"
)

// endpoint groups are sorted by these methods
//  for api documentation
const (
	GET Method = iota
	POST
	PUT
	DELETE
	HEAD
	OPTIONS
	PATCH
)

const ContentJSON = "application/json"

type Method uint

// Endpoint object is used to setup an api endpoint in the api
// this is used to make the correct request to the handler
// and also used to generate the slate api docs
type Endpoint struct {
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
	Pretty       bool        // output the json string as pretty format when true
	// These are used to define the api documentation
	Name        string  // (api docs) a simple statement for the endpoint
	Description string  // (api docs) The description of the endpoint for api docs
	QueryParams []Param // (api docs) listed query params
	URLParams   []Param // (api docs) listed url's path paramaters i.e., http://mydomain.com/{section}/{id}
	JSONFields  []Param // (api docs) listed JSON fields for POST/PUT request body
}

type Param struct {
	Name        string
	Description string
	Required    string // a value of yes or no here
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
			return fmt.Errorf("params in path do not match URLParams %v %s (%s)",
				e.URLParams, e.Method.String(), e.FullPath)
		}
		if e.RequestBody != nil && len(e.JSONFields) == 0 {
			return fmt.Errorf("json request fields need to be defined in JSONFields %s (%s)",
				e.Method.String(), e.FullPath)
		}
	}

	return nil
}

func Docs(w http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(logger.RequestKey).(*logger.Log)
	if ok {
		req.NoLog = true
	}

	workDir, _ := os.Getwd()
	root := http.Dir(filepath.Join(workDir, "docs/build"))
	rctx := chi.RouteContext(r.Context())
	prefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	fs := http.StripPrefix(prefix, http.FileServer(root))
	fs.ServeHTTP(w, r)
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

	// api documentation file serving
	apiConfig.mux.Get("/docs", Docs)
	apiConfig.mux.Get("/docs/*", Docs)

	// add the endpoints to the chi mux router
	for _, e := range apiConfig.Routes {
		// drop any trailing forward slashes on the path
		if e.FullPath[len(e.FullPath)-1:] == "/" && len(e.FullPath) > 1 {
			e.FullPath = e.FullPath[:len(e.FullPath)-1]
		}

		if apiConfig.Debug {
			log.Printf("adding route %4s  path: %s", e.Method, e.FullPath)
		}

		apiConfig.mux.Method(e.Method.String(), e.FullPath, e.HandlerFunc)
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
		_, found := apiConfig.Routes[e.Method.String()+" "+e.FullPath]
		if found {
			log.Fatalf("duplicate endpoint (%s) %s", e.FullPath, e.Description)
		}

		apiConfig.Routes[e.Method.String()+" "+e.FullPath] = e
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

var json = jsoniter.ConfigFastest

func (ep Endpoint) MarshalResp() string {
	b, err := json.MarshalIndent(ep.ResponseBody, "", " ")
	if err != nil {
		log.Fatalf("could not marshal_resp %s", err.Error())
	}
	return string(b)
}

func (ep Endpoint) MarshalReq() string {
	b, err := json.MarshalIndent(ep.RequestBody, "", " ")
	if err != nil {
		log.Fatalf("could not marshal_resp %s", err.Error())
	}
	return string(b)
}
