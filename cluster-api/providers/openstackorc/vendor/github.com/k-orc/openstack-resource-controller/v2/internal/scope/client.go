package scope

import (
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/controller"
)

// RoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default http client RoundTripper
type RoundTripper struct {
	// Default http.RoundTripper
	http.RoundTripper
}

// RoundTrip performs a round-trip HTTP request, injecting the OpenStack
// Request ID header when appropriate
func (rt *RoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	reconcileID := controller.ReconcileIDFromContext(request.Context())
	if reconcileID != "" {
		request.Header.Set("X-OpenStack-Request-ID", "req-"+string(reconcileID))
	}

	return rt.RoundTripper.RoundTrip(request)
}
