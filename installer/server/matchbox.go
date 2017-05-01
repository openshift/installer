package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/server/asset"
)

var (
	defaultTimeout = 30 * time.Second
)

// MatchboxConfig configures a matchbox client with PEM encoded TLS credentials.
type MatchboxConfig struct {
	Endpoint   string
	CA         []byte
	ClientCert []byte
	ClientKey  []byte
}

// MatchboxClient allows Cluster manifests to be written to the matchbox service.
type MatchboxClient struct {
	client *matchbox.Client
}

// NewMatchboxClient returns a new MatchboxClient.
func NewMatchboxClient(config *MatchboxConfig) (*MatchboxClient, error) {
	tlscfg, err := tlsConfig(config.CA, config.ClientCert, config.ClientKey)
	if err != nil {
		return nil, err
	}
	client, err := matchbox.New(&matchbox.Config{
		Endpoints:   []string{config.Endpoint},
		DialTimeout: defaultTimeout,
		TLS:         tlscfg,
	})
	if err != nil {
		return nil, err
	}
	return &MatchboxClient{
		client: client,
	}, nil
}

// Close closes the client's connections.
func (c *MatchboxClient) Close() error {
	return c.client.Close()
}

// Push writes machine profiles, groups, and Ignition templates to the matchbox
// service. Repeated writes are idempotent.
func (c *MatchboxClient) Push(ctx context.Context, groups []*storagepb.Group, profiles []*storagepb.Profile, ignitions []asset.Asset) error {
	// TODO: Parallelize
	ctx, _ = context.WithTimeout(ctx, defaultTimeout)
	for _, profile := range profiles {
		_, err := c.client.Profiles.ProfilePut(ctx, &serverpb.ProfilePutRequest{Profile: profile})
		if err != nil {
			return err
		}
	}

	for _, group := range groups {
		_, err := c.client.Groups.GroupPut(ctx, &serverpb.GroupPutRequest{Group: group})
		if err != nil {
			return err
		}
	}

	for _, asset := range ignitions {
		_, err := c.client.Ignition.IgnitionPut(ctx, &serverpb.IgnitionPutRequest{
			Name:   asset.Name(),
			Config: asset.Data(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// tlsConfig returns a matchbox client TLS.Config.
func tlsConfig(ca, clientCert, clientKey []byte) (*tls.Config, error) {
	// certificate authority for verifying the server
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(ca)
	if !ok {
		return nil, errors.New("no PEM certificates were parsed")
	}

	// client certificate for authentication
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		// CA bundle the client should trust when verifying the server
		RootCAs: pool,
		// Client certificate to authenticate to the server
		Certificates: []tls.Certificate{cert},
	}, nil
}

// ignitionPath returns the Ignition endpoint client machines should use as a
// kernel argument.
func ignitionPath(endpoint string) string {
	return fmt.Sprintf("http://%s/ignition?uuid=${uuid}&mac=${net0/mac:hexhyp}", endpoint)
}

// coreosAssetsPath returns the a matchbox service's CoreOS assets endpoint,
// suitable for use as a baseurl by the CoreOS installer.
func coreosAssetsPath(endpoint string) string {
	return fmt.Sprintf("http://%s/assets/coreos", endpoint)
}
