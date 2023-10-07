/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file contains the implementation of the object that selects the HTTP client to use to
// connect to servers using TCP or Unix sockets.

package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"

	"golang.org/x/net/http2"

	"github.com/openshift-online/ocm-sdk-go/logging"
)

// ClientSelectorBuilder contains the information and logic needed to create an HTTP client
// selector. Don't create instances of this type directly, use the NewClientSelector function.
type ClientSelectorBuilder struct {
	logger            logging.Logger
	trustedCAs        []interface{}
	insecure          bool
	disableKeepAlives bool
	transportWrappers []func(http.RoundTripper) http.RoundTripper
}

// ClientSelector contains the information needed to create select the HTTP client to use to connect
// to servers using TCP or Unix sockets.
type ClientSelector struct {
	logger            logging.Logger
	trustedCAs        *x509.CertPool
	insecure          bool
	disableKeepAlives bool
	transportWrappers []func(http.RoundTripper) http.RoundTripper
	cookieJar         http.CookieJar
	clientsMutex      *sync.Mutex
	clientsTable      map[string]*http.Client
}

// NewClientSelector creates a builder that can then be used to configure and create an HTTP client
// selector.
func NewClientSelector() *ClientSelectorBuilder {
	return &ClientSelectorBuilder{}
}

// Logger sets the logger that will be used by the selector and by the created HTTP clients to write
// messages to the log. This is mandatory.
func (b *ClientSelectorBuilder) Logger(value logging.Logger) *ClientSelectorBuilder {
	b.logger = value
	return b
}

// TrustedCA sets a source that contains he certificate authorities that will be trusted by the HTTP
// clients. If this isn't explicitly specified then the clients will trust the certificate
// authorities trusted by default by the system. The value can be a *x509.CertPool or a string,
// anything else will cause an error when Build method is called. If it is a *x509.CertPool then the
// value will replace any other source given before. If it is a string then it should be the name of
// a PEM file. The contents of that file will be added to the previously given sources.
func (b *ClientSelectorBuilder) TrustedCA(value interface{}) *ClientSelectorBuilder {
	if value != nil {
		b.trustedCAs = append(b.trustedCAs, value)
	}
	return b
}

// TrustedCAs sets a list of sources that contains he certificate authorities that will be trusted
// by the HTTP clients. See the documentation of the TrustedCA method for more information about the
// accepted values.
func (b *ClientSelectorBuilder) TrustedCAs(values ...interface{}) *ClientSelectorBuilder {
	for _, value := range values {
		b.TrustedCA(value)
	}
	return b
}

// Insecure enables insecure communication with the servers. This disables verification of TLS
// certificates and host names and it isn't recommended for a production environment.
func (b *ClientSelectorBuilder) Insecure(flag bool) *ClientSelectorBuilder {
	b.insecure = flag
	return b
}

// DisableKeepAlives disables HTTP keep-alives with the serviers. This is unrelated to similarly
// named TCP keep-alives.
func (b *ClientSelectorBuilder) DisableKeepAlives(flag bool) *ClientSelectorBuilder {
	b.disableKeepAlives = flag
	return b
}

// TransportWrapper adds a function that will be used to wrap the transports of the HTTP clients. If
// used multiple times the transport wrappers will be called in the same order that they are added.
func (b *ClientSelectorBuilder) TransportWrapper(
	value func(http.RoundTripper) http.RoundTripper) *ClientSelectorBuilder {
	if value != nil {
		b.transportWrappers = append(b.transportWrappers, value)
	}
	return b
}

// TransportWrappers adds a list of functions that will be used to wrap the transports of the HTTP clients.
func (b *ClientSelectorBuilder) TransportWrappers(
	values ...func(http.RoundTripper) http.RoundTripper) *ClientSelectorBuilder {
	for _, value := range values {
		b.TransportWrapper(value)
	}
	return b
}

// Build uses the information stored in the builder to create a new HTTP client selector.
func (b *ClientSelectorBuilder) Build(ctx context.Context) (result *ClientSelector, err error) {
	// Check parameters:
	if b.logger == nil {
		err = fmt.Errorf("logger is mandatory")
		return
	}

	// Create the cookie jar:
	cookieJar, err := b.createCookieJar()
	if err != nil {
		return
	}

	// Load trusted CAs:
	trustedCAs, err := b.loadTrustedCAs(ctx)
	if err != nil {
		return
	}

	// Create and populate the object:
	result = &ClientSelector{
		logger:            b.logger,
		trustedCAs:        trustedCAs,
		insecure:          b.insecure,
		disableKeepAlives: b.disableKeepAlives,
		transportWrappers: b.transportWrappers,
		cookieJar:         cookieJar,
		clientsMutex:      &sync.Mutex{},
		clientsTable:      map[string]*http.Client{},
	}

	return
}

func (b *ClientSelectorBuilder) loadTrustedCAs(ctx context.Context) (result *x509.CertPool,
	err error) {
	result, err = loadSystemCAs()
	if err != nil {
		return
	}
	for _, ca := range b.trustedCAs {
		switch source := ca.(type) {
		case *x509.CertPool:
			b.logger.Debug(
				ctx,
				"Default trusted CA certificates have been explicitly replaced",
			)
			result = source
		case string:
			b.logger.Debug(
				ctx,
				"Loading trusted CA certificates from file '%s'",
				source,
			)
			var buffer []byte
			buffer, err = os.ReadFile(source) // #nosec G304
			if err != nil {
				result = nil
				err = fmt.Errorf(
					"can't read trusted CA certificates from file '%s': %w",
					source, err,
				)
				return
			}
			if !result.AppendCertsFromPEM(buffer) {
				result = nil
				err = fmt.Errorf(
					"file '%s' doesn't contain any certificate",
					source,
				)
				return
			}
		default:
			result = nil
			err = fmt.Errorf(
				"don't know how to load trusted CA from source of type '%T'",
				source,
			)
			return
		}
	}
	return
}

func (b *ClientSelectorBuilder) createCookieJar() (result http.CookieJar, err error) {
	result, err = cookiejar.New(nil)
	return
}

// Select returns an HTTP client to use to connect to the given server address. If a client has been
// created previously for the server address it will be reused, otherwise it will be created.
func (s *ClientSelector) Select(ctx context.Context, address *ServerAddress) (client *http.Client,
	err error) {
	// We will be modifiying the clients table so we need to acquire the lock before proceeding:
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	// Get an existing client, or create a new one if it doesn't exist yet:
	key := s.key(address)
	client, ok := s.clientsTable[key]
	if ok {
		return
	}
	s.logger.Debug(ctx, "Client for key '%s' doesn't exist, will create it", key)
	client, err = s.create(ctx, address)
	if err != nil {
		return
	}
	s.clientsTable[key] = client

	return
}

// Forget forgets the client for the given server address. This is intended for situations where a
// client is missbehaving, for example when it is generating protocol errors. In those situations
// connections may be still open but already unusable. To avoid additional errors is beter to
// discard the client and create a new one.
func (s *ClientSelector) Forget(ctx context.Context, address *ServerAddress) error {
	// We will be modifiying the clients table so we need to acquire the lock before proceeding:
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	// Close the client and delete it from the table:
	key := s.key(address)
	client, ok := s.clientsTable[key]
	if ok {
		delete(s.clientsTable, key)
		client.CloseIdleConnections()
	}
	s.logger.Debug(ctx, "Discarded client for key '%s'", key)

	return nil
}

// key calculates from the given server address the key that is used to store clients in the table.
func (s *ClientSelector) key(address *ServerAddress) string {
	// We need to use a different client for each TCP host name and each Unix socket because we
	// explicitly set the TLS server name to the host name. For example, if the first request is
	// for the SSO service (it will usually be) then we would set the TLS server name to
	// `sso.redhat.com`. The next API request would then use the same client and therefore it
	// will use `sso.redhat.com` as the TLS server name. If the server uses SNI to select the
	// certificates it will then fail because the API server doesn't have any certificate for
	// `sso.redhat.com`, it will return the default certificates, and then the validation would
	// fail with an error message like this:
	//
	//      x509: certificate is valid for *.apps.app-sre-prod-04.i5h0.p1.openshiftapps.com,
	//      api.app-sre-prod-04.i5h0.p1.openshiftapps.com,
	//      rh-api.app-sre-prod-04.i5h0.p1.openshiftapps.com, not sso.redhat.com
	//
	// To avoid this we add the host name or socket path as a suffix to the key.
	key := address.Network
	switch address.Network {
	case UnixNetwork:
		key = fmt.Sprintf("%s:%s", key, address.Socket)
	case TCPNetwork:
		key = fmt.Sprintf("%s:%s", key, address.Host)
	}
	return key
}

// create creates a new HTTP client to use to connect to the given address.
func (s *ClientSelector) create(ctx context.Context, address *ServerAddress) (result *http.Client,
	err error) {
	// Create the transport:
	transport, err := s.createTransport(ctx, address)
	if err != nil {
		return
	}

	// Create the client:
	result = &http.Client{
		Jar:       s.cookieJar,
		Transport: transport,
	}
	if s.logger.DebugEnabled() {
		result.CheckRedirect = func(request *http.Request, via []*http.Request) error {
			s.logger.Info(
				request.Context(),
				"Following redirect from '%s' to '%s'",
				via[0].URL,
				request.URL,
			)
			return nil
		}
	}

	return
}

// createTransport creates a new HTTP transport to use to connect to the given server address.
func (s *ClientSelector) createTransport(ctx context.Context,
	address *ServerAddress) (result http.RoundTripper, err error) {
	// Prepare the TLS configuration:
	// #nosec 402
	config := &tls.Config{
		ServerName:         address.Host,
		InsecureSkipVerify: s.insecure,
		RootCAs:            s.trustedCAs,
	}

	// Create the transport:
	if address.Protocol != H2CProtocol {
		// Create a regular transport. Note that this does support HTTP/2 with TLS, but
		// not h2c:
		transport := &http.Transport{
			TLSClientConfig:    config,
			Proxy:              http.ProxyFromEnvironment,
			DisableKeepAlives:  s.disableKeepAlives,
			DisableCompression: false,
			ForceAttemptHTTP2:  true,
		}

		// In order to use Unix sockets we need to explicitly set dialers that use `unix` as
		// network and the socket file as address, otherwise the HTTP client will always use
		// `tcp` as the network and the host name from the request as the address:
		if address.Network == UnixNetwork {
			transport.DialContext = func(ctx context.Context, _, _ string) (net.Conn,
				error) {
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, UnixNetwork, address.Socket)
			}
			transport.DialTLSContext = func(ctx context.Context, _, _ string) (net.Conn,
				error) {
				dialer := tls.Dialer{
					Config: config,
				}
				return dialer.DialContext(ctx, UnixNetwork, address.Socket)
			}
		}

		// Prepare the result:
		result = transport
	} else {
		// In order to use h2c we need to tell the transport to allow the `http` scheme:
		transport := &http2.Transport{
			AllowHTTP:          true,
			DisableCompression: false,
		}

		// We also need to ignore TLS configuration when dialing, and explicitly set the
		// network and socket when using Unix sockets:
		if address.Network == UnixNetwork {
			transport.DialTLS = func(_, _ string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(UnixNetwork, address.Socket)
			}
		} else {
			transport.DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn,
				error) {
				return net.Dial(network, addr)
			}
		}

		// Prepare the result:
		result = transport
	}

	// Transport wrappers are stored in the order that the round trippers that they create
	// should be called. That means that we need to call them in reverse order.
	for i := len(s.transportWrappers) - 1; i >= 0; i-- {
		result = s.transportWrappers[i](result)
	}

	return
}

// TrustedCAs sets returns the certificate pool that contains the certificate authorities that are
// trusted by the HTTP clients.
func (s *ClientSelector) TrustedCAs() *x509.CertPool {
	return s.trustedCAs
}

// Insecure returns the flag that indicates if insecure communication with the server is enabled.
func (s *ClientSelector) Insecure() bool {
	return s.insecure
}

// DisableKeepAlives retursnt the flag that indicates if HTTP keep alive is disabled.
func (s *ClientSelector) DisableKeepAlives() bool {
	return s.disableKeepAlives
}

// Close closes all the connections used by all the clients created by the selector.
func (s *ClientSelector) Close() error {
	for _, client := range s.clientsTable {
		client.CloseIdleConnections()
	}
	return nil
}
