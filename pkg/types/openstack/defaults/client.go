package defaults

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"

	"github.com/openshift/installer/pkg/version"
)

// getUserAgent generates a Gophercloud UserAgent to help cloud operators
// disambiguate openshift-installer requests.
func getUserAgent() (gophercloud.UserAgent, error) {
	ua := gophercloud.UserAgent{}

	version, err := version.Version()
	if err != nil {
		return ua, err
	}

	ua.Prepend(fmt.Sprintf("openshift-installer/%s", version))
	return ua, nil
}

// NewServiceClient is a wrapper around Gophercloud's NewServiceClient that
// ensures we consistently set a user-agent.
func NewServiceClient(ctx context.Context, service string, opts *clientconfig.ClientOpts) (*gophercloud.ServiceClient, error) {
	ua, err := getUserAgent()
	if err != nil {
		return nil, err
	}

	client, err := clientconfig.NewServiceClient(ctx, service, opts)
	if err != nil {
		return nil, err
	}

	client.UserAgent = ua

	return client, nil
}
