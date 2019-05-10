package srv

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

func responseWasNotFound(resp autorest.Response) bool {
	return responseWasStatusCode(resp, http.StatusNotFound)
}

func responseWasStatusCode(resp autorest.Response, statusCode int) bool {
	if r := resp.Response; r != nil {
		if r.StatusCode == statusCode {
			return true
		}
	}

	return false
}
