/*
Copyright 2020 The Kubernetes Authors.

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

package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	capifeature "sigs.k8s.io/cluster-api/feature"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/groups"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/coalescing"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const (
	spIdentityWarning = "You are using Service Principal authentication for Cloud Provider Azure which is less secure than Managed Identity. " +
		"Your Service Principal credentials will be written to a file on the disk of each VM in order to be accessible by Cloud Provider. " +
		"To learn more, see https://capz.sigs.k8s.io/topics/identities-use-cases.html#azure-host-identity "
	deprecatedManagerCredsWarning = "You're using deprecated functionality: " +
		"Using Azure credentials from the manager environment is deprecated and will be removed in future releases. " +
		"Please specify an AzureClusterIdentity for the AzureCluster instead, see: https://capz.sigs.k8s.io/topics/multitenancy.html "
)

type (
	// Options are controller options extended.
	Options struct {
		controller.Options
		Cache *coalescing.ReconcileCache
	}

	// ClusterScoper is a interface used by AzureMachinePools that can be owned by either an AzureManagedCluster or AzureCluster.
	ClusterScoper interface {
		azure.ClusterScoper
		groups.GroupScope
	}
)

// AzureClusterToAzureMachinesMapper creates a mapping handler to transform AzureClusters into AzureMachines. The transform
// requires AzureCluster to map to the owning Cluster, then from the Cluster, collect the Machines belonging to the cluster,
// then finally projecting the infrastructure reference to the AzureMachine.
func AzureClusterToAzureMachinesMapper(_ context.Context, c client.Client, obj runtime.Object, scheme *runtime.Scheme, log logr.Logger) (handler.MapFunc, error) {
	gvk, err := apiutil.GVKForObject(obj, scheme)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find GVK for AzureMachine")
	}

	return func(ctx context.Context, o client.Object) []ctrl.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		azCluster, ok := o.(*infrav1.AzureCluster)
		if !ok {
			log.Error(errors.Errorf("expected an AzureCluster, got %T instead", o), "failed to map AzureCluster")
			return nil
		}

		log := log.WithValues("AzureCluster", azCluster.Name, "Namespace", azCluster.Namespace)

		// Don't handle deleted AzureClusters
		if !azCluster.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("AzureCluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		clusterName, ok := GetOwnerClusterName(azCluster.ObjectMeta)
		if !ok {
			log.Info("unable to get the owner cluster")
			return nil
		}

		machineList := &clusterv1.MachineList{}
		machineList.SetGroupVersionKind(gvk)
		// list all of the requested objects within the cluster namespace with the cluster name label
		if err := c.List(ctx, machineList, client.InNamespace(azCluster.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterName}); err != nil {
			return nil
		}

		mapFunc := util.MachineToInfrastructureMapFunc(gvk)
		var results []ctrl.Request
		for _, machine := range machineList.Items {
			m := machine
			azureMachines := mapFunc(ctx, &m)
			results = append(results, azureMachines...)
		}

		return results
	}, nil
}

// GetOwnerClusterName returns the name of the owning Cluster by finding a clusterv1.Cluster in the ownership references.
func GetOwnerClusterName(obj metav1.ObjectMeta) (string, bool) {
	for _, ref := range obj.OwnerReferences {
		if ref.Kind != "Cluster" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return "", false
		}
		if gv.Group == clusterv1.GroupVersion.Group {
			return ref.Name, true
		}
	}
	return "", false
}

// GetObjectsToRequestsByNamespaceAndClusterName returns the slice of ctrl.Requests consisting the list items contained in the unstructured list.
func GetObjectsToRequestsByNamespaceAndClusterName(ctx context.Context, c client.Client, clusterKey client.ObjectKey, list *unstructured.UnstructuredList) []ctrl.Request {
	// list all of the requested objects within the cluster namespace with the cluster name label
	if err := c.List(ctx, list, client.InNamespace(clusterKey.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterKey.Name}); err != nil {
		return nil
	}

	results := make([]ctrl.Request, len(list.Items))
	for i, obj := range list.Items {
		results[i] = ctrl.Request{
			NamespacedName: client.ObjectKey{Namespace: obj.GetNamespace(), Name: obj.GetName()},
		}
	}
	return results
}

// referSameObject returns true if a and b point to the same object.
func referSameObject(a, b metav1.OwnerReference) bool {
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}

	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}

	return aGV.Group == bGV.Group && a.Kind == b.Kind && a.Name == b.Name
}

// GetCloudProviderSecret returns the required azure json secret for the provided parameters.
func GetCloudProviderSecret(d azure.ClusterScoper, namespace, name string, owner metav1.OwnerReference, identityType infrav1.VMIdentity, userIdentityID string) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s-azure-json", name),
			Labels: map[string]string{
				d.ClusterName(): string(infrav1.ResourceLifecycleOwned),
			},
			OwnerReferences: []metav1.OwnerReference{owner},
		},
	}

	var controlPlaneConfig, workerNodeConfig *CloudProviderConfig

	switch identityType {
	case infrav1.VMIdentitySystemAssigned:
		controlPlaneConfig, workerNodeConfig = systemAssignedIdentityCloudProviderConfig(d)
	case infrav1.VMIdentityUserAssigned:
		if len(userIdentityID) < 1 {
			return nil, errors.New("expected a non-empty userIdentityID")
		}
		controlPlaneConfig, workerNodeConfig = userAssignedIdentityCloudProviderConfig(d, userIdentityID)
	case infrav1.VMIdentityNone:
		controlPlaneConfig, workerNodeConfig = newCloudProviderConfig(d)
	}

	// Enable VMSS Flexible nodes if MachinePools are enabled
	if feature.Gates.Enabled(capifeature.MachinePool) {
		if controlPlaneConfig != nil && controlPlaneConfig.VMType == "vmss" {
			controlPlaneConfig.EnableVmssFlexNodes = true
		}
		if workerNodeConfig != nil && workerNodeConfig.VMType == "vmss" {
			workerNodeConfig.EnableVmssFlexNodes = true
		}
	}

	controlPlaneData, err := json.MarshalIndent(controlPlaneConfig, "", "    ")
	if err != nil {
		return nil, errors.Wrap(err, "failed control plane json marshal")
	}
	workerNodeData, err := json.MarshalIndent(workerNodeConfig, "", "    ")
	if err != nil {
		return nil, errors.Wrap(err, "failed worker node json marshal")
	}

	secret.Data = map[string][]byte{
		"control-plane-azure.json": controlPlaneData,
		"worker-node-azure.json":   workerNodeData,
		// added for backwards compatibility
		"azure.json": controlPlaneData,
	}

	return secret, nil
}

func systemAssignedIdentityCloudProviderConfig(d azure.ClusterScoper) (cpConfig *CloudProviderConfig, wkConfig *CloudProviderConfig) {
	controlPlaneConfig, workerConfig := newCloudProviderConfig(d)
	controlPlaneConfig.AadClientID = ""
	controlPlaneConfig.AadClientSecret = ""
	controlPlaneConfig.UseManagedIdentityExtension = true
	workerConfig.AadClientID = ""
	workerConfig.AadClientSecret = ""
	workerConfig.UseManagedIdentityExtension = true
	return controlPlaneConfig, workerConfig
}

func userAssignedIdentityCloudProviderConfig(d azure.ClusterScoper, identityID string) (cpConfig *CloudProviderConfig, wkConfig *CloudProviderConfig) {
	controlPlaneConfig, workerConfig := newCloudProviderConfig(d)
	controlPlaneConfig.AadClientID = ""
	controlPlaneConfig.AadClientSecret = ""
	controlPlaneConfig.UseManagedIdentityExtension = true
	controlPlaneConfig.UserAssignedIdentityID = identityID
	workerConfig.AadClientID = ""
	workerConfig.AadClientSecret = ""
	workerConfig.UseManagedIdentityExtension = true
	workerConfig.UserAssignedIdentityID = identityID
	return controlPlaneConfig, workerConfig
}

func newCloudProviderConfig(d azure.ClusterScoper) (controlPlaneConfig *CloudProviderConfig, workerConfig *CloudProviderConfig) {
	subnet := getOneNodeSubnet(d)
	return (&CloudProviderConfig{
			Cloud:                        d.CloudEnvironment(),
			AadClientID:                  d.ClientID(),
			AadClientSecret:              d.ClientSecret(),
			TenantID:                     d.TenantID(),
			SubscriptionID:               d.SubscriptionID(),
			ResourceGroup:                d.ResourceGroup(),
			SecurityGroupName:            subnet.SecurityGroup.Name,
			SecurityGroupResourceGroup:   d.Vnet().ResourceGroup,
			Location:                     d.Location(),
			ExtendedLocationType:         d.ExtendedLocationType(),
			ExtendedLocationName:         d.ExtendedLocationName(),
			VMType:                       "vmss",
			VnetName:                     d.Vnet().Name,
			VnetResourceGroup:            d.Vnet().ResourceGroup,
			SubnetName:                   subnet.Name,
			RouteTableName:               subnet.RouteTable.Name,
			LoadBalancerSku:              "Standard",
			LoadBalancerName:             d.OutboundLBName(infrav1.Node),
			MaximumLoadBalancerRuleCount: 250,
			UseManagedIdentityExtension:  false,
			UseInstanceMetadata:          true,
		}).overrideFromSpec(d),
		(&CloudProviderConfig{
			Cloud:                        d.CloudEnvironment(),
			AadClientID:                  d.ClientID(),
			AadClientSecret:              d.ClientSecret(),
			TenantID:                     d.TenantID(),
			SubscriptionID:               d.SubscriptionID(),
			ResourceGroup:                d.ResourceGroup(),
			SecurityGroupName:            subnet.SecurityGroup.Name,
			SecurityGroupResourceGroup:   d.Vnet().ResourceGroup,
			Location:                     d.Location(),
			ExtendedLocationType:         d.ExtendedLocationType(),
			ExtendedLocationName:         d.ExtendedLocationName(),
			VMType:                       "vmss",
			VnetName:                     d.Vnet().Name,
			VnetResourceGroup:            d.Vnet().ResourceGroup,
			SubnetName:                   subnet.Name,
			RouteTableName:               subnet.RouteTable.Name,
			LoadBalancerSku:              "Standard",
			LoadBalancerName:             d.OutboundLBName(infrav1.Node),
			MaximumLoadBalancerRuleCount: 250,
			UseManagedIdentityExtension:  false,
			UseInstanceMetadata:          true,
		}).overrideFromSpec(d)
}

// getOneNodeSubnet returns one of the subnets for the node role.
func getOneNodeSubnet(d azure.ClusterScoper) infrav1.SubnetSpec {
	for _, subnet := range d.Subnets() {
		if subnet.Role == infrav1.SubnetNode || subnet.Role == infrav1.SubnetCluster {
			return subnet
		}
	}
	return infrav1.SubnetSpec{}
}

// CloudProviderConfig is an abbreviated version of the same struct in k/k.
type CloudProviderConfig struct {
	Cloud                        string `json:"cloud"`
	TenantID                     string `json:"tenantId"`
	SubscriptionID               string `json:"subscriptionId"`
	AadClientID                  string `json:"aadClientId,omitempty"`
	AadClientSecret              string `json:"aadClientSecret,omitempty"`
	ResourceGroup                string `json:"resourceGroup"`
	SecurityGroupName            string `json:"securityGroupName"`
	SecurityGroupResourceGroup   string `json:"securityGroupResourceGroup"`
	Location                     string `json:"location"`
	ExtendedLocationType         string `json:"extendedLocationType,omitempty"`
	ExtendedLocationName         string `json:"extendedLocationName,omitempty"`
	VMType                       string `json:"vmType"`
	VnetName                     string `json:"vnetName"`
	VnetResourceGroup            string `json:"vnetResourceGroup"`
	SubnetName                   string `json:"subnetName"`
	RouteTableName               string `json:"routeTableName"`
	LoadBalancerSku              string `json:"loadBalancerSku"`
	LoadBalancerName             string `json:"loadBalancerName"`
	MaximumLoadBalancerRuleCount int    `json:"maximumLoadBalancerRuleCount"`
	UseManagedIdentityExtension  bool   `json:"useManagedIdentityExtension"`
	UseInstanceMetadata          bool   `json:"useInstanceMetadata"`
	EnableVmssFlexNodes          bool   `json:"enableVmssFlexNodes,omitempty"`
	UserAssignedIdentityID       string `json:"userAssignedIdentityID,omitempty"`
	CloudProviderRateLimitConfig
	BackOffConfig
}

// overrideFromSpec overrides cloud provider config with the values provided in cluster spec.
func (cpc *CloudProviderConfig) overrideFromSpec(d azure.ClusterScoper) *CloudProviderConfig {
	if d.CloudProviderConfigOverrides() == nil {
		return cpc
	}

	for _, rateLimit := range d.CloudProviderConfigOverrides().RateLimits {
		switch rateLimit.Name {
		case infrav1.DefaultRateLimit:
			cpc.RateLimitConfig = *toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.RouteRateLimit:
			cpc.RouteRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.SubnetsRateLimit:
			cpc.SubnetsRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.InterfaceRateLimit:
			cpc.InterfaceRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.RouteTableRateLimit:
			cpc.RouteTableRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.LoadBalancerRateLimit:
			cpc.LoadBalancerRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.PublicIPAddressRateLimit:
			cpc.PublicIPAddressRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.SecurityGroupRateLimit:
			cpc.SecurityGroupRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.VirtualMachineRateLimit:
			cpc.VirtualMachineRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.StorageAccountRateLimit:
			cpc.StorageAccountRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.DiskRateLimit:
			cpc.DiskRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.SnapshotRateLimit:
			cpc.SnapshotRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.VirtualMachineScaleSetRateLimit:
			cpc.VirtualMachineScaleSetRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.VirtualMachineSizesRateLimit:
			cpc.VirtualMachineSizeRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		case infrav1.AvailabilitySetRateLimit:
			cpc.AvailabilitySetRateLimit = toCloudProviderRateLimitConfig(rateLimit.Config)
		}
	}

	cpc.BackOffConfig = toCloudProviderBackOffConfig(d.CloudProviderConfigOverrides().BackOffs)
	return cpc
}

// toCloudProviderRateLimitConfig returns converts infrav1.RateLimitConfig to RateLimitConfig that is required with the cloud provider.
func toCloudProviderRateLimitConfig(source infrav1.RateLimitConfig) *RateLimitConfig {
	rateLimitConfig := RateLimitConfig{}
	rateLimitConfig.CloudProviderRateLimit = source.CloudProviderRateLimit
	if source.CloudProviderRateLimitQPS != nil {
		rateLimitConfig.CloudProviderRateLimitQPS = float32(source.CloudProviderRateLimitQPS.AsApproximateFloat64())
	}
	rateLimitConfig.CloudProviderRateLimitBucket = source.CloudProviderRateLimitBucket
	if source.CloudProviderRateLimitQPSWrite != nil {
		rateLimitConfig.CloudProviderRateLimitQPSWrite = float32(source.CloudProviderRateLimitQPSWrite.AsApproximateFloat64())
	}
	rateLimitConfig.CloudProviderRateLimitBucketWrite = source.CloudProviderRateLimitBucketWrite
	return &rateLimitConfig
}

// CloudProviderRateLimitConfig represents the rate limiting configurations in azure cloud provider config.
// See: https://cloud-provider-azure.sigs.k8s.io/install/configs/#per-client-rate-limiting.
// This is a copy of the struct used in cloud-provider-azure: https://github.com/kubernetes-sigs/cloud-provider-azure/blob/d585c2031925b39c925624302f22f8856e29e352/pkg/provider/azure_ratelimit.go#L25
type CloudProviderRateLimitConfig struct {
	RateLimitConfig

	RouteRateLimit                  *RateLimitConfig `json:"routeRateLimit,omitempty"`
	SubnetsRateLimit                *RateLimitConfig `json:"subnetsRateLimit,omitempty"`
	InterfaceRateLimit              *RateLimitConfig `json:"interfaceRateLimit,omitempty"`
	RouteTableRateLimit             *RateLimitConfig `json:"routeTableRateLimit,omitempty"`
	LoadBalancerRateLimit           *RateLimitConfig `json:"loadBalancerRateLimit,omitempty"`
	PublicIPAddressRateLimit        *RateLimitConfig `json:"publicIPAddressRateLimit,omitempty"`
	SecurityGroupRateLimit          *RateLimitConfig `json:"securityGroupRateLimit,omitempty"`
	VirtualMachineRateLimit         *RateLimitConfig `json:"virtualMachineRateLimit,omitempty"`
	StorageAccountRateLimit         *RateLimitConfig `json:"storageAccountRateLimit,omitempty"`
	DiskRateLimit                   *RateLimitConfig `json:"diskRateLimit,omitempty"`
	SnapshotRateLimit               *RateLimitConfig `json:"snapshotRateLimit,omitempty"`
	VirtualMachineScaleSetRateLimit *RateLimitConfig `json:"virtualMachineScaleSetRateLimit,omitempty"`
	VirtualMachineSizeRateLimit     *RateLimitConfig `json:"virtualMachineSizesRateLimit,omitempty"`
	AvailabilitySetRateLimit        *RateLimitConfig `json:"availabilitySetRateLimit,omitempty"`
}

// RateLimitConfig indicates the rate limit config options.
// This is a copy of the struct used in cloud-provider-azure: https://github.com/kubernetes-sigs/cloud-provider-azure/blob/d585c2031925b39c925624302f22f8856e29e352/pkg/azureclients/azure_client_config.go#L48
type RateLimitConfig struct {
	CloudProviderRateLimit            bool    `json:"cloudProviderRateLimit,omitempty"`
	CloudProviderRateLimitQPS         float32 `json:"cloudProviderRateLimitQPS,omitempty"`
	CloudProviderRateLimitBucket      int     `json:"cloudProviderRateLimitBucket,omitempty"`
	CloudProviderRateLimitQPSWrite    float32 `json:"cloudProviderRateLimitQPSWrite,omitempty"`
	CloudProviderRateLimitBucketWrite int     `json:"cloudProviderRateLimitBucketWrite,omitempty"`
}

// BackOffConfig indicates the back-off config options.
// This is a copy of the struct used in cloud-provider-azure: https://github.com/kubernetes-sigs/cloud-provider-azure/blob/d585c2031925b39c925624302f22f8856e29e352/pkg/azureclients/azure_client_config.go#L48
type BackOffConfig struct {
	CloudProviderBackoff         bool    `json:"cloudProviderBackoff,omitempty"`
	CloudProviderBackoffRetries  int     `json:"cloudProviderBackoffRetries,omitempty"`
	CloudProviderBackoffExponent float64 `json:"cloudProviderBackoffExponent,omitempty"`
	CloudProviderBackoffDuration int     `json:"cloudProviderBackoffDuration,omitempty"`
	CloudProviderBackoffJitter   float64 `json:"cloudProviderBackoffJitter,omitempty"`
}

// toCloudProviderBackOffConfig returns converts infrav1.BackOffConfig to BackOffConfig that is required with the cloud provider.
func toCloudProviderBackOffConfig(source infrav1.BackOffConfig) BackOffConfig {
	backOffConfig := BackOffConfig{}
	backOffConfig.CloudProviderBackoff = source.CloudProviderBackoff
	if source.CloudProviderBackoffExponent != nil {
		backOffConfig.CloudProviderBackoffExponent = source.CloudProviderBackoffExponent.AsApproximateFloat64()
	}
	backOffConfig.CloudProviderBackoffRetries = source.CloudProviderBackoffRetries
	if source.CloudProviderBackoffJitter != nil {
		backOffConfig.CloudProviderBackoffJitter = source.CloudProviderBackoffJitter.AsApproximateFloat64()
	}
	backOffConfig.CloudProviderBackoffDuration = source.CloudProviderBackoffDuration
	return backOffConfig
}

func reconcileAzureSecret(ctx context.Context, kubeclient client.Client, owner metav1.OwnerReference, newSecret *corev1.Secret, clusterName string) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.reconcileAzureSecret")
	defer done()

	// Fetch previous secret, if it exists
	key := types.NamespacedName{
		Namespace: newSecret.Namespace,
		Name:      newSecret.Name,
	}
	old := &corev1.Secret{}
	err := kubeclient.Get(ctx, key, old)
	if err != nil && !apierrors.IsNotFound(err) {
		return errors.Wrap(err, "failed to fetch existing secret")
	}

	// Create if it wasn't found
	if apierrors.IsNotFound(err) {
		if err := kubeclient.Create(ctx, newSecret); err != nil && !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "failed to create secret")
		}
		return nil
	}

	tag, exists := old.Labels[clusterName]

	if !exists || tag != string(infrav1.ResourceLifecycleOwned) {
		log.V(2).Info("returning early from secret reconcile, user provided secret already exists")
		return nil
	}

	// Otherwise, check ownership and data freshness. Update as necessary
	hasOwner := false
	for _, ownerRef := range old.OwnerReferences {
		if referSameObject(ownerRef, owner) {
			hasOwner = true
			break
		}
	}

	hasData := equality.Semantic.DeepEqual(old.Data, newSecret.Data)
	if hasData && hasOwner {
		// no update required
		log.V(2).Info("returning early from secret reconcile, no update needed")
		return nil
	}

	if !hasOwner {
		old.OwnerReferences = append(old.OwnerReferences, owner)
	}

	if !hasData {
		old.Data = newSecret.Data
	}

	log.V(2).Info("updating azure secret")
	if err := kubeclient.Update(ctx, old); err != nil {
		return errors.Wrap(err, "failed to update secret when diff was required")
	}

	log.V(2).Info("done updating secret")

	return nil
}

// GetOwnerMachinePool returns the MachinePool object owning the current resource.
func GetOwnerMachinePool(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*expv1.MachinePool, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.GetOwnerMachinePool")
	defer done()

	for _, ref := range obj.OwnerReferences {
		if ref.Kind != "MachinePool" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if gv.Group == expv1.GroupVersion.Group {
			return GetMachinePoolByName(ctx, c, obj.Namespace, ref.Name)
		}
	}
	return nil, nil
}

// GetOwnerAzureMachinePool returns the AzureMachinePool object owning the current resource.
func GetOwnerAzureMachinePool(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*infrav1exp.AzureMachinePool, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.GetOwnerAzureMachinePool")
	defer done()

	for _, ref := range obj.OwnerReferences {
		if ref.Kind != infrav1.AzureMachinePoolKind {
			continue
		}

		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if gv.Group == infrav1exp.GroupVersion.Group {
			return GetAzureMachinePoolByName(ctx, c, obj.Namespace, ref.Name)
		}
	}
	return nil, nil
}

// GetMachinePoolByName finds and return a MachinePool object using the specified params.
func GetMachinePoolByName(ctx context.Context, c client.Client, namespace, name string) (*expv1.MachinePool, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.GetMachinePoolByName")
	defer done()

	m := &expv1.MachinePool{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetAzureMachinePoolByName finds and return an AzureMachinePool object using the specified params.
func GetAzureMachinePoolByName(ctx context.Context, c client.Client, namespace, name string) (*infrav1exp.AzureMachinePool, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.GetAzureMachinePoolByName")
	defer done()

	m := &infrav1exp.AzureMachinePool{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, m); err != nil {
		return nil, err
	}
	return m, nil
}

// ShouldDeleteIndividualResources returns false if the resource group is managed and the whole cluster is being deleted
// meaning that we can rely on a single resource group delete operation as opposed to deleting every individual VM resource.
func ShouldDeleteIndividualResources(ctx context.Context, cluster ClusterScoper) bool {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.ShouldDeleteIndividualResources")
	defer done()

	if cluster.GetDeletionTimestamp().IsZero() {
		return true
	}

	managed, err := groups.New(cluster).IsManaged(ctx)
	// Since this is a best effort attempt to speed up delete, we don't fail the delete if we can't get the RG status.
	// Instead, take the long way and delete all resources one by one.
	return err != nil || !managed
}

// GetClusterIdentityFromRef returns the AzureClusterIdentity referenced by the AzureCluster.
func GetClusterIdentityFromRef(ctx context.Context, c client.Client, azureClusterNamespace string, ref *corev1.ObjectReference) (*infrav1.AzureClusterIdentity, error) {
	identity := &infrav1.AzureClusterIdentity{}
	if ref != nil {
		namespace := ref.Namespace
		if namespace == "" {
			namespace = azureClusterNamespace
		}
		key := client.ObjectKey{Name: ref.Name, Namespace: namespace}
		if err := c.Get(ctx, key, identity); err != nil {
			return nil, err
		}
		return identity, nil
	}
	return nil, nil
}

// deprecatedClusterIdentityFinalizer is was briefly used to compute a finalizer without a hash in releases v1.5.1 and v1.4.4.
// It is kept here to ensure that we can remove it from existing clusters for backwards compatibility.
// This function should be removed in a future release.
func deprecatedClusterIdentityFinalizer(prefix, clusterNamespace, clusterName string) string {
	return fmt.Sprintf("%s/%s-%s", prefix, clusterNamespace, clusterName)
}

// clusterIdentityFinalizer generates a finalizer key.
// The finalizer key is a combination of the prefix and a hash of the cluster name and namespace.
// We use a hash to ensure that the finalizer key name is not longer than 63 characters.
func clusterIdentityFinalizer(prefix, clusterNamespace, clusterName string) string {
	hash := sha256.Sum224([]byte(fmt.Sprintf("%s-%s", clusterNamespace, clusterName)))
	return fmt.Sprintf("%s/%s", prefix, hex.EncodeToString(hash[:]))
}

// EnsureClusterIdentity ensures that the identity ref is allowed in the namespace and sets a finalizer.
func EnsureClusterIdentity(ctx context.Context, c client.Client, object conditions.Setter, identityRef *corev1.ObjectReference, finalizerPrefix string) error {
	name := object.GetName()
	namespace := object.GetNamespace()
	identity, err := GetClusterIdentityFromRef(ctx, c, namespace, identityRef)
	if err != nil {
		return err
	}

	if !scope.IsClusterNamespaceAllowed(ctx, c, identity.Spec.AllowedNamespaces, namespace) {
		conditions.MarkFalse(object, infrav1.NetworkInfrastructureReadyCondition, infrav1.NamespaceNotAllowedByIdentity, clusterv1.ConditionSeverityError, "")
		return errors.New("AzureClusterIdentity list of allowed namespaces doesn't include current cluster namespace")
	}

	// Remove deprecated finalizer if it exists, Register the finalizer immediately to avoid orphaning Azure resources on delete.
	needsPatch := controllerutil.RemoveFinalizer(identity, deprecatedClusterIdentityFinalizer(finalizerPrefix, namespace, name))
	needsPatch = controllerutil.AddFinalizer(identity, clusterIdentityFinalizer(finalizerPrefix, namespace, name)) || needsPatch
	if needsPatch {
		// finalizers are added/removed then patch the object
		identityHelper, err := patch.NewHelper(identity, c)
		if err != nil {
			return errors.Wrap(err, "failed to init patch helper")
		}
		return identityHelper.Patch(ctx, identity)
	}

	return nil
}

// RemoveClusterIdentityFinalizer removes the finalizer on an AzureClusterIdentity.
func RemoveClusterIdentityFinalizer(ctx context.Context, c client.Client, object client.Object, identityRef *corev1.ObjectReference, finalizerPrefix string) error {
	name := object.GetName()
	namespace := object.GetNamespace()
	identity, err := GetClusterIdentityFromRef(ctx, c, namespace, identityRef)
	if err != nil {
		return err
	}
	identityHelper, err := patch.NewHelper(identity, c)
	if err != nil {
		return errors.Wrap(err, "failed to init patch helper")
	}
	controllerutil.RemoveFinalizer(identity, clusterIdentityFinalizer(finalizerPrefix, namespace, name))
	err = identityHelper.Patch(ctx, identity)
	if err != nil {
		return errors.Wrap(err, "failed to patch AzureClusterIdentity")
	}
	return nil
}

// MachinePoolToInfrastructureMapFunc returns a handler.MapFunc that watches for
// MachinePool events and returns reconciliation requests for an infrastructure provider object.
func MachinePoolToInfrastructureMapFunc(gvk schema.GroupVersionKind, log logr.Logger) handler.MapFunc {
	return func(_ context.Context, o client.Object) []reconcile.Request {
		m, ok := o.(*expv1.MachinePool)
		if !ok {
			log.V(4).Info("attempt to map incorrect type", "type", fmt.Sprintf("%T", o))
			return nil
		}

		gk := gvk.GroupKind()
		ref := m.Spec.Template.Spec.InfrastructureRef
		// Return early if the GroupKind doesn't match what we expect.
		infraGK := ref.GroupVersionKind().GroupKind()
		if gk != infraGK {
			log.V(4).Info("gk does not match", "gk", gk, "infraGK", infraGK)
			return nil
		}

		return []reconcile.Request{
			{
				NamespacedName: client.ObjectKey{
					Namespace: m.Namespace,
					Name:      ref.Name,
				},
			},
		}
	}
}

// AzureManagedClusterToAzureManagedMachinePoolsMapper creates a mapping handler to transform AzureManagedClusters into
// AzureManagedMachinePools. The transform requires AzureManagedCluster to map to the owning Cluster, then from the
// Cluster, collect the MachinePools belonging to the cluster, then finally projecting the infrastructure reference
// to the AzureManagedMachinePools.
func AzureManagedClusterToAzureManagedMachinePoolsMapper(_ context.Context, c client.Client, scheme *runtime.Scheme, log logr.Logger) (handler.MapFunc, error) {
	gvk, err := apiutil.GVKForObject(new(infrav1.AzureManagedMachinePool), scheme)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find GVK for AzureManagedMachinePool")
	}

	return func(ctx context.Context, o client.Object) []ctrl.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		azCluster, ok := o.(*infrav1.AzureManagedCluster)
		if !ok {
			log.Error(errors.Errorf("expected an AzureManagedCluster, got %T instead", o.GetObjectKind()), "failed to map AzureManagedCluster")
			return nil
		}

		log := log.WithValues("AzureManagedCluster", azCluster.Name, "Namespace", azCluster.Namespace)

		// Don't handle deleted AzureManagedClusters
		if !azCluster.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("AzureManagedCluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		clusterName, ok := GetOwnerClusterName(azCluster.ObjectMeta)
		if !ok {
			log.V(4).Info("unable to get the owner cluster")
			return nil
		}

		machineList := &expv1.MachinePoolList{}
		machineList.SetGroupVersionKind(gvk)
		// list all of the requested objects within the cluster namespace with the cluster name label
		if err := c.List(ctx, machineList, client.InNamespace(azCluster.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterName}); err != nil {
			return nil
		}

		mapFunc := MachinePoolToInfrastructureMapFunc(gvk, log)
		var results []ctrl.Request
		for _, machine := range machineList.Items {
			m := machine
			azureMachines := mapFunc(ctx, &m)
			results = append(results, azureMachines...)
		}

		return results
	}, nil
}

// AzureManagedControlPlaneToAzureManagedMachinePoolsMapper creates a mapping handler to transform AzureManagedControlPlanes into
// AzureManagedMachinePools. The transform requires AzureManagedControlPlane to map to the owning Cluster, then from the
// Cluster, collect the MachinePools belonging to the cluster, then finally projecting the infrastructure reference
// to the AzureManagedMachinePools.
func AzureManagedControlPlaneToAzureManagedMachinePoolsMapper(_ context.Context, c client.Client, scheme *runtime.Scheme, log logr.Logger) (handler.MapFunc, error) {
	gvk, err := apiutil.GVKForObject(new(infrav1.AzureManagedMachinePool), scheme)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find GVK for AzureManagedMachinePool")
	}

	return func(ctx context.Context, o client.Object) []ctrl.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		azControlPlane, ok := o.(*infrav1.AzureManagedControlPlane)
		if !ok {
			log.Error(errors.Errorf("expected an AzureManagedControlPlane, got %T instead", o.GetObjectKind()), "failed to map AzureManagedControlPlane")
			return nil
		}

		log := log.WithValues("AzureManagedControlPlane", azControlPlane.Name, "Namespace", azControlPlane.Namespace)

		// Don't handle deleted AzureManagedControlPlanes
		if !azControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("AzureManagedControlPlane has a deletion timestamp, skipping mapping.")
			return nil
		}

		clusterName, ok := GetOwnerClusterName(azControlPlane.ObjectMeta)
		if !ok {
			log.Info("unable to get the owner cluster")
			return nil
		}

		machineList := &expv1.MachinePoolList{}
		machineList.SetGroupVersionKind(gvk)
		// list all of the requested objects within the cluster namespace with the cluster name label
		if err := c.List(ctx, machineList, client.InNamespace(azControlPlane.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterName}); err != nil {
			return nil
		}

		mapFunc := MachinePoolToInfrastructureMapFunc(gvk, log)
		var results []ctrl.Request
		for _, machine := range machineList.Items {
			m := machine
			azureMachines := mapFunc(ctx, &m)
			results = append(results, azureMachines...)
		}

		return results
	}, nil
}

// AzureManagedClusterToAzureManagedControlPlaneMapper creates a mapping handler to transform AzureManagedClusters into
// AzureManagedControlPlane. The transform requires AzureManagedCluster to map to the owning Cluster, then from the
// Cluster, collect the control plane infrastructure reference.
func AzureManagedClusterToAzureManagedControlPlaneMapper(_ context.Context, c client.Client, log logr.Logger) (handler.MapFunc, error) {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		azCluster, ok := o.(*infrav1.AzureManagedCluster)
		if !ok {
			log.Error(errors.Errorf("expected an AzureManagedCluster, got %T instead", o), "failed to map AzureManagedCluster")
			return nil
		}

		log := log.WithValues("AzureManagedCluster", azCluster.Name, "Namespace", azCluster.Namespace)

		// Don't handle deleted AzureManagedClusters
		if !azCluster.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("AzureManagedCluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, c, azCluster.ObjectMeta)
		if err != nil {
			log.Error(err, "failed to get the owning cluster")
			return nil
		}

		if cluster == nil {
			log.Error(err, "cluster has not set owner ref yet")
			return nil
		}

		ref := cluster.Spec.ControlPlaneRef
		if ref == nil || ref.Name == "" {
			return nil
		}

		return []ctrl.Request{
			{
				NamespacedName: types.NamespacedName{
					Namespace: ref.Namespace,
					Name:      ref.Name,
				},
			},
		}
	}, nil
}

// AzureManagedControlPlaneToAzureManagedClusterMapper creates a mapping handler to transform AzureManagedClusters into
// AzureManagedControlPlane. The transform requires AzureManagedCluster to map to the owning Cluster, then from the
// Cluster, collect the control plane infrastructure reference.
func AzureManagedControlPlaneToAzureManagedClusterMapper(_ context.Context, c client.Client, log logr.Logger) (handler.MapFunc, error) {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		azManagedControlPlane, ok := o.(*infrav1.AzureManagedControlPlane)
		if !ok {
			log.Error(errors.Errorf("expected an AzureManagedControlPlane, got %T instead", o), "failed to map AzureManagedControlPlane")
			return nil
		}

		log := log.WithValues("AzureManagedControlPlane", azManagedControlPlane.Name, "Namespace", azManagedControlPlane.Namespace)

		// Don't handle deleted AzureManagedControlPlanes
		if !azManagedControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("AzureManagedControlPlane has a deletion timestamp, skipping mapping.")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, c, azManagedControlPlane.ObjectMeta)
		if err != nil {
			log.Error(err, "failed to get the owning cluster")
			return nil
		}

		if cluster == nil {
			log.Error(err, "cluster has not set owner ref yet")
			return nil
		}

		ref := cluster.Spec.InfrastructureRef
		if ref == nil || ref.Name == "" {
			return nil
		}

		return []ctrl.Request{
			{
				NamespacedName: types.NamespacedName{
					Namespace: ref.Namespace,
					Name:      ref.Name,
				},
			},
		}
	}, nil
}

// MachinePoolToAzureManagedControlPlaneMapFunc returns a handler.MapFunc that watches for
// MachinePool events and returns reconciliation requests for a control plane object.
func MachinePoolToAzureManagedControlPlaneMapFunc(_ context.Context, c client.Client, gvk schema.GroupVersionKind, log logr.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultMappingTimeout)
		defer cancel()

		machinePool, ok := o.(*expv1.MachinePool)
		if !ok {
			log.Info("expected a MachinePool, got wrong type", "type", fmt.Sprintf("%T", o))
			return nil
		}

		cluster, err := util.GetClusterByName(ctx, c, machinePool.ObjectMeta.Namespace, machinePool.Spec.ClusterName)
		if err != nil {
			log.Error(err, "failed to get the owning cluster")
			return nil
		}

		gk := gvk.GroupKind()
		ref := cluster.Spec.ControlPlaneRef
		if ref == nil || ref.Name == "" {
			log.Info("control plane ref is nil or empty: control plane ref not found")
			return nil
		}
		// Return early if the GroupKind doesn't match what we expect.
		controlPlaneGK := ref.GroupVersionKind().GroupKind()
		if gk != controlPlaneGK {
			// MachinePool does not correlate to a AzureManagedControlPlane, nothing to do
			return nil
		}

		controlPlaneKey := client.ObjectKey{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		}
		controlPlane := &infrav1.AzureManagedControlPlane{}
		if err := c.Get(ctx, controlPlaneKey, controlPlane); err != nil {
			log.Error(err, "failed to fetch default pool reference")
			// If we get here, we might want to reconcile but aren't sure.
			// Do it anyway to be safe. Worst case we reconcile a few extra times with no-ops.
			return []reconcile.Request{
				{
					NamespacedName: client.ObjectKey{
						Namespace: ref.Namespace,
						Name:      ref.Name,
					},
				},
			}
		}

		infraMachinePoolRef := machinePool.Spec.Template.Spec.InfrastructureRef

		gv, err := schema.ParseGroupVersion(infraMachinePoolRef.APIVersion)
		if err != nil {
			log.Error(err, "failed to parse group version")
			// If we get here, we might want to reconcile but aren't sure.
			// Do it anyway to be safe. Worst case we reconcile a few extra times with no-ops.
			return []reconcile.Request{
				{
					NamespacedName: client.ObjectKey{
						Namespace: ref.Namespace,
						Name:      ref.Name,
					},
				},
			}
		}

		kindMatches := infraMachinePoolRef.Kind == "AzureManagedMachinePool"
		groupMatches := controlPlaneGK.Group == gv.Group

		ammp := &infrav1.AzureManagedMachinePool{}
		key := types.NamespacedName{Namespace: infraMachinePoolRef.Namespace, Name: infraMachinePoolRef.Name}
		if err := c.Get(ctx, key, ammp); err != nil {
			log.Error(err, fmt.Sprintf("failed to fetch azure managed machine pool for Machinepool: %s", infraMachinePoolRef.Name))
			// If we get here, we might want to reconcile but aren't sure.
			// Do it anyway to be safe. Worst case we reconcile a few extra times with no-ops.
			return []reconcile.Request{
				{
					NamespacedName: client.ObjectKey{
						Namespace: ref.Namespace,
						Name:      ref.Name,
					},
				},
			}
		}

		isSystemNodePool := ammp.Spec.Mode == string(infrav1.NodePoolModeSystem)

		if groupMatches && kindMatches && isSystemNodePool {
			return []reconcile.Request{
				{
					NamespacedName: client.ObjectKey{
						Namespace: ref.Namespace,
						Name:      ref.Name,
					},
				},
			}
		}

		// By default, return nothing for a machine pool which is not the default pool for a control plane.
		return nil
	}
}

// ClusterUpdatePauseChange returns a predicate that returns true for an update event when a cluster's
// Spec.Paused changes between any two distinct values.
func ClusterUpdatePauseChange(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := logger.WithValues("predicate", "ClusterUpdatePauseChange", "eventType", "update")

			oldCluster, ok := e.ObjectOld.(*clusterv1.Cluster)
			if !ok {
				log.V(4).Info("Expected Cluster", "type", fmt.Sprintf("%T", e.ObjectOld))
				return false
			}
			log = log.WithValues("Cluster", klog.KObj(oldCluster))

			newCluster := e.ObjectNew.(*clusterv1.Cluster)

			if oldCluster.Spec.Paused != newCluster.Spec.Paused {
				log.V(4).Info("Cluster paused status changed, allowing further processing")
				return true
			}

			log.V(6).Info("Cluster paused status remained the same, blocking further processing")
			return false
		},
		CreateFunc:  func(_ event.CreateEvent) bool { return false },
		DeleteFunc:  func(_ event.DeleteEvent) bool { return false },
		GenericFunc: func(_ event.GenericEvent) bool { return false },
	}
}

// ClusterPauseChangeAndInfrastructureReady is based on ClusterUnpausedAndInfrastructureReady, but
// additionally accepts Cluster pause events.
func ClusterPauseChangeAndInfrastructureReady(scheme *runtime.Scheme, log logr.Logger) predicate.Funcs {
	return predicates.Any(scheme, log, predicates.ClusterCreateInfraReady(scheme, log), predicates.ClusterUpdateInfraReady(scheme, log), ClusterUpdatePauseChange(log)) //nolint:staticcheck
}

// GetClusterScoper returns a ClusterScoper for the given cluster using the infra ref pointing to either an AzureCluster or an AzureManagedCluster.
func GetClusterScoper(ctx context.Context, logger logr.Logger, c client.Client, cluster *clusterv1.Cluster, timeouts reconciler.Timeouts, credCache azure.CredentialCache) (ClusterScoper, error) {
	infraRef := cluster.Spec.InfrastructureRef
	switch infraRef.Kind {
	case "AzureCluster":
		logger = logger.WithValues("AzureCluster", infraRef.Name)
		azureClusterName := client.ObjectKey{
			Namespace: infraRef.Namespace,
			Name:      infraRef.Name,
		}
		azureCluster := &infrav1.AzureCluster{}
		if err := c.Get(ctx, azureClusterName, azureCluster); err != nil {
			logger.V(2).Info("AzureCluster is not available yet")
			return nil, err
		}

		// Create the cluster scope
		return scope.NewClusterScope(ctx, scope.ClusterScopeParams{
			Client:          c,
			Cluster:         cluster,
			AzureCluster:    azureCluster,
			Timeouts:        timeouts,
			CredentialCache: credCache,
		})

	case "AzureManagedCluster":
		logger = logger.WithValues("AzureManagedCluster", infraRef.Name)
		azureManagedControlPlaneName := client.ObjectKey{
			Namespace: infraRef.Namespace,
			Name:      cluster.Spec.ControlPlaneRef.Name,
		}
		azureManagedControlPlane := &infrav1.AzureManagedControlPlane{}
		if err := c.Get(ctx, azureManagedControlPlaneName, azureManagedControlPlane); err != nil {
			logger.V(2).Info("AzureManagedControlPlane is not available yet")
			return nil, err
		}

		// Create the control plane scope
		return scope.NewManagedControlPlaneScope(ctx, scope.ManagedControlPlaneScopeParams{
			Client:          c,
			Cluster:         cluster,
			ControlPlane:    azureManagedControlPlane,
			Timeouts:        timeouts,
			CredentialCache: credCache,
		})
	}

	return nil, errors.Errorf("unsupported infrastructure type %q, should be AzureCluster or AzureManagedCluster", cluster.Spec.InfrastructureRef.Kind)
}

// AddBlockMoveAnnotation adds CAPI's block-move annotation and returns whether or not the annotation was added.
func AddBlockMoveAnnotation(obj metav1.Object) bool {
	annotations := obj.GetAnnotations()

	if _, exists := annotations[clusterctlv1.BlockMoveAnnotation]; exists {
		return false
	}

	if annotations == nil {
		annotations = make(map[string]string)
	}

	// this value doesn't mean anything, only the presence of the annotation matters.
	annotations[clusterctlv1.BlockMoveAnnotation] = "true"
	obj.SetAnnotations(annotations)

	return true
}

// RemoveBlockMoveAnnotation removes CAPI's block-move annotation from the object.
func RemoveBlockMoveAnnotation(obj metav1.Object) {
	azClusterAnnotations := obj.GetAnnotations()
	delete(azClusterAnnotations, clusterctlv1.BlockMoveAnnotation)
	obj.SetAnnotations(azClusterAnnotations)
}
