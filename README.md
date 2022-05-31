# REST API base template

To add endpoints

- Add a new file.go in the routes that will define your group of endpoints
- The cleanest way I found is to write an endpoint func that will return an endpoint object
- Then write a Handler func to be setup in the Endpoint, which will do the work for the endpoint
- You need to give your endpoint a path, method, content-type, and if needed, a request body and/or the response body
  - The path must start with a forward slash '/' paths are build as follows
  - /{version}/{group}/{endpoint.path} (trailing slashes are dropped)
- You need to a setup func to add the endpoints to the api config
- You will also need to add your setup to the InitRoutes() that the config setup calls

### setup func example to add endpoints to the setup
```go
func Setup() {
	setup.AddEndpoints(
		GetKittensEP(),
		KittenEP(),
		RMKittenEP(),
		AddKittnEP(),
	)
}
```

### routes func that gets run by the config setup
```go
// Add all route initializations here
func InitRoutes() {
	root.Setup()
	kittns.Setup()
	// ... add setup functions here
}
```

- Write a func that returns an endpoint object to setup a new endpoint
- Path should always start with a forward slash '/' 
- Name, Path, and Method are required for every endpoint
- With no Version, or Group the root path is used
- If you have url params i.e., /{url_path_value}
	- URLParams are required to define the param (for the api documentation)

### endpoint func example
```go
func MyEndpoint() Endpoint {
	return setup.Endpoint{
		Name:         "Get a Specific Kittn",
		Version:      "v1",
		Group:        "kittns",
		Path:         "/{id}",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		ResponseBody: Kitten{ID: 1, Name: "Fluffums", Breed: "calico", Fluffy: 6, Cute: 7}, // example for docs
		Description:  "This endpoint retrieves a specific kittn",
		HandlerFunc:  GetKitten,
		URLParams: []setup.Param{
			{Name: "id", Description: "the id for a kittn"},
		},
	}
}
```

- Write a HandlerFunc to handle your request
  - then return the data needed for the endpoint
	- return an error if there is a problem
	  - the error handler will then return the error to the user and log the actual error

```go
// Home is an setup.Endpoint handler func example for the root / path
func Home(w http.ResponseWriter, r *http.Request) error {
	Respond(w, struct {
		AppName string         `json:"app_name"`
		Version version.Struct `json:"build_info"`
	}{
		AppName: "rest-api",
		Version: version.JSON(),
	}, true)
	return nil
}
```

# API Documentation
- the api docs are built from a markdown file
- the app will use a template to re-write the markdown
- slate will then generate api docs based on the markdown file
- run `make docs` for linux and `.\make.ps1 -docs` for windows
  - powershell is required on windows to run the windows script
  - docker compose is required to run the slate build
- after the docs are built, running the api will serve the docs
  - http://localhost:9876/docs
