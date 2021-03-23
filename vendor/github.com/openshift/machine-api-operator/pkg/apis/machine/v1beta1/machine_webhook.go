package v1beta1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	osconfigv1 "github.com/openshift/api/config/v1"
	osclientset "github.com/openshift/client-go/config/clientset/versioned"
	gcp "github.com/openshift/cluster-api-provider-gcp/pkg/apis/gcpprovider/v1beta1"
	"github.com/openshift/machine-api-operator/pkg/apis/machine"
	vsphere "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider/v1beta1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	aws "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1beta1"
	azure "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	yaml "sigs.k8s.io/yaml"
)

var (
	// Azure Defaults
	defaultAzureVnet = func(clusterID string) string {
		return fmt.Sprintf("%s-vnet", clusterID)
	}
	defaultAzureSubnet = func(clusterID string) string {
		return fmt.Sprintf("%s-worker-subnet", clusterID)
	}
	defaultAzureNetworkResourceGroup = func(clusterID string) string {
		return fmt.Sprintf("%s-rg", clusterID)
	}
	defaultAzureImageResourceID = func(clusterID string) string {
		return fmt.Sprintf("/resourceGroups/%s/providers/Microsoft.Compute/images/%s", clusterID+"-rg", clusterID)
	}
	defaultAzureManagedIdentiy = func(clusterID string) string {
		return fmt.Sprintf("%s-identity", clusterID)
	}
	defaultAzureResourceGroup = func(clusterID string) string {
		return fmt.Sprintf("%s-rg", clusterID)
	}

	// GCP Defaults
	defaultGCPNetwork = func(clusterID string) string {
		return fmt.Sprintf("%s-network", clusterID)
	}
	defaultGCPSubnetwork = func(clusterID string) string {
		return fmt.Sprintf("%s-worker-subnet", clusterID)
	}
	defaultGCPTags = func(clusterID string) []string {
		return []string{fmt.Sprintf("%s-worker", clusterID)}
	}
)

const (
	DefaultMachineMutatingHookPath      = "/mutate-machine-openshift-io-v1beta1-machine"
	DefaultMachineValidatingHookPath    = "/validate-machine-openshift-io-v1beta1-machine"
	DefaultMachineSetMutatingHookPath   = "/mutate-machine-openshift-io-v1beta1-machineset"
	DefaultMachineSetValidatingHookPath = "/validate-machine-openshift-io-v1beta1-machineset"

	defaultWebhookConfigurationName = "machine-api"
	defaultWebhookServiceName       = "machine-api-operator-webhook"
	defaultWebhookServiceNamespace  = "openshift-machine-api"
	defaultWebhookServicePort       = 443

	defaultUserDataSecret  = "worker-user-data"
	defaultSecretNamespace = "openshift-machine-api"

	// AWS Defaults
	defaultAWSCredentialsSecret = "aws-cloud-credentials"
	defaultAWSInstanceType      = "m4.large"

	// Azure Defaults
	defaultAzureVMSize            = "Standard_D4s_V3"
	defaultAzureCredentialsSecret = "azure-cloud-credentials"
	defaultAzureOSDiskOSType      = "Linux"
	defaultAzureOSDiskStorageType = "Premium_LRS"
	azureMaxDiskSizeGB            = 32768

	// GCP Defaults
	defaultGCPMachineType       = "n1-standard-4"
	defaultGCPCredentialsSecret = "gcp-cloud-credentials"
	defaultGCPDiskSizeGb        = 128
	defaultGCPDiskType          = "pd-standard"
	// https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.6/46.82.202007212240-0/x86_64/meta.json
	// https://github.com/openshift/installer/pull/3808
	// https://github.com/openshift/installer/blob/d75bf7ad98124b901ae7e22b5595e0392ed6ea3c/data/data/rhcos.json
	defaultGCPDiskImage = "projects/rhcos-cloud/global/images/rhcos-46-82-202007212240-0-gcp-x86-64"

	// vSphere Defaults
	defaultVSphereCredentialsSecret = "vsphere-cloud-credentials"
	// Minimum vSphere values taken from vSphere reconciler
	minVSphereCPU       = 2
	minVSphereMemoryMiB = 2048
	// https://docs.openshift.com/container-platform/4.1/installing/installing_vsphere/installing-vsphere.html#minimum-resource-requirements_installing-vsphere
	minVSphereDiskGiB = 120
)

var (
	// webhookFailurePolicy is ignore so we don't want to block machine lifecycle on the webhook operational aspects.
	// This would be particularly problematic for chicken egg issues when bootstrapping a cluster.
	webhookFailurePolicy = admissionregistrationv1.Ignore
	webhookSideEffects   = admissionregistrationv1.SideEffectClassNone
)

func secretExists(c client.Client, name, namespace string) (bool, error) {
	key := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}
	obj := &corev1.Secret{}

	if err := c.Get(context.Background(), key, obj); err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func credentialsSecretExists(c client.Client, name, namespace string) []string {
	secretExists, err := secretExists(c, name, namespace)
	if err != nil {
		return []string{
			field.Invalid(
				field.NewPath("providerSpec", "credentialsSecret"),
				name,
				fmt.Sprintf("failed to get credentialsSecret: %v", err),
			).Error(),
		}
	}

	if !secretExists {
		return []string{
			field.Invalid(
				field.NewPath("providerSpec", "credentialsSecret"),
				name,
				"not found. Expected CredentialsSecret to exist",
			).Error(),
		}
	}

	return []string{}
}

func getInfra() (*osconfigv1.Infrastructure, error) {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := osclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	infra, err := client.ConfigV1().Infrastructures().Get(context.Background(), "cluster", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return infra, nil
}

func getDNS() (*osconfigv1.DNS, error) {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := osclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	dns, err := client.ConfigV1().DNSes().Get(context.Background(), "cluster", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return dns, nil
}

type machineAdmissionFn func(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate)

type admissionConfig struct {
	clusterID       string
	platformStatus  *osconfigv1.PlatformStatus
	dnsDisconnected bool
	client          client.Client
}

type admissionHandler struct {
	*admissionConfig
	webhookOperations machineAdmissionFn
	decoder           *admission.Decoder
}

// InjectDecoder injects the decoder.
func (a *admissionHandler) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}

// machineValidatorHandler validates Machine API resources.
// implements type Handler interface.
// https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/webhook/admission#Handler
type machineValidatorHandler struct {
	*admissionHandler
}

// machineDefaulterHandler defaults Machine API resources.
// implements type Handler interface.
// https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/webhook/admission#Handler
type machineDefaulterHandler struct {
	*admissionHandler
}

// NewValidator returns a new machineValidatorHandler.
func NewMachineValidator() (*machineValidatorHandler, error) {
	infra, err := getInfra()
	if err != nil {
		return nil, err
	}

	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}
	c, err := client.New(cfg, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to build kubernetes client: %v", err)
	}

	dns, err := getDNS()
	if err != nil {
		return nil, err
	}

	return createMachineValidator(infra, c, dns), nil
}

func createMachineValidator(infra *osconfigv1.Infrastructure, client client.Client, dns *osconfigv1.DNS) *machineValidatorHandler {
	admissionConfig := &admissionConfig{
		dnsDisconnected: dns.Spec.PublicZone == nil,
		clusterID:       infra.Status.InfrastructureName,
		platformStatus:  infra.Status.PlatformStatus,
		client:          client,
	}
	return &machineValidatorHandler{
		admissionHandler: &admissionHandler{
			admissionConfig:   admissionConfig,
			webhookOperations: getMachineValidatorOperation(infra.Status.PlatformStatus.Type),
		},
	}
}

func getMachineValidatorOperation(platform osconfigv1.PlatformType) machineAdmissionFn {
	switch platform {
	case osconfigv1.AWSPlatformType:
		return validateAWS
	case osconfigv1.AzurePlatformType:
		return validateAzure
	case osconfigv1.GCPPlatformType:
		return validateGCP
	case osconfigv1.VSpherePlatformType:
		return validateVSphere
	default:
		// just no-op
		return func(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
			return true, []string{}, nil
		}
	}
}

// NewDefaulter returns a new machineDefaulterHandler.
func NewMachineDefaulter() (*machineDefaulterHandler, error) {
	infra, err := getInfra()
	if err != nil {
		return nil, err
	}

	return createMachineDefaulter(infra.Status.PlatformStatus, infra.Status.InfrastructureName), nil
}

func createMachineDefaulter(platformStatus *osconfigv1.PlatformStatus, clusterID string) *machineDefaulterHandler {
	return &machineDefaulterHandler{
		admissionHandler: &admissionHandler{
			admissionConfig:   &admissionConfig{clusterID: clusterID},
			webhookOperations: getMachineDefaulterOperation(platformStatus),
		},
	}
}

func getMachineDefaulterOperation(platformStatus *osconfigv1.PlatformStatus) machineAdmissionFn {
	switch platformStatus.Type {
	case osconfigv1.AWSPlatformType:
		region := ""
		if platformStatus.AWS != nil {
			region = platformStatus.AWS.Region
		}
		return awsDefaulter{region: region}.defaultAWS
	case osconfigv1.AzurePlatformType:
		return defaultAzure
	case osconfigv1.GCPPlatformType:
		return defaultGCP
	case osconfigv1.VSpherePlatformType:
		return defaultVSphere
	default:
		// just no-op
		return func(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
			return true, []string{}, nil
		}
	}
}

// NewValidatingWebhookConfiguration creates a validation webhook configuration with configured Machine and MachineSet webhooks
func NewValidatingWebhookConfiguration() *admissionregistrationv1.ValidatingWebhookConfiguration {
	validatingWebhookConfiguration := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultWebhookConfigurationName,
			Annotations: map[string]string{
				"service.beta.openshift.io/inject-cabundle": "true",
			},
		},
		Webhooks: []admissionregistrationv1.ValidatingWebhook{
			MachineValidatingWebhook(),
			MachineSetValidatingWebhook(),
		},
	}

	// Setting group version is required for testEnv to create unstructured objects, as the new structure sets it on empty strings
	// Usual way to populate those values, is to create the resource in the cluster first, which we can't yet do.
	validatingWebhookConfiguration.SetGroupVersionKind(admissionregistrationv1.SchemeGroupVersion.WithKind("ValidatingWebhookConfiguration"))
	return validatingWebhookConfiguration
}

// MachineValidatingWebhook returns validating webhooks for machine to populate the configuration
func MachineValidatingWebhook() admissionregistrationv1.ValidatingWebhook {
	serviceReference := admissionregistrationv1.ServiceReference{
		Namespace: defaultWebhookServiceNamespace,
		Name:      defaultWebhookServiceName,
		Path:      pointer.StringPtr(DefaultMachineValidatingHookPath),
		Port:      pointer.Int32Ptr(defaultWebhookServicePort),
	}
	return admissionregistrationv1.ValidatingWebhook{
		AdmissionReviewVersions: []string{"v1beta1"},
		Name:                    "validation.machine.machine.openshift.io",
		FailurePolicy:           &webhookFailurePolicy,
		SideEffects:             &webhookSideEffects,
		ClientConfig: admissionregistrationv1.WebhookClientConfig{
			Service: &serviceReference,
		},
		Rules: []admissionregistrationv1.RuleWithOperations{
			{
				Rule: admissionregistrationv1.Rule{
					APIGroups:   []string{machine.GroupName},
					APIVersions: []string{SchemeGroupVersion.Version},
					Resources:   []string{"machines"},
				},
				Operations: []admissionregistrationv1.OperationType{
					admissionregistrationv1.Create,
					admissionregistrationv1.Update,
				},
			},
		},
	}
}

// MachineSetValidatingWebhook returns validating webhooks for machineSet to populate the configuration
func MachineSetValidatingWebhook() admissionregistrationv1.ValidatingWebhook {
	machinesetServiceReference := admissionregistrationv1.ServiceReference{
		Namespace: defaultWebhookServiceNamespace,
		Name:      defaultWebhookServiceName,
		Path:      pointer.StringPtr(DefaultMachineSetValidatingHookPath),
		Port:      pointer.Int32Ptr(defaultWebhookServicePort),
	}
	return admissionregistrationv1.ValidatingWebhook{
		AdmissionReviewVersions: []string{"v1beta1"},
		Name:                    "validation.machineset.machine.openshift.io",
		FailurePolicy:           &webhookFailurePolicy,
		SideEffects:             &webhookSideEffects,
		ClientConfig: admissionregistrationv1.WebhookClientConfig{
			Service: &machinesetServiceReference,
		},
		Rules: []admissionregistrationv1.RuleWithOperations{
			{
				Rule: admissionregistrationv1.Rule{
					APIGroups:   []string{machine.GroupName},
					APIVersions: []string{SchemeGroupVersion.Version},
					Resources:   []string{"machinesets"},
				},
				Operations: []admissionregistrationv1.OperationType{
					admissionregistrationv1.Create,
					admissionregistrationv1.Update,
				},
			},
		},
	}
}

// NewMutatingWebhookConfiguration creates a mutating webhook configuration with configured Machine and MachineSet webhooks
func NewMutatingWebhookConfiguration() *admissionregistrationv1.MutatingWebhookConfiguration {
	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultWebhookConfigurationName,
			Annotations: map[string]string{
				"service.beta.openshift.io/inject-cabundle": "true",
			},
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			MachineMutatingWebhook(),
			MachineSetMutatingWebhook(),
		},
	}

	// Setting group version is required for testEnv to create unstructured objects, as the new structure sets it on empty strings
	// Usual way to populate those values, is to create the resource in the cluster first, which we can't yet do.
	mutatingWebhookConfiguration.SetGroupVersionKind(admissionregistrationv1.SchemeGroupVersion.WithKind("MutatingWebhookConfiguration"))
	return mutatingWebhookConfiguration
}

// MachineMutatingWebhook returns mutating webhooks for machine to apply in configuration
func MachineMutatingWebhook() admissionregistrationv1.MutatingWebhook {
	machineServiceReference := admissionregistrationv1.ServiceReference{
		Namespace: defaultWebhookServiceNamespace,
		Name:      defaultWebhookServiceName,
		Path:      pointer.StringPtr(DefaultMachineMutatingHookPath),
		Port:      pointer.Int32Ptr(defaultWebhookServicePort),
	}
	return admissionregistrationv1.MutatingWebhook{
		AdmissionReviewVersions: []string{"v1beta1"},
		Name:                    "default.machine.machine.openshift.io",
		FailurePolicy:           &webhookFailurePolicy,
		SideEffects:             &webhookSideEffects,
		ClientConfig: admissionregistrationv1.WebhookClientConfig{
			Service: &machineServiceReference,
		},
		Rules: []admissionregistrationv1.RuleWithOperations{
			{
				Rule: admissionregistrationv1.Rule{
					APIGroups:   []string{machine.GroupName},
					APIVersions: []string{SchemeGroupVersion.Version},
					Resources:   []string{"machines"},
				},
				Operations: []admissionregistrationv1.OperationType{
					admissionregistrationv1.Create,
				},
			},
		},
	}
}

// MachineSetMutatingWebhook returns mutating webhook for machineSet to apply in configuration
func MachineSetMutatingWebhook() admissionregistrationv1.MutatingWebhook {
	machineSetServiceReference := admissionregistrationv1.ServiceReference{
		Namespace: defaultWebhookServiceNamespace,
		Name:      defaultWebhookServiceName,
		Path:      pointer.StringPtr(DefaultMachineSetMutatingHookPath),
		Port:      pointer.Int32Ptr(defaultWebhookServicePort),
	}
	return admissionregistrationv1.MutatingWebhook{
		AdmissionReviewVersions: []string{"v1beta1"},
		Name:                    "default.machineset.machine.openshift.io",
		FailurePolicy:           &webhookFailurePolicy,
		SideEffects:             &webhookSideEffects,
		ClientConfig: admissionregistrationv1.WebhookClientConfig{
			Service: &machineSetServiceReference,
		},
		Rules: []admissionregistrationv1.RuleWithOperations{
			{
				Rule: admissionregistrationv1.Rule{
					APIGroups:   []string{machine.GroupName},
					APIVersions: []string{SchemeGroupVersion.Version},
					Resources:   []string{"machinesets"},
				},
				Operations: []admissionregistrationv1.OperationType{
					admissionregistrationv1.Create,
				},
			},
		},
	}
}

// Handle handles HTTP requests for admission webhook servers.
func (h *machineValidatorHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	m := &Machine{}

	if err := h.decoder.Decode(req, m); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	klog.V(3).Infof("Validate webhook called for Machine: %s", m.GetName())

	ok, warnings, errs := h.webhookOperations(m, h.admissionConfig)
	if !ok {
		return admission.Denied(errs.Error()).WithWarnings(warnings...)
	}

	return admission.Allowed("Machine valid").WithWarnings(warnings...)
}

// Handle handles HTTP requests for admission webhook servers.
func (h *machineDefaulterHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	m := &Machine{}

	if err := h.decoder.Decode(req, m); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	klog.V(3).Infof("Mutate webhook called for Machine: %s", m.GetName())

	// Only enforce the clusterID if it's not set.
	// Otherwise a discrepancy on the value would leave the machine orphan
	// and would trigger a new machine creation by the machineSet.
	// https://bugzilla.redhat.com/show_bug.cgi?id=1857175
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}
	if _, ok := m.Labels[MachineClusterIDLabel]; !ok {
		m.Labels[MachineClusterIDLabel] = h.clusterID
	}

	ok, warnings, errs := h.webhookOperations(m, h.admissionConfig)
	if !ok {
		return admission.Denied(errs.Error()).WithWarnings(warnings...)
	}

	marshaledMachine, err := json.Marshal(m)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err).WithWarnings(warnings...)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledMachine).WithWarnings(warnings...)
}

type awsDefaulter struct {
	region string
}

func (a awsDefaulter) defaultAWS(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
	klog.V(3).Infof("Defaulting AWS providerSpec")

	var errs []error
	var warnings []string
	providerSpec := new(aws.AWSMachineProviderConfig)
	if err := unmarshalInto(m, providerSpec); err != nil {
		errs = append(errs, err)
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	if providerSpec.InstanceType == "" {
		providerSpec.InstanceType = defaultAWSInstanceType
	}

	if providerSpec.Placement.Region == "" {
		providerSpec.Placement.Region = a.region
	}

	if providerSpec.UserDataSecret == nil {
		providerSpec.UserDataSecret = &corev1.LocalObjectReference{Name: defaultUserDataSecret}
	}

	if providerSpec.CredentialsSecret == nil {
		providerSpec.CredentialsSecret = &corev1.LocalObjectReference{Name: defaultAWSCredentialsSecret}
	}

	rawBytes, err := json.Marshal(providerSpec)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	m.Spec.ProviderSpec.Value = &runtime.RawExtension{Raw: rawBytes}
	return true, warnings, nil
}

func unmarshalInto(m *Machine, providerSpec interface{}) error {
	if m.Spec.ProviderSpec.Value == nil {
		return field.Required(field.NewPath("providerSpec", "value"), "a value must be provided")
	}

	if err := yaml.Unmarshal(m.Spec.ProviderSpec.Value.Raw, &providerSpec); err != nil {
		return field.Invalid(field.NewPath("providerSpec", "value"), providerSpec, err.Error())
	}
	return nil
}

func validateAWS(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
	klog.V(3).Infof("Validating AWS providerSpec")

	var errs []error
	var warnings []string
	providerSpec := new(aws.AWSMachineProviderConfig)
	if err := unmarshalInto(m, providerSpec); err != nil {
		errs = append(errs, err)
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	if providerSpec.AMI.ARN == nil && providerSpec.AMI.Filters == nil && providerSpec.AMI.ID == nil {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "ami"),
				"expected either providerSpec.ami.arn or providerSpec.ami.filters or providerSpec.ami.id to be populated",
			),
		)
	}

	if providerSpec.Placement.Region == "" {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "placement", "region"),
				"expected providerSpec.placement.region to be populated",
			),
		)
	}

	if providerSpec.InstanceType == "" {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "instanceType"),
				"expected providerSpec.instanceType to be populated",
			),
		)
	}

	if providerSpec.UserDataSecret == nil {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "userDataSecret"),
				"expected providerSpec.userDataSecret to be populated",
			),
		)
	}

	if providerSpec.CredentialsSecret == nil {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "credentialsSecret"),
				"expected providerSpec.credentialsSecret to be populated",
			),
		)
	} else {
		warnings = append(warnings, credentialsSecretExists(config.client, providerSpec.CredentialsSecret.Name, m.GetNamespace())...)
	}

	if providerSpec.Subnet.ARN == nil && providerSpec.Subnet.ID == nil && providerSpec.Subnet.Filters == nil {
		warnings = append(
			warnings,
			"providerSpec.subnet: No subnet has been provided. Instances may be created in an unexpected subnet and may not join the cluster.",
		)
	}
	// TODO(alberto): Validate providerSpec.BlockDevices.
	// https://github.com/openshift/cluster-api-provider-aws/pull/299#discussion_r433920532

	switch providerSpec.Placement.Tenancy {
	case "", aws.DefaultTenancy, aws.DedicatedTenancy, aws.HostTenancy:
		// Do nothing, valid values
	default:
		errs = append(
			errs,
			field.Invalid(
				field.NewPath("providerSpec", "tenancy"),
				providerSpec.Placement.Tenancy,
				fmt.Sprintf("Invalid providerSpec.tenancy, the only allowed options are: %s, %s, %s", aws.DefaultTenancy, aws.DedicatedTenancy, aws.HostTenancy),
			),
		)
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	return true, warnings, nil
}

func defaultAzure(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
	klog.V(3).Infof("Defaulting Azure providerSpec")

	var errs []error
	var warnings []string
	providerSpec := new(azure.AzureMachineProviderSpec)
	if err := unmarshalInto(m, providerSpec); err != nil {
		errs = append(errs, err)
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	if providerSpec.VMSize == "" {
		providerSpec.VMSize = defaultAzureVMSize
	}

	// Vnet and Subnet need to be provided together by the user
	if providerSpec.Vnet == "" && providerSpec.Subnet == "" {
		providerSpec.Vnet = defaultAzureVnet(config.clusterID)
		providerSpec.Subnet = defaultAzureSubnet(config.clusterID)
	}

	if providerSpec.Image == (azure.Image{}) {
		providerSpec.Image.ResourceID = defaultAzureImageResourceID(config.clusterID)
	}

	if providerSpec.UserDataSecret == nil {
		providerSpec.UserDataSecret = &corev1.SecretReference{Name: defaultUserDataSecret}
	} else if providerSpec.UserDataSecret.Name == "" {
		providerSpec.UserDataSecret.Name = defaultUserDataSecret
	}

	if providerSpec.CredentialsSecret == nil {
		providerSpec.CredentialsSecret = &corev1.SecretReference{Name: defaultAzureCredentialsSecret, Namespace: defaultSecretNamespace}
	} else {
		if providerSpec.CredentialsSecret.Namespace == "" {
			providerSpec.CredentialsSecret.Namespace = defaultSecretNamespace
		}
		if providerSpec.CredentialsSecret.Name == "" {
			providerSpec.CredentialsSecret.Name = defaultAzureCredentialsSecret
		}
	}

	rawBytes, err := json.Marshal(providerSpec)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	m.Spec.ProviderSpec.Value = &runtime.RawExtension{Raw: rawBytes}
	return true, warnings, nil
}

func validateAzure(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
	klog.V(3).Infof("Validating Azure providerSpec")

	var errs []error
	var warnings []string
	providerSpec := new(azure.AzureMachineProviderSpec)
	if err := unmarshalInto(m, providerSpec); err != nil {
		errs = append(errs, err)
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	if providerSpec.VMSize == "" {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "vmSize"), "vmSize should be set to one of the supported Azure VM sizes"))
	}

	if providerSpec.PublicIP && config.dnsDisconnected {
		errs = append(errs, field.Forbidden(field.NewPath("providerSpec", "publicIP"), "publicIP is not allowed in Azure disconnected installation"))
	}
	// Vnet requires Subnet
	if providerSpec.Vnet != "" && providerSpec.Subnet == "" {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "subnet"), "must provide a subnet when a virtual network is specified"))
	}

	// Subnet requires Vnet
	if providerSpec.Subnet != "" && providerSpec.Vnet == "" {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "vnet"), "must provide a virtual network when supplying subnets"))
	}

	errs = append(errs, validateAzureImage(providerSpec.Image)...)

	if providerSpec.UserDataSecret == nil {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "userDataSecret"), "userDataSecret must be provided"))
	} else if providerSpec.UserDataSecret.Name == "" {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "userDataSecret", "name"), "name must be provided"))
	}

	if providerSpec.CredentialsSecret == nil {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "credentialsSecret"), "credentialsSecret must be provided"))
	} else {
		if providerSpec.CredentialsSecret.Namespace == "" {
			errs = append(errs, field.Required(field.NewPath("providerSpec", "credentialsSecret", "namespace"), "namespace must be provided"))
		}
		if providerSpec.CredentialsSecret.Name == "" {
			errs = append(errs, field.Required(field.NewPath("providerSpec", "credentialsSecret", "name"), "name must be provided"))
		}
		if providerSpec.CredentialsSecret.Name != "" && providerSpec.CredentialsSecret.Namespace != "" {
			warnings = append(warnings, credentialsSecretExists(config.client, providerSpec.CredentialsSecret.Name, providerSpec.CredentialsSecret.Namespace)...)
		}
	}

	if providerSpec.OSDisk.DiskSizeGB <= 0 || providerSpec.OSDisk.DiskSizeGB >= azureMaxDiskSizeGB {
		errs = append(errs, field.Invalid(field.NewPath("providerSpec", "osDisk", "diskSizeGB"), providerSpec.OSDisk.DiskSizeGB, "diskSizeGB must be greater than zero and less than 32768"))
	}

	if isAzureGovCloud(config.platformStatus) && providerSpec.SpotVMOptions != nil {
		warnings = append(warnings, "spot VMs may not be supported when using GovCloud region")
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}
	return true, warnings, nil
}

func validateAzureImage(image azure.Image) []error {
	errors := []error{}
	if image == (azure.Image{}) {
		return append(errors, field.Required(field.NewPath("providerSpec", "image"), "an image reference must be provided"))
	}

	if image.ResourceID != "" {
		if image != (azure.Image{ResourceID: image.ResourceID}) {
			return append(errors, field.Required(field.NewPath("providerSpec", "image", "resourceID"), "resourceID is already specified, other fields such as [Offer, Publisher, SKU, Version] should not be set"))
		}
		return errors
	}

	// Resource ID not provided, so Offer, Publisher, SKU and Version are required
	if image.Offer == "" {
		errors = append(errors, field.Required(field.NewPath("providerSpec", "image", "Offer"), "Offer must be provided"))
	}
	if image.Publisher == "" {
		errors = append(errors, field.Required(field.NewPath("providerSpec", "image", "Publisher"), "Publisher must be provided"))
	}
	if image.SKU == "" {
		errors = append(errors, field.Required(field.NewPath("providerSpec", "image", "SKU"), "SKU must be provided"))
	}
	if image.Version == "" {
		errors = append(errors, field.Required(field.NewPath("providerSpec", "image", "Version"), "Version must be provided"))
	}

	return errors
}

func defaultGCP(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
	klog.V(3).Infof("Defaulting GCP providerSpec")

	var errs []error
	var warnings []string
	providerSpec := new(gcp.GCPMachineProviderSpec)
	if err := unmarshalInto(m, providerSpec); err != nil {
		errs = append(errs, err)
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	if providerSpec.MachineType == "" {
		providerSpec.MachineType = defaultGCPMachineType
	}

	if len(providerSpec.NetworkInterfaces) == 0 {
		providerSpec.NetworkInterfaces = append(providerSpec.NetworkInterfaces, &gcp.GCPNetworkInterface{
			Network:    defaultGCPNetwork(config.clusterID),
			Subnetwork: defaultGCPSubnetwork(config.clusterID),
		})
	}

	providerSpec.Disks = defaultGCPDisks(providerSpec.Disks, config.clusterID)

	if len(providerSpec.Tags) == 0 {
		providerSpec.Tags = defaultGCPTags(config.clusterID)
	}

	if providerSpec.UserDataSecret == nil {
		providerSpec.UserDataSecret = &corev1.LocalObjectReference{Name: defaultUserDataSecret}
	}

	if providerSpec.CredentialsSecret == nil {
		providerSpec.CredentialsSecret = &corev1.LocalObjectReference{Name: defaultGCPCredentialsSecret}
	}

	rawBytes, err := json.Marshal(providerSpec)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	m.Spec.ProviderSpec.Value = &runtime.RawExtension{Raw: rawBytes}
	return true, warnings, nil
}

func defaultGCPDisks(disks []*gcp.GCPDisk, clusterID string) []*gcp.GCPDisk {
	if len(disks) == 0 {
		return []*gcp.GCPDisk{
			{
				AutoDelete: true,
				Boot:       true,
				SizeGb:     defaultGCPDiskSizeGb,
				Type:       defaultGCPDiskType,
				Image:      defaultGCPDiskImage,
			},
		}
	}

	for _, disk := range disks {
		if disk.Type == "" {
			disk.Type = defaultGCPDiskType
		}

		if disk.Image == "" {
			disk.Image = defaultGCPDiskImage
		}
	}

	return disks
}

func validateGCP(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
	klog.V(3).Infof("Validating GCP providerSpec")

	var errs []error
	var warnings []string
	providerSpec := new(gcp.GCPMachineProviderSpec)
	if err := unmarshalInto(m, providerSpec); err != nil {
		errs = append(errs, err)
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	if providerSpec.Region == "" {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "region"), "region is required"))
	}

	if !strings.HasPrefix(providerSpec.Zone, providerSpec.Region) {
		errs = append(errs, field.Invalid(field.NewPath("providerSpec", "zone"), providerSpec.Zone, fmt.Sprintf("zone not in configured region (%s)", providerSpec.Region)))
	}

	if providerSpec.MachineType == "" {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "machineType"), "machineType should be set to one of the supported GCP machine types"))
	}

	errs = append(errs, validateGCPNetworkInterfaces(providerSpec.NetworkInterfaces, field.NewPath("providerSpec", "networkInterfaces"))...)
	errs = append(errs, validateGCPDisks(providerSpec.Disks, field.NewPath("providerSpec", "disks"))...)

	if len(providerSpec.ServiceAccounts) == 0 {
		warnings = append(warnings, "providerSpec.serviceAccounts: no service account provided: nodes may be unable to join the cluster")
	} else {
		errs = append(errs, validateGCPServiceAccounts(providerSpec.ServiceAccounts, field.NewPath("providerSpec", "serviceAccounts"))...)
	}

	if providerSpec.UserDataSecret == nil {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "userDataSecret"), "userDataSecret must be provided"))
	} else {
		if providerSpec.UserDataSecret.Name == "" {
			errs = append(errs, field.Required(field.NewPath("providerSpec", "userDataSecret", "name"), "name must be provided"))
		}
	}

	if providerSpec.CredentialsSecret == nil {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "credentialsSecret"), "credentialsSecret must be provided"))
	} else {
		if providerSpec.CredentialsSecret.Name == "" {
			errs = append(errs, field.Required(field.NewPath("providerSpec", "credentialsSecret", "name"), "name must be provided"))
		} else {
			warnings = append(warnings, credentialsSecretExists(config.client, providerSpec.CredentialsSecret.Name, m.GetNamespace())...)
		}
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}
	return true, warnings, nil
}

func validateGCPNetworkInterfaces(networkInterfaces []*gcp.GCPNetworkInterface, parentPath *field.Path) []error {
	if len(networkInterfaces) == 0 {
		return []error{field.Required(parentPath, "at least 1 network interface is required")}
	}

	var errs []error
	for i, ni := range networkInterfaces {
		fldPath := parentPath.Index(i)

		if ni.Network == "" {
			errs = append(errs, field.Required(fldPath.Child("network"), "network is required"))
		}

		if ni.Subnetwork == "" {
			errs = append(errs, field.Required(fldPath.Child("subnetwork"), "subnetwork is required"))
		}
	}

	return errs
}

func validateGCPDisks(disks []*gcp.GCPDisk, parentPath *field.Path) []error {
	if len(disks) == 0 {
		return []error{field.Required(parentPath, "at least 1 disk is required")}
	}

	var errs []error
	for i, disk := range disks {
		fldPath := parentPath.Index(i)

		if disk.SizeGb != 0 {
			if disk.SizeGb < 16 {
				errs = append(errs, field.Invalid(fldPath.Child("sizeGb"), disk.SizeGb, "must be at least 16GB in size"))
			} else if disk.SizeGb > 65536 {
				errs = append(errs, field.Invalid(fldPath.Child("sizeGb"), disk.SizeGb, "exceeding maximum GCP disk size limit, must be below 65536"))
			}
		}

		if disk.Type != "" {
			diskTypes := sets.NewString("pd-standard", "pd-ssd")
			if !diskTypes.Has(disk.Type) {
				errs = append(errs, field.NotSupported(fldPath.Child("type"), disk.Type, diskTypes.List()))
			}
		}
	}

	return errs
}

func validateGCPServiceAccounts(serviceAccounts []gcp.GCPServiceAccount, parentPath *field.Path) []error {
	if len(serviceAccounts) != 1 {
		return []error{field.Invalid(parentPath, fmt.Sprintf("%d service accounts supplied", len(serviceAccounts)), "exactly 1 service account must be supplied")}
	}

	var errs []error
	for i, serviceAccount := range serviceAccounts {
		fldPath := parentPath.Index(i)

		if serviceAccount.Email == "" {
			errs = append(errs, field.Required(fldPath.Child("email"), "email is required"))
		}

		if len(serviceAccount.Scopes) == 0 {
			errs = append(errs, field.Required(fldPath.Child("scopes"), "at least 1 scope is required"))
		}
	}
	return errs
}

func defaultVSphere(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
	klog.V(3).Infof("Defaulting vSphere providerSpec")

	var errs []error
	var warnings []string
	providerSpec := new(vsphere.VSphereMachineProviderSpec)
	if err := unmarshalInto(m, providerSpec); err != nil {
		errs = append(errs, err)
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	if providerSpec.UserDataSecret == nil {
		providerSpec.UserDataSecret = &corev1.LocalObjectReference{Name: defaultUserDataSecret}
	}

	if providerSpec.CredentialsSecret == nil {
		providerSpec.CredentialsSecret = &corev1.LocalObjectReference{Name: defaultVSphereCredentialsSecret}
	}

	rawBytes, err := json.Marshal(providerSpec)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	m.Spec.ProviderSpec.Value = &runtime.RawExtension{Raw: rawBytes}
	return true, warnings, nil
}

func validateVSphere(m *Machine, config *admissionConfig) (bool, []string, utilerrors.Aggregate) {
	klog.V(3).Infof("Validating vSphere providerSpec")

	var errs []error
	var warnings []string
	providerSpec := new(vsphere.VSphereMachineProviderSpec)
	if err := unmarshalInto(m, providerSpec); err != nil {
		errs = append(errs, err)
		return false, warnings, utilerrors.NewAggregate(errs)
	}

	if providerSpec.Template == "" {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "template"), "template must be provided"))
	}

	workspaceWarnings, workspaceErrors := validateVSphereWorkspace(providerSpec.Workspace, field.NewPath("providerSpec", "workspace"))
	warnings = append(warnings, workspaceWarnings...)
	errs = append(errs, workspaceErrors...)

	errs = append(errs, validateVSphereNetwork(providerSpec.Network, field.NewPath("providerSpec", "network"))...)

	if providerSpec.NumCPUs < minVSphereCPU {
		warnings = append(warnings, fmt.Sprintf("providerSpec.numCPUs: %d is missing or less than the minimum value (%d): nodes may not boot correctly", providerSpec.NumCPUs, minVSphereCPU))
	}
	if providerSpec.MemoryMiB < minVSphereMemoryMiB {
		warnings = append(warnings, fmt.Sprintf("providerSpec.memoryMiB: %d is missing or less than the recommended minimum value (%d): nodes may not boot correctly", providerSpec.MemoryMiB, minVSphereMemoryMiB))
	}
	if providerSpec.DiskGiB < minVSphereDiskGiB {
		warnings = append(warnings, fmt.Sprintf("providerSpec.diskGiB: %d is missing or less than the recommended minimum (%d): nodes may fail to start if disk size is too low", providerSpec.DiskGiB, minVSphereDiskGiB))
	}

	if providerSpec.UserDataSecret == nil {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "userDataSecret"), "userDataSecret must be provided"))
	} else {
		if providerSpec.UserDataSecret.Name == "" {
			errs = append(errs, field.Required(field.NewPath("providerSpec", "userDataSecret", "name"), "name must be provided"))
		}
	}

	if providerSpec.CredentialsSecret == nil {
		errs = append(errs, field.Required(field.NewPath("providerSpec", "credentialsSecret"), "credentialsSecret must be provided"))
	} else {
		if providerSpec.CredentialsSecret.Name == "" {
			errs = append(errs, field.Required(field.NewPath("providerSpec", "credentialsSecret", "name"), "name must be provided"))
		} else {
			warnings = append(warnings, credentialsSecretExists(config.client, providerSpec.CredentialsSecret.Name, m.GetNamespace())...)
		}
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}
	return true, warnings, nil
}

func validateVSphereWorkspace(workspace *vsphere.Workspace, parentPath *field.Path) ([]string, []error) {
	if workspace == nil {
		return []string{}, []error{field.Required(parentPath, "workspace must be provided")}
	}

	var errs []error
	var warnings []string
	if workspace.Server == "" {
		errs = append(errs, field.Required(parentPath.Child("server"), "server must be provided"))
	}
	if workspace.Datacenter == "" {
		warnings = append(warnings, fmt.Sprintf("%s: datacenter is unset: if more than one datacenter is present, VMs cannot be created", parentPath.Child("datacenter")))
	}
	if workspace.Folder != "" {
		expectedPrefix := fmt.Sprintf("/%s/vm/", workspace.Datacenter)
		if !strings.HasPrefix(workspace.Folder, expectedPrefix) {
			errMsg := fmt.Sprintf("folder must be absolute path: expected prefix %q", expectedPrefix)
			errs = append(errs, field.Invalid(parentPath.Child("folder"), workspace.Folder, errMsg))
		}
	}

	return warnings, errs
}

func validateVSphereNetwork(network vsphere.NetworkSpec, parentPath *field.Path) []error {
	if len(network.Devices) == 0 {
		return []error{field.Required(parentPath.Child("devices"), "at least 1 network device must be provided")}
	}

	var errs []error
	for i, spec := range network.Devices {
		fldPath := parentPath.Child("devices").Index(i)
		if spec.NetworkName == "" {
			errs = append(errs, field.Required(fldPath.Child("networkName"), "networkName must be provided"))
		}
	}

	return errs
}

func isAzureGovCloud(platformStatus *osconfigv1.PlatformStatus) bool {
	return platformStatus != nil && platformStatus.Azure != nil &&
		platformStatus.Azure.CloudName != osconfigv1.AzurePublicCloud
}
