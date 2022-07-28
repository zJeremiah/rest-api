package setup

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rest-api/internal/logger"
)

func Docs(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.String()

	// redirect if full path is not given
	if !strings.Contains(urlStr, "/docs/") {
		http.Redirect(w, r, "/docs/#introduction", http.StatusMovedPermanently)
		return
	}

	// do not log api docs requests
	req, ok := r.Context().Value(logger.RequestKey).(*logger.Log)
	if ok {
		req.NoLog = true
	}

	workDir, _ := os.Getwd()
	root := http.Dir(filepath.Join(workDir, "docs/build"))
	rctx := chi.RouteContext(r.Context())
	rp := rctx.RoutePattern()
	prefix := strings.TrimSuffix(rp, "/*")
	fs := http.StripPrefix(prefix, http.FileServer(root))
	fs.ServeHTTP(w, r)
}
