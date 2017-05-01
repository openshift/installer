package server

import (
	"errors"
	"net/http"

	"github.com/dghubble/sessions"

	"github.com/coreos/tectonic-installer/installer/server/ctxh"
)

const (
	installerSessionName = "tectonic-installer"
)

var (
	errEmptyCookieSigningSecret = errors.New("installer: Empty cookie signing secret")
)

// Config configures a server.
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

// NewServer returns a new server handler.
func NewServer(config *Config) (http.Handler, error) {
	if config.CookieSigningSecret == "" {
		return nil, errEmptyCookieSigningSecret
	}

	// client-side cookie sessions
	sessionProvider := sessions.NewCookieStore([]byte(config.CookieSigningSecret), nil)
	// allow Javascript access to cookies
	sessionProvider.Config.HTTPOnly = false
	// allow cookies to be sent over HTTP
	if config.DisableSecureCookie {
		sessionProvider.Config.Secure = false
	} else {
		sessionProvider.Config.Secure = true
	}

	mux := http.NewServeMux()
	mux.Handle("/", logRequests(frontendHandler(config.AssetDir, config.Platforms, config.DevMode)))
	mux.Handle("/images", logRequests(ctxh.NewHandler(listImagesHandler())))
	mux.Handle("/cluster/create", logRequests(syncHandler(ctxh.NewHandler(createHandler(sessionProvider)))))
	mux.Handle("/cluster/status", logRequests(ctxh.NewHandler(statusHandler(sessionProvider))))
	mux.Handle("/cluster/done", logRequests(doneHandler(sessionProvider)))
	mux.Handle("/proxy", logRequests(ctxh.NewHandler(proxy())))
	mux.Handle("/aws/regions", logRequests(ctxh.NewHandler(awsDescribeRegionsHandler())))
	mux.Handle("/aws/default-subnets", logRequests(ctxh.NewHandler(awsDefaultSubnetsHandler())))
	mux.Handle("/aws/subnets/validate", logRequests(ctxh.NewHandler(awsValidateSubnetsHandler())))
	mux.Handle("/aws/kms/create", logRequests(ctxh.NewHandler(awsCreateKMSHandler())))
	mux.Handle("/aws/kms", logRequests(ctxh.NewHandler(awsGetKMSHandler())))
	mux.Handle("/aws/vpcs", logRequests(ctxh.NewHandler(awsGetVPCsHandler())))
	mux.Handle("/aws/vpcs/subnets", logRequests(ctxh.NewHandler(awsGetVPCsSubnetsHandler())))
	mux.Handle("/aws/ssh-key-pairs", logRequests(ctxh.NewHandler(awsGetKeyPairsHandler())))
	mux.Handle("/aws/zones", logRequests(ctxh.NewHandler(awsGetZonesHandler())))
	mux.Handle("/aws/domain", logRequests(ctxh.NewHandler(awsGetDomainInfoHandler())))
	mux.Handle("/terraform/apply", logRequests(ctxh.NewHandler(terraformApplyHandler(sessionProvider))))
	mux.Handle("/terraform/status", logRequests(ctxh.NewHandler(terraformStatusHandler(sessionProvider))))
	mux.Handle("/terraform/assets", logRequests(ctxh.NewHandler(terraformAssetsHandler(sessionProvider))))
	mux.Handle("/terraform/destroy", logRequests(ctxh.NewHandler(terraformDestroyHandler(sessionProvider))))

	return mux, nil
}
