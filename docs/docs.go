package openapi

import "github.com/rest-api/internal/version"

type OpenApi map[string]interface{}

var openApi = make(OpenApi)

func init() {
	openApi["openapi"] = "3.1.0"
	openApi["info"] = map[string]string{
		"title":   "rest-api",
		"version": version.Version,
		"summary": "an api for requestion data"}

}
