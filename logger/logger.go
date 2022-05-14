package logger

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pcelvng/task-tools/file"
	"github.com/rest-api/internal/apierr"
)

// File log options
type Options struct {
	StdOut        io.Writer
	Writer        file.Writer
	Rotation      string `toml:"rotation" json:"rotation"`
	FilePath      string `toml:"file_path" json:"file_path"`
	*file.Options `toml:"file_options" json:"file_options"`
}

// Request represents an external request to the api
type Request struct {
	Host        string           `json:"host"`
	URI         string           `json:"request_uri"`
	Time        time.Time        `json:"request_time"`
	Body        interface{}      `json:"request_body,omitempty"`
	ContentLen  int64            `json:"content_length,omitempty"`
	Method      string           `json:"method"`
	RemoteAddr  string           `json:"remote_address"`
	UserAgent   string           `json:"user_agent,omitempty"`
	ContentType string           `json:"content_type,omitempty"`
	APIError    *apierr.APIError `json:"api_error,omitempty"`
	Latency     float64          `json:"latency"`
}

var json = jsoniter.ConfigFastest

func (o *Options) WriteRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		req := &Request{
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

		r = r.WithContext(context.WithValue(r.Context(), "request", req))
		next.ServeHTTP(rw, r)

		req.Latency = time.Since(start).Seconds()
		b, err := json.Marshal(req)
		if err != nil {
			log.Printf("error marshal request object %v", err)
		}

		_, err = o.StdOut.Write(b)
		if err != nil {
			log.Printf("error writing to stdout %v", err)
		}
		o.StdOut.Write([]byte("\n"))
	})
}

func retrieveBody(req *http.Request) (i interface{}) {
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
