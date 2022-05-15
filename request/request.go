package request

type PostReq struct {
	ID    int               `json:"id"`
	Name  string            `json:"name"`
	Float float64           `json:"float"`
	Map   map[string]string `json:"map"`
	Array []interface{}     `json:"array"`
}
