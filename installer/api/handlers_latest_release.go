package api

import (
	"io"
	"net/http"
)

// latestReleaseHandler gets the tectonic's latest release from coreos.com.
func latestReleaseHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	res, err := http.Get("https://coreos.com/tectonic/api/releases/latest")
	if err != nil {
		return newBadRequestError("Failed to get a response from coreos.com: %s", err)
	}
	io.Copy(w, res.Body)
	res.Body.Close()
	return nil
}
