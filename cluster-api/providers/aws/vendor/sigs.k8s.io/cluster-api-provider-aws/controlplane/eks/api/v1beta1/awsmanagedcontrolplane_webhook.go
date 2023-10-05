/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/version"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/eks"
)

const (
	minAddonVersion      = "v1.18.0"
	maxClusterNameLength = 100
)

// log is for logging in this package.
var mcpLog = logf.Log.WithName("awsmanagedcontrolplane-resource")

const (
	cidrSizeMax    = 65536
	cidrSizeMin    = 16
	vpcCniAddon    = "vpc-cni"
	kubeProxyAddon = "kube-proxy"
)

// SetupWebhookWithManager will setup the webhooks for the AWSManagedControlPlane.
func (r *AWSManagedControlPlane) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-controlplane-cluster-x-k8s-io-v1beta1-awsmanagedcontrolplane,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,versions=v1beta1,name=validation.awsmanagedcontrolplanes.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-controlplane-cluster-x-k8s-io-v1beta1-awsmanagedcontrolplane,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,versions=v1beta1,name=default.awsmanagedcontrolplanes.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Defaulter = &AWSManagedControlPlane{}
var _ webhook.Validator = &AWSManagedControlPlane{}

func parseEKSVersion(raw string) (*version.Version, error) {
	v, err := version.ParseGeneric(raw)
	if err != nil {
		return nil, err
	}
	return version.MustParseGeneric(fmt.Sprintf("%d.%d", v.Major(), v.Minor())), nil
}

func normalizeVersion(raw string) (string, error) {
	// Normalize version (i.e. remove patch, add "v" prefix) if necessary
	eksV, err := parseEKSVersion(raw)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("v%d.%d", eksV.Major(), eksV.Minor()), nil
}

// ValidateCreate will do any extra validation when creating a AWSManagedControlPlane.
func (r *AWSManagedControlPlane) ValidateCreate() error {
	mcpLog.Info("AWSManagedControlPlane validate create", "name", r.Name)

	var allErrs field.ErrorList

	if r.Spec.EKSClusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.eksClusterName"), "eksClusterName is required"))
	}

	allErrs = append(allErrs, r.validateEKSVersion(nil)...)
	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateIAMAuthConfig()...)
	allErrs = append(allErrs, r.validateSecondaryCIDR()...)
	allErrs = append(allErrs, r.validateEKSAddons()...)
	allErrs = append(allErrs, r.validateDisableVPCCNI()...)
	allErrs = append(allErrs, r.validateKubeProxy()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when updating a AWSManagedControlPlane.
func (r *AWSManagedControlPlane) ValidateUpdate(old runtime.Object) error {
	mcpLog.Info("AWSManagedControlPlane validate update", "name", r.Name)
	oldAWSManagedControlplane, ok := old.(*AWSManagedControlPlane)
	if !ok {
		return apierrors.NewInvalid(GroupVersion.WithKind("AWSManagedControlPlane").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.New("failed to convert old AWSManagedControlPlane to object")),
		})
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateEKSClusterName()...)
	allErrs = append(allErrs, r.validateEKSClusterNameSame(oldAWSManagedControlplane)...)
	allErrs = append(allErrs, r.validateEKSVersion(oldAWSManagedControlplane)...)
	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateIAMAuthConfig()...)
	allErrs = append(allErrs, r.validateSecondaryCIDR()...)
	allErrs = append(allErrs, r.validateEKSAddons()...)
	allErrs = append(allErrs, r.validateDisableVPCCNI()...)
	allErrs = append(allErrs, r.validateKubeProxy()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

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

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete allows you to add any extra validation when deleting.
func (r *AWSManagedControlPlane) ValidateDelete() error {
	mcpLog.Info("AWSManagedControlPlane validate delete", "name", r.Name)

	return nil
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
	var allErrs field.ErrorList

	if r.Spec.Version == nil {
		return allErrs
	}

	v, err := parseEKSVersion(*r.Spec.Version)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(path, *r.Spec.Version, err.Error()))
	}

	if old != nil {
		oldV, err := parseEKSVersion(*old.Spec.Version)
		if err == nil && (v.Major() < oldV.Major() || v.Minor() < oldV.Minor()) {
			allErrs = append(allErrs, field.Invalid(path, *r.Spec.Version, "new version less than old version"))
		}
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateEKSAddons() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.Addons == nil || len(*r.Spec.Addons) == 0 {
		return allErrs
	}

	path := field.NewPath("spec.version")
	v, err := parseEKSVersion(*r.Spec.Version)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(path, *r.Spec.Version, err.Error()))
	}

	minVersion, _ := version.ParseSemantic(minAddonVersion)

	addonsPath := field.NewPath("spec.addons")

	if v.LessThan(minVersion) {
		message := fmt.Sprintf("addons requires Kubernetes %s or greater", minAddonVersion)
		allErrs = append(allErrs, field.Invalid(addonsPath, *r.Spec.Version, message))
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateIAMAuthConfig() field.ErrorList {
	var allErrs field.ErrorList

	parentPath := field.NewPath("spec.iamAuthenticatorConfig")

	cfg := r.Spec.IAMAuthenticatorConfig
	if cfg == nil {
		return allErrs
	}

	for i, userMapping := range cfg.UserMappings {
		usersPathName := fmt.Sprintf("mapUsers[%d]", i)
		usersPath := parentPath.Child(usersPathName)
		errs := userMapping.Validate()
		for _, validErr := range errs {
			allErrs = append(allErrs, field.Invalid(usersPath, userMapping, validErr.Error()))
		}
	}

	for i, roleMapping := range cfg.RoleMappings {
		rolePathName := fmt.Sprintf("mapRoles[%d]", i)
		rolePath := parentPath.Child(rolePathName)
		errs := roleMapping.Validate()
		for _, validErr := range errs {
			allErrs = append(allErrs, field.Invalid(rolePath, roleMapping, validErr.Error()))
		}
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateSecondaryCIDR() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.SecondaryCidrBlock != nil {
		cidrField := field.NewPath("spec", "secondaryCidrBlock")
		_, validRange1, _ := net.ParseCIDR("100.64.0.0/10")
		_, validRange2, _ := net.ParseCIDR("198.19.0.0/16")

		_, ipv4Net, err := net.ParseCIDR(*r.Spec.SecondaryCidrBlock)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(cidrField, *r.Spec.SecondaryCidrBlock, "must be valid CIDR range"))
			return allErrs
		}

		cidrSize := cidr.AddressCount(ipv4Net)
		if cidrSize > cidrSizeMax || cidrSize < cidrSizeMin {
			allErrs = append(allErrs, field.Invalid(cidrField, *r.Spec.SecondaryCidrBlock, "CIDR block sizes must be between a /16 netmask and /28 netmask"))
		}

		start, end := cidr.AddressRange(ipv4Net)
		if (!validRange1.Contains(start) || !validRange1.Contains(end)) && (!validRange2.Contains(start) || !validRange2.Contains(end)) {
			allErrs = append(allErrs, field.Invalid(cidrField, *r.Spec.SecondaryCidrBlock, "must be within the 100.64.0.0/10 or 198.19.0.0/16 range"))
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func (r *AWSManagedControlPlane) validateKubeProxy() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.KubeProxy.Disable {
		disableField := field.NewPath("spec", "kubeProxy", "disable")

		if r.Spec.Addons != nil {
			for _, addon := range *r.Spec.Addons {
				if addon.Name == kubeProxyAddon {
					allErrs = append(allErrs, field.Invalid(disableField, r.Spec.KubeProxy.Disable, "cannot disable kube-proxy if the kube-proxy addon is specified"))
					break
				}
			}
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func (r *AWSManagedControlPlane) validateDisableVPCCNI() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.DisableVPCCNI {
		disableField := field.NewPath("spec", "disableVPCCNI")

		if r.Spec.Addons != nil {
			for _, addon := range *r.Spec.Addons {
				if addon.Name == vpcCniAddon {
					allErrs = append(allErrs, field.Invalid(disableField, r.Spec.DisableVPCCNI, "cannot disable vpc cni if the vpc-cni addon is specified"))
					break
				}
			}
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

// Default will set default values for the AWSManagedControlPlane.
func (r *AWSManagedControlPlane) Default() {
	mcpLog.Info("AWSManagedControlPlane setting defaults", "name", r.Name)

	if r.Spec.EKSClusterName == "" {
		mcpLog.Info("EKSClusterName is empty, generating name")
		name, err := eks.GenerateEKSName(r.Name, r.Namespace, maxClusterNameLength)
		if err != nil {
			mcpLog.Error(err, "failed to create EKS cluster name")
			return
		}

		mcpLog.Info("defaulting EKS cluster name", "cluster-name", name)
		r.Spec.EKSClusterName = name
	}

	if r.Spec.IdentityRef == nil {
		r.Spec.IdentityRef = &infrav1.AWSIdentityReference{
			Kind: infrav1.ControllerIdentityKind,
			Name: infrav1.AWSClusterControllerIdentityName,
		}
	}

	// Normalize version (i.e. remove patch, add "v" prefix) if necessary
	if r.Spec.Version != nil {
		normalizedV, err := normalizeVersion(*r.Spec.Version)
		if err != nil {
			mcpLog.Error(err, "couldn't parse version")
			return
		}
		r.Spec.Version = &normalizedV
	}

	infrav1.SetDefaults_Bastion(&r.Spec.Bastion)
	infrav1.SetDefaults_NetworkSpec(&r.Spec.NetworkSpec)
}
