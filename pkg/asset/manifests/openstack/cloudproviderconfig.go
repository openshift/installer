package openstack

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	networkutils "github.com/gophercloud/utils/v2/openstack/networking/v2/networks"

	"github.com/openshift/installer/pkg/asset/installconfig/openstack"
	"github.com/openshift/installer/pkg/types"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
	"github.com/openshift/installer/pkg/types/powervc"
)

// Error represents a failure while generating OpenStack provider
// configuration.
type Error struct {
	err error
	msg string
}

func (e Error) Error() string { return e.msg + ": " + e.err.Error() }
func (e Error) Unwrap() error { return e.err }

// isOctaviaAvailable checks if the Octavia (load-balancer) endpoint exists
// in the OpenStack service catalog.
func isOctaviaAvailable(ctx context.Context, clientOpts *clientconfig.ClientOpts) bool {
	if clientOpts == nil {
		// If no client options provided, assume Octavia is available
		// to maintain backward compatibility
		return true
	}
	_, err := openstackdefaults.NewServiceClient(ctx, "load-balancer", clientOpts)
	if err != nil {
		var gerr *gophercloud.ErrEndpointNotFound
		if errors.As(err, &gerr) {
			return false
		}
		// For other errors, assume Octavia might be available
		// to avoid incorrectly disabling it
		return true
	}
	return true
}

func generateCloudProviderConfig(ctx context.Context, networkClient *gophercloud.ServiceClient, cloudConfig *clientconfig.Cloud, clientOpts *clientconfig.ClientOpts, installConfig types.InstallConfig) (cloudProviderConfigData, cloudProviderConfigCABundleData string, err error) {
	cloudProviderConfigData = `[Global]
secret-name = openstack-credentials
secret-namespace = kube-system
`
	if regionName := cloudConfig.RegionName; regionName != "" {
		cloudProviderConfigData += "region = " + regionName + "\n"
	}

	if caCertFile := cloudConfig.CACertFile; caCertFile != "" {
		cloudProviderConfigData += "ca-file = /etc/kubernetes/static-pod-resources/configmaps/cloud-config/ca-bundle.pem\n"
		caFile, err := os.ReadFile(caCertFile)
		if err != nil {
			return "", "", Error{err, "failed to read clouds.yaml ca-cert from disk"}
		}
		cloudProviderConfigCABundleData = string(caFile)
	}

	switch {
	case installConfig.Platform.Name() == powervc.Name:
		if installConfig.OpenStack.ExternalNetwork != "" {
			return "", "", fmt.Errorf("powervc does not support external network")
		}
		// powervc does not provide an equivalent to Octavia
		cloudProviderConfigData += "\n[LoadBalancer]\nenabled = false\n"
	case !isOctaviaAvailable(ctx, clientOpts):
		// Explicitly disable LoadBalancer when Octavia is not available
		// to prevent CCM from crashing on startup.
		// See: https://issues.redhat.com/browse/OCPBUGS-64842
		cloudProviderConfigData += "\n[LoadBalancer]\nenabled = false\n"
	case installConfig.OpenStack.ExternalNetwork != "":
		networkName := installConfig.OpenStack.ExternalNetwork // Yes, we use a name in install-config.yaml :/
		networkID, err := networkutils.IDFromName(ctx, networkClient, networkName)
		if err != nil {
			return "", "", Error{err, "failed to fetch external network " + networkName}
		}
		// If set get the ID and configure CCM to use that network for LB FIPs.
		cloudProviderConfigData += "\n[LoadBalancer]\n"
		cloudProviderConfigData += "floating-network-id = " + networkID + "\n"
	}

	return cloudProviderConfigData, cloudProviderConfigCABundleData, nil
}

// GenerateCloudProviderConfig adds the cloud provider config for the OpenStack
// platform in the provided configmap.
func GenerateCloudProviderConfig(ctx context.Context, installConfig types.InstallConfig) (cloudProviderConfigData, cloudProviderConfigCABundleData string, err error) {
	session, err := openstack.GetSession(installConfig.Platform.OpenStack.Cloud)
	if err != nil {
		return "", "", Error{err, "failed to get cloud config for openstack"}
	}

	networkClient, err := openstackdefaults.NewServiceClient(ctx, "network", session.ClientOpts)
	if err != nil {
		return "", "", Error{err, "failed to create a network client"}
	}

	return generateCloudProviderConfig(ctx, networkClient, session.CloudConfig, session.ClientOpts, installConfig)
}
