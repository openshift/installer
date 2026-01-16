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
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
)

// log is for logging in this package.
var log = ctrl.Log.WithName("awsmachine-resource")

func (r *AWSMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsMachineWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachines,versions=v1beta2,name=validation.awsmachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachine,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=awsmachines,versions=v1beta2,name=mawsmachine.kb.io,name=mutation.awsmachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsMachineWebhook struct{}

var (
	_ webhook.CustomValidator = &awsMachineWebhook{}
	_ webhook.CustomDefaulter = &awsMachineWebhook{}
)

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (*awsMachineWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSMachine)
	if !ok {
		return nil, fmt.Errorf("expected an AWSMachine object but got %T", r)
	}

	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateCloudInitSecret()...)
	allErrs = append(allErrs, r.validateIgnitionAndCloudInit()...)
	allErrs = append(allErrs, r.validateRootVolume()...)
	allErrs = append(allErrs, r.validateNonRootVolumes()...)
	allErrs = append(allErrs, r.validateSSHKeyName()...)
	allErrs = append(allErrs, r.validateAdditionalSecurityGroups()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.validateNetworkElasticIPPool()...)
	allErrs = append(allErrs, r.validateInstanceMarketType()...)
	allErrs = append(allErrs, r.validateCapacityReservation()...)
	allErrs = append(allErrs, r.validateHostAllocation()...)

	return nil, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (*awsMachineWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSMachine)
	if !ok {
		return nil, fmt.Errorf("expected an AWSMachine object but got %T", r)
	}

	newAWSMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	if err != nil {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AWSMachine").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert new AWSMachine to unstructured object")),
		})
	}
	oldAWSMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldObj)
	if err != nil {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AWSMachine").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert old AWSMachine to unstructured object")),
		})
	}

	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateCloudInitSecret()...)
	allErrs = append(allErrs, r.validateAdditionalSecurityGroups()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.validateHostAllocation()...)

	newAWSMachineSpec := newAWSMachine["spec"].(map[string]interface{})
	oldAWSMachineSpec := oldAWSMachine["spec"].(map[string]interface{})

	// allow changes to providerID
	delete(oldAWSMachineSpec, "providerID")
	delete(newAWSMachineSpec, "providerID")

	// allow changes to instanceID
	delete(oldAWSMachineSpec, "instanceID")
	delete(newAWSMachineSpec, "instanceID")

	// allow changes to additionalTags
	delete(oldAWSMachineSpec, "additionalTags")
	delete(newAWSMachineSpec, "additionalTags")

	// allow changes to additionalSecurityGroups
	delete(oldAWSMachineSpec, "additionalSecurityGroups")
	delete(newAWSMachineSpec, "additionalSecurityGroups")

	// allow changes to secretPrefix, secretCount, and secureSecretsBackend
	if cloudInit, ok := oldAWSMachineSpec["cloudInit"].(map[string]interface{}); ok {
		delete(cloudInit, "secretPrefix")
		delete(cloudInit, "secretCount")
		delete(cloudInit, "secureSecretsBackend")
	}

	if cloudInit, ok := newAWSMachineSpec["cloudInit"].(map[string]interface{}); ok {
		delete(cloudInit, "secretPrefix")
		delete(cloudInit, "secretCount")
		delete(cloudInit, "secureSecretsBackend")
	}

	// allow changes to enableResourceNameDNSAAAARecord and enableResourceNameDNSARecord
	if privateDNSName, ok := oldAWSMachineSpec["privateDnsName"].(map[string]interface{}); ok {
		delete(privateDNSName, "enableResourceNameDnsAAAARecord")
		delete(privateDNSName, "enableResourceNameDnsARecord")
	}

	if privateDNSName, ok := newAWSMachineSpec["privateDnsName"].(map[string]interface{}); ok {
		delete(privateDNSName, "enableResourceNameDnsAAAARecord")
		delete(privateDNSName, "enableResourceNameDnsARecord")
	}

	if !cmp.Equal(oldAWSMachineSpec, newAWSMachineSpec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "cannot be modified"))
	}

	return nil, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

func (r *AWSMachine) validateCloudInitSecret() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.CloudInit.InsecureSkipSecretsManager {
		if r.Spec.CloudInit.SecretPrefix != "" {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretPrefix"), "cannot be set if spec.cloudInit.insecureSkipSecretsManager is true"))
		}
		if r.Spec.CloudInit.SecretCount != 0 {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretCount"), "cannot be set if spec.cloudInit.insecureSkipSecretsManager is true"))
		}
		if r.Spec.CloudInit.SecureSecretsBackend != "" {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secureSecretsBackend"), "cannot be set if spec.cloudInit.insecureSkipSecretsManager is true"))
		}
	}

	if (r.Spec.CloudInit.SecretPrefix != "") != (r.Spec.CloudInit.SecretCount != 0) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretCount"), "must be set together with spec.CloudInit.SecretPrefix"))
	}

	return allErrs
}

func (r *AWSMachine) cloudInitConfigured() bool {
	configured := false

	configured = configured || r.Spec.CloudInit.SecretPrefix != ""
	configured = configured || r.Spec.CloudInit.SecretCount != 0
	configured = configured || r.Spec.CloudInit.SecureSecretsBackend != ""
	configured = configured || r.Spec.CloudInit.InsecureSkipSecretsManager

	return configured
}

func (r *AWSMachine) ignitionEnabled() bool {
	return r.Spec.Ignition != nil
}

func (r *AWSMachine) validateIgnitionAndCloudInit() field.ErrorList {
	var allErrs field.ErrorList
	if !r.ignitionEnabled() {
		return allErrs
	}

	// Feature gate is not enabled but ignition is enabled then send a forbidden error.
	if !feature.Gates.Enabled(feature.BootstrapFormatIgnition) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "ignition"),
			"can be set only if the BootstrapFormatIgnition feature gate is enabled"))
	}

	// If ignition is enabled, cloudInit should not be configured.
	if r.cloudInitConfigured() {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit"), "cannot be set if spec.ignition is set"))
	}

	// Proxy and TLS are only valid for Ignition versions >= 3.1.
	if r.Spec.Ignition.Version == "2.3" || r.Spec.Ignition.Version == "3.0" {
		if r.Spec.Ignition.Proxy != nil {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "ignition", "proxy"), "cannot be set if spec.ignition.version is 2.3 or 3.0"))
		}
		if r.Spec.Ignition.TLS != nil {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "ignition", "tls"), "cannot be set if spec.ignition.version is 2.3 or 3.0"))
		}
	}

	allErrs = append(allErrs, r.validateIgnitionProxy()...)
	allErrs = append(allErrs, r.validateIgnitionTLS()...)

	return allErrs
}

func (r *AWSMachine) validateIgnitionProxy() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.Ignition.Proxy == nil {
		return allErrs
	}

	// Validate HTTPProxy.
	if r.Spec.Ignition.Proxy.HTTPProxy != nil {
		// Parse the url to check if it is valid.
		_, err := url.Parse(*r.Spec.Ignition.Proxy.HTTPProxy)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "ignition", "proxy", "httpProxy"), *r.Spec.Ignition.Proxy.HTTPProxy, "invalid URL"))
		}
	}

	// Validate HTTPSProxy.
	if r.Spec.Ignition.Proxy.HTTPSProxy != nil {
		// Parse the url to check if it is valid.
		_, err := url.Parse(*r.Spec.Ignition.Proxy.HTTPSProxy)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "ignition", "proxy", "httpsProxy"), *r.Spec.Ignition.Proxy.HTTPSProxy, "invalid URL"))
		}
	}

	// Validate NoProxy.
	for _, noProxy := range r.Spec.Ignition.Proxy.NoProxy {
		noProxy := string(noProxy)
		// Validate here that the value `noProxy` is:
		// - A domain name
		//   - A domain name matches that name and all subdomains
		//   - A domain name with a leading . matches subdomains only

		// A special DNS label (*).
		if noProxy == "*" {
			continue
		}
		// An IP address prefix (1.2.3.4).
		if ip := net.ParseIP(noProxy); ip != nil {
			continue
		}
		// An IP address prefix in CIDR notation (1.2.3.4/8).
		if _, _, err := net.ParseCIDR(noProxy); err == nil {
			continue
		}
		// An IP or domain name with a port.
		if _, _, err := net.SplitHostPort(noProxy); err == nil {
			continue
		}
		// A domain name.
		if noProxy[0] == '.' {
			// If it starts with a dot, it should be a domain name.
			noProxy = noProxy[1:]
		}
		// Validate that the value matches DNS 1123.
		if errs := validation.IsDNS1123Subdomain(noProxy); len(errs) > 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "ignition", "proxy", "noProxy"), noProxy, fmt.Sprintf("invalid noProxy value, please refer to the field documentation: %s", strings.Join(errs, "; "))))
		}
	}

	return allErrs
}

func (r *AWSMachine) validateIgnitionTLS() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.Ignition.TLS == nil {
		return allErrs
	}

	for _, source := range r.Spec.Ignition.TLS.CASources {
		// Validate that source is RFC 2397 data URL.
		u, err := url.Parse(string(source))
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "ignition", "tls", "caSources"), source, "invalid URL"))
		}

		switch u.Scheme {
		case "http", "https", "tftp", "s3", "arn", "gs":
			// Valid schemes.
		case "data":
			// Validate that the data URL is base64 encoded.
			i := strings.Index(u.Opaque, ",")
			if i < 0 {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "ignition", "tls", "caSources"), source, "invalid data URL"))
			}
			// Validate that the data URL is base64 encoded.
			if _, err := base64.StdEncoding.DecodeString(u.Opaque[i+1:]); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "ignition", "tls", "caSources"), source, "invalid base64 encoding for data url"))
			}
		default:
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "ignition", "tls", "caSources"), source, "unsupported URL scheme"))
		}
	}

	return allErrs
}

func (r *AWSMachine) validateRootVolume() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.RootVolume == nil {
		return allErrs
	}

	if VolumeTypesProvisioned.Has(string(r.Spec.RootVolume.Type)) && r.Spec.RootVolume.IOPS == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.rootVolume.iops"), "iops required if type is 'io1' or 'io2'"))
	}

	if r.Spec.RootVolume.Throughput != nil {
		if r.Spec.RootVolume.Type != VolumeTypeGP3 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.rootVolume.throughput"), "throughput is valid only for type 'gp3'"))
		}
		if *r.Spec.RootVolume.Throughput < 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.rootVolume.throughput"), "throughput must be nonnegative"))
		}
	}

	if r.Spec.RootVolume.DeviceName != "" {
		log.Info("root volume shouldn't have a device name (this can be ignored if performing a `clusterctl move`)")
	}

	return allErrs
}

func (r *AWSMachine) validateNetworkElasticIPPool() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.ElasticIPPool == nil {
		return allErrs
	}
	if !ptr.Deref(r.Spec.PublicIP, false) {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.elasticIpPool"), "publicIp must be set to 'true' to assign custom public IPv4 pools with elasticIpPool"))
	}
	eipp := r.Spec.ElasticIPPool
	if eipp.PublicIpv4Pool != nil {
		if eipp.PublicIpv4PoolFallBackOrder == nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.elasticIpPool.publicIpv4PoolFallbackOrder"), r.Spec.ElasticIPPool, "publicIpv4PoolFallbackOrder must be set when publicIpv4Pool is defined."))
		}
		awsPublicIpv4PoolPrefix := "ipv4pool-ec2-"
		if !strings.HasPrefix(*eipp.PublicIpv4Pool, awsPublicIpv4PoolPrefix) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.elasticIpPool.publicIpv4Pool"), r.Spec.ElasticIPPool, fmt.Sprintf("publicIpv4Pool must start with %s.", awsPublicIpv4PoolPrefix)))
		}
	} else if eipp.PublicIpv4PoolFallBackOrder != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.elasticIpPool.publicIpv4PoolFallbackOrder"), r.Spec.ElasticIPPool, "publicIpv4Pool must be set when publicIpv4PoolFallbackOrder is defined."))
	}

	return allErrs
}

func (r *AWSMachine) validateCapacityReservation() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.CapacityReservationID != nil && r.Spec.CapacityReservationPreference != CapacityReservationPreferenceOnly && r.Spec.CapacityReservationPreference != "" {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "capacityReservationPreference"), "when capacityReservationId is specified, capacityReservationPreference may only be 'CapacityReservationsOnly' or empty"))
	}
	if r.Spec.CapacityReservationPreference == CapacityReservationPreferenceOnly && r.Spec.MarketType == MarketTypeSpot {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "capacityReservationPreference"), "when marketType is set to 'Spot', capacityReservationPreference cannot be set to 'CapacityReservationsOnly'"))
	}
	if r.Spec.CapacityReservationPreference == CapacityReservationPreferenceOnly && r.Spec.SpotMarketOptions != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "capacityReservationPreference"), "when capacityReservationPreference is 'CapacityReservationsOnly', spotMarketOptions cannot be set (which implies marketType: 'Spot')"))
	}
	return allErrs
}

func (r *AWSMachine) validateInstanceMarketType() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.MarketType == MarketTypeCapacityBlock && r.Spec.SpotMarketOptions != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "marketType"), "marketType set to CapacityBlock and spotMarketOptions cannot be used together"))
	}
	if r.Spec.MarketType == MarketTypeOnDemand && r.Spec.SpotMarketOptions != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "marketType"), "setting marketType to OnDemand and spotMarketOptions cannot be used together"))
	}
	if r.Spec.MarketType == MarketTypeCapacityBlock && r.Spec.CapacityReservationID == nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "capacityReservationID"), "is required when CapacityBlock is provided"))
	}
	return allErrs
}

func (r *AWSMachine) validateNonRootVolumes() field.ErrorList {
	var allErrs field.ErrorList

	for _, volume := range r.Spec.NonRootVolumes {
		if VolumeTypesProvisioned.Has(string(volume.Type)) && volume.IOPS == 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.nonRootVolumes.iops"), "iops required if type is 'io1' or 'io2'"))
		}

		if volume.Throughput != nil {
			if volume.Type != VolumeTypeGP3 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec.nonRootVolumes.throughput"), "throughput is valid only for type 'gp3'"))
			}
			if *volume.Throughput < 0 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec.nonRootVolumes.throughput"), "throughput must be nonnegative"))
			}
		}

		if volume.DeviceName == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.nonRootVolumes.deviceName"), "non root volume should have device name"))
		}
	}

	return allErrs
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*awsMachineWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// Default implements webhook.Defaulter such that an empty CloudInit will be defined with a default
// SecureSecretsBackend as SecretBackendSecretsManager iff InsecureSkipSecretsManager is unset.
func (*awsMachineWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSMachine)
	if !ok {
		return fmt.Errorf("expected an AWSMachine object but got %T", r)
	}

	if !r.Spec.CloudInit.InsecureSkipSecretsManager && r.Spec.CloudInit.SecureSecretsBackend == "" && !r.ignitionEnabled() {
		r.Spec.CloudInit.SecureSecretsBackend = SecretBackendSecretsManager
	}

	if r.ignitionEnabled() && r.Spec.Ignition.StorageType == "" {
		r.Spec.Ignition.StorageType = DefaultIgnitionStorageType
	}
	// Defaults the version field if StorageType is not set to `UnencryptedUserData`.
	// When using `UnencryptedUserData` the version field is ignored because the userdata defines its version itself.
	if r.ignitionEnabled() && r.Spec.Ignition.Version == "" && r.Spec.Ignition.StorageType != IgnitionStorageTypeOptionUnencryptedUserData {
		r.Spec.Ignition.Version = DefaultIgnitionVersion
	}

	return nil
}

func (r *AWSMachine) validateAdditionalSecurityGroups() field.ErrorList {
	var allErrs field.ErrorList

	for _, additionalSecurityGroup := range r.Spec.AdditionalSecurityGroups {
		if len(additionalSecurityGroup.Filters) > 0 && additionalSecurityGroup.ID != nil {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.additionalSecurityGroups"), "only one of ID or Filters may be specified, specifying both is forbidden"))
		}
	}
	return allErrs
}

func (r *AWSMachine) validateHostAllocation() field.ErrorList {
	var allErrs field.ErrorList

	// Check if both hostID and dynamicHostAllocation are specified
	hasHostID := r.Spec.HostID != nil && len(*r.Spec.HostID) > 0
	hasDynamicHostAllocation := r.Spec.DynamicHostAllocation != nil

	// If both hostID and dynamicHostAllocation are specified, return an error
	if hasHostID && hasDynamicHostAllocation {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.hostID"), "hostID and dynamicHostAllocation are mutually exclusive"), field.Forbidden(field.NewPath("spec.dynamicHostAllocation"), "hostID and dynamicHostAllocation are mutually exclusive"))
	}

	return allErrs
}

func (r *AWSMachine) validateSSHKeyName() field.ErrorList {
	return validateSSHKeyName(r.Spec.SSHKeyName)
}
