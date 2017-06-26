package api

import (
	"encoding/json"
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

	response := struct {
		TectonicConsole tectonic.ServiceStatus `json:"tectonicConsole"`
	}{
		TectonicConsole: tectonic.ConsoleHealth(nil, input.TectonicDomain),
	}

	return writeJSONResponse(w, req, http.StatusOK, response)
}
