package conversion

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilsslice "k8s.io/utils/strings/slices"

	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	powervcconversion "github.com/openshift/installer/pkg/types/conversion/powervc"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervc"
	"github.com/openshift/installer/pkg/types/vsphere"
	vsphereconversion "github.com/openshift/installer/pkg/types/vsphere/conversion"
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
	case powervc.Name:
		// The first thing on the agenda is to convert the PowerVC install config
		// to an underlying OpenStack install config.
		if err := powervcconversion.ConvertPowerVCInstallConfig(config); err != nil {
			return err
		}
		if err := convertOpenStack(config); err != nil {
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
	case azure.Name:
		if err := convertAzure(config); err != nil {
			return err
		}
	case vsphere.Name:
		if err := vsphereconversion.ConvertInstallConfig(config); err != nil {
			return err
		}
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

	// Recognize the network plugin name regardless of capitalization, for
	// backward compatibility
	if strings.EqualFold(netconf.NetworkType, string(operv1.NetworkTypeOpenShiftSDN)) {
		netconf.NetworkType = string(operv1.NetworkTypeOpenShiftSDN)
	}
	if strings.EqualFold(netconf.NetworkType, string(operv1.NetworkTypeOVNKubernetes)) {
		netconf.NetworkType = string(operv1.NetworkTypeOVNKubernetes)
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

	// type has been deprecated in favor of types in the machinePools.
	if config.ControlPlane != nil &&
		config.ControlPlane.Platform.OpenStack != nil &&
		config.ControlPlane.Platform.OpenStack.RootVolume != nil &&
		config.ControlPlane.Platform.OpenStack.RootVolume.DeprecatedType != "" {
		if len(config.ControlPlane.Platform.OpenStack.RootVolume.Types) > 0 {
			// Return error if both type and types of rootVolume are specified in the config
			return field.Forbidden(field.NewPath("controlPlane").Child("platform").Child("openstack").Child("rootVolume").Child("type"), "cannot specify type and types in rootVolume together")
		}
		config.ControlPlane.Platform.OpenStack.RootVolume.Types = []string{config.ControlPlane.Platform.OpenStack.RootVolume.DeprecatedType}
		config.ControlPlane.Platform.OpenStack.RootVolume.DeprecatedType = ""
	}
	for _, pool := range config.Compute {
		mpool := pool.Platform.OpenStack
		if mpool != nil && mpool.RootVolume != nil && mpool.RootVolume.DeprecatedType != "" {
			if mpool.RootVolume.Types != nil && len(mpool.RootVolume.Types) > 0 {
				// Return error if both type and types of rootVolume are specified in the config
				return field.Forbidden(field.NewPath("compute").Child("platform").Child("openstack").Child("rootVolume").Child("type"), "cannot specify type and types in rootVolume together")
			}
			mpool.RootVolume.Types = []string{mpool.RootVolume.DeprecatedType}
			mpool.RootVolume.DeprecatedType = ""
		}
	}
	if config.Platform.OpenStack.DefaultMachinePlatform != nil && config.Platform.OpenStack.DefaultMachinePlatform.RootVolume != nil && config.Platform.OpenStack.DefaultMachinePlatform.RootVolume.DeprecatedType != "" {
		if len(config.Platform.OpenStack.DefaultMachinePlatform.RootVolume.Types) > 0 {
			// Return error if both type and types of defaultMachinePlatform are specified in the config
			return field.Forbidden(field.NewPath("platform").Child("openstack").Child("type"), "cannot specify type and types in defaultMachinePlatform together")
		}
		config.Platform.OpenStack.DefaultMachinePlatform.RootVolume.Types = []string{config.Platform.OpenStack.DefaultMachinePlatform.RootVolume.DeprecatedType}
		config.Platform.OpenStack.DefaultMachinePlatform.RootVolume.DeprecatedType = ""
	}

	if err := upconvertVIP(&config.Platform.OpenStack.APIVIPs, config.Platform.OpenStack.DeprecatedAPIVIP, "apiVIP", "apiVIPs", field.NewPath("platform").Child("openstack")); err != nil {
		return err
	}

	if err := upconvertVIP(&config.Platform.OpenStack.IngressVIPs, config.Platform.OpenStack.DeprecatedIngressVIP, "ingressVIP", "ingressVIPs", field.NewPath("platform").Child("openstack")); err != nil {
		return err
	}

	// machinesSubnet has been deprecated in favor of ControlPlanePort
	controlPlanePort := config.Platform.OpenStack.ControlPlanePort
	deprecatedMachinesSubnet := config.Platform.OpenStack.DeprecatedMachinesSubnet
	if deprecatedMachinesSubnet != "" && controlPlanePort == nil {
		fixedIPs := []openstack.FixedIP{{Subnet: openstack.SubnetFilter{ID: deprecatedMachinesSubnet}}}
		config.Platform.OpenStack.ControlPlanePort = &openstack.PortTarget{FixedIPs: fixedIPs}
	} else if deprecatedMachinesSubnet != "" &&
		controlPlanePort != nil {
		if !(len(controlPlanePort.FixedIPs) == 1 && controlPlanePort.FixedIPs[0].Subnet.ID == deprecatedMachinesSubnet) {
			return field.Invalid(field.NewPath("platform").Child("openstack").Child("machinesSubnet"), deprecatedMachinesSubnet, fmt.Sprintf("%s is deprecated; only %s needs to be specified", "machinesSubnet", "controlPlanePort"))
		}
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
	// BestEffortDeleteIgnition takes precedence when set
	if !config.AWS.BestEffortDeleteIgnition {
		config.AWS.BestEffortDeleteIgnition = config.AWS.PreserveBootstrapIgnition
	}
	if ami := config.AWS.AMIID; len(ami) > 0 {
		if config.AWS.DefaultMachinePlatform == nil {
			config.AWS.DefaultMachinePlatform = &aws.MachinePool{}
		}
		// DefaultMachinePlatform.AMIID takes precedence in the machine manifest code anyway
		if len(config.AWS.DefaultMachinePlatform.AMIID) == 0 {
			config.AWS.DefaultMachinePlatform.AMIID = ami
		}
	}

	// Subnets field is deprecated in favor of VPC.Subnets.
	fldPath := field.NewPath("platform", "aws")
	if len(config.AWS.DeprecatedSubnets) > 0 && len(config.AWS.VPC.Subnets) > 0 { // nolint: staticcheck
		return field.Forbidden(fldPath.Child("subnets"), fmt.Sprintf("cannot specify %s and %s together", fldPath.Child("subnets"), fldPath.Child("vpc", "subnets")))
	} else if len(config.AWS.DeprecatedSubnets) > 0 { // nolint: staticcheck
		var subnets []aws.Subnet
		for _, subnetID := range config.AWS.DeprecatedSubnets { // nolint: staticcheck
			subnets = append(subnets, aws.Subnet{
				ID: aws.AWSSubnetID(subnetID),
			})
		}
		config.AWS.VPC.Subnets = subnets
		logrus.Warnf("%s is deprecated. Converted to %s", fldPath.Child("subnets"), fldPath.Child("vpc", "subnets"))
	}

	return nil
}

func convertAzure(config *types.InstallConfig) error {
	subnets := config.Azure.Subnets
	if len(subnets) == 0 {
		subnets = []azure.SubnetSpec{}
	}
	if config.Azure.DeprecatedControlPlaneSubnet != "" { // nolint: staticcheck
		subnets = append(subnets, azure.SubnetSpec{
			Name: config.Azure.DeprecatedControlPlaneSubnet, // nolint: staticcheck
			Role: azure.SubnetControlPlane,
		})
	}
	if config.Azure.DeprecatedComputeSubnet != "" { // nolint: staticcheck
		subnets = append(subnets, azure.SubnetSpec{
			Name: config.Azure.DeprecatedComputeSubnet, // nolint: staticcheck
			Role: azure.SubnetNode,
		})
	}
	config.Azure.Subnets = subnets
	return nil
}
