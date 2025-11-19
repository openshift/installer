/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains the implementations of the Builder and Connection objects.

package sdk

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/openshift-online/ocm-sdk-go/accesstransparency"
	"github.com/openshift-online/ocm-sdk-go/accountsmgmt"
	"github.com/openshift-online/ocm-sdk-go/addonsmgmt"
	"github.com/openshift-online/ocm-sdk-go/arohcp"
	"github.com/openshift-online/ocm-sdk-go/authentication"
	"github.com/openshift-online/ocm-sdk-go/authorizations"
	"github.com/openshift-online/ocm-sdk-go/clustersmgmt"
	"github.com/openshift-online/ocm-sdk-go/configuration"
	"github.com/openshift-online/ocm-sdk-go/internal"
	"github.com/openshift-online/ocm-sdk-go/jobqueue"
	"github.com/openshift-online/ocm-sdk-go/logging"
	"github.com/openshift-online/ocm-sdk-go/metrics"
	"github.com/openshift-online/ocm-sdk-go/osdfleetmgmt"
	"github.com/openshift-online/ocm-sdk-go/retry"
	"github.com/openshift-online/ocm-sdk-go/servicelogs"
	"github.com/openshift-online/ocm-sdk-go/servicemgmt"
	"github.com/openshift-online/ocm-sdk-go/statusboard"
	"github.com/openshift-online/ocm-sdk-go/webrca"
)

// Default values:
const (
	// #nosec G101
	DefaultTokenURL     = authentication.DefaultTokenURL
	DefaultClientID     = authentication.DefaultClientID
	DefaultClientSecret = authentication.DefaultClientSecret
	DefaultURL          = "https://api.openshift.com"
	DefaultAgent        = "OCM-SDK/" + Version
	FedRAMPURL          = "https://api.openshiftusgov.com"
)

// DefaultScopes is the ser of scopes used by default:
var DefaultScopes = []string{
	"openid",
}

// ConnectionBuilder contains the configuration and logic needed to create connections to
// `api.openshift.com`. Don't create instances of this type directly, use the NewConnectionBuilder
// function instead.
type ConnectionBuilder struct {
	// Basic attributes:
	logger            logging.Logger
	trustedCAs        []interface{}
	insecure          bool
	disableKeepAlives bool
	tokenURL          string
	clientID          string
	clientSecret      string
	urlTable          map[string]string
	agent             string
	user              string
	password          string
	tokens            []string
	scopes            []string
	retryLimit        int
	retryInterval     time.Duration
	retryJitter       float64
	transportWrappers []func(http.RoundTripper) http.RoundTripper

	includeDefaultAuthnTransportWrapper bool

	// Metrics:
	metricsSubsystem  string
	metricsRegisterer prometheus.Registerer

	// Error detected while populating the builder. Once set calls to methods to
	// set other builder parameters will be ignored and the Build method will
	// exit inmediately returning this error.
	err error
}

// TransportWrapper is a wrapper for a transport of type http.RoundTripper. Creating a transport
// wrapper, enables to preform actions and manipulations on the transport request and response.
type TransportWrapper func(http.RoundTripper) http.RoundTripper

// Connection contains the data needed to connect to the `api.openshift.com`. Don't create instances
// of this type directly, use the builder instead.
type Connection struct {
	// Basic attributes:
	closed         bool
	logger         logging.Logger
	authnWrapper   *authentication.TransportWrapper
	retryWrapper   *retry.TransportWrapper
	clientSelector *internal.ClientSelector
	urlTable       []urlTableEntry
	agent          string

	// Metrics:
	metricsSubsystem  string
	metricsRegisterer prometheus.Registerer
}

// urlTableEntry is used to store one entry of the table that contains the correspondence between
// path prefixes and base URLs.
type urlTableEntry struct {
	prefix string
	re     *regexp.Regexp
	url    *internal.ServerAddress
}

// NewConnectionBuilder creates an builder that knows how to create connections with the default
// configuration.
func NewConnectionBuilder() *ConnectionBuilder {
	return &ConnectionBuilder{
		urlTable: map[string]string{
			"": DefaultURL,
		},
		retryLimit:                          retry.DefaultLimit,
		retryInterval:                       retry.DefaultInterval,
		retryJitter:                         retry.DefaultJitter,
		metricsRegisterer:                   prometheus.DefaultRegisterer,
		includeDefaultAuthnTransportWrapper: true,
	}
}

// NewConnectionBuilder creates a Builder that knows how to create connections
// without authentication
func NewUnauthenticatedConnectionBuilder() *ConnectionBuilder {
	connectionBuilder := NewConnectionBuilder()
	connectionBuilder.includeDefaultAuthnTransportWrapper = false
	return connectionBuilder
}

// Logger sets the logger that will be used by the connection. By default it uses the Go `log`
// package, and with the debug level disabled and the rest enabled. If you need to change that you
// can create a logger and pass it to this method. For example:
//
//	// Create a logger with the debug level enabled:
//	logger, err := logging.NewGoLoggerBuilder().
//		Debug(true).
//		Build()
//	if err != nil {
//		panic(err)
//	}
//
//	// Create the connection:
//	cl, err := client.NewConnectionBuilder().
//		Logger(logger).
//		Build()
//	if err != nil {
//		panic(err)
//	}
//
// You can also build your own logger, implementing the Logger interface.
func (b *ConnectionBuilder) Logger(logger logging.Logger) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.logger = logger
	return b
}

// TokenURL sets the URL that will be used to request OpenID access tokens. The default is
// `https://sso.redhat.com/auth/realms/cloud-services/protocol/openid-connect/token`.
func (b *ConnectionBuilder) TokenURL(url string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.tokenURL = url
	return b
}

// Client sets OpenID client identifier and secret that will be used to request OpenID tokens. The
// default identifier is `cloud-services`. The default secret is the empty string. When these two
// values are provided and no user name and password is provided, the connection will use the client
// credentials grant to obtain the token. For example, to create a connection using the client
// credentials grant do the following:
//
//	// Use the client credentials grant:
//	connection, err := sdk.NewConnectionBuilder().
//		Client("myclientid", "myclientsecret").
//		Build()
//
// Note that some OpenID providers (Keycloak, for example) require the client identifier also for
// the resource owner password grant. In that case use the set only the identifier, and let the
// secret blank. For example:
//
//	// Use the resource owner password grant:
//	connection, err := sdk.NewConnectionBuilder().
//		User("myuser", "mypassword").
//		Client("myclientid", "").
//		Build()
//
// Note the empty client secret.
func (b *ConnectionBuilder) Client(id string, secret string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.clientID = id
	b.clientSecret = secret
	return b
}

// URL sets the base URL of the API gateway. The default is `https://api.openshift.com`.
//
// To connect using a Unix sockets and HTTP use the `unix` URL scheme and put the name of socket file
// in the URL path:
//
//	connection, err := sdk.NewConnectionBuilder().
//		URL("unix://my.server.com/tmp/api.socket").
//		Build()
//
// To connect using Unix sockets and HTTPS use `unix+https://my.server.com/tmp/api.socket`.
//
// To force use of HTTP/2 without TLS use `h2c://...`. This can also be combined with Unix sockets,
// for example `unix+h2c://...`.
//
// Note that the host name is mandatory even when using Unix sockets because it is used to populate
// the `Host` header sent to the server.
func (b *ConnectionBuilder) URL(url string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	return b.AlternativeURL("", url)
}

// AlternativeURL sets an alternative base URL for the given path prefix. For example, to configure
// the connection so that it sends the requests for the clusters management service to
// `https://my.server.com`:
//
//	connection, err := client.NewConnectionBuilder().
//		URL("https://api.example.com").
//		AlternativeURL("/api/clusters_mgmt", "https://my.server.com").
//		Build()
//
// Requests for other paths that don't start with the given prefix will still be sent to the default
// base URL.
//
// This method can be called multiple times to set alternative URLs for multiple prefixes.
func (b *ConnectionBuilder) AlternativeURL(prefix, base string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.urlTable[prefix] = base
	return b
}

// AlternativeURLs sets an collection of alternative base URLs. For example, to configure the
// connection so that it sends the requests for the clusters management service to
// `https://my.server.com` and the requests for the accounts management service to
// `https://your.server.com`:
//
//	connection, err := client.NewConnectionBuilder().
//		URL("https://api.example.com").
//		AlternativeURLs(map[string]string{
//			"/api/clusters_mgmt": "https://my.server.com",
//			"/api/accounts_mgmt": "https://your.server.com",
//		}).
//		Build()
//
// The effect is the same as calling the AlternativeURL multiple times.
func (b *ConnectionBuilder) AlternativeURLs(entries map[string]string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	for prefix, base := range entries {
		b.urlTable[prefix] = base
	}
	return b
}

// Agent sets the `User-Agent` header that the client will use in all the HTTP requests. The default
// is `OCM` followed by an slash and the version of the client, for example `OCM/0.0.0`.
func (b *ConnectionBuilder) Agent(agent string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.agent = agent
	return b
}

// User sets the user name and password that will be used to request OpenID access tokens. When
// these two values are provided the connection will use the resource owner password grant type to
// obtain the token. For example:
//
//	// Use the resource owner password grant:
//	connection, err := sdk.NewConnectionBuilder().
//		User("myuser", "mypassword").
//		Build()
//
// Note that some OpenID providers (Keycloak, for example) require the client identifier also for
// the resource owner password grant. In that case use the set only the identifier, and let the
// secret blank. For example:
//
//	// Use the resource owner password grant:
//	connection, err := sdk.NewConnectionBuilder().
//		User("myuser", "mypassword").
//		Client("myclientid", "").
//		Build()
//
// Note the empty client secret.
func (b *ConnectionBuilder) User(name string, password string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.user = name
	b.password = password
	return b
}

// Scopes sets the OpenID scopes that will be included in the token request. The default is to use
// the `openid` scope. If this method is used then that default will be completely replaced, so you
// will need to specify it explicitly if you want to use it. For example, if you want to add the
// scope 'myscope' without loosing the default you will have to do something like this:
//
//	// Create a connection with the default 'openid' scope and some additional scopes:
//	connection, err := sdk.NewConnectionBuilder().
//		User("myuser", "mypassword").
//		Scopes("openid", "myscope", "yourscope").
//		Build()
//
// If you just want to use the default 'openid' then there is no need to use this method.
func (b *ConnectionBuilder) Scopes(values ...string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.scopes = make([]string, len(values))
	copy(b.scopes, values)
	return b
}

// Tokens sets the OpenID tokens that will be used to authenticate. Multiple types of tokens are
// accepted, and used according to their type. For example, you can pass a single access token, or
// an access token and a refresh token, or just a refresh token. If no token is provided then the
// connection will the user name and password or the client identifier and client secret (see the
// User and Client methods) to request new ones.
//
// If the connection is created with these tokens and no user or client credentials, it will
// stop working when both tokens expire. That can happen, for example, if the connection isn't used
// for a period of time longer than the life of the refresh token.
func (b *ConnectionBuilder) Tokens(tokens ...string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.tokens = append(b.tokens, tokens...)
	return b
}

// TrustedCAs sets the certificate pool that contains the certificate authorities that will be
// trusted by the connection. If this isn't explicitly specified then the client will trust the
// certificate authorities trusted by default by the system.
func (b *ConnectionBuilder) TrustedCAs(value *x509.CertPool) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.trustedCAs = append(b.trustedCAs, value)
	return b
}

// TrustedCAFile sets the name of a file that contains the certificate authorities that will be
// trusted by the connection. If this isn't explicitly specified then the client will trust the
// certificate authorities trusted by default by the system.
func (b *ConnectionBuilder) TrustedCAFile(value string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.trustedCAs = append(b.trustedCAs, value)
	return b
}

// Insecure enables insecure communication with the server. This disables verification of TLS
// certificates and host names and it isn't recommended for a production environment.
func (b *ConnectionBuilder) Insecure(flag bool) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.insecure = flag
	return b
}

// DisableKeepAlives disables HTTP keep-alives with the server. This is unrelated to similarly
// named TCP keep-alives.
func (b *ConnectionBuilder) DisableKeepAlives(flag bool) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.disableKeepAlives = flag
	return b
}

// RetryLimit sets the maximum number of retries for a request. When this is zero no retries will be
// performed. The default value is two.
func (b *ConnectionBuilder) RetryLimit(value int) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.retryLimit = value
	return b
}

// RetryInterval sets the time to wait before the first retry. The interval time will be doubled for
// each retry. For example, if this is set to one second then the first retry will happen
// approximately one second after the failure of the initial request, the second retry will happen
// affer four seconds, the third will happen after eitght seconds, so on.
func (b *ConnectionBuilder) RetryInterval(value time.Duration) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.retryInterval = value
	return b
}

// RetryJitter sets a factor that will be used to randomize the retry intervals. For example, if
// this is set to 0.1 then a random adjustment between -10% and +10% will be done to the interval
// for each retry.  This is intended to reduce simultaneous retries by clients when a server starts
// failing.  The default value is 0.2.
func (b *ConnectionBuilder) RetryJitter(value float64) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.retryJitter = value
	return b
}

// TransportWrapper allows setting a transport layer into the connection for capturing and
// manipulating the request or response.
func (b *ConnectionBuilder) TransportWrapper(value TransportWrapper) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.transportWrappers = append(b.transportWrappers, value)
	return b
}

// MetricsSubsystem sets the name of the subsystem that will be used by the connection to register
// metrics with Prometheus. If this isn't explicitly specified, or if it is an empty string, then no
// metrics will be registered.  For example, if the value is `api_outbound` then the following
// metrics will be registered:
//
//	api_outbound_request_count - Number of API requests sent.
//	api_outbound_request_duration_sum - Total time to send API requests, in seconds.
//	api_outbound_request_duration_count - Total number of API requests measured.
//	api_outbound_request_duration_bucket - Number of API requests organized in buckets.
//	api_outbound_token_request_count - Number of token requests sent.
//	api_outbound_token_request_duration_sum - Total time to send token requests, in seconds.
//	api_outbound_token_request_duration_count - Total number of token requests measured.
//	api_outbound_token_request_duration_bucket - Number of token requests organized in buckets.
//
// The duration buckets metrics contain an `le` label that indicates the upper bound. For example if
// the `le` label is `1` then the value will be the number of requests that were processed in less
// than one second.
//
// The API request metrics have the following labels:
//
//	method - Name of the HTTP method, for example GET or POST.
//	path - Request path, for example /api/clusters_mgmt/v1/clusters.
//	code - HTTP response code, for example 200 or 500.
//
// To calculate the average request duration during the last 10 minutes, for example, use a
// Prometheus expression like this:
//
//	rate(api_outbound_request_duration_sum[10m]) / rate(api_outbound_request_duration_count[10m])
//
// In order to reduce the cardinality of the metrics the path label is modified to remove the
// identifiers of the objects. For example, if the original path is .../clusters/123 then it will
// be replaced by .../clusters/-, and the values will be accumulated. The line returned by the
// metrics server will be like this:
//
//	api_outbound_request_count{code="200",method="GET",path="/api/clusters_mgmt/v1/clusters/-"} 56
//
// The meaning of that is that there were a total of 56 requests to get specific clusters,
// independently of the specific identifier of the cluster.
//
// The token request metrics will contain the following labels:
//
//	code - HTTP response code, for example 200 or 500.
//
// The value of the `code` label will be zero when sending the request failed without a response
// code, for example if it wasn't possible to open the connection, or if there was a timeout waiting
// for the response.
//
// Note that setting this attribute is not enough to have metrics published, you also need to
// create and start a metrics server, as described in the documentation of the Prometheus library.
func (b *ConnectionBuilder) MetricsSubsystem(value string) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	b.metricsSubsystem = value
	return b
}

// MetricsRegisterer sets the Prometheus registerer that will be used to register the metrics. The
// default is to use the default Prometheus registerer and there is usually no need to change that.
// This is intended for unit tests, where it is convenient to have a registerer that doesn't
// interfere with the rest of the system.
func (b *ConnectionBuilder) MetricsRegisterer(value prometheus.Registerer) *ConnectionBuilder {
	if b.err != nil {
		return b
	}
	if value == nil {
		value = prometheus.DefaultRegisterer
	}
	b.metricsRegisterer = value
	return b
}

// Metrics sets the name of the subsystem that will be used by the connection to register metrics
// with Prometheus.
//
// Deprecated: has been replaced by MetricsSubsystem.
func (b *ConnectionBuilder) Metrics(value string) *ConnectionBuilder {
	return b.MetricsSubsystem(value)
}

// Load loads the connection configuration from the given source. The source must be a YAML
// document with content similar to this:
//
//	url: https://my.server.com
//	alternative_urls:
//	- /api/clusters_mgmt: https://your.server.com
//	- /api/accounts_mgmt: https://her.server.com
//	token_url: https://openid.server.com
//	user: myuser
//	password: mypassword
//	client_id: myclient
//	client_secret: mysecret
//	tokens:
//	- eY...
//	- eY...
//	scopes:
//	- openid
//	insecure: false
//	trusted_cas:
//	- /my/ca.pem
//	- /your/ca.pem
//	agent: myagent
//	retry: true
//	retry_limit: 1
//
// Setting any of these fields in the file has the same effect that calling the corresponding method
// of the builder.
//
// For details of the supported syntax see the documentation of the configuration package.
func (b *ConnectionBuilder) Load(source interface{}) *ConnectionBuilder {
	if b.err != nil {
		return b
	}

	// Load the configuration:
	var config *configuration.Object
	config, b.err = configuration.New().
		Load(source).
		Build()
	if b.err != nil {
		return b
	}
	var view struct {
		URL              *string           `yaml:"url"`
		AlternativeURLs  map[string]string `yaml:"alternative_urls"`
		TokenURL         *string           `yaml:"token_url"`
		User             *string           `yaml:"user"`
		Password         *string           `yaml:"password"`
		ClientID         *string           `yaml:"client_id"`
		ClientSecret     *string           `yaml:"client_secret"`
		Tokens           []string          `yaml:"tokens"`
		Insecure         *bool             `yaml:"insecure"`
		TrustedCAs       []string          `yaml:"trusted_cas"`
		Scopes           []string          `yaml:"scopes"`
		Agent            *string           `yaml:"agent"`
		Retry            *bool             `yaml:"retry"`
		RetryLimit       *int              `yaml:"retry_limit"`
		MetricsSubsystem *string           `yaml:"metrics_subsystem"`
	}
	b.err = config.Populate(&view)
	if b.err != nil {
		return b
	}

	// URL:
	if view.URL != nil {
		b.URL(*view.URL)
	}
	if view.TokenURL != nil {
		b.TokenURL(*view.TokenURL)
	}

	// Alternative URLs:
	if view.AlternativeURLs != nil {
		for prefix, base := range view.AlternativeURLs {
			b.AlternativeURL(prefix, base)
		}
	}

	// User and password:
	var user string
	var password string
	if view.User != nil {
		user = *view.User
	}
	if view.Password != nil {
		password = *view.Password
	}
	if user != "" || password != "" {
		b.User(user, password)
	}

	// Client identifier and secret:
	var clientID string
	var clientSecret string
	if view.ClientID != nil {
		clientID = *view.ClientID
	}
	if view.ClientSecret != nil {
		clientSecret = *view.ClientSecret
	}
	if clientID != "" || clientSecret != "" {
		b.Client(clientID, clientSecret)
	}

	// Tokens:
	if view.Tokens != nil {
		b.Tokens(view.Tokens...)
	}

	// Scopes:
	if view.Scopes != nil {
		b.Scopes(view.Scopes...)
	}

	// Insecure:
	if view.Insecure != nil {
		b.Insecure(*view.Insecure)
	}

	// Trusted CAs:
	for _, trustedCA := range view.TrustedCAs {
		b.TrustedCAFile(trustedCA)
	}

	// Agent:
	if view.Agent != nil {
		b.Agent(*view.Agent)
	}

	// Retry:
	if view.RetryLimit != nil {
		b.RetryLimit(*view.RetryLimit)
	}

	// Metrics subsystem:
	if view.MetricsSubsystem != nil {
		b.MetricsSubsystem(*view.MetricsSubsystem)
	}

	return b
}

// Build uses the configuration stored in the builder to create a new connection. The builder can be
// reused to create multiple connections with the same configuration. It returns a pointer to the
// connection, and an error if something fails when trying to create it.
//
// This operation is potentially lengthy, as it may require network communications. Consider using a
// context and the BuildContext method.
func (b *ConnectionBuilder) Build() (connection *Connection, err error) {
	return b.BuildContext(context.Background())
}

// BuildContext uses the configuration stored in the builder to create a new connection. The builder
// can be reused to create multiple connections with the same configuration. It returns a pointer to
// the connection, and an error if something fails when trying to create it.
func (b *ConnectionBuilder) BuildContext(ctx context.Context) (connection *Connection, err error) {
	// If an error has been detected while populating the builder then return it and finish:
	if b.err != nil {
		err = b.err
		return
	}

	// Create the default logger, if needed:
	if b.logger == nil {
		b.logger, err = logging.NewGoLoggerBuilder().
			Debug(false).
			Info(true).
			Warn(true).
			Error(true).
			Build()
		if err != nil {
			err = fmt.Errorf("can't create default logger: %w", err)
			return
		}
		b.logger.Debug(ctx, "Logger wasn't provided, will use Go log")
	}

	// Create the URL table:
	urlTable, err := b.createURLTable(ctx)
	if err != nil {
		return
	}

	// Set the default agent, if needed:
	agent := b.agent
	if b.agent == "" {
		agent = DefaultAgent
	}

	// Create the metrics wrapper:
	var metricsWrapper func(http.RoundTripper) http.RoundTripper
	if b.metricsSubsystem != "" {
		var parsed *url.URL
		parsed, err = url.Parse(b.tokenURL)
		if err != nil {
			return
		}
		var wrapper *metrics.TransportWrapper
		wrapper, err = metrics.NewTransportWrapper().
			Path(parsed.Path).
			Subsystem(b.metricsSubsystem).
			Registerer(b.metricsRegisterer).
			Build()
		if err != nil {
			return
		}
		metricsWrapper = wrapper.Wrap
	}

	// Create the logging wrapper:
	var loggingWrapper func(http.RoundTripper) http.RoundTripper
	if b.logger.DebugEnabled() {
		wrapper := &dumpTransportWrapper{
			logger: b.logger,
		}
		loggingWrapper = wrapper.Wrap
	}

	// Initialize the client selector builder:
	clientSelectorBuilder := internal.NewClientSelector().
		Logger(b.logger).
		TrustedCAs(b.trustedCAs...).
		Insecure(b.insecure)

	var authnWrapper *authentication.TransportWrapper
	if b.includeDefaultAuthnTransportWrapper {
		// Create the authentication wrapper:
		authnWrapper, err = authentication.NewTransportWrapper().
			Logger(b.logger).
			TokenURL(b.tokenURL).
			User(b.user, b.password).
			Client(b.clientID, b.clientSecret).
			Tokens(b.tokens...).
			Scopes(b.scopes...).
			TrustedCAs(b.trustedCAs...).
			Insecure(b.insecure).
			TransportWrapper(metricsWrapper).
			TransportWrapper(loggingWrapper).
			TransportWrappers(b.transportWrappers...).
			MetricsSubsystem(b.metricsSubsystem).
			MetricsRegisterer(b.metricsRegisterer).
			Build(ctx)
		if err != nil {
			return
		}
		clientSelectorBuilder.TransportWrapper(authnWrapper.Wrap)
	}

	// Create the retry wrapper:
	retryWrapper, err := retry.NewTransportWrapper().
		Logger(b.logger).
		Limit(b.retryLimit).
		Interval(b.retryInterval).
		Jitter(b.retryJitter).
		Build(ctx)
	if err != nil {
		return
	}

	// Create the client selector:
	clientSelector, err := clientSelectorBuilder.
		TransportWrapper(metricsWrapper).
		TransportWrapper(retryWrapper.Wrap).
		TransportWrapper(loggingWrapper).
		TransportWrappers(b.transportWrappers...).
		Build(ctx)
	if err != nil {
		return
	}

	// Allocate and populate the connection object:
	connection = &Connection{
		logger:            b.logger,
		authnWrapper:      authnWrapper,
		retryWrapper:      retryWrapper,
		clientSelector:    clientSelector,
		urlTable:          urlTable,
		agent:             agent,
		metricsSubsystem:  b.metricsSubsystem,
		metricsRegisterer: b.metricsRegisterer,
	}

	return
}

func (b *ConnectionBuilder) createURLTable(ctx context.Context) (table []urlTableEntry, err error) {
	// Check that all the prefixes are acceptable:
	for prefix, base := range b.urlTable {
		if !validPrefixRE.MatchString(prefix) {
			err = fmt.Errorf(
				"prefix '%s' for URL '%s' isn't valid; it must start with a "+
					"slash and be composed of slash separated segments "+
					"containing only digits, letters, dashes and undercores",
				prefix, base,
			)
			return
		}
	}

	// Allocate space for the table:
	table = make([]urlTableEntry, len(b.urlTable))

	// For each alternative URL create the regular expression that will be used to check if
	// paths match it, and parse the base URL:
	i := 0
	for prefix, base := range b.urlTable {
		entry := &table[i]
		entry.prefix = prefix
		pattern := fmt.Sprintf("^%s(/.*)?$", regexp.QuoteMeta(prefix))
		entry.re, err = regexp.Compile(pattern)
		if err != nil {
			err = fmt.Errorf(
				"can't compile regular expression '%s' for URL with "+
					"prefix '%s' and URL '%s': %v",
				pattern, prefix, base, err,
			)
			return
		}
		entry.url, err = internal.ParseServerAddress(ctx, base)
		if err != nil {
			err = fmt.Errorf(
				"can't parse URL '%s' for prefix '%s': %w",
				base, prefix, err,
			)
			return
		}
		i++
	}

	// Sort the entries in descending order of the length of the prefix, so that later
	// when matching it will be easier to select the longest prefix that matches:
	sort.Slice(table, func(i, j int) bool {
		lenI := len(table[i].prefix)
		lenJ := len(table[j].prefix)
		return lenI > lenJ
	})

	// Write to the log the resulting table:
	if b.logger.DebugEnabled() {
		for _, entry := range table {
			b.logger.Debug(
				ctx,
				"Added URL with prefix '%s', regular expression "+
					"'%s' and URL '%s'",
				entry.prefix, entry.re, entry.url.Text,
			)
		}
	}

	return
}

// Logger returns the logger that is used by the connection.
func (c *Connection) Logger() logging.Logger {
	return c.logger
}

// TokenURL returns the URL that the connection is using request OpenID access tokens.
// An empty string is returned if the connection does not use authentication.
func (c *Connection) TokenURL() string {
	if c.authnWrapper == nil {
		return ""
	}
	return c.authnWrapper.TokenURL()
}

// Client returns OpenID client identifier and secret that the connection is using to request OpenID
// access tokens.
// Empty strings are returned if the connection does not use authentication.
func (c *Connection) Client() (id, secret string) {
	if c.authnWrapper != nil {
		id, secret = c.authnWrapper.Client()
	}
	return
}

// User returns the user name and password that the is using to request OpenID access tokens.
// Empty strings are returned if the connection does not use authentication.
func (c *Connection) User() (user, password string) {
	if c.authnWrapper != nil {
		user, password = c.authnWrapper.User()
	}
	return
}

// Scopes returns the OpenID scopes that the connection is using to request OpenID access tokens.
// An empty slice is returned if the connection does not use authentication.
func (c *Connection) Scopes() []string {
	if c.authnWrapper == nil {
		return []string{}
	}
	return c.authnWrapper.Scopes()
}

// URL returns the base URL of the API gateway.
func (c *Connection) URL() string {
	// The base URL will most likely be the last in the URL table because it is sorted in
	// descending order of the prefix length, so it is faster to traverse the table in
	// reverse order.
	for i := len(c.urlTable) - 1; i >= 0; i-- {
		entry := &c.urlTable[i]
		if entry.prefix == "" {
			return entry.url.Text
		}
	}
	return ""
}

// Agent returns the `User-Agent` header that the client is using for all HTTP requests.
func (c *Connection) Agent() string {
	return c.agent
}

// TrustedCAs sets returns the certificate pool that contains the certificate authorities that are
// trusted by the connection.
func (c *Connection) TrustedCAs() *x509.CertPool {
	return c.clientSelector.TrustedCAs()
}

// Insecure returns the flag that indicates if insecure communication with the server is enabled.
func (c *Connection) Insecure() bool {
	return c.clientSelector.Insecure()
}

// DisableKeepAlives returns the flag that indicates if HTTP keep alive is disabled.
func (c *Connection) DisableKeepAlives() bool {
	return c.clientSelector.DisableKeepAlives()
}

// RetryLimit gets the maximum number of retries for a request.
func (c *Connection) RetryLimit() int {
	return c.retryWrapper.Limit()
}

// RetryInteval returns the initial retry interval.
func (c *Connection) RetryInterval() time.Duration {
	return c.retryWrapper.Interval()
}

// RetryJitter returns the retry interval jitter factor.
func (c *Connection) RetryJitter() float64 {
	return c.retryWrapper.Jitter()
}

// MetricsSubsystem returns the name of the subsystem that is used by the connection to register
// metrics with Prometheus. An empty string means that no metrics are registered.
func (c *Connection) MetricsSubsystem() string {
	return c.metricsSubsystem
}

// AlternativeURLs returns the alternative URLs in use by the connection. Note that the map returned
// is a copy of the data used internally, so changing it will have no effect on the connection.
func (c *Connection) AlternativeURLs() map[string]string {
	// Copy all the entries of the URL table except the one corresponding to the empty prefix, as
	// that isn't usually set via the alternative URLs mechanism:
	result := map[string]string{}
	for _, entry := range c.urlTable {
		if entry.prefix != "" {
			result[entry.prefix] = entry.url.Text
		}
	}
	return result
}

// AccessTransparency returns the client for the access transparency service.
func (c *Connection) AccessTransparency() *accesstransparency.Client {
	return accesstransparency.NewClient(c, "/api/access_transparency")
}

// AccountsMgmt returns the client for the accounts management service.
func (c *Connection) AccountsMgmt() *accountsmgmt.Client {
	return accountsmgmt.NewClient(c, "/api/accounts_mgmt")
}

// AccountsMgmt returns the client for the accounts management service.
func (c *Connection) AddonsMgmt() *addonsmgmt.Client {
	return addonsmgmt.NewClient(c, "/api/addons_mgmt")
}

// ClustersMgmt returns the client for the clusters management service.
func (c *Connection) ClustersMgmt() *clustersmgmt.Client {
	return clustersmgmt.NewClient(c, "/api/clusters_mgmt")
}

// AroHCP returns the client for the ARO-HCP clusters management service.
func (c *Connection) AroHCP() *arohcp.Client {
	return arohcp.NewClient(c, "/api/aro_hcp")
}

// OSDFleetMgmt returns the client for the OSD management service.
func (c *Connection) OSDFleetMgmt() *osdfleetmgmt.Client {

	return osdfleetmgmt.NewClient(c, "/api/osd_fleet_mgmt")
}

// Authorizations returns the client for the authorizations service.
func (c *Connection) Authorizations() *authorizations.Client {
	return authorizations.NewClient(c, "/api/authorizations")
}

// ServiceLogs returns the client for the logs service.
func (c *Connection) ServiceLogs() *servicelogs.Client {
	return servicelogs.NewClient(c, "/api/service_logs")
}

// JobQueue returns the client for the Job Queues service.
func (c *Connection) JobQueue() *jobqueue.Client {
	return jobqueue.NewClient(c, "/api/job_queue")
}

// Status board returns the client for the status board service.
func (c *Connection) StatusBoard() *statusboard.Client {
	return statusboard.NewClient(c, "/api/status-board")
}

// ServiceMgmt returns the client for the service management service.
func (c *Connection) ServiceMgmt() *servicemgmt.Client {
	return servicemgmt.NewClient(c, "/api/service_mgmt")
}

// WebRCA returns the client for the web RCA service.
func (c *Connection) WebRCA() *webrca.Client {
	return webrca.NewClient(c, "/api/web-rca")
}

// Close releases all the resources used by the connection. It is very important to always close it
// once it is no longer needed, as otherwise those resources may be leaked. Trying to use a
// connection that has been closed will result in a error.
func (c *Connection) Close() error {
	var err error

	// in case the connection is already closed, return instead of printing an error message
	if c.closed {
		return nil
	}

	// Close the HTTP clients:
	err = c.clientSelector.Close()
	if err != nil {
		return err
	}

	// If the default authentication wrapper is set close it
	if c.authnWrapper != nil {
		// Close the authentication wrapper:
		err = c.authnWrapper.Close()
		if err != nil {
			return err
		}
	}

	// Mark the connection as closed, so that further attempts to use it will fail:
	c.closed = true
	return nil
}

func (c *Connection) checkClosed() error {
	if c.closed {
		return fmt.Errorf("connection is closed")
	}
	return nil
}

// validPrefixRE is the regular expression used to check patch prefixes.
var validPrefixRE = regexp.MustCompile(`^((/\w+)*)?$`)
