package routes

import (
	"github.com/rest-api/routes/kittens"
	"github.com/rest-api/routes/root"
)

// Add all route initializations here
func InitRoutes() {
	root.Setup()
	kittens.Setup()
}
