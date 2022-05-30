package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hydronica/go-config"
	"github.com/rest-api/docs"
	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
	"github.com/rest-api/internal/version"
	"github.com/rest-api/routes"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := &setup.Config{
		Log:    &logger.Options{},
		Port:   9876, // default port
		Routes: make(setup.Endpoints),
	}

	c.Log.StdOut(os.Stdout)
	config.New(c).Version(version.Get()).LoadOrDie()

	setup.InitConfig(c) // initialize the config
	routes.InitRoutes() // adds the routes to the setup

	// custom handler methods to avoid logging
	setup.Mux().MethodNotAllowed(NotAllowed)
	setup.Mux().NotFound(NotFound)

	// order matters, middleware is called in the order added
	setup.Mux().Use(middleware.Recoverer)
	setup.Mux().Use(middleware.Timeout(time.Minute))
	setup.Mux().Use(middleware.RequestID)
	setup.Mux().Use(Cors())
	setup.Mux().Use(middleware.StripSlashes)
	setup.Mux().Use(c.Log.WriteRequest)

	setup.AddRoutes()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "./docs/build"))
	docs.FileServer(setup.Mux(), "/docs", filesDir)

	c.Log.Pretty = c.PrettyLog

	if c.BuildDocs {
		err := docs.ParseTemplate()
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	log.Printf("running api on port %d", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), setup.Mux())
}

func NotAllowed(rw http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(logger.RequestKey).(*logger.Log)
	if ok {
		req.NoLog = true
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func NotFound(rw http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(logger.RequestKey).(*logger.Log)
	if ok {
		req.NoLog = true
	}
	rw.WriteHeader(http.StatusNotFound)
}

// Basic CORS
// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
func Cors() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{
			setup.GET.String(),
			setup.POST.String(),
			setup.PUT.String(),
			setup.DELETE.String(),
			setup.OPTIONS.String(),
			setup.HEAD.String(),
		},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
