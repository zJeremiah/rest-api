package routes

import (
	"github.com/rest-api/routes/kittns"
	"github.com/rest-api/routes/root"
)

// Add all route initializations here
func InitRoutes() {
	root.Setup()
	kittns.Setup()
}
