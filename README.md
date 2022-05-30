# REST API base template

To add endpoints

- Add a new file.go in the routes that will define your group of endpoints
- The cleanest way I found is to write an endpoint func that will return an endpoint object
- You need to give your endpoint a path, method, content-type, and the response body
  - The request body is only needed on requests that have a request body
  - The path must start with a forward slash '/' paths are build as follows
  - /{version}/{group}/{endpoint.path} (trailing slashes are dropped)

```go
func MyRoute() Endpoint {
	e := Endpoint{
		Path:         "/",
		Method:       GET,
		ResponseType: ContentJSON,
		ResponseBody: Home{
			AppName: "rest-api",
			Version: version.JSON(),
		},
	}

	e.HandlerFunc = e.TestError
	return e
}
```

- Write a HandlerFunc to handle your request

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
