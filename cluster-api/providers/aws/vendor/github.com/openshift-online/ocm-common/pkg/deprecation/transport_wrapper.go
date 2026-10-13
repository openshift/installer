package deprecation

//go:generate mockgen -destination=test/mock_roundtripper.go -package=test net/http RoundTripper

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
)

// TransportWrapper creates an HTTP transport wrapper that automatically detects and handles
// deprecation headers in all OCM API responses.
type TransportWrapper struct {
	transport http.RoundTripper
}

// NewTransportWrapper creates a new deprecation transport wrapper.
func NewTransportWrapper() func(http.RoundTripper) http.RoundTripper {
	return func(transport http.RoundTripper) http.RoundTripper {
		return &TransportWrapper{
			transport: transport,
		}
	}
}

// RoundTrip implements the http.RoundTripper interface and intercepts responses
// to check for deprecation headers.
func (w *TransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Execute the original request
	resp, err := w.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	w.handleDeprecationHeaders(resp)

	return resp, nil
}

// handleDeprecationHeaders checks for deprecation headers and prints warnings
func (w *TransportWrapper) handleDeprecationHeaders(resp *http.Response) {
	if resp == nil || resp.Header == nil {
		return
	}

	deprecationHeader := resp.Header.Get(consts.DeprecationHeader)
	messageHeader := resp.Header.Get(consts.OcmDeprecationMessage)

	if deprecationHeader != "" || messageHeader != "" {
		w.printDeprecationWarning(deprecationHeader, messageHeader)
	}
}

// printDeprecationWarning prints a deprecation warning to stderr
func (w *TransportWrapper) printDeprecationWarning(deprecationHeader, messageHeader string) {
	var message strings.Builder
	message.WriteString("WARNING: You are using a deprecated OCM API\n")

	if deprecationHeader != "" {
		// Try to parse the date from the header
		if deprecationTime, err := time.Parse(time.RFC3339, deprecationHeader); err == nil {
			message.WriteString(fmt.Sprintf("This endpoint will be removed on: %s\n",
				deprecationTime.Format(time.RFC3339)))
		} else if deprecationTime, err := time.Parse(time.RFC1123Z, deprecationHeader); err == nil {
			message.WriteString(fmt.Sprintf("This endpoint will be removed on: %s\n",
				deprecationTime.Format(time.RFC1123Z)))
		} else {
			message.WriteString(fmt.Sprintf("Deprecation: %s\n", deprecationHeader))
		}
	}

	if messageHeader != "" {
		message.WriteString(fmt.Sprintf("Details: %s\n", messageHeader))
	}

	message.WriteString("Please update your usage to avoid issues when this endpoint is removed\n")

	fmt.Fprint(os.Stderr, message.String())
}
