package clusterapi

import (
	"context"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/cluster/tfvars"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
)

// Provider is the base interface that cloud platforms
// should implement for the CAPI infrastructure provider.
type Provider interface {
	// Name provides the name for the cloud platform.
	Name() string

	// BootstrapHasPublicIP indicates whether a public IP address
	// is expected on the bootstrap node in a public cluster.
	// When BootstrapHasPublicIP returns true, the machine ready checks
	// wait for an ExternalIP address to be populated in the machine status.
	BootstrapHasPublicIP() bool
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

// IgnitionProvider handles preconditions for bootstrap ignition and
// generates ignition data for the CAPI bootstrap ignition secret.
//
// WARNING! Low-level primitive. Use only if absolutely necessary.
type IgnitionProvider interface {
	Ignition(ctx context.Context, in IgnitionInput) ([]byte, error)
}

// IgnitionInput collects the args passed to the IgnitionProvider call.
type IgnitionInput struct {
	Client           client.Client
	BootstrapIgnData []byte
	InfraID          string
	InstallConfig    *installconfig.InstallConfig
	TFVarsAsset      *tfvars.TerraformVariables
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
