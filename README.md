# REST API base boiler plate

To add endpoints

- Add a new file.go in the routes that will define your group of endpoints
- The cleanest way I found is to write an endpoint func that will return an endpoint object
- You need to give your endpoint a path, method, content-type, and the response body
  - The request body is only needed on requests that have a request body
  - The path must start with a forward slash '/' paths are build as follows
  - /{version}/{group}/{endpoint.path} (trailing slashes are dropped)
- Write a HandlerFunc to handle your request

```go
func MyRoute() Endpoint {
	e := Endpoint{
		Name:         "test error",
		Path:         "/error",
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

- In the router.go file add a group where you will list your endpoints

```go
Group{
			Name:        "test",
			Version:     "v2",
			Description: "v2 test routes",
			// each endpoint function is listed here to build the group routes
			Routes: Endpoints{
				MyRoute(),
				// ...
			},
		})
```
