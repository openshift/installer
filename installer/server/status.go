package server

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sessions"
	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/server/ctxh"
)

// statusHandler returns the status of the created cluster based on the
// cluster monitoring cookie.
func statusHandler(sessionProvider sessions.Store) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		// read cluster details from the session cookie
		session, err := sessionProvider.Get(req, installerSessionName)
		if err != nil {
			// Error directly (rather than NewAppError, which logs) since the
			// frontend periodically calls this endpoint to advance screens
			http.Error(w, fmt.Sprintf("failed to get session: %v", err), http.StatusNotFound)
			return nil
		}
		kind, ok := session.Values["kind"].(string)
		if !ok {
			return ctxh.NewAppError(nil, "malformed 'kind' in session", http.StatusBadRequest)
		}

		// deserialize the StatusChecker from the session cookie
		var checker StatusChecker

		switch kind {
		case "tectonic-metal":
			metalChecker := new(TectonicMetalChecker)
			if metalChecker, ok = session.Values["checker"].(*TectonicMetalChecker); ok {
				checker = metalChecker
			} else {
				return ctxh.NewAppError(nil, "malformed 'checker' in session", http.StatusBadRequest)
			}
		case "tectonic-aws":
			awsChecker := new(TectonicAWSChecker)
			if awsChecker, ok = session.Values["checker"].(*TectonicAWSChecker); ok {
				checker = awsChecker
			} else {
				return ctxh.NewAppError(nil, "malformed 'checker' in session", http.StatusBadRequest)
			}

		default:
			return ctxh.NewAppError(nil, "missing 'checker' in session", http.StatusBadRequest)
		}

		// Return the status JSON description for the cluster
		b, err := checker.Status()
		if err != nil {
			return ctxh.NewAppError(nil, "failed to check cluster health", http.StatusInternalServerError)
		}
		writeJSON(w, b)
		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}
