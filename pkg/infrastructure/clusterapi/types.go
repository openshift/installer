package clusterapi

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/cluster/tfvars"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
)

// Provider is the base interface that cloud platforms
// should implement for the CAPI infrastructure provider.
type Provider interface {
	// Name provides the name for the cloud platform.
	Name() string

	// PublicGatherEndpoint returns how the cloud platform expects the installer
	// to connect to the bootstrap node for log gathering. CAPI providers are not
	// consistent in how public IP addresses are represented in the machine status.
	// Furthermore, Azure cannot attach a public IP to the bootstrap node, so SSH
	// must be performed through the API load balancer.
	// When a platform returns ExternalIP, the installer will require an ExternalIP
	// to be present in the status, before it declares the machine ready.
	PublicGatherEndpoint() GatherEndpoint
}

// PreProvider defines the PreProvision hook, which is called prior to
// CAPI infrastructure provisioning.
type PreProvider interface {
	// PreProvision is called before provisioning using CAPI controllers has begun
	// and should be used to create dependencies needed for CAPI provisioning,
	// such as IAM roles or policies.
	PreProvision(ctx context.Context, in PreProvisionInput) error
}

// PreProvisionInput collects the args passed to the PreProvision call.
type PreProvisionInput struct {
	InfraID          string
	InstallConfig    *installconfig.InstallConfig
	RhcosImage       *rhcos.Image
	ManifestsAsset   *manifests.Manifests
	MachineManifests []client.Object
	WorkersAsset     *machines.Worker
}

// IgnitionProvider handles preconditions for bootstrap ignition,
// such as pushing to cloud storage. Returns bootstrap and master
// ignition secrets.
//
// WARNING! Low-level primitive. Use only if absolutely necessary.
type IgnitionProvider interface {
	Ignition(ctx context.Context, in IgnitionInput) ([]*corev1.Secret, error)
}

// IgnitionInput collects the args passed to the IgnitionProvider call.
type IgnitionInput struct {
	Client           client.Client
	BootstrapIgnData []byte
	MasterIgnData    []byte
	WorkerIgnData    []byte
	InfraID          string
	InstallConfig    *installconfig.InstallConfig
	TFVarsAsset      *tfvars.TerraformVariables
	RootCA           *tls.RootCA
}

// SetIgnFromOutput set the ignition data for input from an output.
// This is used to chain multiple ignition editions.
func (ignIn *IgnitionInput) SetIgnFromOutput(outIgn *IgnitionOutput) {
	if ignIn == nil || outIgn == nil {
		return
	}
	ignIn.BootstrapIgnData = outIgn.UpdatedBootstrapIgn
	ignIn.MasterIgnData = outIgn.UpdatedMasterIgn
	ignIn.WorkerIgnData = outIgn.UpdatedWorkerIgn
}

// IgnitionOutput collects updated Ignition Data for Bootstrap, Master and Worker nodes.
type IgnitionOutput struct {
	UpdatedBootstrapIgn []byte
	UpdatedMasterIgn    []byte
	UpdatedWorkerIgn    []byte
}

// Set set the ignition data from another output.
// This is used for chained multiple ignition editions.
func (ignOut *IgnitionOutput) Set(other *IgnitionOutput) {
	if ignOut == nil || other == nil {
		return
	}
	ignOut.UpdatedBootstrapIgn = other.UpdatedBootstrapIgn
	ignOut.UpdatedMasterIgn = other.UpdatedMasterIgn
	ignOut.UpdatedWorkerIgn = other.UpdatedWorkerIgn
}

// InfraReadyProvider defines the InfraReady hook, which is
// called after the initial infrastructure manifests have been created
// and InfrastructureReady == true on the cluster status, and before
// IgnitionProvider hook and creation of the control-plane machines.
type InfraReadyProvider interface {
	// InfraReady is called once cluster.Status.InfrastructureReady
	// is true, typically after load balancers have been provisioned. It can be used
	// to create DNS records.
	InfraReady(ctx context.Context, in InfraReadyInput) error
}

// InfraReadyInput collects the args passed to the InfraReady call.
type InfraReadyInput struct {
	// Client is the client for kube-apiserver running locally on the installer host.
	// It can be used to read the status of the cluster object on the local control plane.
	Client        client.Client
	InstallConfig *installconfig.InstallConfig
	InfraID       string
}

// PostProvider defines the PostProvision hook, which is called after
// machine provisioning has completed.
type PostProvider interface {
	PostProvision(ctx context.Context, in PostProvisionInput) error
}

// PostProvisionInput collects the args passed to the PostProvision hook.
type PostProvisionInput struct {
	Client        client.Client
	InstallConfig *installconfig.InstallConfig
	InfraID       string
}

// BootstrapDestroyer allows platform-specific behavior when
// destroying bootstrap resources.
type BootstrapDestroyer interface {
	DestroyBootstrap(ctx context.Context, in BootstrapDestroyInput) error
}

// BootstrapDestroyInput collects args passed to the DestroyBootstrap hook.
type BootstrapDestroyInput struct {
	Client   client.Client
	Metadata types.ClusterMetadata
}

// PostDestroyer allows platform-specific behavior after bootstrap has been destroyed and
// ClusterAPI has stopped running.
type PostDestroyer interface {
	PostDestroy(ctx context.Context, in PostDestroyerInput) error
}

// PostDestroyerInput collects args passed to the PostDestroyer hook.
type PostDestroyerInput struct {
	Metadata types.ClusterMetadata
}

// Timeouts allows platform provider to override the timeouts for certain phases.
type Timeouts interface {
	// When waiting for the network infrastructure to become ready.
	NetworkTimeout() time.Duration
	// When waiting for the machines to provision.
	ProvisionTimeout() time.Duration
}

// GatherEndpoint represents the valid values for connecting to the bootstrap nude
// in a public cluster to gather logs.
type GatherEndpoint string

const (
	// ExternalIP indicates that the machine status will include an ExternalIP that can be used for gather.
	ExternalIP GatherEndpoint = "ExternalIP"

	// InternalIP indicates that the machine status will only include InternalIPs.
	InternalIP GatherEndpoint = "InternalIP"

	// APILoadBalancer indicates that gather bootstrap should connect to the API load balancer.
	APILoadBalancer GatherEndpoint = "APILoadBalancer"
)
