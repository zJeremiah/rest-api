package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/hydronica/go-config"
	"github.com/rest-api/internal/setup"
	"github.com/rest-api/internal/version"
	"github.com/rest-api/logger"
	"github.com/rest-api/router"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := &setup.Config{
		Router: chi.NewRouter(),
		Log: logger.Options{
			StdOut: os.Stdout,
		},
		Endpoints: router.InitRoutes(),
		Port:      8080, // default port
	}

	config.New(c).Version(version.Get()).LoadOrDie()

	c.Router.Use(c.Log.WriteRequest)

	router.InitRoutes().SetupRoutes(c.Router)

	log.Printf("running api on port %d", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), c.Router)
}
