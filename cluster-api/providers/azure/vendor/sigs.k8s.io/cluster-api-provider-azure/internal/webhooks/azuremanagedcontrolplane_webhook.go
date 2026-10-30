/*
Copyright The Kubernetes Authors.

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

package webhooks

import (
	"context"
	"regexp"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	apiinternal "sigs.k8s.io/cluster-api-provider-azure/internal/api/v1beta1"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

var (
	kubeSemver                 = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)([-0-9a-zA-Z_\.+]*)?$`)
	rMaxNodeProvisionTime      = regexp.MustCompile(`^(\d+)m$`)
	rScaleDownTime             = regexp.MustCompile(`^(\d+)m$`)
	rScaleDownDelayAfterDelete = regexp.MustCompile(`^(\d+)s$`)
	rScanInterval              = regexp.MustCompile(`^(\d+)s$`)
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (mw *AzureManagedControlPlaneWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mw.client = mgr.GetClient()
	mw.logger = mgr.GetLogger().WithName("AzureManagedControlPlane")

	return ctrl.NewWebhookManagedBy(mgr, &infrav1.AzureManagedControlPlane{}).
		WithDefaulter(mw).
		WithValidator(mw).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplane,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanes,verbs=create;update,versions=v1beta1,name=default.azuremanagedcontrolplanes.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// AzureManagedControlPlaneWebhook implements a validating and defaulting webhook for AzureManagedControlPlane.
type AzureManagedControlPlaneWebhook struct {
	client client.Client
	logger logr.Logger
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (mw *AzureManagedControlPlaneWebhook) Default(_ context.Context, m *infrav1.AzureManagedControlPlane) error {
	m.Spec.Version = apiinternal.NormalizeVersion(m.Spec.Version)
	m.Spec.SKU = apiinternal.DefaultSku(mw.logger, m.Spec.SKU)
	m.Spec.FleetsMember = apiinternal.DefaultFleetsMember(m.Spec.FleetsMember, m.Labels)

	if err := apiinternal.SetDefaultAzureManagedControlPlaneSSHPublicKey(m); err != nil {
		ctrl.Log.WithName("AzureManagedControlPlaneWebHookLogger").Error(err, "setDefaultSSHPublicKey failed")
	}

	apiinternal.SetDefaultAzureManagedControlPlaneResourceGroupName(m)
	apiinternal.SetDefaultAzureManagedControlPlaneNodeResourceGroupName(m)
	apiinternal.SetDefaultAzureManagedControlPlaneVirtualNetwork(m)
	apiinternal.SetDefaultAzureManagedControlPlaneSubnet(m)
	apiinternal.SetDefaultAzureManagedControlPlaneOIDCIssuerProfile(m)
	apiinternal.SetDefaultAzureManagedControlPlaneDNSPrefix(m)
	apiinternal.SetDefaultAzureManagedControlPlaneAKSExtensions(m)

	return nil
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplane,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanes,versions=v1beta1,name=validation.azuremanagedcontrolplanes.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureManagedControlPlaneWebhook) ValidateCreate(_ context.Context, m *infrav1.AzureManagedControlPlane) (admission.Warnings, error) {
	return nil, validateAzureManagedControlPlane(m, mw.client)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureManagedControlPlaneWebhook) ValidateUpdate(_ context.Context, old, m *infrav1.AzureManagedControlPlane) (admission.Warnings, error) {
	var allErrs field.ErrorList

	immutableFields := []struct {
		path *field.Path
		old  any
		new  any
	}{
		{field.NewPath("spec", "subscriptionID"), old.Spec.SubscriptionID, m.Spec.SubscriptionID},
		{field.NewPath("spec", "resourceGroupName"), old.Spec.ResourceGroupName, m.Spec.ResourceGroupName},
		{field.NewPath("spec", "nodeResourceGroupName"), old.Spec.NodeResourceGroupName, m.Spec.NodeResourceGroupName},
		{field.NewPath("spec", "location"), old.Spec.Location, m.Spec.Location},
		{field.NewPath("spec", "sshPublicKey"), old.Spec.SSHPublicKey, m.Spec.SSHPublicKey},
		{field.NewPath("spec", "dnsServiceIP"), old.Spec.DNSServiceIP, m.Spec.DNSServiceIP},
		{field.NewPath("spec", "networkPlugin"), old.Spec.NetworkPlugin, m.Spec.NetworkPlugin},
		{field.NewPath("spec", "networkPolicy"), old.Spec.NetworkPolicy, m.Spec.NetworkPolicy},
		{field.NewPath("spec", "networkDataplane"), old.Spec.NetworkDataplane, m.Spec.NetworkDataplane},
		{field.NewPath("spec", "loadBalancerSKU"), old.Spec.LoadBalancerSKU, m.Spec.LoadBalancerSKU},
		{field.NewPath("spec", "httpProxyConfig"), old.Spec.HTTPProxyConfig, m.Spec.HTTPProxyConfig},
		{field.NewPath("spec", "azureEnvironment"), old.Spec.AzureEnvironment, m.Spec.AzureEnvironment},
	}

	for _, f := range immutableFields {
		if err := webhookutils.ValidateImmutable(f.path, f.old, f.new); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	// This nil check is only to streamline tests from having to define this correctly in every test case.
	// Normally, the defaulting webhooks will always set the new DNSPrefix so users can never entirely unset it.
	if m.Spec.DNSPrefix != nil {
		// Pre-1.12 versions of CAPZ do not set this field while 1.12+ defaults it, so emulate the current
		// defaulting here to avoid unrelated updates from failing this immutability check due to the
		// nil -> non-nil transition.
		oldDNSPrefix := old.Spec.DNSPrefix
		if oldDNSPrefix == nil {
			oldDNSPrefix = ptr.To(old.Name)
		}
		if err := webhookutils.ValidateImmutable(
			field.NewPath("spec", "dnsPrefix"),
			oldDNSPrefix,
			m.Spec.DNSPrefix,
		); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	// Consider removing this once moves out of preview
	// Updating outboundType after cluster creation (PREVIEW)
	// https://learn.microsoft.com/en-us/azure/aks/egress-outboundtype#updating-outboundtype-after-cluster-creation-preview
	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "outboundType"),
		old.Spec.OutboundType,
		m.Spec.OutboundType); err != nil {
		allErrs = append(allErrs, err)
	}

	if errs := validateAzureManagedControlPlaneVirtualNetworkUpdate(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneAddonProfilesUpdate(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneAPIServerAccessProfileUpdate(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneNetworkPluginModeUpdate(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneAADProfileUpdateAndLocalAccounts(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneAutoUpgradeProfile(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneK8sVersionUpdate(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneOIDCIssuerProfileUpdate(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneFleetsMemberUpdate(m, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAKSExtensionsUpdate(old.Spec.Extensions, m.Spec.Extensions); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneClassSpecSecurityProfileUpdate(&m.Spec.AzureManagedControlPlaneClassSpec, &old.Spec.AzureManagedControlPlaneClassSpec); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, validateAzureManagedControlPlane(m, mw.client)
	}

	return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind(infrav1.AzureManagedControlPlaneKind).GroupKind(), m.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureManagedControlPlaneWebhook) ValidateDelete(_ context.Context, _ *infrav1.AzureManagedControlPlane) (admission.Warnings, error) {
	return nil, nil
}
