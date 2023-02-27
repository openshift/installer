package ovirtclient

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// ExtraSettings are the optional settings for the oVirt connection.
//
// For future development, an interface named ExtraSettingsV2, V3, etc. will be added that incorporate this interface.
// This is done for backwards compatibility.
type ExtraSettings interface {
	// ExtraHeaders adds headers to the request.
	ExtraHeaders() map[string]string
	// Compression enables GZIP or DEFLATE compression on HTTP queries
	Compression() bool
	// Proxy returns the proxy server to use for connecting the oVirt engine. If none is set, the system settings are
	// used. If an empty string is passed, no proxy is used.
	Proxy() *string
}

// ExtraSettingsBuilder is a buildable version of ExtraSettings.
type ExtraSettingsBuilder interface {
	ExtraSettings

	// WithExtraHeaders adds extra headers to send along with each request.
	WithExtraHeaders(map[string]string) ExtraSettingsBuilder
	// WithCompression enables compression on HTTP requests.
	WithCompression() ExtraSettingsBuilder
	// WithProxy explicitly sets a proxy server to use for requests.
	WithProxy(string) ExtraSettingsBuilder
}

// NewExtraSettings creates a builder for ExtraSettings.
func NewExtraSettings() ExtraSettingsBuilder {
	return &extraSettings{}
}

type extraSettings struct {
	headers     map[string]string
	compression bool
	proxy       *string
}

func (e *extraSettings) ExtraHeaders() map[string]string {
	return e.headers
}

func (e *extraSettings) Compression() bool {
	return e.compression
}

func (e *extraSettings) Proxy() *string {
	return e.proxy
}

func (e *extraSettings) WithExtraHeaders(m map[string]string) ExtraSettingsBuilder {
	e.headers = m
	return e
}

func (e *extraSettings) WithCompression() ExtraSettingsBuilder {
	e.compression = true
	return e
}

func (e *extraSettings) WithProxy(addr string) ExtraSettingsBuilder {
	e.proxy = &addr
	return e
}

// New creates a new copy of the enhanced oVirt client. It accepts the following options:
//
//	url
//
// This is the oVirt engine URL. This must start with http:// or https:// and typically ends with /ovirt-engine/.
//
//	username
//
// This is the username for the oVirt engine. This must contain the profile separated with an @ sign. For example,
// admin@internal.
//
//	password
//
// This is the password for the oVirt engine. Other authentication mechanisms are not supported.
//
//	tls
//
// This is a TLSProvider responsible for supplying TLS configuration to the client. See below for a simple example.
//
//	logger
//
// This is an implementation of ovirtclientlog.Logger to provide logging.
//
//	extraSettings
//
// This is an implementation of the ExtraSettings interface, allowing for customization of headers and turning on
// compression.
//
// # TLS
//
// This library tries to follow best practices when it comes to connection security. Therefore, you will need to pass
// a valid implementation of the TLSProvider interface in the tls parameter. The easiest way to do this is calling
// the ovirtclient.TLS() function and then configuring the resulting variable with the following functions:
//
//	tls := ovirtclient.TLS()
//
//	// Add certificates from an in-memory byte slice. Certificates must be in PEM format.
//	tls.CACertsFromMemory(caCerts)
//
//	// Add certificates from a single file. Certificates must be in PEM format.
//	tls.CACertsFromFile("/path/to/file.pem")
//
//	// Add certificates from a directory. Optionally, regular expressions can be passed that must match the file
//	// names.
//	tls.CACertsFromDir("/path/to/certs", regexp.MustCompile(`\.pem`))
//
//	// Add system certificates
//	tls.CACertsFromSystem()
//
//	// Disable certificate verification. This is a bad idea.
//	tls.Insecure()
//
//	client, err := ovirtclient.New(
//	    url, username, password,
//	    tls,
//	    logger, extraSettings
//	)
//
// # Extra settings
//
// This library also supports customizing the connection settings. In order to stay backwards compatible the
// extraSettings parameter must implement the ovirtclient.ExtraSettings interface. Future versions of this library will
// add new interfaces (e.g. ExtraSettingsV2) to add new features without breaking compatibility.
func New(
	url string,
	username string,
	password string,
	tls TLSProvider,
	logger Logger,
	extraSettings ExtraSettings,
) (ClientWithLegacySupport, error) {
	return NewWithVerify(url, username, password, tls, logger, extraSettings, testConnection)
}

// NewWithVerify is equivalent to New, but allows customizing the verification function for the connection.
// Alternatively, a nil can be passed to disable connection verification.
func NewWithVerify(
	u string,
	username string,
	password string,
	tls TLSProvider,
	logger Logger,
	extraSettings ExtraSettings,
	verify func(connection Client) error,
) (ClientWithLegacySupport, error) {
	if err := validateURL(u); err != nil {
		return nil, wrap(err, EBadArgument, "invalid URL: %s", u)
	}
	if err := validateUsername(username); err != nil {
		return nil, wrap(err, EBadArgument, "invalid username: %s", username)
	}
	tlsConfig, err := tls.CreateTLSConfig()
	if err != nil {
		return nil, wrap(err, ETLSError, "failed to create TLS configuration")
	}

	proxyFunc, err := getProxyFunc(extraSettings)
	if err != nil {
		return nil, err
	}
	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           proxyFunc,
		},
	}

	client := &oVirtClient{
		&sync.Mutex{},
		nil,
		nil,
		httpClient,
		logger,
		u,
		username,
		password,
		tlsConfig,
		extraSettings,
		rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
		verify,
	}

	if err := client.Reconnect(); err != nil {
		return nil, err
	}

	return client, nil
}

func getProxyFunc(extraSettings ExtraSettings) (func(req *http.Request) (*url.URL, error), error) {
	proxyFunc := http.ProxyFromEnvironment
	if extraSettings == nil {
		return proxyFunc, nil
	}
	proxy := extraSettings.Proxy()
	if proxy == nil {
		return proxyFunc, nil
	}
	if *proxy != "" {
		u, err := url.Parse(*proxy)
		if err != nil {
			return nil, wrap(err, EBadArgument, "failed to parse proxy URL: %s", *proxy)
		}
		proxyFunc = func(req *http.Request) (*url.URL, error) {
			return u, nil
		}
	} else {
		proxyFunc = nil
	}
	return proxyFunc, nil
}

func processExtraSettings(
	extraSettings ExtraSettings,
	connBuilder *ovirtsdk4.ConnectionBuilder,
) error {
	if extraSettings == nil {
		connBuilder.ProxyFromEnvironment()
		return nil
	}
	if len(extraSettings.ExtraHeaders()) > 0 {
		connBuilder.Headers(extraSettings.ExtraHeaders())
	}
	if extraSettings.Compression() {
		connBuilder.Compress(true)
	}
	proxy := extraSettings.Proxy()
	if proxy == nil {
		connBuilder.ProxyFromEnvironment()
		return nil
	}
	if *proxy != "" {
		u, err := url.Parse(*proxy)
		if err != nil {
			return wrap(err, EBadArgument, "failed to parse proxy URL: %s", *proxy)
		}
		connBuilder.Proxy(u)
	}
	return nil
}

func testConnection(conn Client) error {
	return conn.Test()
}

func validateUsername(username string) error {
	usernameParts := strings.Split(username, "@")
	if len(usernameParts) < 2 {
		return newError(EBadArgument, "username must contain at least one @ sign (format should be admin@internal)")
	}
	user := strings.Join(usernameParts[:len(usernameParts)-1], "@")
	scope := usernameParts[len(usernameParts)-1]

	if len(user) == 0 {
		return newError(EBadArgument, "no user supplied before the @ sign in username (format should be admin@internal)")
	}
	if len(scope) == 0 {
		return newError(EBadArgument, "no user supplied after the @ sign in username (format should be admin@internal)")
	}
	return nil
}

func validateURL(url string) error {
	//goland:noinspection HttpUrlsUsage
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return newError(EBadArgument, "URL must start with http:// or https://")
	}
	return nil
}
