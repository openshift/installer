package validation

import (
	"fmt"
	"net"

	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsvalidation "github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/openshift/installer/pkg/types/libvirt"
	libvirtvalidation "github.com/openshift/installer/pkg/types/libvirt/validation"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackvalidation "github.com/openshift/installer/pkg/types/openstack/validation"
	"github.com/openshift/installer/pkg/validate"
)

// ValidateInstallConfig checks that the specified install config is valid.
func ValidateInstallConfig(c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	if c.ObjectMeta.Name == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("metadata", "name"), "cluster name required"))
	}
	if c.ClusterID == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("clusterID"), "cluster ID required"))
	}
	if c.SSHKey != "" {
		if err := validate.SSHPublicKey(c.SSHKey); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("sshKey"), c.SSHKey, err.Error()))
		}
	}
	if err := validate.DomainName(c.BaseDomain); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("baseDomain"), c.BaseDomain, err.Error()))
	}
	allErrs = append(allErrs, validateNetworking(&c.Networking, field.NewPath("networking"))...)
	allErrs = append(allErrs, validateMachinePools(c.Machines, field.NewPath("machines"), c.Platform.Name())...)
	allErrs = append(allErrs, validatePlatform(&c.Platform, field.NewPath("platform"))...)
	if err := validate.ImagePullSecret(c.PullSecret); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("pullSecret"), c.PullSecret, err.Error()))
	}
	return allErrs
}

func validateNetworking(n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if !validate.ValidNetworkTypes[n.Type] {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("type"), n.Type, validate.ValidNetworkTypeValues))
	}
	if err := validate.SubnetCIDR(&n.ServiceCIDR.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceCIDR"), n.ServiceCIDR, err.Error()))
	}
	for i, cn := range n.ClusterNetworks {
		allErrs = append(allErrs, validateClusterNetwork(&cn, fldPath.Child("clusterNetworks").Index(i), &n.ServiceCIDR.IPNet)...)
	}
	if n.PodCIDR != nil {
		if err := validate.SubnetCIDR(&n.PodCIDR.IPNet); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("podCIDR"), n.PodCIDR, err.Error()))
		}
		if validate.DoCIDRsOverlap(&n.ServiceCIDR.IPNet, &n.PodCIDR.IPNet) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("podCIDR"), n.PodCIDR, "podCIDR must not overlap with serviceCIDR"))
		}
	}
	if len(n.ClusterNetworks) == 0 && n.PodCIDR == nil {
		allErrs = append(allErrs, field.Invalid(fldPath, n, "either clusterNetworks or podCIDR is required"))
	}
	if len(n.ClusterNetworks) != 0 && n.PodCIDR != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("podCIDR"), n.PodCIDR, "cannot use podCIDR when clusterNetworks is used"))
	}
	return allErrs
}

func validateClusterNetwork(cn *netopv1.ClusterNetwork, fldPath *field.Path, serviceCIDR *net.IPNet) field.ErrorList {
	allErrs := field.ErrorList{}
	_, cidr, err := net.ParseCIDR(cn.CIDR)
	if err == nil {
		if validate.DoCIDRsOverlap(cidr, serviceCIDR) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR, "cluster network CIDR must not overlap with serviceCIDR"))
		}
		if ones, bits := cidr.Mask.Size(); cn.HostSubnetLength > uint32(bits-ones) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostSubnetLength"), cn.HostSubnetLength, "cluster network host subnet length must not be greater than CIDR length"))
		}
	} else {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR, err.Error()))
	}
	return allErrs
}

func validateMachinePools(pools []types.MachinePool, fldPath *field.Path, platform string) field.ErrorList {
	allErrs := field.ErrorList{}
	poolNames := map[string]bool{}
	for i, p := range pools {
		if poolNames[p.Name] {
			allErrs = append(allErrs, field.Duplicate(fldPath.Index(i), p))
		}
		poolNames[p.Name] = true
		allErrs = append(allErrs, ValidateMachinePool(&p, fldPath.Index(i), platform)...)
	}
	if !poolNames["master"] {
		allErrs = append(allErrs, field.Required(fldPath, "must specify a machine pool with a name of 'master'"))
	}
	if !poolNames["worker"] {
		allErrs = append(allErrs, field.Required(fldPath, "must specify a machine pool with a name of 'worker'"))
	}
	return allErrs
}

func validatePlatform(platform *types.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	activePlatform := ""
	validate := func(n string, value interface{}, validation func(*field.Path) field.ErrorList) {
		if activePlatform != "" {
			allErrs = append(allErrs, field.Invalid(fldPath, platform, fmt.Sprintf("must only specify a single type of platform; cannot use both %q and %q", activePlatform, n)))
		} else {
			activePlatform = n
		}
		allErrs = append(allErrs, validation(fldPath.Child(n))...)
	}
	if platform.AWS != nil {
		validate(aws.Name, platform.AWS, func(f *field.Path) field.ErrorList { return awsvalidation.ValidatePlatform(platform.AWS, f) })
	}
	if platform.Libvirt != nil {
		validate(libvirt.Name, platform.Libvirt, func(f *field.Path) field.ErrorList { return libvirtvalidation.ValidatePlatform(platform.Libvirt, f) })
	}
	if platform.OpenStack != nil {
		validate(openstack.Name, platform.OpenStack, func(f *field.Path) field.ErrorList {
			return openstackvalidation.ValidatePlatform(platform.OpenStack, f)
		})
	}
	if activePlatform == "" {
		allErrs = append(allErrs, field.Invalid(fldPath, platform, "must specify one of the platforms"))
	}
	return allErrs
}
