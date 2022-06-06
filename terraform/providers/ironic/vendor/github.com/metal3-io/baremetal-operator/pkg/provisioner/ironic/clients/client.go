package clients

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/httpbasic"
	"github.com/gophercloud/gophercloud/openstack/baremetal/noauth"
	httpbasicintrospection "github.com/gophercloud/gophercloud/openstack/baremetalintrospection/httpbasic"
	noauthintrospection "github.com/gophercloud/gophercloud/openstack/baremetalintrospection/noauth"
	"go.etcd.io/etcd/pkg/transport"
)

var tlsConnectionTimeout = time.Second * 30

// TLSConfig contains the TLS configuration for the Ironic connection.
// Using Go default values for this will result in no additional trusted
// CA certificates and a secure connection.
// When specifying Certificate and Private key, TLS connection will use
// client certificate authentication.
type TLSConfig struct {
	TrustedCAFile         string
	ClientCertificateFile string
	ClientPrivateKeyFile  string
	InsecureSkipVerify    bool
	SkipClientSANVerify   bool
}

func updateHTTPClient(client *gophercloud.ServiceClient, tlsConf TLSConfig) (*gophercloud.ServiceClient, error) {
	tlsInfo := transport.TLSInfo{
		TrustedCAFile:       tlsConf.TrustedCAFile,
		CertFile:            tlsConf.ClientCertificateFile,
		KeyFile:             tlsConf.ClientPrivateKeyFile,
		InsecureSkipVerify:  tlsConf.InsecureSkipVerify,
		SkipClientSANVerify: tlsConf.SkipClientSANVerify,
	}
	if _, err := os.Stat(tlsConf.TrustedCAFile); err != nil {
		if os.IsNotExist(err) {
			tlsInfo.TrustedCAFile = ""
		} else {
			return client, err
		}
	}
	if _, err := os.Stat(tlsConf.ClientCertificateFile); err != nil {
		if os.IsNotExist(err) {
			tlsInfo.CertFile = ""
		} else {
			return client, err
		}
	}
	if _, err := os.Stat(tlsConf.ClientPrivateKeyFile); err != nil {
		if os.IsNotExist(err) {
			tlsInfo.KeyFile = ""
		} else {
			return client, err
		}
	}
	if tlsInfo.CertFile != "" && tlsInfo.KeyFile != "" {
		tlsInfo.ClientCertAuth = true
	}

	tlsTransport, err := transport.NewTransport(tlsInfo, tlsConnectionTimeout)
	if err != nil {
		return client, err
	}
	c := http.Client{
		Transport: tlsTransport,
	}
	client.HTTPClient = c
	return client, nil
}

// IronicClient creates a client for Ironic
func IronicClient(ironicEndpoint string, auth AuthConfig, tls TLSConfig) (client *gophercloud.ServiceClient, err error) {
	switch auth.Type {
	case NoAuth:
		client, err = noauth.NewBareMetalNoAuth(noauth.EndpointOpts{
			IronicEndpoint: ironicEndpoint,
		})
	case HTTPBasicAuth:
		client, err = httpbasic.NewBareMetalHTTPBasic(httpbasic.EndpointOpts{
			IronicEndpoint:     ironicEndpoint,
			IronicUser:         auth.Username,
			IronicUserPassword: auth.Password,
		})
	default:
		err = fmt.Errorf("Unknown auth type %s", auth.Type)
	}
	if err != nil {
		return
	}

	// Ensure we have a microversion high enough to get the features
	// we need. Update docs/configuration.md when updating the version.
	// Version 1.74 allows retrival of the BIOS Registry
	client.Microversion = "1.74"

	return updateHTTPClient(client, tls)
}

// InspectorClient creates a client for Ironic Inspector
func InspectorClient(inspectorEndpoint string, auth AuthConfig, tls TLSConfig) (client *gophercloud.ServiceClient, err error) {
	switch auth.Type {
	case NoAuth:
		client, err = noauthintrospection.NewBareMetalIntrospectionNoAuth(
			noauthintrospection.EndpointOpts{
				IronicInspectorEndpoint: inspectorEndpoint,
			})
	case HTTPBasicAuth:
		client, err = httpbasicintrospection.NewBareMetalIntrospectionHTTPBasic(httpbasicintrospection.EndpointOpts{
			IronicInspectorEndpoint:     inspectorEndpoint,
			IronicInspectorUser:         auth.Username,
			IronicInspectorUserPassword: auth.Password,
		})
	default:
		err = fmt.Errorf("Unknown auth type %s", auth.Type)
	}
	if err != nil {
		return
	}
	return updateHTTPClient(client, tls)
}
