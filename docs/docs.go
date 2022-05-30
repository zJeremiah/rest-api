package docs

import (
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
)

type EP struct {
	setup.Endpoint
	FullPath    string
	JsonResp    string // json string of the response body
	JsonReq     string // json string ofthe request body
	QueryParams []Param
	URLParams   []Param
}

type Param struct {
	Name        string
	Default     string
	Descritpion string
}

type Group struct {
	Name   string
	EPList []EP
}

func ParseTemplate() {
	t := template.New("api-docs")   // Create a template.
	t.ParseFiles("index.html.tmpl") // Parse template file.
}

func EPGroups() map[string][]EP {
	groups := make(map[string][]EP)
	// first organize the endpoints into groups
	for route, ep := range setup.EndpointsList() {
		g := path.Clean(ep.Version + "/" + ep.Group)

		epDoc := EP{
			Endpoint: ep,
			FullPath: route,
		}

		if ep.RespFunc != nil {
			epDoc.JsonResp = epDoc.RespFunc()
		}

		if ep.ReqFunc != nil {
			epDoc.JsonReq = epDoc.ReqFunc()
		}

		groups[g] = append(groups[g], epDoc)
	}

	return groups
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		req, ok := r.Context().Value(logger.RequestKey).(*logger.Log)
		if ok {
			req.NoLog = true
		}
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
