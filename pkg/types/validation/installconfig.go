package validation

import (
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
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

const (
	masterPoolName = "master"
)

// ClusterDomain returns the cluster domain for a cluster with the specified
// base domain and cluster name.
func ClusterDomain(baseDomain, clusterName string) string {
	return fmt.Sprintf("%s.%s", clusterName, baseDomain)
}

// ValidateInstallConfig checks that the specified install config is valid.
func ValidateInstallConfig(c *types.InstallConfig, openStackValidValuesFetcher openstackvalidation.ValidValuesFetcher) field.ErrorList {
	allErrs := field.ErrorList{}
	if c.TypeMeta.APIVersion == "" {
		return field.ErrorList{field.Required(field.NewPath("apiVersion"), "install-config version required")}
	}
	switch v := c.APIVersion; v {
	case types.InstallConfigVersion:
		// Current version
	case "v1beta1", "v1beta2":
		logrus.Warnf("install-config.yaml is using a deprecated version %q. The expected version is %q.", v, types.InstallConfigVersion)
	default:
		return field.ErrorList{field.Invalid(field.NewPath("apiVersion"), c.TypeMeta.APIVersion, fmt.Sprintf("install-config version must be %q", types.InstallConfigVersion))}
	}
	if c.SSHKey != "" {
		if err := validate.SSHPublicKey(c.SSHKey); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("sshKey"), c.SSHKey, err.Error()))
		}
	}
	nameErr := validate.DomainName(c.ObjectMeta.Name, false)
	if nameErr != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), c.ObjectMeta.Name, nameErr.Error()))
	}
	baseDomainErr := validate.DomainName(c.BaseDomain, true)
	if baseDomainErr != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("baseDomain"), c.BaseDomain, baseDomainErr.Error()))
	}
	if nameErr == nil && baseDomainErr == nil {
		clusterDomain := ClusterDomain(c.BaseDomain, c.ObjectMeta.Name)
		if err := validate.DomainName(clusterDomain, true); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("baseDomain"), clusterDomain, err.Error()))
		}
	}
	if c.Networking != nil {
		allErrs = append(allErrs, validateNetworking(c.Networking, field.NewPath("networking"))...)
	} else {
		allErrs = append(allErrs, field.Required(field.NewPath("networking"), "networking is required"))
	}
	if c.ControlPlane != nil {
		allErrs = append(allErrs, validateControlPlane(c.ControlPlane, field.NewPath("controlPlane"), c.Platform.Name())...)
	} else {
		allErrs = append(allErrs, field.Required(field.NewPath("controlPlane"), "controlPlane is required"))
	}
	allErrs = append(allErrs, validateCompute(c.Compute, field.NewPath("compute"), c.Platform.Name())...)
	allErrs = append(allErrs, validatePlatform(&c.Platform, field.NewPath("platform"), openStackValidValuesFetcher)...)
	if err := validate.ImagePullSecret(c.PullSecret); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("pullSecret"), c.PullSecret, err.Error()))
	}
	return allErrs
}

func validateNetworking(n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if n.Type == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("type"), "network provider type required"))
	}
	if n.MachineCIDR != nil {
		if err := validate.SubnetCIDR(&n.MachineCIDR.IPNet); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("machineCIDR"), n.MachineCIDR.String(), err.Error()))
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("machineCIDR"), "a machine CIDR is required"))
	}
	if n.ServiceCIDR != nil {
		if err := validate.SubnetCIDR(&n.ServiceCIDR.IPNet); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceCIDR"), n.ServiceCIDR.String(), err.Error()))
		}
		for i, cn := range n.ClusterNetworks {
			allErrs = append(allErrs, validateClusterNetwork(&cn, fldPath.Child("clusterNetworks").Index(i), &n.ServiceCIDR.IPNet)...)
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("serviceCIDR"), "a service CIDR is required"))
	}
	if len(n.ClusterNetworks) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("clusterNetworks"), "cluster network required"))
	}
	return allErrs
}

func validateClusterNetwork(cn *types.ClusterNetworkEntry, fldPath *field.Path, serviceCIDR *net.IPNet) field.ErrorList {
	allErrs := field.ErrorList{}
	if err := validate.SubnetCIDR(&cn.CIDR.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR.IPNet.String(), err.Error()))
	}
	if validate.DoCIDRsOverlap(&cn.CIDR.IPNet, serviceCIDR) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cidr"), cn.CIDR.String(), "cluster network CIDR must not overlap with serviceCIDR"))
	}
	if cn.HostSubnetLength < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("hostSubnetLength"), cn.HostSubnetLength, "hostSubnetLength must be positive"))
	}
	if ones, bits := cn.CIDR.Mask.Size(); cn.HostSubnetLength > int32(bits-ones) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("hostSubnetLength"), cn.HostSubnetLength, "cluster network host subnet must not be larger than CIDR "+cn.CIDR.String()))
	}
	return allErrs
}

func validateControlPlane(pool *types.MachinePool, fldPath *field.Path, platform string) field.ErrorList {
	allErrs := field.ErrorList{}
	if pool.Name != masterPoolName {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("name"), pool.Name, []string{masterPoolName}))
	}
	if pool.Replicas != nil && *pool.Replicas == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), pool.Replicas, "number of control plane replicas must be positive"))
	}
	allErrs = append(allErrs, ValidateMachinePool(pool, fldPath, platform)...)
	return allErrs
}

func validateCompute(pools []types.MachinePool, fldPath *field.Path, platform string) field.ErrorList {
	allErrs := field.ErrorList{}
	poolNames := map[string]bool{}
	foundPositiveReplicas := false
	for i, p := range pools {
		poolFldPath := fldPath.Index(i)
		if p.Name != "worker" {
			allErrs = append(allErrs, field.NotSupported(poolFldPath.Child("name"), p.Name, []string{"worker"}))
		}
		if poolNames[p.Name] {
			allErrs = append(allErrs, field.Duplicate(poolFldPath.Child("name"), p.Name))
		}
		poolNames[p.Name] = true
		if p.Replicas != nil && *p.Replicas > 0 {
			foundPositiveReplicas = true
		}
		allErrs = append(allErrs, ValidateMachinePool(&p, poolFldPath, platform)...)
	}
	if !foundPositiveReplicas {
		logrus.Warnf("There are no compute nodes specified. The cluster will not fully initialize without compute nodes.")
	}
	return allErrs
}

func validatePlatform(platform *types.Platform, fldPath *field.Path, openStackValidValuesFetcher openstackvalidation.ValidValuesFetcher) field.ErrorList {
	allErrs := field.ErrorList{}
	activePlatform := platform.Name()
	platforms := make([]string, len(types.PlatformNames))
	copy(platforms, types.PlatformNames)
	platforms = append(platforms, types.HiddenPlatformNames...)
	sort.Strings(platforms)
	i := sort.SearchStrings(platforms, activePlatform)
	if i == len(platforms) || platforms[i] != activePlatform {
		allErrs = append(allErrs, field.Invalid(fldPath, platform, fmt.Sprintf("must specify one of the platforms (%s)", strings.Join(platforms, ", "))))
	}
	validate := func(n string, value interface{}, validation func(*field.Path) field.ErrorList) {
		if n != activePlatform {
			allErrs = append(allErrs, field.Invalid(fldPath, platform, fmt.Sprintf("must only specify a single type of platform; cannot use both %q and %q", activePlatform, n)))
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
			return openstackvalidation.ValidatePlatform(platform.OpenStack, f, openStackValidValuesFetcher)
		})
	}
	return allErrs
}
