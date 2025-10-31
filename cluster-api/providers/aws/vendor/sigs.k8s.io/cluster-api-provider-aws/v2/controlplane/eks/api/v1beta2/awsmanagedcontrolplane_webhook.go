/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta2

import (
	"context"
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
)

const (
	minAddonVersion          = "v1.18.0"
	minKubeVersionForIPv6    = "v1.21.0"
	minVpcCniVersionForIPv6  = "1.10.2"
	maxClusterNameLength     = 100
	hostnameTypeResourceName = "resource-name"
)

// log is for logging in this package.
var mcpLog = ctrl.Log.WithName("awsmanagedcontrolplane-resource")

const (
	cidrSizeMax    = 65536
	cidrSizeMin    = 16
	vpcCniAddon    = "vpc-cni"
	kubeProxyAddon = "kube-proxy"
)

// SetupWebhookWithManager will setup the webhooks for the AWSManagedControlPlane.
func (r *AWSManagedControlPlane) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsManagedControlPlaneWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-controlplane-cluster-x-k8s-io-v1beta2-awsmanagedcontrolplane,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,versions=v1beta2,name=validation.awsmanagedcontrolplanes.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-controlplane-cluster-x-k8s-io-v1beta2-awsmanagedcontrolplane,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,versions=v1beta2,name=default.awsmanagedcontrolplanes.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsManagedControlPlaneWebhook struct{}

var _ webhook.CustomDefaulter = &awsManagedControlPlaneWebhook{}
var _ webhook.CustomValidator = &awsManagedControlPlaneWebhook{}

func parseEKSVersion(raw string) (*version.Version, error) {
	v, err := version.ParseGeneric(raw)
	if err != nil {
		return nil, err
	}
	return version.MustParseGeneric(fmt.Sprintf("%d.%d", v.Major(), v.Minor())), nil
}

// ValidateCreate will do any extra validation when creating a AWSManagedControlPlane.
func (*awsManagedControlPlaneWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSManagedControlPlane)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedControlPlane object but got %T", r)
	}

	mcpLog.Info("AWSManagedControlPlane validate create", "control-plane", klog.KObj(r))

	var allErrs field.ErrorList

	if r.Spec.EKSClusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.eksClusterName"), "eksClusterName is required"))
	}

	// TODO: Add ipv6 validation things in these validations.
	allErrs = append(allErrs, r.validateEKSVersion(nil)...)
	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateIAMAuthConfig()...)
	allErrs = append(allErrs, r.validateSecondaryCIDR()...)
	allErrs = append(allErrs, r.validateEKSAddons()...)
	allErrs = append(allErrs, r.validateDisableVPCCNI()...)
	allErrs = append(allErrs, r.validateRestrictPrivateSubnets()...)
	allErrs = append(allErrs, r.validateKubeProxy()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.validateNetwork()...)
	allErrs = append(allErrs, r.validatePrivateDNSHostnameTypeOnLaunch()...)
	allErrs = append(allErrs, r.validateAccessConfigCreate()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when updating a AWSManagedControlPlane.
func (*awsManagedControlPlaneWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSManagedControlPlane)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedControlPlane object but got %T", r)
	}

	mcpLog.Info("AWSManagedControlPlane validate update", "control-plane", klog.KObj(r))

	oldAWSManagedControlplane, ok := oldObj.(*AWSManagedControlPlane)
	if !ok {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AWSManagedControlPlane").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.New("failed to convert old AWSManagedControlPlane to object")),
		})
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateEKSClusterName()...)
	allErrs = append(allErrs, r.validateEKSClusterNameSame(oldAWSManagedControlplane)...)
	allErrs = append(allErrs, r.validateEKSVersion(oldAWSManagedControlplane)...)
	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateAccessConfigUpdate(oldAWSManagedControlplane)...)
	allErrs = append(allErrs, r.validateIAMAuthConfig()...)
	allErrs = append(allErrs, r.validateSecondaryCIDR()...)
	allErrs = append(allErrs, r.validateEKSAddons()...)
	allErrs = append(allErrs, r.validateDisableVPCCNI()...)
	allErrs = append(allErrs, r.validateRestrictPrivateSubnets()...)
	allErrs = append(allErrs, r.validateKubeProxy()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.validatePrivateDNSHostnameTypeOnLaunch()...)

	if r.Spec.Region != oldAWSManagedControlplane.Spec.Region {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "region"), r.Spec.Region, "field is immutable"),
		)
	}

	// If encryptionConfig is already set, do not allow removal of it.
	if oldAWSManagedControlplane.Spec.EncryptionConfig != nil && r.Spec.EncryptionConfig == nil {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "encryptionConfig"), r.Spec.EncryptionConfig, "disabling EKS encryption is not allowed after it has been enabled"),
		)
	}

	// If encryptionConfig is already set, do not allow change in provider
	if r.Spec.EncryptionConfig != nil &&
		r.Spec.EncryptionConfig.Provider != nil &&
		oldAWSManagedControlplane.Spec.EncryptionConfig != nil &&
		oldAWSManagedControlplane.Spec.EncryptionConfig.Provider != nil &&
		*r.Spec.EncryptionConfig.Provider != *oldAWSManagedControlplane.Spec.EncryptionConfig.Provider {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "encryptionConfig", "provider"), r.Spec.EncryptionConfig.Provider, "changing EKS encryption is not allowed after it has been enabled"),
		)
	}

	// If a identityRef is already set, do not allow removal of it.
	if oldAWSManagedControlplane.Spec.IdentityRef != nil && r.Spec.IdentityRef == nil {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "identityRef"),
				r.Spec.IdentityRef, "field cannot be set to nil"),
		)
	}

	if oldAWSManagedControlplane.Spec.NetworkSpec.VPC.IsIPv6Enabled() != r.Spec.NetworkSpec.VPC.IsIPv6Enabled() {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "network", "vpc", "enableIPv6"), r.Spec.NetworkSpec.VPC.IsIPv6Enabled(), "changing IP family is not allowed after it has been set"))
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete allows you to add any extra validation when deleting.
func (*awsManagedControlPlaneWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (r *AWSManagedControlPlane) validateEKSClusterName() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.EKSClusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.eksClusterName"), "eksClusterName is required"))
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateEKSClusterNameSame(old *AWSManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList
	if old.Spec.EKSClusterName != "" && r.Spec.EKSClusterName != old.Spec.EKSClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.eksClusterName"), r.Spec.EKSClusterName, "eksClusterName is different to current cluster name"))
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateEKSVersion(old *AWSManagedControlPlane) field.ErrorList {
	path := field.NewPath("spec.version")
	var oldVersion *string
	if old != nil {
		oldVersion = old.Spec.Version
	}
	return validateEKSVersion(r.Spec.Version, oldVersion, r.Spec.NetworkSpec, path)
}

func validateEKSVersion(eksVersion *string, oldVersion *string, networkSpec infrav1.NetworkSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if eksVersion == nil {
		return allErrs
	}

	v, err := parseEKSVersion(*eksVersion)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(path, *eksVersion, err.Error()))
	}

	if oldVersion != nil {
		oldV, err := parseEKSVersion(*oldVersion)
		if err == nil && (v.Major() < oldV.Major() || v.Minor() < oldV.Minor()) {
			allErrs = append(allErrs, field.Invalid(path, *eksVersion, "new version less than old version"))
		}
	}

	if networkSpec.VPC.IsIPv6Enabled() {
		minIPv6, _ := version.ParseSemantic(minKubeVersionForIPv6)
		if v.LessThan(minIPv6) {
			allErrs = append(allErrs, field.Invalid(path, *eksVersion, fmt.Sprintf("IPv6 requires Kubernetes %s or greater", minKubeVersionForIPv6)))
		}
	}
	return allErrs
}

func (r *AWSManagedControlPlane) validateEKSAddons() field.ErrorList {
	return validateEKSAddons(r.Spec.Version, r.Spec.NetworkSpec, r.Spec.Addons, field.NewPath("spec"))
}

func validateEKSAddons(eksVersion *string, networkSpec infrav1.NetworkSpec, addons *[]Addon, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// If not using IPv6 and no addons are specified, return no errors
	if !networkSpec.VPC.IsIPv6Enabled() && (addons == nil || len(*addons) == 0) {
		return allErrs
	}

	// Version is required for addon validation
	if eksVersion == nil {
		return allErrs
	}

	versionPath := path.Child("version")
	v, err := parseEKSVersion(*eksVersion)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(versionPath, *eksVersion, err.Error()))
	}

	minVersion, _ := version.ParseSemantic(minAddonVersion)

	addonsPath := path.Child("addons")

	if v.LessThan(minVersion) {
		message := fmt.Sprintf("addons require Kubernetes %s or greater", minAddonVersion)
		allErrs = append(allErrs, field.Invalid(addonsPath, *eksVersion, message))
	}

	// validations for IPv6:
	// - addons have to be defined in case IPv6 is enabled
	// - minimum version requirement for VPC-CNI using IPv6 ipFamily is 1.10.2
	if networkSpec.VPC.IsIPv6Enabled() {
		if addons == nil || len(*addons) == 0 {
			allErrs = append(allErrs, field.Invalid(addonsPath, "", "addons are required to be set explicitly if IPv6 is enabled"))
			return allErrs
		}

		for _, addon := range *addons {
			if addon.Name == vpcCniAddon {
				v, err := version.ParseGeneric(addon.Version)
				if err != nil {
					allErrs = append(allErrs, field.Invalid(addonsPath, addon.Version, err.Error()))
					break
				}
				minCniVersion, _ := version.ParseSemantic(minVpcCniVersionForIPv6)
				if v.LessThan(minCniVersion) {
					allErrs = append(allErrs, field.Invalid(addonsPath, addon.Version, fmt.Sprintf("vpc-cni version must be above or equal to %s for IPv6", minVpcCniVersionForIPv6)))
					break
				}
			}
		}
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateAccessConfigUpdate(old *AWSManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList

	// If accessConfig is already set, do not allow removal of it.
	if old.Spec.AccessConfig != nil && r.Spec.AccessConfig == nil {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "accessConfig"), r.Spec.AccessConfig, "removing AccessConfig is not allowed after it has been enabled"),
		)
	}

	// AuthenticationMode is ratcheting - do not allow downgrades
	if old.Spec.AccessConfig != nil && r.Spec.AccessConfig != nil &&
		old.Spec.AccessConfig.AuthenticationMode != r.Spec.AccessConfig.AuthenticationMode &&
		((old.Spec.AccessConfig.AuthenticationMode == EKSAuthenticationModeAPIAndConfigMap && r.Spec.AccessConfig.AuthenticationMode == EKSAuthenticationModeConfigMap) ||
			old.Spec.AccessConfig.AuthenticationMode == EKSAuthenticationModeAPI) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "accessConfig", "authenticationMode"), r.Spec.AccessConfig.AuthenticationMode, "downgrading authentication mode is not allowed after it has been enabled"),
		)
	}

	// BootstrapClusterCreatorAdminPermissions only applies on create, but changes should not invalidate updates
	if old.Spec.AccessConfig != nil && r.Spec.AccessConfig != nil &&
		old.Spec.AccessConfig.BootstrapClusterCreatorAdminPermissions != r.Spec.AccessConfig.BootstrapClusterCreatorAdminPermissions {
		mcpLog.Info("Ignoring changes to BootstrapClusterCreatorAdminPermissions on cluster update", "old", old.Spec.AccessConfig.BootstrapClusterCreatorAdminPermissions, "new", r.Spec.AccessConfig.BootstrapClusterCreatorAdminPermissions)
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateAccessConfigCreate() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.AccessConfig != nil {
		if r.Spec.AccessConfig.AuthenticationMode == EKSAuthenticationModeConfigMap &&
			r.Spec.AccessConfig.BootstrapClusterCreatorAdminPermissions != nil &&
			!*r.Spec.AccessConfig.BootstrapClusterCreatorAdminPermissions {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "accessConfig", "bootstrapClusterCreatorAdminPermissions"),
					*r.Spec.AccessConfig.BootstrapClusterCreatorAdminPermissions,
					"bootstrapClusterCreatorAdminPermissions must be true if cluster authentication mode is set to config_map"),
			)
		}
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateIAMAuthConfig() field.ErrorList {
	return validateIAMAuthConfig(r.Spec.IAMAuthenticatorConfig, field.NewPath("spec.iamAuthenticatorConfig"))
}

func validateIAMAuthConfig(cfg *IAMAuthenticatorConfig, parentPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if cfg == nil {
		return allErrs
	}

	for i, userMapping := range cfg.UserMappings {
		usersPath := parentPath.Child(fmt.Sprintf("mapUsers[%d]", i))
		errs := userMapping.Validate()
		for _, validErr := range errs {
			allErrs = append(allErrs, field.Invalid(usersPath, userMapping, validErr.Error()))
		}
	}

	for i, roleMapping := range cfg.RoleMappings {
		rolePath := parentPath.Child(fmt.Sprintf("mapRoles[%d]", i))
		errs := roleMapping.Validate()
		for _, validErr := range errs {
			allErrs = append(allErrs, field.Invalid(rolePath, roleMapping, validErr.Error()))
		}
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateSecondaryCIDR() field.ErrorList {
	return validateSecondaryCIDR(r.Spec.SecondaryCidrBlock, field.NewPath("spec", "secondaryCidrBlock"))
}

func validateSecondaryCIDR(secondaryCidrBlock *string, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if secondaryCidrBlock != nil {
		_, validRange1, _ := net.ParseCIDR("100.64.0.0/10")
		_, validRange2, _ := net.ParseCIDR("198.19.0.0/16")

		_, ipv4Net, err := net.ParseCIDR(*secondaryCidrBlock)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(path, *secondaryCidrBlock, "must be a valid CIDR range"))
			return allErrs
		}

		cidrSize := cidr.AddressCount(ipv4Net)
		if cidrSize > cidrSizeMax || cidrSize < cidrSizeMin {
			allErrs = append(allErrs, field.Invalid(path, *secondaryCidrBlock, "CIDR block sizes must be between a /16 netmask and /28 netmask"))
		}

		start, end := cidr.AddressRange(ipv4Net)
		if (!validRange1.Contains(start) || !validRange1.Contains(end)) && (!validRange2.Contains(start) || !validRange2.Contains(end)) {
			allErrs = append(allErrs, field.Invalid(path, *secondaryCidrBlock, "must be within the 100.64.0.0/10 or 198.19.0.0/16 range"))
		}
	}
	return allErrs
}

func (r *AWSManagedControlPlane) validateKubeProxy() field.ErrorList {
	return validateKubeProxy(r.Spec.KubeProxy, r.Spec.Addons, field.NewPath("spec"))
}

func validateKubeProxy(kubeProxy KubeProxy, addons *[]Addon, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if kubeProxy.Disable {
		disableField := path.Child("kubeProxy", "disable")

		if addons != nil {
			for _, addon := range *addons {
				if addon.Name == kubeProxyAddon {
					allErrs = append(allErrs, field.Invalid(disableField, kubeProxy.Disable, "cannot disable kube-proxy if the kube-proxy addon is specified"))
					break
				}
			}
		}
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateDisableVPCCNI() field.ErrorList {
	return validateDisableVPCCNI(r.Spec.VpcCni, r.Spec.Addons, field.NewPath("spec"))
}

func validateDisableVPCCNI(vpcCni VpcCni, addons *[]Addon, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if vpcCni.Disable {
		disableField := path.Child("vpcCni", "disable")

		if addons != nil {
			for _, addon := range *addons {
				if addon.Name == vpcCniAddon {
					allErrs = append(allErrs, field.Invalid(disableField, vpcCni.Disable, "cannot disable vpc cni if the vpc-cni addon is specified"))
					break
				}
			}
		}
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateRestrictPrivateSubnets() field.ErrorList {
	return validateRestrictPrivateSubnets(r.Spec.RestrictPrivateSubnets, r.Spec.NetworkSpec, r.Spec.EKSClusterName, field.NewPath("spec"))
}

func validateRestrictPrivateSubnets(restrictPrivateSubnets bool, networkSpec infrav1.NetworkSpec, eksClusterName string, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if restrictPrivateSubnets && networkSpec.VPC.IsUnmanaged(eksClusterName) {
		boolField := path.Child("restrictPrivateSubnets")
		if len(networkSpec.Subnets.FilterPrivate()) == 0 {
			allErrs = append(allErrs, field.Invalid(boolField, restrictPrivateSubnets, "cannot enable private subnets restriction when no private subnets are specified"))
		}
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validatePrivateDNSHostnameTypeOnLaunch() field.ErrorList {
	return validatePrivateDNSHostnameTypeOnLaunch(r.Spec.NetworkSpec, field.NewPath("spec"))
}

func validatePrivateDNSHostnameTypeOnLaunch(networkSpec infrav1.NetworkSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if networkSpec.VPC.IsIPv6Enabled() && networkSpec.VPC.PrivateDNSHostnameTypeOnLaunch != nil && *networkSpec.VPC.PrivateDNSHostnameTypeOnLaunch != hostnameTypeResourceName {
		privateDNSHostnameTypeOnLaunchPath := path.Child("networkSpec", "vpc", "privateDNSHostnameTypeOnLaunch")
		allErrs = append(allErrs, field.Invalid(
			privateDNSHostnameTypeOnLaunchPath, networkSpec.VPC.PrivateDNSHostnameTypeOnLaunch,
			fmt.Sprintf("only %s HostnameType can be used in IPv6 mode", hostnameTypeResourceName),
		))
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateNetwork() field.ErrorList {
	return validateNetwork("AWSManagedControlPlane", r.Spec.NetworkSpec, r.Spec.SecondaryCidrBlock, field.NewPath("spec"))
}

func validateNetwork(resourceName string, networkSpec infrav1.NetworkSpec, secondaryCidrBlock *string, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// If only `AWSManagedControlPlane.spec.secondaryCidrBlock` is set, no additional checks are done to remain
	// backward-compatible. The `VPCSpec.SecondaryCidrBlocks` field was added later - if that list is not empty, we
	// require `AWSManagedControlPlane.spec.secondaryCidrBlock` to be listed in there as well. This may allow merging
	// the fields later on.
	secondaryCidrBlocks := networkSpec.VPC.SecondaryCidrBlocks
	secondaryCidrBlocksField := path.Child("network", "vpc", "secondaryCidrBlocks")
	if secondaryCidrBlock != nil && len(secondaryCidrBlocks) > 0 {
		found := false
		for _, cidrBlock := range secondaryCidrBlocks {
			if cidrBlock.IPv4CidrBlock == *secondaryCidrBlock {
				found = true
				break
			}
		}
		if !found {
			allErrs = append(allErrs, field.Invalid(
				secondaryCidrBlocksField, secondaryCidrBlocks,
				fmt.Sprintf("%s.spec.secondaryCidrBlock %v must be listed in %s.spec.network.vpc.secondaryCidrBlocks (required if both fields are filled)", resourceName, *secondaryCidrBlock, resourceName),
			))
		}
	}

	if secondaryCidrBlock != nil && networkSpec.VPC.CidrBlock != "" && networkSpec.VPC.CidrBlock == *secondaryCidrBlock {
		secondaryCidrBlockField := path.Child("vpc", "secondaryCidrBlock")
		allErrs = append(allErrs, field.Invalid(
			secondaryCidrBlockField, secondaryCidrBlocks,
			fmt.Sprintf("%s.spec.secondaryCidrBlock %v must not be equal to the primary %s.spec.network.vpc.cidrBlock", resourceName, *secondaryCidrBlock, resourceName),
		))
	}

	for _, cidrBlock := range secondaryCidrBlocks {
		if networkSpec.VPC.CidrBlock != "" && networkSpec.VPC.CidrBlock == cidrBlock.IPv4CidrBlock {
			allErrs = append(allErrs, field.Invalid(
				secondaryCidrBlocksField, secondaryCidrBlocks,
				fmt.Sprintf("%s.spec.network.vpc.secondaryCidrBlocks must not contain the primary %s.spec.network.vpc.cidrBlock %v", resourceName, resourceName, networkSpec.VPC.CidrBlock),
			))
		}
	}

	// IPv6 validations
	if networkSpec.VPC.IsIPv6Enabled() {
		ipv6Path := path.Child("network", "vpc", "ipv6")

		if networkSpec.VPC.IPv6.CidrBlock != "" && networkSpec.VPC.IPv6.PoolID == "" {
			allErrs = append(allErrs, field.Invalid(
				ipv6Path.Child("poolId"), networkSpec.VPC.IPv6.PoolID,
				"poolId cannot be empty if cidrBlock is set",
			))
		}

		if networkSpec.VPC.IPv6.PoolID != "" && networkSpec.VPC.IPv6.IPAMPool != nil {
			allErrs = append(allErrs, field.Invalid(
				ipv6Path.Child("poolId"), networkSpec.VPC.IPv6.PoolID,
				"poolId and ipamPool cannot be used together",
			))
		}

		if networkSpec.VPC.IPv6.CidrBlock != "" && networkSpec.VPC.IPv6.IPAMPool != nil {
			allErrs = append(allErrs, field.Invalid(
				ipv6Path.Child("cidrBlock"), networkSpec.VPC.IPv6.CidrBlock,
				"cidrBlock and ipamPool cannot be used together",
			))
		}

		if networkSpec.VPC.IPv6.IPAMPool != nil && networkSpec.VPC.IPv6.IPAMPool.ID == "" && networkSpec.VPC.IPv6.IPAMPool.Name == "" {
			allErrs = append(allErrs, field.Invalid(
				ipv6Path.Child("ipamPool"), networkSpec.VPC.IPv6.IPAMPool,
				"ipamPool must have either id or name",
			))
		}
	}

	return allErrs
}

// Default will set default values for the AWSManagedControlPlane.
func (*awsManagedControlPlaneWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSManagedControlPlane)
	if !ok {
		return fmt.Errorf("expected an AWSManagedControlPlane object but got %T", r)
	}

	mcpLog.Info("AWSManagedControlPlane setting defaults", "control-plane", klog.KObj(r))

	if r.Spec.EKSClusterName == "" {
		mcpLog.Info("EKSClusterName is empty, generating name")
		name, err := eks.GenerateEKSName(r.Name, r.Namespace, maxClusterNameLength)
		if err != nil {
			mcpLog.Error(err, "failed to create EKS cluster name")
			return nil
		}

		mcpLog.Info("defaulting EKS cluster name", "cluster", klog.KRef(r.Namespace, name))
		r.Spec.EKSClusterName = name
	}

	if r.Spec.IdentityRef == nil {
		r.Spec.IdentityRef = &infrav1.AWSIdentityReference{
			Kind: infrav1.ControllerIdentityKind,
			Name: infrav1.AWSClusterControllerIdentityName,
		}
	}

	infrav1.SetDefaults_Bastion(&r.Spec.Bastion)
	infrav1.SetDefaults_NetworkSpec(&r.Spec.NetworkSpec)

	// Set default value for BootstrapSelfManagedAddons
	r.Spec.BootstrapSelfManagedAddons = true
	return nil
}
