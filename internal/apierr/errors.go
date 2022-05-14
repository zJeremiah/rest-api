package apierr

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type APIError struct {
	Internal `json:"internal,omitempty"`
	RespBody `json:"response_body,omitempty"`
}

type RespBody struct {
	Msg    string `json:"message,omitempty"`
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
}
type Internal struct {
	Msg     string `json:"message,omitempty"`
	Err     error  `json:"error,omitempty"`
	ErrText string `json:"error_text,omitempty"`
}

var json = jsoniter.ConfigFastest

func (a *APIError) Error() string {
	body, _ := json.MarshalToString(a.RespBody)
	return body
}

func AddCTX(ctx context.Context, a *APIError) context.Context {
	return context.WithValue(ctx, "error", a)
}

func FromContext(ctx context.Context) (found bool, a APIError) {
	v := ctx.Value("error")
	if v == nil {
		return false, APIError{}
	}
	return true, v.(APIError)
}

func NewError(internal, external string, code int, err error) *APIError {
	a := &APIError{
		Internal: Internal{
			Msg: internal,
			Err: err,
		},
		RespBody: RespBody{
			Msg:    external,
			Code:   code,
			Status: http.StatusText(code),
		},
	}

	if err != nil {
		a.ErrText = err.Error()
	}
	return a
}

func (a *APIError) Write(w http.ResponseWriter) {
	body, _ := json.Marshal(a.RespBody)
	w.WriteHeader(a.Code)
	w.Write(body)
}
