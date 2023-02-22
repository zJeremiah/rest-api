package setup

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/hydronica/go-openapi"
)

type OpenAPI struct {
	*openapi.OpenAPI
}

var docsSkipPaths = Skip{
	{Contains: "docs"},
	{Path: "/stats"},
	{Path: "/status"},
	{Contains: "version"},
	{Path: "/"},
}

//go:embed base.json
var base string

type SkipRoute struct {
	Contains string
	Path     string
}

type Skip []SkipRoute

func (sk Skip) Skipper(path string) bool {
	for _, s := range sk {
		if s.Contains != "" && strings.Contains(path, s.Contains) {
			return true
		}
		if s.Path != "" && s.Path == path {
			return true
		}
	}

	return false
}

func BuildDocs() error {
	o, err := openapi.NewFromJson(base)
	if err != nil {
		return err
	}
	oa := OpenAPI{o}

	for _, ep := range EndpointsList() {
		if docsSkipPaths.Skipper(ep.Path) {
			continue
		}
		for _, m := range ep.Methods {
			if ep.Group != "" {
				oa.AddTag(ep.Group, "")
			}

			ur, err := oa.AddRoute(ep.FullPath, strings.ToLower(m.String()), ep.Group, ep.Description, "")
			if err != nil {
				return fmt.Errorf("error adding route to docs %w", err)
			}
			for _, p := range ep.QueryParams {
				oa.AddParam(ur, openapi.RouteParam{
					Name:     p.Name,
					Desc:     p.Description,
					Location: "query",
					Type:     openapi.String,
					Required: p.Required,
				})
			}
			for _, p := range ep.PathParams {
				oa.AddParam(ur, openapi.RouteParam{
					Name:     p.Name,
					Desc:     p.Description,
					Location: "path",
					Type:     openapi.String,
					Required: p.Required,
				})
			}
			if ep.ResponseBody != nil {
				oa.AddRequest(ur, openapi.NewReqBody(openapi.Json, ep.Description, []openapi.ExampleObject{
					{Name: ep.Name, Example: ep.ResponseBody, Desc: ep.Description},
				}))
			}
		}

	}

	os.WriteFile("./swagger.json", oa.JSON(), os.ModePerm)

	return nil
}

func Docs(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.String()
	if strings.Contains(urlStr, "/swagger.json") {
		b, _ := os.ReadFile("swagger.json")
		w.Header().Set("Content-Type", ContentJSON)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}

	proxyURL, _ := url.Parse("http://swagger_ui:8080")
	proxyURL, _ = url.Parse("http://localhost:8080")

	proxy := httputil.ReverseProxy{Director: func(rNew *http.Request) {
		rNew.URL.Scheme = proxyURL.Scheme
		rNew.URL.Host = proxyURL.Host
		rNew.URL.Path = proxyURL.Path + r.URL.Path
		rNew.Host = proxyURL.Host
	}}

	proxy.ServeHTTP(w, r)

	// // redirect if full path is not given
	// if !strings.Contains(urlStr, "/docs/") {
	// 	http.Redirect(w, r, "/docs/#introduction", http.StatusMovedPermanently)
	// 	return
	// }

	// // do not log api docs requests
	// req, ok := r.Context().Value(logger.RequestKey).(*logger.Log)
	// if ok {
	// 	req.NoLog = true
	// }

	// workDir, _ := os.Getwd()
	// root := http.Dir(filepath.Join(workDir, "docs/build"))
	// rctx := chi.RouteContext(r.Context())
	// rp := rctx.RoutePattern()
	// prefix := strings.TrimSuffix(rp, "/*")
	// fs := http.StripPrefix(prefix, http.FileServer(root))
	// fs.ServeHTTP(w, r)
}
