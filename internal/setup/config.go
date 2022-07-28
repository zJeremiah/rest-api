package setup

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/logger"
)

// endpoint groups are sorted by these methods
//  for api documentation
const (
	GET Method = iota + 1
	HEAD
	POST
	PUT
	PATCH
	DELETE
	CONNECT
	OPTIONS
	TRACE
	ANY
)

const ContentJSON = "application/json"

type Method uint
type Methods []Method

// Endpoint object is used to setup an api endpoint in the api
// this is used to make the correct request to the handler
// and also used to generate the slate api docs
type Endpoint struct {
	Version      string      // the version pathing v1, v2, etc...
	Group        string      // the name of the group path, is part of the path /v1/name/{path}
	Path         string      // the URI path for the just the endpoint must start with a forward slash /
	FullPath     string      // (auto built) the full uri path for the endpoint /{version}/{group}/{path}
	Methods      Methods     // the HTTP methods to use for the endpoint GET, POST etc...
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
	ColorLog  bool            `flag:"color" toml:"color" json:"color" comment:"use linux coloring for request logs"`
	PrettyLog bool            `flag:"pretty" comment:"will pretty print request logs"`
	BuildDocs bool            `flag:"docs" comment:"flag to build the docs md file for slate"`
	mux       *chi.Mux
	Routes    Endpoints
}

var apiConfig *Config

func (ms Methods) First() (m Method) {
	if len(ms) > 0 {
		m = ms[0]
	}
	return m
}

func (rm Method) String() string {
	switch rm {
	case GET:
		return "GET"
	case HEAD:
		return "HEAD"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	case CONNECT:
		return "CONNECT"
	case OPTIONS:
		return "OPTIONS"
	case TRACE:
		return "TRACE"

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

		a.Status = http.StatusText(a.Code)
		respBody, _ = json.Marshal(a.RespBody)
		req.APIError = &a.Internal

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
			return fmt.Errorf("the path must start with a forward slash / %s %v (%s)",
				e.Name, e.Methods, e.Path)
		}

		urlParams := 0
		pStart := false
		for _, c := range e.Path {
			if c == '{' {
				pStart = true
			}
			if c == '}' {
				if !pStart {
					return fmt.Errorf("mis-matched url params %s %v (%s)",
						e.Name, e.Methods, e.Path)
				}
				pStart = false
				urlParams++
			}
		}
		if len(e.URLParams) != urlParams {
			return fmt.Errorf("params in path do not match URLParams %v %v (%s)",
				e.URLParams, e.Methods, e.FullPath)
		}
		if e.RequestBody != nil && len(e.JSONFields) == 0 {
			return fmt.Errorf("json request fields need to be defined in JSONFields %v (%s)",
				e.Methods, e.FullPath)
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

	// api documentation file serving
	apiConfig.mux.Get("/docs", Docs)
	apiConfig.mux.Get("/docs/*", Docs)

	// add the endpoints to the chi mux router
	for x, e := range apiConfig.Routes {
		// drop any trailing forward slashes on the path
		if e.FullPath[len(e.FullPath)-1:] == "/" && len(e.FullPath) > 1 {
			e.FullPath = e.FullPath[:len(e.FullPath)-1]
		}

		any := false
		for _, m := range e.Methods {
			if m == ANY {
				any = true
				e.Methods = make(Methods, 0)
				for i := 1; i < 10; i++ {
					method := Method(i).String()
					if method != "HEAD" && method != "DELETE" {
						e.Methods = append(e.Methods, Method(i))
						apiConfig.mux.Method(method, e.FullPath, e.HandlerFunc)
					}
				}
				apiConfig.Routes[x] = e
				break
			}
		}
		if !any {
			for _, m := range e.Methods {
				apiConfig.mux.Method(m.String(), e.FullPath, e.HandlerFunc)
			}
		}

		if apiConfig.Debug {
			log.Printf("adding route %v path: %s", e.Methods, e.FullPath)
		}
	}
}

// Mux returns the chi.mux (router)
func Mux() *chi.Mux {
	return apiConfig.mux
}

// AddEndpoint will add the endpoint list by {version}/{group}
// a full route is determined by {domain}/{version}/{group}/{endpoint.path}
// an error is returned if the endpoint has already been created
func AddEndpoints(ep ...Endpoint) error {
	for _, e := range ep {
		e.FullPath = path.Clean("/" + e.Version + "/" + e.Group + e.Path)
		id := fmt.Sprintf("%v %s", e.Methods, e.FullPath)
		_, found := apiConfig.Routes[id]
		if found {
			log.Fatalf("duplicate endpoint (%s) %s", e.FullPath, e.Description)
		}

		apiConfig.Routes[id] = e
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
