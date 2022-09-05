package conversion

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/vsphere"
	utilsslice "k8s.io/utils/strings/slices"
)

// ConvertInstallConfig is modeled after the k8s conversion schemes, which is
// how deprecated values are upconverted.
// This updates the APIVersion to reflect the fact that we've internally
// upconverted.
func ConvertInstallConfig(config *types.InstallConfig) error {
	// check that the version is convertible
	switch config.APIVersion {
	case types.InstallConfigVersion, "v1beta3", "v1beta4":
		// works
	case "":
		return field.Required(field.NewPath("apiVersion"), "no version was provided")
	default:
		return field.Invalid(field.NewPath("apiVersion"), config.APIVersion, fmt.Sprintf("cannot upconvert from version %s", config.APIVersion))
	}
	convertNetworking(config)

	switch config.Platform.Name() {
	case baremetal.Name:
		if err := convertBaremetal(config); err != nil {
			return err
		}
	case nutanix.Name:
		if err := convertNutanix(config); err != nil {
			return err
		}
	case openstack.Name:
		if err := convertOpenStack(config); err != nil {
			return err
		}
	case aws.Name:
		if err := convertAWS(config); err != nil {
			return err
		}
	case vsphere.Name:
		if err := convertVSphere(config); err != nil {
			return err
		}
	case ovirt.Name:
		if err := convertOVirt(config); err != nil {
			return err
		}
	}

	config.APIVersion = types.InstallConfigVersion
	return nil
}

// convertNetworking upconverts deprecated fields in networking
func convertNetworking(config *types.InstallConfig) {
	if config.Networking == nil {
		return
	}

	netconf := config.Networking

	if len(netconf.ClusterNetwork) == 0 {
		netconf.ClusterNetwork = netconf.DeprecatedClusterNetworks
	}

	if len(netconf.MachineNetwork) == 0 && netconf.DeprecatedMachineCIDR != nil {
		netconf.MachineNetwork = []types.MachineNetworkEntry{
			{CIDR: *netconf.DeprecatedMachineCIDR},
		}
	}

	if len(netconf.ServiceNetwork) == 0 && netconf.DeprecatedServiceCIDR != nil {
		netconf.ServiceNetwork = []ipnet.IPNet{*netconf.DeprecatedServiceCIDR}
	}

	// Convert type to networkType if the latter is missing
	if netconf.NetworkType == "" {
		netconf.NetworkType = netconf.DeprecatedType
	}

	// Recognize the OpenShiftSDN network plugin name regardless of capitalization, for
	// backward compatibility
	if strings.ToLower(netconf.NetworkType) == strings.ToLower(string(operv1.NetworkTypeOpenShiftSDN)) {
		netconf.NetworkType = string(operv1.NetworkTypeOpenShiftSDN)
	}

	// Convert hostSubnetLength to hostPrefix
	for i, entry := range netconf.ClusterNetwork {
		if entry.HostPrefix == 0 && entry.DeprecatedHostSubnetLength != 0 {
			_, size := entry.CIDR.Mask.Size()
			netconf.ClusterNetwork[i].HostPrefix = int32(size) - entry.DeprecatedHostSubnetLength
		}
	}
}

// convertBaremetal upconverts deprecated fields in the baremetal platform.
// ProvisioningDHCPExternal has been replaced by setting the ProvisioningNetwork
// field to "Unmanaged", ProvisioningHostIP has been replaced by
// ClusterProvisioningIP, apiVIP has been replaced by apiVIPs and ingressVIP has
// been replaced by ingressVIPs.
func convertBaremetal(config *types.InstallConfig) error {
	if config.Platform.BareMetal.DeprecatedProvisioningDHCPExternal && config.Platform.BareMetal.ProvisioningNetwork == "" {
		config.Platform.BareMetal.ProvisioningNetwork = baremetal.UnmanagedProvisioningNetwork
	}

	if config.Platform.BareMetal.DeprecatedProvisioningHostIP != "" && config.Platform.BareMetal.ClusterProvisioningIP == "" {
		config.Platform.BareMetal.ClusterProvisioningIP = config.Platform.BareMetal.DeprecatedProvisioningHostIP
	}

	// If user specified both, but they aren't equal, let them know they are the same field
	if config.Platform.BareMetal.DeprecatedProvisioningHostIP != "" &&
		config.Platform.BareMetal.DeprecatedProvisioningHostIP != config.Platform.BareMetal.ClusterProvisioningIP {
		return field.Invalid(field.NewPath("platform").Child("baremetal").Child("provisioningHostIP"),
			config.Platform.BareMetal.DeprecatedProvisioningHostIP, "provisioningHostIP is deprecated; only clusterProvisioningIP needs to be specified")
	}

	if err := upconvertVIP(&config.Platform.BareMetal.APIVIPs, config.Platform.BareMetal.DeprecatedAPIVIP, "apiVIP", "apiVIPs", field.NewPath("platform").Child("baremetal")); err != nil {
		return err
	}

	if err := upconvertVIP(&config.Platform.BareMetal.IngressVIPs, config.Platform.BareMetal.DeprecatedIngressVIP, "ingressVIP", "ingressVIPs", field.NewPath("platform").Child("baremetal")); err != nil {
		return err
	}

	return nil
}

// convertOpenStack upconverts deprecated fields in the OpenStack platform.
func convertOpenStack(config *types.InstallConfig) error {
	// LbFloatingIP has been renamed to APIFloatingIP
	if config.Platform.OpenStack.DeprecatedLbFloatingIP != "" {
		if config.Platform.OpenStack.APIFloatingIP == "" {
			config.Platform.OpenStack.APIFloatingIP = config.Platform.OpenStack.DeprecatedLbFloatingIP
		} else if config.Platform.OpenStack.DeprecatedLbFloatingIP != config.Platform.OpenStack.APIFloatingIP {
			// Return error if both LbFloatingIP and APIFloatingIP are specified in the config
			return field.Forbidden(field.NewPath("platform").Child("openstack").Child("lbFloatingIP"), "cannot specify lbFloatingIP and apiFloatingIP together")
		}
	}

	// computeFlavor has been deprecated in favor of type in defaultMachinePlatform.
	if config.Platform.OpenStack.DeprecatedFlavorName != "" {
		if config.Platform.OpenStack.DefaultMachinePlatform == nil {
			config.Platform.OpenStack.DefaultMachinePlatform = &openstack.MachinePool{}
		}

		if config.Platform.OpenStack.DefaultMachinePlatform.FlavorName != "" && config.Platform.OpenStack.DefaultMachinePlatform.FlavorName != config.Platform.OpenStack.DeprecatedFlavorName {
			// Return error if both computeFlavor and type of defaultMachinePlatform are specified in the config
			return field.Forbidden(field.NewPath("platform").Child("openstack").Child("computeFlavor"), "cannot specify computeFlavor and type in defaultMachinePlatform together")
		}

		config.Platform.OpenStack.DefaultMachinePlatform.FlavorName = config.Platform.OpenStack.DeprecatedFlavorName
	}

	if err := upconvertVIP(&config.Platform.OpenStack.APIVIPs, config.Platform.OpenStack.DeprecatedAPIVIP, "apiVIP", "apiVIPs", field.NewPath("platform").Child("openstack")); err != nil {
		return err
	}

	if err := upconvertVIP(&config.Platform.OpenStack.IngressVIPs, config.Platform.OpenStack.DeprecatedIngressVIP, "ingressVIP", "ingressVIPs", field.NewPath("platform").Child("openstack")); err != nil {
		return err
	}

	return nil
}

// convertNutanix upconverts deprecated fields in the Nutanix platform.
func convertNutanix(config *types.InstallConfig) error {
	if err := upconvertVIP(&config.Platform.Nutanix.APIVIPs, config.Platform.Nutanix.DeprecatedAPIVIP, "apiVIP", "apiVIPs", field.NewPath("platform").Child("nutanix")); err != nil {
		return err
	}

	if err := upconvertVIP(&config.Platform.Nutanix.IngressVIPs, config.Platform.Nutanix.DeprecatedIngressVIP, "ingressVIP", "ingressVIPs", field.NewPath("platform").Child("nutanix")); err != nil {
		return err
	}

	return nil
}

// convertVSphere upconverts deprecated fields in the VSphere platform.
func convertVSphere(config *types.InstallConfig) error {
	if err := upconvertVIP(&config.Platform.VSphere.APIVIPs, config.Platform.VSphere.DeprecatedAPIVIP, "apiVIP", "apiVIPs", field.NewPath("platform").Child("vsphere")); err != nil {
		return err
	}

	if err := upconvertVIP(&config.Platform.VSphere.IngressVIPs, config.Platform.VSphere.DeprecatedIngressVIP, "ingressVIP", "ingressVIPs", field.NewPath("platform").Child("vsphere")); err != nil {
		return err
	}

	return nil
}

// convertOVirt upconverts deprecated fields in the OVirt platform.
func convertOVirt(config *types.InstallConfig) error {
	if err := upconvertVIP(&config.Platform.Ovirt.APIVIPs, config.Platform.Ovirt.DeprecatedAPIVIP, "api_vip", "api_vips", field.NewPath("platform").Child("ovirt")); err != nil {
		return err
	}

	if err := upconvertVIP(&config.Platform.Ovirt.IngressVIPs, config.Platform.Ovirt.DeprecatedIngressVIP, "ingress_vip", "ingress_vips", field.NewPath("platform").Child("ovirt")); err != nil {
		return err
	}

	return nil
}

// upconvertVIP upconverts the deprecated VIP (oldVIPValue) to the new VIPs
// slice (newVIPValues). It returns errors, if both fields are set and all
// contain unique values
func upconvertVIP(newVIPValues *[]string, oldVIPValue, newFieldName, oldFieldName string, fldPath *field.Path) error {
	if oldVIPValue != "" && len(*newVIPValues) == 0 {
		*newVIPValues = []string{oldVIPValue}
	} else if oldVIPValue != "" &&
		len(*newVIPValues) > 0 &&
		!utilsslice.Contains(*newVIPValues, oldVIPValue) {

		return field.Invalid(fldPath.Child(oldFieldName), oldVIPValue, fmt.Sprintf("%s is deprecated; only %s needs to be specified", oldFieldName, newFieldName))
	}

	return nil
}

// convertAWS upconverts deprecated fields in the AWS platform.
func convertAWS(config *types.InstallConfig) error {
	// Deprecated ExperimentalPropagateUserTag takes precedence when set
	if config.Platform.AWS.ExperimentalPropagateUserTag != nil {
		config.Platform.AWS.PropagateUserTag = *config.Platform.AWS.ExperimentalPropagateUserTag
	}
	return nil
}
