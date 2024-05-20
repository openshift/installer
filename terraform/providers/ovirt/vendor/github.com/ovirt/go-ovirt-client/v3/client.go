package ovirtclient

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net/http"
	"sync"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// Client is a simplified client for the oVirt API.
//
//goland:noinspection GoDeprecation
type Client interface {
	// GetURL returns the oVirt engine base URL.
	GetURL() string
	// Reconnect triggers the client to reauthenticate against the oVirt Engine.
	Reconnect() (err error)
	// WithContext creates a subclient with the specified context applied.
	WithContext(ctx context.Context) Client
	// GetContext returns the current context of the client. May be nil.
	GetContext() context.Context

	AffinityGroupClient
	DiskClient
	DiskAttachmentClient
	VMClient
	NICClient
	VNICProfileClient
	NetworkClient
	DatacenterClient
	ClusterClient
	StorageDomainClient
	HostClient
	TemplateClient
	TemplateDiskClient
	TestConnectionClient
	TagClient
	FeatureClient
	InstanceTypeClient
	GraphicsConsoleClient
}

// ClientWithLegacySupport is an extension of Client that also offers the ability to retrieve the underlying
// SDK connection or a configured HTTP client.
type ClientWithLegacySupport interface {
	// GetSDKClient returns a configured oVirt SDK client for the use cases that are not covered by goVirt.
	GetSDKClient() *ovirtsdk4.Connection

	// GetHTTPClient returns a configured HTTP client for the oVirt engine. This can be used to send manual
	// HTTP requests to the oVirt engine.
	GetHTTPClient() http.Client

	Client
}

type oVirtClient struct {
	reconnectLock   *sync.Mutex
	conn            *ovirtsdk4.Connection
	ctx             context.Context
	httpClient      http.Client
	logger          Logger
	url             string
	username        string
	password        string
	tlsConfig       *tls.Config
	extraSettings   ExtraSettings
	nonSecureRandom *rand.Rand
	verify          func(connection Client) error
}

func (o *oVirtClient) WithContext(ctx context.Context) Client {
	return &oVirtClient{
		o.reconnectLock,
		o.conn,
		ctx,
		o.httpClient,
		o.logger.WithContext(ctx),
		o.url,
		o.username,
		o.password,
		o.tlsConfig,
		o.extraSettings,
		o.nonSecureRandom,
		o.verify,
	}
}

func (o *oVirtClient) GetContext() context.Context {
	return o.ctx
}

func (o *oVirtClient) Reconnect() error {
	o.reconnectLock.Lock()
	defer o.reconnectLock.Unlock()
	connBuilder := ovirtsdk4.NewConnectionBuilder().
		URL(o.url).
		Username(o.username).
		Password(o.password).
		TLSConfig(o.tlsConfig)
	if err := processExtraSettings(o.extraSettings, connBuilder); err != nil {
		return err
	}

	conn, err := connBuilder.Build()
	if err != nil {
		return wrap(err, EUnidentified, "failed to create underlying oVirt connection")
	}
	if o.conn == nil {
		o.conn = conn
	} else {
		// Replace the structure under the pointer to make all instances update.
		*o.conn = *conn
	}

	if o.verify != nil {
		if err := o.verify(o); err != nil {
			return err
		}
	}
	return nil
}

func (o *oVirtClient) GetSDKClient() *ovirtsdk4.Connection {
	return o.conn
}

func (o *oVirtClient) GetHTTPClient() http.Client {
	return o.httpClient
}

func (o *oVirtClient) GetURL() string {
	return o.url
}
