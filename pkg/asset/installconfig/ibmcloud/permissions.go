package ibmcloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/types"
)

const (
	serviceIAMIdentity    = "IAM Identity Service"
	serviceResourceGroups = "Resource Groups"
	serviceVPC            = "VPC Infrastructure Services"
	serviceCIS            = "Internet Services"
	serviceDNS            = "DNS Services"
	serviceCOS            = "Cloud Object Storage"

	// cosPermsProbeName is used only to exercise COS list permission. A missing
	// instance is success; HTTP 403 is not.
	cosPermsProbeName = "openshift-install-perms-probe"
)

// ValidatePerms probes IBM Cloud APIs with the installer API key and fails
// fast if the key cannot access services required to create a cluster.
// Empty LIST results are success. HTTP 403 / unauthorized is a permission error.
//
// IBM Cloud IPI always uses credentialsMode: Manual (CCO does not mint). Manual
// only skips in-cluster credential minting; the installer API key still
// provisions infrastructure and must be checked.
func ValidatePerms(ctx context.Context, client API, ic *types.InstallConfig) error {
	if ic == nil || ic.IBMCloud == nil {
		return fmt.Errorf("ibmcloud platform configuration is required")
	}

	var missing []string

	if err := probeIAMIdentity(ctx, client); err != nil {
		missing = append(missing, formatMissing(serviceIAMIdentity, err))
	}

	if err := probeResourceGroups(ctx, client, ic.IBMCloud.ResourceGroupName); err != nil {
		missing = append(missing, formatMissing(serviceResourceGroups, err))
	}

	if err := probeVPCs(ctx, client, ic.IBMCloud.Region); err != nil {
		missing = append(missing, formatMissing(serviceVPC, err))
	}

	dnsService := serviceCIS
	if ic.Publish == types.InternalPublishingStrategy {
		dnsService = serviceDNS
	}
	if err := probeDNS(ctx, client, ic.Publish); err != nil {
		missing = append(missing, formatMissing(dnsService, err))
	}

	// IPI always uploads the RHCOS image to COS.
	if err := probeCOS(ctx, client); err != nil {
		missing = append(missing, formatMissing(serviceCOS, err))
	}

	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("IBM Cloud API key lacks required permissions: %s", strings.Join(missing, "; "))
}

func probeIAMIdentity(ctx context.Context, client API) error {
	_, err := client.GetAuthenticatorAPIKeyDetails(ctx)
	if isForbidden(err) {
		return err
	}
	return nil
}

func probeResourceGroups(ctx context.Context, client API, resourceGroupName string) error {
	_, err := client.GetResourceGroups(ctx)
	if isForbidden(err) {
		return err
	}
	if resourceGroupName == "" {
		return nil
	}
	_, err = client.GetResourceGroup(ctx, resourceGroupName)
	if isForbidden(err) {
		return err
	}
	return nil
}

func probeVPCs(ctx context.Context, client API, region string) error {
	_, err := client.GetVPCs(ctx, region)
	if isForbidden(err) {
		return err
	}
	return nil
}

func probeDNS(ctx context.Context, client API, publish types.PublishingStrategy) error {
	_, err := client.GetDNSZones(ctx, publish)
	if isForbidden(err) {
		return err
	}
	return nil
}

func probeCOS(ctx context.Context, client API) error {
	_, err := client.GetCOSInstanceByName(ctx, cosPermsProbeName)
	if isForbidden(err) {
		return err
	}
	return nil
}

func formatMissing(service string, err error) string {
	return fmt.Sprintf("%s (%v)", service, err)
}

func isForbidden(err error) bool {
	if err == nil {
		return false
	}
	var cosNotFound *COSResourceNotFoundError
	if errors.As(err, &cosNotFound) {
		return false
	}
	var vpcNotFound *VPCResourceNotFoundError
	if errors.As(err, &vpcNotFound) {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, needle := range []string{
		"status 403",
		"status code: 403",
		"statuscode: 403",
		"forbidden",
		"not authorized",
		"unauthorized",
		"access is denied",
		"user is not authorized",
	} {
		if strings.Contains(msg, needle) {
			return true
		}
	}
	return false
}
