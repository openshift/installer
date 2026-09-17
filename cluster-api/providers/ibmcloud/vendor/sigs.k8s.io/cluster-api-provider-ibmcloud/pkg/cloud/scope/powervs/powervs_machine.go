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

package powervs

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	ignV3Types "github.com/coreos/ignition/v2/config/v3_4/types"
	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
	ignV2Types "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/bootstrap/ignition"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/cos"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/vpc"
)

const cosURLDomain = "cloud-object-storage.appdomain.cloud"

// presignExpiry is the lifetime of the pre-signed ignition URL embedded in user-data.
// One hour is enough for any reasonable bootstrap window; after expiry the URL is worthless.
const presignExpiry = time.Hour

// zoneCacheEntry holds the supported system types and the exact time they were fetched for a specific zone.
type zoneCacheEntry struct {
	supportedTypes []string
	lastFetch      time.Time
}

// systemTypeCache stores supported system types per datacenter to avoid frequent API calls.
type systemTypeCache struct {
	mu       sync.RWMutex
	zonesMap map[string]zoneCacheEntry
	ttl      time.Duration
}

// Global instance of the cache (TTL set to 6 hours).
var sysCache = &systemTypeCache{
	zonesMap: make(map[string]zoneCacheEntry),
	ttl:      6 * time.Hour,
}

// ConfigurationError represents an error due to invalid machine configuration.
type ConfigurationError struct {
	message string
}

func (e *ConfigurationError) Error() string {
	return e.message
}

// NewConfigurationError creates a new ConfigurationError.
func NewConfigurationError(message string) *ConfigurationError {
	return &ConfigurationError{message: message}
}

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	Client            client.Client
	Cluster           *clusterv1.Cluster
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
	Machine           *clusterv1.Machine
	IBMPowerVSMachine *infrav1.IBMPowerVSMachine
	IBMPowerVSImage   *infrav1.IBMPowerVSImage
	ServiceEndpoint   []endpoints.ServiceEndpoint
	DHCPIPCacheStore  cache.Store

	ClientBuilder ClientBuilder
	Recorder      cgrecord.EventRecorder
}

// MachineScope defines a scope defined around a Power VS Machine.
type MachineScope struct {
	Client client.Client

	IBMPowerVSClient powervs.PowerVS
	IBMVPCClient     vpc.Vpc
	ResourceClient   resourcecontroller.ResourceController
	// COSClient is an HMAC-credential COS client (built via GetHMACCOSClient).
	// It supports SigV4 pre-signed URL generation; the IAM OAuth client does not.
	COSClient cos.Cos

	// Recorder is the event recorder for emitting Kubernetes events.
	Recorder cgrecord.EventRecorder

	Cluster           *clusterv1.Cluster
	Machine           *clusterv1.Machine
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
	IBMPowerVSMachine *infrav1.IBMPowerVSMachine
	IBMPowerVSImage   *infrav1.IBMPowerVSImage
	ServiceEndpoint   []endpoints.ServiceEndpoint
	DHCPIPCacheStore  cache.Store
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
func NewMachineScope(ctx context.Context, params MachineScopeParams) (*MachineScope, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	scope := &MachineScope{
		Client:            params.Client,
		Cluster:           params.Cluster,
		Machine:           params.Machine,
		IBMPowerVSCluster: params.IBMPowerVSCluster,
		IBMPowerVSMachine: params.IBMPowerVSMachine,
		IBMPowerVSImage:   params.IBMPowerVSImage,
		ServiceEndpoint:   params.ServiceEndpoint,
		DHCPIPCacheStore:  params.DHCPIPCacheStore,
		Recorder:          params.Recorder,
	}

	if err := scope.initClients(ctx, &params); err != nil {
		return nil, fmt.Errorf("failed to initialize IBM Cloud clients for machine: %w", err)
	}

	return scope, nil
}

// validate ensures all required fields are present before scope creation.
func (p *MachineScopeParams) validate() error {
	if p.Client == nil {
		return errors.New("client is required when creating a MachineScope")
	}
	if p.Machine == nil {
		return errors.New("machine is required when creating a MachineScope")
	}
	if p.Cluster == nil {
		return errors.New("cluster is required when creating a MachineScope")
	}
	if p.IBMPowerVSMachine == nil {
		return errors.New("ibmPowerVSMachine is required when creating a MachineScope")
	}
	if p.IBMPowerVSCluster == nil {
		return errors.New("ibmPowerVSCluster is required when creating a MachineScope")
	}
	if p.ClientBuilder == nil {
		return errors.New("ClientBuilder is required when creating a MachineScope")
	}
	return nil
}

func (s *MachineScope) initClients(ctx context.Context, params *MachineScopeParams) error {
	log := ctrl.LoggerFrom(ctx)

	auth, err := params.ClientBuilder.GetAuthenticator(ctx)
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	// 1. Build Base Options
	opts := ClientOptions{
		Authenticator:   auth,
		ServiceEndpoint: s.ServiceEndpoint,
		Debug:           log.V(DEBUGLEVEL).Enabled(),
	}

	// 2. Build ResourceController
	s.ResourceClient, err = params.ClientBuilder.GetResourceControllerClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create Resource Controller client: %w", err)
	}

	// 3. Resolve Workspace and Zone
	workspaceID, workspaceZone, err := s.resolveWorkspace(ctx)
	if err != nil {
		return err
	}

	// Set Scope Region/Zone from the resolved workspace
	s.SetZone(workspaceZone)
	s.SetRegion(endpoints.ConstructRegionFromZone(workspaceZone))

	// 4. Build PowerVS Client with the resolved Workspace ID
	opts.WorkspaceID = workspaceID
	opts.Zone = workspaceZone

	s.IBMPowerVSClient, err = params.ClientBuilder.GetPowerVSClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create PowerVS client: %w", err)
	}

	// 5. Build VPC Client
	vpcRegion := s.IBMPowerVSCluster.Spec.VPC.Region
	if vpcRegion == "" {
		vpcRegion, err = regionUtil.VPCRegionForPowerVSRegion(s.GetRegion())
		if err != nil {
			return fmt.Errorf("failed to determine VPC region from PowerVS region: %w", err)
		}
	}
	opts.VPCRegion = vpcRegion

	s.IBMVPCClient, err = params.ClientBuilder.GetVPCClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create IBM VPC client: %w", err)
	}

	// 6. Build COS Client — only when Ignition is configured.
	if s.useIgnition() {
		s.COSClient, err = params.ClientBuilder.GetHMACCOSClient(ctx, s.hmacCOSClientOptions())
		if err != nil {
			return fmt.Errorf("failed to create HMAC COS client: %w", err)
		}
	}

	return nil
}

// resolveWorkspace figures out which workspace this machine should belong to, validates it, and returns the ID and Zone.
func (s *MachineScope) resolveWorkspace(_ context.Context) (string, string, error) {
	var workspaceID, workspaceName string

	// 1. Check if the Machine explicitly overrides the Workspace
	if s.IBMPowerVSMachine.Spec.Workspace.ID != "" {
		workspaceID = s.IBMPowerVSMachine.Spec.Workspace.ID
	} else if s.IBMPowerVSMachine.Spec.Workspace.Name != "" {
		workspaceName = s.IBMPowerVSMachine.Spec.Workspace.Name
	} else {
		// 2. Inherit from the Cluster
		if s.IBMPowerVSCluster.Status.Workspace.ID != "" {
			workspaceID = s.IBMPowerVSCluster.Status.Workspace.ID
		} else {
			return "", "", errors.New("PowerVS workspace ID is not yet populated in the cluster status")
		}
	}

	// 3. Validate against IBM Cloud
	filter := resourcecontroller.InstanceFilter{
		ID:             workspaceID,
		Name:           workspaceName,
		Zone:           &s.IBMPowerVSCluster.Spec.Zone,
		ResourceID:     resourcecontroller.PowerVSResourceID,
		ResourcePlanID: resourcecontroller.PowerVSResourcePlanID,
	}

	workspace, err := s.ResourceClient.GetResourceInstanceByFilter(filter)
	if err != nil {
		return "", "", fmt.Errorf("failed to get PowerVS workspace details (name: %q, id: %q): %w", workspaceName, workspaceID, err)
	}

	if workspace == nil || workspace.GUID == nil || workspace.RegionID == nil {
		return "", "", fmt.Errorf("PowerVS workspace or GUID or RegionID is nil (name: %q, id: %q)", workspaceName, workspaceID)
	}
	if workspace.State == nil || *workspace.State != string(infrav1.WorkspaceStateActive) {
		return "", "", fmt.Errorf("PowerVS workspace (name: %q, id: %q) is not in active state", workspaceName, workspaceID)
	}

	return *workspace.GUID, *workspace.RegionID, nil
}

// CreateMachine creates a PowerVS machine.
//
//nolint:gocyclo
func (s *MachineScope) CreateMachine(ctx context.Context) (*models.PVMInstanceReference, error) {
	log := ctrl.LoggerFrom(ctx)
	machineSpec := s.IBMPowerVSMachine.Spec

	// 1. Idempotency Check: check if the instance already exist
	instanceReply, err := s.ensureInstanceUnique(ctx, s.IBMPowerVSMachine.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to verify instance uniqueness: %w", err)
	} else if instanceReply != nil {
		log.Info("PowerVS instance already exists", "instanceID", *instanceReply.PvmInstanceID)
		return instanceReply, nil
	}

	// 2. Prevent Duplicate API Calls:
	// If the creation request was just sent and K8s status is pending/unknown, wait.
	for _, con := range s.IBMPowerVSMachine.Status.Conditions {
		if con.Type == infrav1.InstanceReadyCondition && con.Status == metav1.ConditionUnknown {
			log.Info("Instance creation already triggered, waiting for cloud status update")
			return nil, nil
		}
	}

	// 3. Validate SystemType Capabilities
	if machineSpec.SystemType != "" {
		valid, supportedTypes, err := s.validateSystemType(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to validate systemType: %w", err)
		}
		if !valid {
			return nil, NewConfigurationError(fmt.Sprintf("systemType '%s' is not supported in this zone. Supported types: %v", machineSpec.SystemType, supportedTypes))
		}
	}

	// 4. Resolve UserData (Ignition / Cloud-init)
	userData, err := s.resolveUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve userdata: %w", err)
	}

	// 5. Parse Processors
	var processors float64
	switch machineSpec.Processors.Type {
	case intstr.Int:
		processors = float64(machineSpec.Processors.IntVal)
	case intstr.String:
		processors, err = strconv.ParseFloat(machineSpec.Processors.StrVal, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Processors (%s) to float64: %w", machineSpec.Processors.StrVal, err)
		}
	}

	// 6. Resolve Image ID
	var imageID string
	if machineSpec.Image.Type == infrav1.ImageSourceTypeImport {
		if s.IBMPowerVSImage == nil || s.IBMPowerVSImage.Status.ImageID == "" {
			return nil, fmt.Errorf("imported image is not ready yet")
		}
		imageID = s.IBMPowerVSImage.Status.ImageID
	} else {
		imageID, err = s.getImageID(ctx, machineSpec.Image.Reference)
		if err != nil {
			s.Recorder.Eventf(s.IBMPowerVSMachine, corev1.EventTypeWarning, "FailedRetrieveImage", "Failed image retrieval: %v", err)
			return nil, fmt.Errorf("error getting image ID from reference: %w", err)
		}
	}
	log.V(3).Info("Resolved image ID", "imageID", imageID)

	// 7. Resolve Network ID
	network := machineSpec.Network

	// Fallback to cluster network if explicitly omitted on the machine
	if network.ID == "" && network.Name == "" {
		networkID := s.IBMPowerVSCluster.Status.Network.ID
		if networkID == "" {
			return nil, fmt.Errorf("network ID is not yet resolved in cluster status and was not specified on machine")
		}
		network.ID = networkID
	}

	networkID, err := s.getNetworkID(ctx, network)
	if err != nil {
		s.Recorder.Eventf(s.IBMPowerVSMachine, corev1.EventTypeWarning, "FailedRetrieveNetwork", "Failed network retrieval: %v", err)
		return nil, fmt.Errorf("error getting network ID: %w", err)
	}
	log.V(3).Info("Retrieved network id", "networkID", *networkID)

	// 8. Construct IBM Cloud SDK Payload
	procType := strings.ToLower(string(machineSpec.ProcessorType))

	payload := &models.PVMInstanceCreate{
		ServerName: ptr.To(s.IBMPowerVSMachine.Name),
		ImageID:    ptr.To(imageID),
		Memory:     ptr.To(float64(machineSpec.MemoryGiB)),
		Processors: ptr.To(processors),
		ProcType:   ptr.To(procType),
		SysType:    machineSpec.SystemType,
		UserData:   userData,
		Networks: []*models.PVMInstanceAddNetwork{
			{NetworkID: networkID},
		},
	}

	if machineSpec.SSHKey != "" {
		payload.KeyPairName = machineSpec.SSHKey
	}

	// 9. Execute Instance Creation
	log.Info("Triggering PowerVS instance creation", "machine", s.IBMPowerVSMachine.Name)
	if _, err := s.IBMPowerVSClient.CreateInstance(ctx, payload); err != nil {
		s.Recorder.Eventf(s.IBMPowerVSMachine, corev1.EventTypeWarning, "FailedCreateInstance", "Failed instance creation: %v", err)
		return nil, fmt.Errorf("failed to create PowerVS instance via SDK: %w", err)
	}

	s.Recorder.Eventf(s.IBMPowerVSMachine, corev1.EventTypeNormal, "SuccessfulCreateInstance", "Successfully triggered creation for Instance %q", s.IBMPowerVSMachine.Name)

	return nil, nil
}

// DeleteMachine deletes the power vs machine associated with machine instance id and service instance id.
func (s *MachineScope) DeleteMachine(ctx context.Context) error {
	if err := s.IBMPowerVSClient.DeleteInstance(ctx, s.IBMPowerVSMachine.Status.InstanceID); err != nil {
		s.Recorder.Eventf(s.IBMPowerVSMachine, corev1.EventTypeWarning, "FailedDeleteInstance", "Failed instance deletion - %v", err)
		return err
	}
	s.Recorder.Eventf(s.IBMPowerVSMachine, corev1.EventTypeNormal, "SuccessfulDeleteInstance", "Deleted Instance %q", s.IBMPowerVSMachine.Name)
	return nil
}

// DeleteMachineIgnition deletes the ignition data associated with the machine.
func (s *MachineScope) DeleteMachineIgnition(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)

	// 1. Guard check: If the machine isn't using Ignition, skip teardown.
	if !s.useIgnition() {
		log.V(3).Info("Machine is not using user data of type ignition")
		return nil
	}

	// 2. Fetch the bucket name strictly from the Status field
	bucket := s.IBMPowerVSCluster.Status.COSInstance.BucketName
	if bucket == "" {
		log.Info("COS bucket name is not populated in cluster status, skipping ignition deletion")
		return nil
	}

	// 3. Delete the bootstrap key
	key := s.bootstrapDataKey()

	if _, err := s.COSClient.DeleteObject(&s3.DeleteObjectInput{
		Bucket: ptr.To(bucket),
		Key:    ptr.To(key),
	}); err != nil {
		s.Recorder.Eventf(s.IBMPowerVSMachine, corev1.EventTypeWarning, "FailedDeleteMachineIgnition", "Failed machine ignition deletion - %v", err)
		return fmt.Errorf("failed to delete COS object %s from bucket %s: %w", key, bucket, err)
	}

	s.Recorder.Eventf(s.IBMPowerVSMachine, corev1.EventTypeNormal, "SuccessfulDeleteMachineIgnition", "Deleted machine ignition %q", s.IBMPowerVSMachine.Name)
	return nil
}

// CreateVPCLoadBalancerPoolMember registers the current machine's IP in every eligible load balancer
// pool. It is designed for parallel progress across multiple load balancers:
//
//   - A transiently-busy load balancer (update_pending, create_pending) is skipped — not an error —
//     so the other load balancer can still make progress in the same reconcile pass.
//   - A non-transient non-active state (e.g. delete_pending, or any unknown state) is returned as
//     an error so that a genuinely broken load balancer surfaces through conditions/events.
//   - At most one CreateLoadBalancerPoolMember (POST) is issued per load balancer per pass. After a
//     write the load balancer goes update_pending, so subsequent pools on the same load balancer are
//     checked for membership only; if any are still unregistered the pending flag is set.
//
// Returns (pendingUpdate, err):
//   - pendingUpdate: true when any registration is still incomplete.
//   - err: only set for real API failures or a non-transient load balancer state.
func (s *MachineScope) CreateVPCLoadBalancerPoolMember(ctx context.Context) (bool, error) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)
	loadBalancers := make([]infrav1.LoadBalancerSource, 0)
	if len(s.IBMPowerVSCluster.Spec.LoadBalancers) == 0 {
		loadBalancers = append(loadBalancers, infrav1.LoadBalancerSource{
			Type: infrav1.SourceTypeProvision,
			Provision: infrav1.LoadBalancerProvision{
				Name: ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeLBPublic, ""),
				Type: infrav1.LoadBalancerTypePublic,
			},
		})
	}
	for index, loadBalancer := range s.IBMPowerVSCluster.Spec.LoadBalancers {
		if loadBalancer.Type == infrav1.SourceTypeProvision && loadBalancer.Provision.Name == "" {
			lbType := ResourceTypeLBPublic
			if loadBalancer.Provision.Type == infrav1.LoadBalancerTypePrivate {
				lbType = ResourceTypeLBPrivate
			}
			qualifier := ""
			if index > 0 {
				qualifier = strconv.Itoa(index)
			}
			loadBalancer.Provision.Name = ResourceName(s.IBMPowerVSCluster.Name, lbType, qualifier)
		}
		loadBalancers = append(loadBalancers, loadBalancer)
	}

	pendingUpdate := false

	for _, lb := range loadBalancers {
		var lbName string
		switch lb.Type {
		case infrav1.SourceTypeReference:
			lbName = lb.Reference.Name
		case infrav1.SourceTypeProvision:
			lbName = lb.Provision.Name
		}
		if lbName == "" {
			return false, fmt.Errorf("failed to determine VPC load balancer name")
		}

		if len(s.IBMPowerVSCluster.Status.LoadBalancers) == 0 {
			return false, fmt.Errorf("failed to find VPC load balancer ID")
		}
		var lbID string
		found := false
		for _, val := range s.IBMPowerVSCluster.Status.LoadBalancers {
			if val.Name == lbName {
				lbID = val.ID
				found = true
				break
			}
		}
		if !found || lbID == "" {
			return false, fmt.Errorf("failed to find VPC load balancer ID")
		}

		loadBalancerDetails, _, err := s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{ID: &lbID})
		if err != nil {
			return false, fmt.Errorf("failed to find VPC load balancer details: %w", err)
		}

		// Only transient states are safe to skip and retry. A non-transient non-active state
		// (e.g. delete_pending, or any unknown state) means the load balancer is genuinely broken
		// and should surface as an error so conditions/events reflect it.
		// A create_pending load balancer may legitimately have zero pools, so the
		// provisioning status is checked before the pool count.
		if loadBalancerDetails.ProvisioningStatus == nil {
			return false, fmt.Errorf("VPC load balancer %s has no provisioning status", lbName)
		}

		loadBalancerState := infrav1.LoadBalancerState(*loadBalancerDetails.ProvisioningStatus)

		switch loadBalancerState {
		case infrav1.LoadBalancerStateActive:
			// Load balancer is active and ready; proceed normally.
			log.V(3).Info("VPC load balancer is active", "loadBalancerName", lbName, "loadBalancerState", loadBalancerState)
		case infrav1.LoadBalancerStateUpdatePending, infrav1.LoadBalancerStateCreatePending:
			log.V(3).Info("VPC load balancer is transiently busy, skipping this pass",
				"loadBalancerName", lbName, "loadBalancerState", loadBalancerState)
			pendingUpdate = true
			continue
		default:
			return false, fmt.Errorf("VPC load balancer %s is in non-recoverable state %q", lbName, loadBalancerState)
		}

		if len(loadBalancerDetails.Pools) == 0 {
			return false, fmt.Errorf("no pools exist for the VPC load balancer %s", lbName)
		}

		internalIP := s.GetMachineInternalIP()

		// lbAdditionalListeners: port-protocol key → spec entry (used to resolve target ports).
		lbAdditionalListeners := map[string]infrav1.AdditionalListener{}
		for _, additionalListener := range lb.Provision.AdditionalListeners {
			protocol := additionalListener.Protocol
			if protocol == "" {
				protocol = infrav1.LoadBalancerListenerProtocolTCP
			}
			lbAdditionalListeners[fmt.Sprintf("%d-%s", additionalListener.Port, protocol)] = additionalListener
		}

		// loadBalancerListeners: pool name → listener spec (populated from the cloud's listener list).
		loadBalancerListeners := map[string]infrav1.AdditionalListener{}
		for _, listenerRef := range loadBalancerDetails.Listeners {
			listenerOpts := &vpcv1.GetLoadBalancerListenerOptions{}
			listenerOpts.SetLoadBalancerID(*loadBalancerDetails.ID)
			listenerOpts.SetID(*listenerRef.ID)
			loadBalancerListener, _, err := s.IBMVPCClient.GetLoadBalancerListener(listenerOpts)
			if err != nil {
				return false, fmt.Errorf("failed to get load balancer listener %s: %w", *listenerRef.ID, err)
			}
			if additionalListener, ok := lbAdditionalListeners[fmt.Sprintf("%d-%s", *loadBalancerListener.Port, *loadBalancerListener.Protocol)]; ok {
				if loadBalancerListener.DefaultPool != nil {
					loadBalancerListeners[*loadBalancerListener.DefaultPool.Name] = additionalListener
				}
			} else if loadBalancerListener.Port != nil && *loadBalancerListener.Port == int64(6443) {
				protocol := infrav1.LoadBalancerListenerProtocol(*loadBalancerListener.Protocol)
				listener := infrav1.AdditionalListener{
					Port:     *loadBalancerListener.Port,
					Protocol: protocol,
				}
				if loadBalancerListener.DefaultPool != nil {
					loadBalancerListeners[*loadBalancerListener.DefaultPool.Name] = listener
				} else {
					log.V(3).Error(fmt.Errorf("unable to get the default pool details"), "default pool is nil", "port", loadBalancerListener.Port)
				}
			}
		}

		// updatedThisLB tracks whether we've already issued a POST to this load balancer in the
		// current pass. After one write the load balancer goes update_pending, so further writes
		// would fail. Instead we continue scanning remaining pools for membership to set pendingUpdate.
		updatedThisLB := false

		for _, pool := range loadBalancerDetails.Pools {
			log.V(3).Info("Checking LoadBalancer pool member", "pool", *pool.Name, "loadBalancerName", *loadBalancerDetails.Name, "IP", internalIP)

			// Determine the target port and evaluate label selectors for this pool.
			var targetPort int64
			if loadBalancerListener, ok := loadBalancerListeners[*pool.Name]; ok {
				targetPort = loadBalancerListener.Port
				log.V(3).Info("Checking if machine label matches with the label selector in listener", "machineLabel", s.IBMPowerVSMachine.Labels, "labelSelector", loadBalancerListener.Selector)
				selector, err := metav1.LabelSelectorAsSelector(&loadBalancerListener.Selector)
				if err != nil {
					log.V(5).Error(err, "Skipping listener addition, failed to get label selector from spec selector")
					continue
				}
				if selector.Empty() && !util.IsControlPlaneMachine(s.Machine) {
					log.V(3).Info("Skipping listener addition as the selector is empty and not a control plane machine")
					continue
				}
				if !selector.Empty() && !selector.Matches(labels.Set(s.IBMPowerVSMachine.Labels)) {
					log.V(3).Info("Skip adding listener, machine label doesn't match with the listener label selector", "pool", *pool.Name, "IP", internalIP)
					continue
				}
			}

			// If we already wrote to this load balancer this pass, it is now update_pending;
			// further writes would fail. Break out without listing remaining pools.
			if updatedThisLB {
				log.V(3).Info("Already updated load balancer this pass, skipping remaining pools",
					"pool", *pool.Name, "loadBalancerName", *loadBalancerDetails.Name)
				pendingUpdate = true
				break
			}

			// Check existing membership.
			listOpts := &vpcv1.ListLoadBalancerPoolMembersOptions{}
			listOpts.SetLoadBalancerID(*loadBalancerDetails.ID)
			listOpts.SetPoolID(*pool.ID)
			existingPoolMembers, _, err := s.IBMVPCClient.ListLoadBalancerPoolMembers(listOpts)
			if err != nil {
				return false, fmt.Errorf("failed to list %s VPC load balancer pool members: %w", *pool.Name, err)
			}

			alreadyRegistered := false
			for _, member := range existingPoolMembers.Members {
				if target, ok := member.Target.(*vpcv1.LoadBalancerPoolMemberTarget); ok && target.Address != nil {
					if *target.Address == internalIP {
						alreadyRegistered = true
						log.V(3).Info("Pool member already exists", "poolName", *pool.Name, "IP", internalIP)
						break
					}
				}
			}
			if alreadyRegistered {
				continue
			}

			// Machine is not yet in this pool — write the membership.
			opts := &vpcv1.CreateLoadBalancerPoolMemberOptions{}
			opts.SetPort(targetPort)
			opts.SetLoadBalancerID(*loadBalancerDetails.ID)
			opts.SetPoolID(*pool.ID)
			opts.SetTarget(&vpcv1.LoadBalancerPoolMemberTargetPrototype{Address: &internalIP})
			log.V(3).Info("Creating VPC load balancer pool member", "pool", *pool.Name, "IP", internalIP, "port", targetPort)
			newMember, _, err := s.IBMVPCClient.CreateLoadBalancerPoolMember(opts)
			if err != nil {
				return false, fmt.Errorf("failed to create VPC load balancer %s pool member: %w", *loadBalancerDetails.Name, err)
			}
			log.Info("Created VPC load balancer pool member", "memberID", *newMember.ID, "pool", *pool.Name)

			updatedThisLB = true
			// If the member is not yet active, the load balancer is settling; signal a requeue.
			if newMember.ProvisioningStatus == nil || *newMember.ProvisioningStatus != string(infrav1.LoadBalancerStateActive) {
				pendingUpdate = true
			}
		}
	}
	return pendingUpdate, nil
}

// SetReady will set the status as ready for the machine.
func (s *MachineScope) SetReady() {
	s.IBMPowerVSMachine.Status.Initialization.Provisioned = ptr.To(true)
}

// SetNotReady will set status as not ready for the machine.
func (s *MachineScope) SetNotReady() {
	s.IBMPowerVSMachine.Status.Initialization.Provisioned = ptr.To(false)
}

// IsReady will return the status for the machine.
func (s *MachineScope) IsReady() bool {
	return ptr.Deref(s.IBMPowerVSMachine.Status.Initialization.Provisioned, false)
}

// SetInstanceID will set the instance id for the machine.
func (s *MachineScope) SetInstanceID(id *string) {
	if id != nil {
		s.IBMPowerVSMachine.Status.InstanceID = *id
	}
}

// GetInstanceID will get the instance id for the machine.
func (s *MachineScope) GetInstanceID() string {
	return s.IBMPowerVSMachine.Status.InstanceID
}

// SetInstanceState will set the state for the machine.
func (s *MachineScope) SetInstanceState(status *string) {
	s.IBMPowerVSMachine.Status.InstanceState = infrav1.PowerVSInstanceState(*status)
}

// GetInstanceState will get the state for the machine.
func (s *MachineScope) GetInstanceState() infrav1.PowerVSInstanceState {
	return s.IBMPowerVSMachine.Status.InstanceState
}

// SetHealth will set the health status for the machine.
func (s *MachineScope) SetHealth(health *models.PVMInstanceHealth) {
	if health != nil {
		s.IBMPowerVSMachine.Status.Health = health.Status
	}
}

// SetAddresses will set the addresses for the machine.
func (s *MachineScope) SetAddresses(ctx context.Context, instance *models.PVMInstance) {
	log := ctrl.LoggerFrom(ctx)

	// 1. Set Base Addresses (DNS and Hostname)
	addresses := []clusterv1.MachineAddress{
		{Type: clusterv1.MachineInternalDNS, Address: *instance.ServerName},
		{Type: clusterv1.MachineHostName, Address: *instance.ServerName},
	}

	// 2. Check if the instance already has IPs assigned
	if instanceIPs := extractIPsFromInstance(instance); len(instanceIPs) > 0 {
		s.IBMPowerVSMachine.Status.Addresses = append(addresses, instanceIPs...)
		return
	}

	// 3. Look for DHCP IP in our local cache
	if ip, found := s.getIPFromCache(ctx, *instance.ServerName); found {
		addresses = append(addresses, clusterv1.MachineAddress{Type: clusterv1.MachineInternalIP, Address: ip})
		s.IBMPowerVSMachine.Status.Addresses = addresses
		return
	}

	// 4. Query the IBM Cloud DHCP Server directly; DHCP errors are non-fatal —
	// the machine may not have received a lease yet, so we settle for base addresses.
	ip, err := s.getIPFromDHCPServer(ctx, instance)
	if err != nil {
		log.Info("Could not retrieve IP from DHCP server, using base addresses only", "reason", err.Error())
		s.IBMPowerVSMachine.Status.Addresses = addresses
		return
	}

	// 5. Update Cache and Status on success
	if err := s.DHCPIPCacheStore.Add(powervs.VMip{Name: *instance.ServerName, IP: ip}); err != nil {
		log.Error(err, "failed to update the DHCP cache store with the IP", "IP", ip)
	}

	log.V(3).Info("Successfully resolved internal IP for VM", "IP", ip)
	addresses = append(addresses, clusterv1.MachineAddress{Type: clusterv1.MachineInternalIP, Address: ip})
	s.IBMPowerVSMachine.Status.Addresses = addresses
}

// SetRegion will set the region for the machine.
func (s *MachineScope) SetRegion(region string) {
	s.IBMPowerVSMachine.Status.Region = region
}

// GetRegion will get the region for the machine.
func (s *MachineScope) GetRegion() string {
	return s.IBMPowerVSMachine.Status.Region
}

// SetZone will set the zone for the machine.
func (s *MachineScope) SetZone(zone string) {
	s.IBMPowerVSMachine.Status.Zone = zone
}

// GetZone will get the zone for the machine.
func (s *MachineScope) GetZone() string {
	return s.IBMPowerVSMachine.Status.Zone
}

// GetWorkspaceID returns the PowerVS workspace ID, evaluating in the following order of precedence:
// 1. Machine Spec explicitly sets Workspace ID
// 2. Machine Spec explicitly sets Workspace Name (requires IBM Cloud lookup)
// 3. Inherit resolved Workspace ID from the Cluster Status.
func (s *MachineScope) GetWorkspaceID() (string, error) {
	// 1. Precedence 1: Machine Spec Workspace ID
	if s.IBMPowerVSMachine.Spec.Workspace.ID != "" {
		return s.IBMPowerVSMachine.Spec.Workspace.ID, nil
	}

	// 2. Precedence 2: Machine Spec Workspace Name (requires lookup)
	if s.IBMPowerVSMachine.Spec.Workspace.Name != "" {
		workspaceName := s.IBMPowerVSMachine.Spec.Workspace.Name

		resourceInstance := resourcecontroller.InstanceFilter{
			Name:           workspaceName,
			Zone:           &s.IBMPowerVSCluster.Spec.Zone,
			ResourceID:     resourcecontroller.PowerVSResourceID,
			ResourcePlanID: resourcecontroller.PowerVSResourcePlanID,
		}

		workspace, err := s.ResourceClient.GetResourceInstanceByFilter(resourceInstance)
		if err != nil {
			return "", fmt.Errorf("failed to lookup PowerVS workspace by name %q: %w", workspaceName, err)
		}
		if workspace == nil || workspace.GUID == nil {
			return "", fmt.Errorf("PowerVS workspace %q not found or GUID is nil in IBM Cloud", workspaceName)
		}

		return *workspace.GUID, nil
	}

	// 3. Precedence 3: Inherit from Cluster Status
	// In v1beta3, the Cluster controller guarantees this is populated during the cluster's reconciliation loop.
	if s.IBMPowerVSCluster.Status.Workspace.ID != "" {
		return s.IBMPowerVSCluster.Status.Workspace.ID, nil
	}

	return "", errors.New("failed to find workspace ID: not specified in Machine spec and not yet populated in Cluster status")
}

// SetProviderID will set the provider id for the machine.
func (s *MachineScope) SetProviderID(instanceID string) error {
	workspaceID, err := s.GetWorkspaceID()
	if err != nil {
		return err
	}
	s.IBMPowerVSMachine.Spec.ProviderID = fmt.Sprintf("ibmpowervs://%s/%s/%s/%s", s.GetRegion(), s.GetZone(), workspaceID, instanceID)
	return nil
}

// GetMachineInternalIP returns the machine's internal IP.
func (s *MachineScope) GetMachineInternalIP() string {
	for _, address := range s.IBMPowerVSMachine.Status.Addresses {
		if address.Type == clusterv1.MachineInternalIP {
			return address.Address
		}
	}
	return ""
}

// resolveUserData fetches raw bootstrap data and, when Ignition is configured,
// uploads it to COS and returns a base64-encoded ignition redirect document.
// For plain cloud-init, it base64-encodes the raw secret value directly.
func (s *MachineScope) resolveUserData(ctx context.Context) (string, error) {
	userData, err := s.getRawBootstrapData()
	if err != nil {
		return "", err
	}

	if s.useIgnition() {
		data, err := s.ignitionUserData(ctx, userData)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(data), nil
	}

	// Explicitly return nil for the error since it is guaranteed to be nil here
	return base64.StdEncoding.EncodeToString(userData), nil
}

// ignitionUserData uploads the raw bootstrap data to COS via createIgnitionData,
// then wraps the resulting pre-signed URL in an Ignition v2 or v3 redirect document.
func (s *MachineScope) ignitionUserData(ctx context.Context, userData []byte) ([]byte, error) {
	presignedURL, err := s.createIgnitionData(ctx, userData)
	if err != nil {
		return nil, fmt.Errorf("failed to create user data object: %w", err)
	}

	ignVersion := s.getIgnitionVersion()
	semver, err := semver.ParseTolerant(ignVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ignition version %q: %w", ignVersion, err)
	}

	switch semver.Major {
	case 2:
		ignData := &ignV2Types.Config{
			Ignition: ignV2Types.Ignition{
				Version: semver.String(),
				Config: ignV2Types.IgnitionConfig{
					Replace: &ignV2Types.ConfigReference{
						Source: presignedURL,
					},
				},
			},
		}
		return json.Marshal(ignData)
	case 3:
		ignData := &ignV3Types.Config{
			Ignition: ignV3Types.Ignition{
				Version: semver.String(),
				Config: ignV3Types.IgnitionConfig{
					Replace: ignV3Types.Resource{
						Source: aws.String(presignedURL),
					},
				},
			},
		}
		return json.Marshal(ignData)
	default:
		return nil, fmt.Errorf("unsupported ignition version %q", ignVersion)
	}
}

// createIgnitionData uploads userData to the COS bucket, then returns a
// pre-signed HTTPS URL that Ignition uses to fetch the bootstrap config.
// The URL is signed with HMAC SigV4 credentials — the signature is embedded as query parameters.
func (s *MachineScope) createIgnitionData(ctx context.Context, data []byte) (string, error) {
	log := ctrl.LoggerFrom(ctx)
	if len(data) == 0 {
		return "", fmt.Errorf("user data is empty")
	}

	key := s.bootstrapDataKey()
	log.V(3).Info("Bootstrap data key name", "key", key)

	bucket := s.IBMPowerVSCluster.Status.COSInstance.BucketName
	if bucket == "" {
		return "", fmt.Errorf("COS bucket name is not yet populated in cluster status. Waiting for cluster reconciler")
	}

	if _, err := s.COSClient.PutObject(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return "", fmt.Errorf("failed to push object to COS bucket: %w", err)
	}

	presignedURL, err := s.COSClient.PresignedURL(bucket, key, presignExpiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL for ignition data: %w", err)
	}

	log.V(3).Info("Generated pre-signed Ignition URL", "bucket", bucket, "key", key)
	return presignedURL, nil
}

// hmacCOSClientOptions builds the HMACCOSClientOptions for the current machine scope,
// drawing region, secret name, namespace, and endpoint overrides from cluster status.
func (s *MachineScope) hmacCOSClientOptions() HMACCOSClientOptions {
	cosStatus := s.IBMPowerVSCluster.Status.COSInstance
	return HMACCOSClientOptions{
		Client:              s.Client,
		HMACSecretName:      cosStatus.HMACSecretName,
		HMACSecretNamespace: s.IBMPowerVSCluster.Namespace,
		BucketRegion:        cosStatus.BucketRegion,
		ServiceEndpoint:     s.ServiceEndpoint,
	}
}

// bootstrapDataKey returns the COS object key for this machine's bootstrap data.
func (s *MachineScope) bootstrapDataKey() string {
	// Use machine name as object key.
	return path.Join(s.role(), s.name())
}

// getIgnitionVersion returns the user-specified Ignition version,
// or falls back to the default "2.3" if it was left unset.
func (s *MachineScope) getIgnitionVersion() string {
	if s.IBMPowerVSCluster.Spec.Ignition.Version == "" {
		return "2.3"
	}
	return s.IBMPowerVSCluster.Spec.Ignition.Version
}

// useIgnition returns true if the user configured a COS Instance,
// which acts as the master switch for the Ignition bootstrap workflow.
func (s *MachineScope) useIgnition() bool {
	return s.IBMPowerVSCluster.Spec.COSInstance.Type != ""
}

// getRawBootstrapData returns the bootstrap data if present.
func (s *MachineScope) getRawBootstrapData() ([]byte, error) {
	if s.Machine == nil || s.Machine.Spec.Bootstrap.DataSecretName == nil {
		return nil, errors.New("failed to retrieve bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: s.Machine.Namespace, Name: *s.Machine.Spec.Bootstrap.DataSecretName}
	if err := s.Client.Get(context.TODO(), key, secret); err != nil {
		return nil, fmt.Errorf("failed to retrieve bootstrap data secret: %v", err)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, errors.New("failed to retrieve bootstrap data: secret value key is missing")
	}

	return value, nil
}

// getImageID resolves an image ResourceIdentifier to a concrete image ID string.
func (s *MachineScope) getImageID(ctx context.Context, image infrav1.ResourceIdentifier) (string, error) {
	if image.ID != "" {
		return image.ID, nil
	}

	if image.Name != "" {
		images, err := s.getImages(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get images from IBM Cloud: %w", err)
		}

		for _, img := range images.Images {
			if image.Name == *img.Name {
				return *img.ImageID, nil
			}
		}

		return "", fmt.Errorf("image with name %q not found", image.Name)
	}

	return "", fmt.Errorf("image reference must contain either an ID or a Name")
}

// getNetworkID resolves a network ResourceIdentifier to a concrete network ID pointer.
func (s *MachineScope) getNetworkID(ctx context.Context, network infrav1.ResourceIdentifier) (*string, error) {
	if network.ID != "" {
		return ptr.To(network.ID), nil
	}

	if network.Name != "" {
		net, err := s.IBMPowerVSClient.GetNetworkByName(ctx, network.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get network by name %q: %w", network.Name, err)
		}
		if net == nil || net.NetworkID == nil {
			return nil, fmt.Errorf("network with name %q not found", network.Name)
		}
		return net.NetworkID, nil
	}
	return nil, fmt.Errorf("network identifier must contain either an ID or a Name")
}

// ensureInstanceUnique returns the existing PVMInstanceReference if an instance
// with the given name already exists, or nil if no such instance is found.
func (s *MachineScope) ensureInstanceUnique(ctx context.Context, instanceName string) (*models.PVMInstanceReference, error) {
	instances, err := s.IBMPowerVSClient.ListInstances(ctx)
	if err != nil {
		return nil, err
	}
	for _, ins := range instances.PvmInstances {
		if *ins.ServerName == instanceName {
			return ins, nil
		}
	}
	return nil, nil
}

// validateSystemType checks whether the machine's configured SystemType is
// supported by the target datacenter zone, using a TTL-backed in-memory cache
// to avoid redundant API calls.
func (s *MachineScope) validateSystemType(ctx context.Context) (bool, []string, error) {
	systemType := s.IBMPowerVSMachine.Spec.SystemType

	if systemType == "" {
		return false, nil, fmt.Errorf("systemType is not set")
	}

	zone := s.GetZone()

	// Read from Cache
	sysCache.mu.RLock()
	entry, exists := sysCache.zonesMap[zone]
	isFresh := time.Since(entry.lastFetch) < sysCache.ttl
	sysCache.mu.RUnlock()

	if exists && isFresh {
		return slices.Contains(entry.supportedTypes, systemType), entry.supportedTypes, nil
	}

	// Cache is expired or empty. Fetch from IBM Cloud.
	sysCache.mu.Lock()
	defer sysCache.mu.Unlock()

	// Double check inside the lock
	entry, exists = sysCache.zonesMap[zone]
	if exists && time.Since(entry.lastFetch) < sysCache.ttl {
		return slices.Contains(entry.supportedTypes, systemType), entry.supportedTypes, nil
	}

	// Fetch the specific datacenter capabilities.
	datacenter, err := s.IBMPowerVSClient.GetDatacenterDetails(ctx, zone)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get datacenter details for zone %s: %w", zone, err)
	}

	if datacenter == nil || datacenter.CapabilitiesDetails == nil || datacenter.CapabilitiesDetails.SupportedSystems == nil {
		return false, nil, fmt.Errorf("system capabilities details are missing for zone %s", zone)
	}

	// Extract the General list.
	systemTypes := datacenter.CapabilitiesDetails.SupportedSystems.General

	if len(systemTypes) == 0 {
		return false, nil, fmt.Errorf("no general system types available in zone %s", zone)
	}

	// Sort once so error messages are always in alphabetical order.
	sort.Strings(systemTypes)

	// Update the cache for the zone
	sysCache.zonesMap[zone] = zoneCacheEntry{
		supportedTypes: systemTypes,
		lastFetch:      time.Now(),
	}

	// Validate against the newly refreshed data.
	return slices.Contains(systemTypes, systemType), systemTypes, nil
}

// role returns the machine role label ("control-plane" or "node").
func (s *MachineScope) role() string {
	if util.IsControlPlaneMachine(s.Machine) {
		return "control-plane"
	}
	return "node"
}

// name returns the IBMPowerVSMachine name.
func (s *MachineScope) name() string {
	return s.IBMPowerVSMachine.Name
}

// getImages will get list of images for the powervs service instance.
func (s *MachineScope) getImages(ctx context.Context) (*models.Images, error) {
	return s.IBMPowerVSClient.ListImages(ctx)
}

// extractIPsFromInstance parses the raw instance data for assigned IP addresses.
func extractIPsFromInstance(instance *models.PVMInstance) []clusterv1.MachineAddress {
	var ips []clusterv1.MachineAddress
	for _, network := range instance.Networks {
		if ip := strings.TrimSpace(network.IPAddress); ip != "" {
			ips = append(ips, clusterv1.MachineAddress{Type: clusterv1.MachineInternalIP, Address: ip})
		}
		if extIP := strings.TrimSpace(network.ExternalIP); extIP != "" {
			ips = append(ips, clusterv1.MachineAddress{Type: clusterv1.MachineExternalIP, Address: extIP})
		}
	}
	return ips
}

// getIPFromCache attempts to retrieve the IP from the fast local cache.
func (s *MachineScope) getIPFromCache(ctx context.Context, serverName string) (string, bool) {
	log := ctrl.LoggerFrom(ctx)

	obj, exists, err := s.DHCPIPCacheStore.GetByKey(serverName)
	if err != nil {
		log.Error(err, "failed to fetch the DHCP IP address from cache store")
		return "", false
	}

	if exists {
		vmIP := obj.(powervs.VMip)
		log.V(3).Info("Found IP for machine from DHCP cache", "IP", vmIP.IP)
		return vmIP.IP, true
	}

	return "", false
}

// getIPFromDHCPServer performs the complex lookup of Network -> DHCP Server -> Lease MAC -> IP.
func (s *MachineScope) getIPFromDHCPServer(ctx context.Context, instance *models.PVMInstance) (string, error) {
	// 1. Resolve Network ID
	network := s.IBMPowerVSMachine.Spec.Network
	if network.ID == "" && network.Name == "" {
		network.ID = s.IBMPowerVSCluster.Status.Network.ID
	}

	networkID, err := s.getNetworkID(ctx, network)
	if err != nil {
		return "", fmt.Errorf("failed to fetch network ID: %w", err)
	}

	// 2. Find the exact network attachment to get the MAC address
	var pvmNetwork *models.PVMInstanceNetwork
	for _, net := range instance.Networks {
		if net.NetworkID == *networkID {
			pvmNetwork = net
			break
		}
	}
	if pvmNetwork == nil {
		return "", fmt.Errorf("failed to find network attached to machine with ID %s", *networkID)
	}

	// 3. Find the DHCP server associated with this network
	dhcpServers, err := s.IBMPowerVSClient.ListDHCPServers(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list DHCP servers: %w", err)
	}

	var dhcpServerID string
	for _, server := range dhcpServers {
		if server.Network != nil && server.Network.ID != nil && *server.Network.ID == *networkID {
			dhcpServerID = *server.ID
			break
		}
	}
	if dhcpServerID == "" {
		return "", fmt.Errorf("no DHCP server found associated with network %s", *networkID)
	}

	// 4. Inspect the DHCP server's active leases for our MAC address
	dhcpServerDetails, err := s.IBMPowerVSClient.GetDHCPServer(ctx, dhcpServerID)
	if err != nil {
		return "", fmt.Errorf("failed to get DHCP server details for %s: %w", dhcpServerID, err)
	}

	for _, lease := range dhcpServerDetails.Leases {
		if *lease.InstanceMacAddress == pvmNetwork.MacAddress {
			return *lease.InstanceIP, nil
		}
	}

	return "", fmt.Errorf("DHCP lease not found for machine with MAC %s in DHCP Server %s", pvmNetwork.MacAddress, dhcpServerID)
}
