package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
)

func tectonicStatusHandler(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	// Read the input from the request's body.
	input := struct {
		TectonicDomain string `json:"tectonicDomain"`
	}{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return newBadRequestError("Could not unmarshal input: %s", err)
	}
	defer req.Body.Close()

	// Determine whether there is an execution environment already in the session
	if _, _, _, errCtx := restoreExecutionFromSession(req, ctx.Sessions, nil); errCtx != nil {
		// Error directly (rather than NewAppError, which logs) since the
		// frontend periodically calls this endpoint to advance screens
		http.Error(w, fmt.Sprintf("Could not find session data: %v", errCtx), http.StatusNotFound)
		return nil
	}

	response := struct {
		TectonicConsole tectonic.ServiceStatus `json:"tectonicConsole"`
	}{
		TectonicConsole: tectonic.ConsoleHealth(nil, input.TectonicDomain),
	}

	return writeJSONResponse(w, req, http.StatusOK, response)
}
