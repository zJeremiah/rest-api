package logger

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	jsoniter "github.com/json-iterator/go"
	"github.com/pcelvng/task-tools/file"
)

// File log options
type Options struct {
	stdOut io.Writer
	writer io.Writer

	//writer        file.Writer
	Rotation      string `toml:"rotation" json:"rotation"`
	FilePath      string `toml:"file_path" json:"file_path"`
	*file.Options `toml:"file_options" json:"file_options"`
	Pretty        bool `toml:"-" json:"-"`
	Color         bool `toml:"-" json:"-"`
	Debug         bool `toml:"-" json:"-"`
}

// Request represents an external request to the api
type APIErr struct {
	Internal `json:"internal,omitempty"`
	RespBody `json:"response_body,omitempty"`
}

type RespBody struct {
	Msg    string `json:"message,omitempty"`
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"` // you don't set the status value, that is derived from the status code
}

type Internal struct {
	Msg     string `json:"message,omitempty"`
	Err     error  `json:"-"`
	ErrText string `json:"error_text,omitempty"`
}

type Log struct {
	ID          string    `json:"id"`
	Host        string    `json:"host"`
	URI         string    `json:"request_uri"`
	Time        time.Time `json:"request_time"`
	Body        any       `json:"request_body,omitempty"`
	ContentLen  int64     `json:"content_length,omitempty"`
	Method      string    `json:"method"`
	RemoteAddr  string    `json:"remote_address"`
	UserAgent   string    `json:"user_agent,omitempty"`
	ContentType string    `json:"content_type,omitempty"`
	APIError    *Internal `json:"error,omitempty"`
	Latency     float64   `json:"latency"`
	RespCode    int       `json:"response_code"`
	Response    string    `json:"response"`
	NoLog       bool      `json:"-"` // will cancel the request log
}

type ctxRequestKey int

const RequestKey ctxRequestKey = 0

var json = jsoniter.ConfigFastest

func (o *Options) InitLogger() {
	if o.FilePath == "" {
		o.FilePath = "nop://"
	}
	w, err := file.NewWriter(o.FilePath, o.Options)
	if err != nil {
		log.Fatalf("could not create request log writer %+v", err)
	}
	o.writer = w

	if o.Debug {
		o.stdOut = os.Stdout
	} else {
		o.stdOut, _ = file.NewWriter("nop://", o.Options)
	}
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (o *Options) WriteRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rec := &StatusRecorder{
			ResponseWriter: rw,
			Status:         200,
		}

		id := r.Context().Value(middleware.RequestIDKey).(string)
		var err error
		start := time.Now()
		req := &Log{
			ID:          id,
			Host:        r.Host,
			URI:         r.RequestURI,
			Time:        time.Now().UTC(),
			Body:        retrieveBody(r),
			ContentLen:  r.ContentLength,
			Method:      r.Method,
			RemoteAddr:  r.RemoteAddr,
			UserAgent:   r.Header.Get("user-agent"),
			ContentType: r.Header.Get("content-type"),
		}

		r = r.WithContext(context.WithValue(r.Context(), RequestKey, req))
		next.ServeHTTP(rec, r)
		if req.NoLog {
			return // return without writing log output
		}
		var body []byte
		req.Latency = time.Since(start).Seconds()
		req.RespCode = rec.Status
		req.Response = http.StatusText(req.RespCode)
		if o.Pretty {
			body, err = json.MarshalIndent(req, "", "  ")
			if err != nil {
				log.Printf("error marshal request object %v", err)
			}
		} else {
			body, err = json.Marshal(req)
			if err != nil {
				log.Printf("error marshal request object %v", err)
			}
		}

		o.writer.Write(body)
		o.writer.Write([]byte("\n"))

		if o.Color {
			if req.APIError != nil || req.RespCode/200 != 1 {
				o.stdOut.Write([]byte("\033[31m"))
				o.stdOut.Write(body)
				o.stdOut.Write([]byte("\033[0m\n"))
			} else {
				o.stdOut.Write([]byte("\033[32m"))
				o.stdOut.Write(body)
				o.stdOut.Write([]byte("\033[0m\n"))
			}
		} else {
			o.stdOut.Write(body)
			o.stdOut.Write([]byte("\n"))
		}
	})
}

func retrieveBody(req *http.Request) (i any) {
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return "could not read request body " + err.Error()
	}

	if len(buf) == 0 {
		return nil
	}

	// one read closer for modifing and one to set back to the request
	b := io.NopCloser(bytes.NewBuffer(buf))
	req.Body = io.NopCloser(bytes.NewBuffer(buf))

	// read the request body
	body, err := io.ReadAll(b)
	if err != nil {
		return "could not read request body copy " + err.Error()
	}

	if len(body) == 0 {
		return ""
	}

	err = json.Unmarshal(body, &i)
	if err != nil {
		return "request body is not valid json"
	}

	return i
}

func (a *APIErr) Error() string {
	body, _ := json.MarshalToString(a.RespBody)
	return body
}

func FromContext(ctx context.Context) (found bool, a APIErr) {
	v := ctx.Value("error")
	if v == nil {
		return false, APIErr{}
	}
	return true, v.(APIErr)
}

// NewError with take the http request and set a the api error in the request log
// internal messages are logged, external messages are sent back in the response
func NewError(internal, external string, code int, err error) *APIErr {
	a := &APIErr{
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

func (a *APIErr) Write(w http.ResponseWriter) {
	body, _ := json.Marshal(a.RespBody)
	w.WriteHeader(a.Code)
	w.Write(body)
}
