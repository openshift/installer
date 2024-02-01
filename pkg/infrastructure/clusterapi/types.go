package clusterapi

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/installconfig"
)

// Provider is the base interface that cloud platforms
// should implement for the CAPI infrastructure provider.
type Provider interface {
	// Name provides the name for the cloud platform.
	Name() string
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
	InfraID       string
	InstallConfig *installconfig.InstallConfig
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
	BootstrapIgnData []byte
	InfraID          string
	InstallConfig    *installconfig.InstallConfig
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
