package response

import "github.com/rest-api/internal/version"

type Home struct {
	AppName string `json:"app_name"`
	Version version.Struct
}

type PostResp struct {
	RID    int               `json:"id,omitempty"`
	RName  string            `json:"name,omitempty"`
	RFloat float64           `json:"field_float,omitempty"`
	RMap   map[string]string `json:"field_map,omitempty"`
	RArray []interface{}     `json:"array,omitempty"`
}
