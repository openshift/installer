package vsphere

import (
	"context"
	"crypto/x509"
	"net/http"
	"net/url"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/master/pkg/session/session.go

// CreateVSphereClients creates the SOAP and REST client to access
// different portions of the vSphere API
// e.g. tags are only available in REST
func CreateVSphereClients(ctx context.Context, vcenter, username, password string, certificates ...x509.Certificate) (*vim25.Client, *rest.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	u, err := soap.ParseURL(vcenter)
	if err != nil {
		return nil, nil, err
	}
	u.User = url.UserPassword(username, password)

	soapClient := soap.NewClient(u, false)

	setHTTPTransportRootCA(soapClient.DefaultTransport(), certificates...)

	vimClient, err := vim25.NewClient(ctx, soapClient)

	if err != nil {
		return nil, nil, err
	}

	c := &govmomi.Client{
		Client:         vimClient,
		SessionManager: session.NewManager(vimClient),
	}

	err = c.SessionManager.Login(ctx, u.User)
	if err != nil {
		return nil, nil, err
	}

	restClient := rest.NewClient(vimClient)

	setHTTPTransportRootCA(restClient.DefaultTransport(), certificates...)

	err = restClient.Login(ctx, u.User)
	if err != nil {
		return nil, nil, err
	}

	return vimClient, restClient, nil
}

func setHTTPTransportRootCA(transport *http.Transport, certificates ...x509.Certificate) {
	if certificates != nil {
		// Use the systemcertpool if available
		certpool, err := x509.SystemCertPool()

		//If not create a newcertpool
		if err != nil {
			certpool = x509.NewCertPool()
		}

		for _, c := range certificates {
			certpool.AddCert(&c)
		}

		transport.TLSClientConfig.RootCAs = certpool
	}
}
