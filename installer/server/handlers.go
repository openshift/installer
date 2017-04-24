package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/server/ctxh"
)

// logRequests logs HTTP requests and calls the next handler.
func logRequests(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		log.Debugf("HTTP %s %v", req.Method, req.URL)
		if next != nil {
			next.ServeHTTP(w, req)
		}
	}
	return http.HandlerFunc(fn)
}

// requireHTTPMethod will respond with HTTP 405, method not allowed if
// the http request is not the HTTP method provided.
func requireHTTPMethod(method string, next ctxh.ContextHandler) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		if req.Method != method {
			return ctxh.NewAppError(nil, fmt.Sprintf("request must be %s", method), http.StatusMethodNotAllowed)
		}
		next.ServeHTTP(ctx, w, req)
		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}
