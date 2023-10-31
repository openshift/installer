/*
Copyright 2023 The Kubernetes Authors.

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
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/maps"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	capifeature "sigs.k8s.io/cluster-api/feature"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var validNodePublicPrefixID = regexp.MustCompile(`(?i)^/?subscriptions/[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}/resourcegroups/[^/]+/providers/microsoft\.network/publicipprefixes/[^/]+$`)

// SetupAzureManagedMachinePoolWebhookWithManager sets up and registers the webhook with the manager.
func SetupAzureManagedMachinePoolWebhookWithManager(mgr ctrl.Manager) error {
	mw := &azureManagedMachinePoolWebhook{Client: mgr.GetClient()}
	return ctrl.NewWebhookManagedBy(mgr).
		For(&AzureManagedMachinePool{}).
		WithDefaulter(mw).
		WithValidator(mw).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedmachinepools,verbs=create;update,versions=v1beta1,name=default.azuremanagedmachinepools.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// azureManagedMachinePoolWebhook implements a validating and defaulting webhook for AzureManagedMachinePool.
type azureManagedMachinePoolWebhook struct {
	Client client.Client
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (mw *azureManagedMachinePoolWebhook) Default(ctx context.Context, obj runtime.Object) error {
	m, ok := obj.(*AzureManagedMachinePool)
	if !ok {
		return apierrors.NewBadRequest("expected an AzureManagedMachinePool")
	}
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}
	m.Labels[LabelAgentPoolMode] = m.Spec.Mode

	if m.Spec.Name == nil || *m.Spec.Name == "" {
		m.Spec.Name = &m.Name
	}

	if m.Spec.OSType == nil {
		m.Spec.OSType = ptr.To(DefaultOSType)
	}

	if ptr.Deref(m.Spec.ScaleSetPriority, "") == string(armcontainerservice.ScaleSetPrioritySpot) && m.Spec.SpotMaxPrice == nil {
		m.Spec.SpotMaxPrice = ptr.To(resource.MustParse("-1"))
	}

	return nil
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedmachinepools,versions=v1beta1,name=validation.azuremanagedmachinepools.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedMachinePoolWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	m, ok := obj.(*AzureManagedMachinePool)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedMachinePool")
	}
	// NOTE: AzureManagedMachinePool relies upon MachinePools, which is behind a feature gate flag.
	// The webhook must prevent creating new objects in case the feature flag is disabled.
	if !feature.Gates.Enabled(capifeature.MachinePool) {
		return nil, field.Forbidden(
			field.NewPath("spec"),
			"can be set only if the Cluster API 'MachinePool' feature flag is enabled",
		)
	}
	validators := []func() error{
		m.validateMaxPods,
		m.validateOSType,
		m.validateName,
		m.validateNodeLabels,
		m.validateNodePublicIPPrefixID,
		m.validateEnableNodePublicIP,
		m.validateKubeletConfig,
		m.validateLinuxOSConfig,
		m.validateSubnetName,
	}

	var errs []error
	for _, validator := range validators {
		if err := validator(); err != nil {
			errs = append(errs, err)
		}
	}

	return nil, kerrors.NewAggregate(errs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedMachinePoolWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	old, ok := oldObj.(*AzureManagedMachinePool)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedMachinePool")
	}
	m, ok := newObj.(*AzureManagedMachinePool)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedMachinePool")
	}
	var allErrs field.ErrorList

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "Name"),
		old.Spec.Name,
		m.Spec.Name); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := m.validateNodeLabels(); err != nil {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "NodeLabels"),
				m.Spec.NodeLabels,
				err.Error()))
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "OSType"),
		old.Spec.OSType,
		m.Spec.OSType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "SKU"),
		old.Spec.SKU,
		m.Spec.SKU); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "OSDiskSizeGB"),
		old.Spec.OSDiskSizeGB,
		m.Spec.OSDiskSizeGB); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "SubnetName"),
		old.Spec.SubnetName,
		m.Spec.SubnetName); err != nil && old.Spec.SubnetName != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "EnableFIPS"),
		old.Spec.EnableFIPS,
		m.Spec.EnableFIPS); err != nil && old.Spec.EnableFIPS != nil {
		allErrs = append(allErrs, err)
	}

	// custom headers are immutable
	oldCustomHeaders := maps.FilterByKeyPrefix(old.ObjectMeta.Annotations, CustomHeaderPrefix)
	newCustomHeaders := maps.FilterByKeyPrefix(m.ObjectMeta.Annotations, CustomHeaderPrefix)
	if !reflect.DeepEqual(oldCustomHeaders, newCustomHeaders) {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("metadata", "annotations"),
				m.ObjectMeta.Annotations,
				fmt.Sprintf("annotations with '%s' prefix are immutable", CustomHeaderPrefix)))
	}

	if !webhookutils.EnsureStringSlicesAreEquivalent(m.Spec.AvailabilityZones, old.Spec.AvailabilityZones) {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "AvailabilityZones"),
				m.Spec.AvailabilityZones,
				"field is immutable"))
	}

	if m.Spec.Mode != string(NodePoolModeSystem) && old.Spec.Mode == string(NodePoolModeSystem) {
		// validate for last system node pool
		if err := m.validateLastSystemNodePool(mw.Client); err != nil {
			allErrs = append(allErrs, field.Forbidden(
				field.NewPath("Spec", "Mode"),
				"Cannot change node pool mode to User, you must have at least one System node pool in your cluster"))
		}
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "MaxPods"),
		old.Spec.MaxPods,
		m.Spec.MaxPods); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "OsDiskType"),
		old.Spec.OsDiskType,
		m.Spec.OsDiskType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "ScaleSetPriority"),
		old.Spec.ScaleSetPriority,
		m.Spec.ScaleSetPriority); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "EnableUltraSSD"),
		old.Spec.EnableUltraSSD,
		m.Spec.EnableUltraSSD); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "EnableNodePublicIP"),
		old.Spec.EnableNodePublicIP,
		m.Spec.EnableNodePublicIP); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "NodePublicIPPrefixID"),
		old.Spec.NodePublicIPPrefixID,
		m.Spec.NodePublicIPPrefixID); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "KubeletConfig"),
		old.Spec.KubeletConfig,
		m.Spec.KubeletConfig); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "KubeletDiskType"),
		old.Spec.KubeletDiskType,
		m.Spec.KubeletDiskType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "LinuxOSConfig"),
		old.Spec.LinuxOSConfig,
		m.Spec.LinuxOSConfig); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AzureManagedMachinePool").GroupKind(), m.Name, allErrs)
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedMachinePoolWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	m, ok := obj.(*AzureManagedMachinePool)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedMachinePool")
	}
	if m.Spec.Mode != string(NodePoolModeSystem) {
		return nil, nil
	}

	return nil, errors.Wrapf(m.validateLastSystemNodePool(mw.Client), "if the delete is triggered via owner MachinePool please refer to trouble shooting section in https://capz.sigs.k8s.io/topics/managedcluster.html")
}

// validateLastSystemNodePool is used to check if the existing system node pool is the last system node pool.
// If it is a last system node pool it cannot be deleted or mutated to user node pool as AKS expects min 1 system node pool.
func (m *AzureManagedMachinePool) validateLastSystemNodePool(cli client.Client) error {
	ctx := context.Background()

	// Fetch the Cluster.
	clusterName, ok := m.Labels[clusterv1.ClusterNameLabel]
	if !ok {
		return nil
	}

	ownerCluster := &clusterv1.Cluster{}
	key := client.ObjectKey{
		Namespace: m.Namespace,
		Name:      clusterName,
	}

	if err := cli.Get(ctx, key, ownerCluster); err != nil {
		return err
	}

	if !ownerCluster.DeletionTimestamp.IsZero() {
		return nil
	}

	// checking if the Cluster is going to be deleted for clusterctl move operation
	if _, found := ownerCluster.Annotations[clusterctlv1.DeleteForMoveAnnotation]; found {
		return nil
	}

	opt1 := client.InNamespace(m.Namespace)
	opt2 := client.MatchingLabels(map[string]string{
		clusterv1.ClusterNameLabel: clusterName,
		LabelAgentPoolMode:         string(NodePoolModeSystem),
	})

	ammpList := &AzureManagedMachinePoolList{}
	if err := cli.List(ctx, ammpList, opt1, opt2); err != nil {
		return err
	}

	if len(ammpList.Items) <= 1 {
		return errors.New("AKS Cluster must have at least one system pool")
	}
	return nil
}

func (m *AzureManagedMachinePool) validateMaxPods() error {
	if m.Spec.MaxPods != nil {
		if ptr.Deref[int32](m.Spec.MaxPods, 0) < 10 || ptr.Deref[int32](m.Spec.MaxPods, 0) > 250 {
			return field.Invalid(
				field.NewPath("Spec", "MaxPods"),
				m.Spec.MaxPods,
				"MaxPods must be between 10 and 250")
		}
	}

	return nil
}

func (m *AzureManagedMachinePool) validateOSType() error {
	if m.Spec.Mode == string(NodePoolModeSystem) {
		if m.Spec.OSType != nil && *m.Spec.OSType != LinuxOS {
			return field.Forbidden(
				field.NewPath("Spec", "OSType"),
				"System node pooll must have OSType 'Linux'")
		}
	}

	return nil
}

func (m *AzureManagedMachinePool) validateName() error {
	var name *string
	var fieldNameMessage string
	if m.Spec.Name == nil || *m.Spec.Name == "" {
		name = &m.Name
		fieldNameMessage = "when spec.name is empty, metadata.name"
	} else {
		name = m.Spec.Name
		fieldNameMessage = "spec.name"
	}

	if err := validateNameLength(m.Spec.OSType, name, fieldNameMessage); err != nil {
		return err
	}
	return validateNamePattern(name, fieldNameMessage)
}

func validateNameLength(osType *string, name *string, fieldNameMessage string) error {
	if osType != nil && *osType == WindowsOS &&
		name != nil && len(*name) > 6 {
		return field.Invalid(
			field.NewPath("Spec", "Name"),
			name,
			fmt.Sprintf("For OSType Windows, %s can not be longer than 6 characters.", fieldNameMessage))
	} else if (osType == nil || *osType == LinuxOS) &&
		(name != nil && len(*name) > 12) {
		return field.Invalid(
			field.NewPath("Spec", "Name"),
			osType,
			fmt.Sprintf("For OSType Linux, %s can not be longer than 12 characters.", fieldNameMessage))
	}
	return nil
}

func validateNamePattern(name *string, fieldNameMessage string) error {
	if name == nil || *name == "" {
		return nil
	}

	if !unicode.IsLower(rune((*name)[0])) {
		return field.Invalid(
			field.NewPath("Spec", "Name"),
			name,
			fmt.Sprintf("%s must begin with a lowercase letter.", fieldNameMessage))
	}

	for _, char := range *name {
		if !(unicode.IsLower(char) || unicode.IsNumber(char)) {
			return field.Invalid(
				field.NewPath("Spec", "Name"),
				name,
				fmt.Sprintf("%s may only contain lowercase alphanumeric characters.", fieldNameMessage))
		}
	}
	return nil
}

func (m *AzureManagedMachinePool) validateNodeLabels() error {
	for key := range m.Spec.NodeLabels {
		if azureutil.IsAzureSystemNodeLabelKey(key) {
			return field.Invalid(
				field.NewPath("Spec", "NodeLabels"),
				key,
				fmt.Sprintf("Node pool label key must not start with %s", azureutil.AzureSystemNodeLabelPrefix))
		}
	}

	return nil
}

func (m *AzureManagedMachinePool) validateNodePublicIPPrefixID() error {
	if m.Spec.NodePublicIPPrefixID != nil && !validNodePublicPrefixID.MatchString(*m.Spec.NodePublicIPPrefixID) {
		return field.Invalid(
			field.NewPath("Spec", "NodePublicIPPrefixID"),
			m.Spec.NodePublicIPPrefixID,
			fmt.Sprintf("resource ID must match %q", validNodePublicPrefixID.String()))
	}
	return nil
}

func (m *AzureManagedMachinePool) validateEnableNodePublicIP() error {
	if (m.Spec.EnableNodePublicIP == nil || !*m.Spec.EnableNodePublicIP) &&
		m.Spec.NodePublicIPPrefixID != nil {
		return field.Invalid(
			field.NewPath("Spec", "EnableNodePublicIP"),
			m.Spec.EnableNodePublicIP,
			"must be set to true when NodePublicIPPrefixID is set")
	}
	return nil
}

func (m *AzureManagedMachinePool) validateSubnetName() error {
	if m.Spec.SubnetName != nil {
		subnetRegex := "^[a-zA-Z0-9][a-zA-Z0-9-]{0,78}[a-zA-Z0-9]$"
		regex := regexp.MustCompile(subnetRegex)
		if success := regex.MatchString(ptr.Deref(m.Spec.SubnetName, "")); !success {
			return field.Invalid(field.NewPath("Spec", "SubnetName"), m.Spec.SubnetName,
				fmt.Sprintf("name of subnet doesn't match regex %s", subnetRegex))
		}
	}
	return nil
}

// validateKubeletConfig enforces the AKS API configuration for KubeletConfig.
// See:  https://learn.microsoft.com/en-us/azure/aks/custom-node-configuration.
func (m *AzureManagedMachinePool) validateKubeletConfig() error {
	var allowedUnsafeSysctlsPatterns = []string{
		`^kernel\.shm.+$`,
		`^kernel\.msg.+$`,
		`^kernel\.sem$`,
		`^fs\.mqueue\..+$`,
		`^net\..+$`,
	}
	if m.Spec.KubeletConfig != nil {
		if m.Spec.KubeletConfig.CPUCfsQuotaPeriod != nil {
			if !strings.HasSuffix(ptr.Deref(m.Spec.KubeletConfig.CPUCfsQuotaPeriod, ""), "ms") {
				return field.Invalid(
					field.NewPath("Spec", "KubeletConfig", "CPUCfsQuotaPeriod"),
					m.Spec.KubeletConfig.CPUCfsQuotaPeriod,
					"must be a string value in milliseconds with a 'ms' suffix, e.g., '100ms'")
			}
		}
		if m.Spec.KubeletConfig.ImageGcHighThreshold != nil && m.Spec.KubeletConfig.ImageGcLowThreshold != nil {
			if ptr.Deref[int32](m.Spec.KubeletConfig.ImageGcLowThreshold, 0) > ptr.Deref[int32](m.Spec.KubeletConfig.ImageGcHighThreshold, 0) {
				return field.Invalid(
					field.NewPath("Spec", "KubeletConfig", "ImageGcLowThreshold"),
					m.Spec.KubeletConfig.ImageGcLowThreshold,
					fmt.Sprintf("must not be greater than ImageGcHighThreshold, ImageGcLowThreshold=%d, ImageGcHighThreshold=%d",
						ptr.Deref[int32](m.Spec.KubeletConfig.ImageGcLowThreshold, 0), ptr.Deref[int32](m.Spec.KubeletConfig.ImageGcHighThreshold, 0)))
			}
		}
		for _, val := range m.Spec.KubeletConfig.AllowedUnsafeSysctls {
			var hasMatch bool
			for _, p := range allowedUnsafeSysctlsPatterns {
				if m, _ := regexp.MatchString(p, val); m {
					hasMatch = true
					break
				}
			}
			if !hasMatch {
				return field.Invalid(
					field.NewPath("Spec", "KubeletConfig", "AllowedUnsafeSysctls"),
					m.Spec.KubeletConfig.AllowedUnsafeSysctls,
					fmt.Sprintf("%s is not a supported AllowedUnsafeSysctls configuration", val))
			}
		}
	}
	return nil
}

// validateLinuxOSConfig enforces AKS API configuration for Linux OS custom configuration
// See: https://learn.microsoft.com/en-us/azure/aks/custom-node-configuration#linux-os-custom-configuration for detailed information.
func (m *AzureManagedMachinePool) validateLinuxOSConfig() error {
	var errs []error
	if m.Spec.LinuxOSConfig == nil {
		return nil
	}

	if m.Spec.LinuxOSConfig.SwapFileSizeMB != nil {
		if m.Spec.KubeletConfig == nil || ptr.Deref(m.Spec.KubeletConfig.FailSwapOn, true) {
			errs = append(errs, field.Invalid(
				field.NewPath("Spec", "LinuxOSConfig", "SwapFileSizeMB"),
				m.Spec.LinuxOSConfig.SwapFileSizeMB,
				"KubeletConfig.FailSwapOn must be set to false to enable swap file on nodes"))
		}
	}

	if m.Spec.LinuxOSConfig.Sysctls != nil && m.Spec.LinuxOSConfig.Sysctls.NetIpv4IPLocalPortRange != nil {
		// match numbers separated by a space
		portRangeRegex := `^[0-9]+ [0-9]+$`
		portRange := *m.Spec.LinuxOSConfig.Sysctls.NetIpv4IPLocalPortRange

		match, matchErr := regexp.MatchString(portRangeRegex, portRange)
		if matchErr != nil {
			errs = append(errs, matchErr)
		}
		if !match {
			errs = append(errs, field.Invalid(
				field.NewPath("Spec", "LinuxOSConfig", "Sysctls", "NetIpv4IpLocalPortRange"),
				m.Spec.LinuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
				"LinuxOSConfig.Sysctls.NetIpv4IpLocalPortRange must be of the format \"<int> <int>\""))
		} else {
			ports := strings.Split(portRange, " ")
			firstPort, _ := strconv.Atoi(ports[0])
			lastPort, _ := strconv.Atoi(ports[1])

			if firstPort < 1024 || firstPort > 60999 {
				errs = append(errs, field.Invalid(
					field.NewPath("Spec", "LinuxOSConfig", "Sysctls", "NetIpv4IpLocalPortRange", "First"),
					m.Spec.LinuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
					fmt.Sprintf("first port of NetIpv4IpLocalPortRange=%d must be in between [1024 - 60999]", firstPort)))
			}

			if lastPort < 32768 || lastPort > 65000 {
				errs = append(errs, field.Invalid(
					field.NewPath("Spec", "LinuxOSConfig", "Sysctls", "NetIpv4IpLocalPortRange", "Last"),
					m.Spec.LinuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
					fmt.Sprintf("last port of NetIpv4IpLocalPortRange=%d must be in between [32768 -65000]", lastPort)))
			}

			if firstPort > lastPort {
				errs = append(errs, field.Invalid(
					field.NewPath("Spec", "LinuxOSConfig", "Sysctls", "NetIpv4IpLocalPortRange", "First"),
					m.Spec.LinuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
					fmt.Sprintf("first port of NetIpv4IpLocalPortRange=%d cannot be greater than last port of NetIpv4IpLocalPortRange=%d", firstPort, lastPort)))
			}
		}
	}
	return kerrors.NewAggregate(errs)
}
