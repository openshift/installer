package server

import (
	"fmt"
	"net/http"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/dghubble/sessions"
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

// syncHandler synchronizes request handling and calls the next handler. Only
// one request may be processed by the handler chain at a time.
func syncHandler(next http.Handler) http.Handler {
	mu := &sync.Mutex{}
	fn := func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		if next != nil {
			next.ServeHTTP(w, req)
		}
		mu.Unlock()
	}
	return http.HandlerFunc(fn)
}

// doneHandler removes (expires) the session cookie, if present.
func doneHandler(sessionProvider sessions.Store) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		sessionProvider.Destroy(w, installerSessionName)
		fmt.Fprintf(w, "ok")
	}
	return http.HandlerFunc(fn)
}
