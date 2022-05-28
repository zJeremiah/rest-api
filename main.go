package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/hydronica/go-config"
	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
	"github.com/rest-api/internal/version"
	"github.com/rest-api/router"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := &setup.Config{
		Router:      chi.NewRouter(),
		Log:         logger.Options{},
		RouteGroups: router.InitGroups(),
		Port:        8080, // default port
	}
	c.Log.StdOut(os.Stdout)
	config.New(c).Version(version.Get()).LoadOrDie()

	if c.Debug {
		log.Println("debug flag enabled")
	}

	c.Router.Use(c.Log.WriteRequest)

	c.RouteGroups.SetupRoutes(c.Router, c.Debug)
	err := router.Validate(c.RouteGroups)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("running api on port %d", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), c.Router)
}
