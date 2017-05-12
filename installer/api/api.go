package api

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/dghubble/sessions"
)

const (
	installerSessionName = "tectonic-installer"
)

// ErrEmptyCookieSigningSecret is an error that occurs when a new API server is
// initialized without being provided a cookie signing secret.
var ErrEmptyCookieSigningSecret = errors.New("Empty cookie signing secret")

// Config configures an API server.
type Config struct {
	// If not "", serve assets from this local directory rather than from binassets
	AssetDir string

	// List of platform names to support
	Platforms []string

	// Whether the server was started with --dev
	DevMode bool

	// Cookie Sessions
	CookieSigningSecret string

	// Allow cookies to be sent over HTTP
	DisableSecureCookie bool
}

// Context contains data that can be used by the API handlers.
type Context struct {
	Sessions sessions.Store
	Config   *Config
}

// New initializes and returns a API server.
// TODO: Most methods should be GET, not POST!
func New(config *Config) (http.Handler, error) {
	if config.CookieSigningSecret == "" {
		return nil, ErrEmptyCookieSigningSecret
	}

	// Create a new client-side cookie-based session provider.
	sessions := sessions.NewCookieStore([]byte(config.CookieSigningSecret), nil)
	sessions.Config.HTTPOnly = false                     // allow Javascript access to cookies
	sessions.Config.Secure = !config.DisableSecureCookie // allow cookies to be sent over HTTP

	// Create a context.
	ctx := &Context{Sessions: sessions, Config: config}

	// Create the router.
	mux := http.NewServeMux()

	// handlers_frontend.go
	mux.Handle("/", frontendHandler(config.AssetDir, config.Platforms, config.DevMode))

	// handlers_aws.go
	mux.Handle("/aws/regions", httpHandler("POST", ctx, awsDescribeRegionsHandler))
	mux.Handle("/aws/default-subnets", httpHandler("POST", ctx, awsDefaultSubnetsHandler))
	mux.Handle("/aws/subnets/validate", httpHandler("POST", ctx, awsValidateSubnetsHandler))
	mux.Handle("/aws/vpcs", httpHandler("POST", ctx, awsGetVPCsHandler))
	mux.Handle("/aws/vpcs/subnets", httpHandler("POST", ctx, awsGetVPCsSubnetsHandler))
	mux.Handle("/aws/ssh-key-pairs", httpHandler("POST", ctx, awsGetKeyPairsHandler))
	mux.Handle("/aws/zones", httpHandler("POST", ctx, awsGetZonesHandler))
	mux.Handle("/aws/domain", httpHandler("POST", ctx, awsGetDomainInfoHandler))

	// handlers_terraform.go
	mux.Handle("/terraform/apply", httpHandler("POST", ctx, terraformApplyHandler))
	mux.Handle("/terraform/status", httpHandler("POST", ctx, terraformStatusHandler))
	mux.Handle("/terraform/assets", httpHandler("GET", ctx, terraformAssetsHandler))
	mux.Handle("/terraform/destroy", httpHandler("POST", ctx, terraformDestroyHandler))

	// handlers_containerlinux.go
	mux.Handle("/containerlinux/images/matchbox", httpHandler("GET", ctx, listMatchboxImagesHandler))
	mux.Handle("/containerlinux/images/amis", httpHandler("GET", ctx, listAMIImagesHandler))

	return mux, nil
}

// httpHandler wraps a handler with logging, and METHOD checker, and inject a
// Context.
func httpHandler(method string, context *Context, handler func(http.ResponseWriter, *http.Request, *Context) error) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logFields := log.Fields{
			"remote addr": request.RemoteAddr,
			"method":      request.Method,
			"request uri": request.RequestURI,
		}

		// Expect the right METHOD, or return http.StatusMethodNotAllowed.
		if request.Method != method {
			logFields["status"] = http.StatusMethodNotAllowed
			log.WithFields(logFields).Warning("Invalid METHOD for HTTP Request")

			writer.Header().Set("Allow", method)
			http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		// Execute the handler and log the request.
		// If an error occurred, write it.
		if err := handler(writer, request, context); err != nil {
			logFields["error"] = err.Error()

			if httpErr, ok := err.(*httpError); ok {
				http.Error(writer, httpErr.message, httpErr.status)
				log.WithFields(logFields).Error("Failed to handle HTTP Request")
			} else {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		log.WithFields(logFields).Info("Handled HTTP Request")
	})
}

// writeJSONResponse is an utilitarian function that writes a JSON response to the
// ResponseWriter for the Request, with the provided status and content.
func writeJSONResponse(writer http.ResponseWriter, _ *http.Request, status int, resp interface{}) error {
	// Headers must be written before the response.
	header := writer.Header()
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("Server", "Tectonic")

	// Write the response.
	writer.WriteHeader(status)
	err := json.NewEncoder(writer).Encode(resp)

	if err != nil {
		switch err.(type) {
		case *json.MarshalerError, *json.UnsupportedTypeError, *json.UnsupportedValueError:
			return newInternalServerError("Failed to marshal response: %s", err)
		default:
			return newInternalServerError("Failed to write response")
		}
	}
	return nil
}
