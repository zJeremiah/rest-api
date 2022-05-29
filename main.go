package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hydronica/go-config"
	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
	"github.com/rest-api/internal/version"
	"github.com/rest-api/router"
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

	setup.InitConfig(c)
	setup.SetRouter(chi.NewRouter())
	setup.Router().MethodNotAllowed(func(rw http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(logger.Key("request")).(*logger.Log)
		req.NoLog = true
		rw.WriteHeader(http.StatusMethodNotAllowed)
	})

	setup.Router().NotFound(func(rw http.ResponseWriter, r *http.Request) {
		req, ok := r.Context().Value(logger.Key("request")).(*logger.Log)
		if ok {
			req.NoLog = true
		}
		rw.WriteHeader(http.StatusNotFound)
	})

	setup.Router().Use(c.Log.WriteRequest)
	setup.Router().Use(middleware.StripSlashes)
	setup.Router().Use(middleware.Timeout(time.Minute))
	setup.Router().Use(Cors())

	router.InitGroups()
	setup.Routes()

	err := setup.Validate()
	if err != nil {
		log.Fatalln(err)
	}

	c.Log.Pretty = c.PrettyLog

	log.Printf("running api on port %d", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), setup.Router())
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
