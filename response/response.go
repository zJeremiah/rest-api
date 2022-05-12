package response

import "github.com/rest-api/internal/version"

type Home struct {
	AppName string `json:"app_name"`
	Version version.Struct
}
