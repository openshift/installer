package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	baremetalapi "github.com/openshift/cluster-api-provider-baremetal/pkg/apis"
	baremetalprovider "github.com/openshift/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	libvirtapi "github.com/openshift/cluster-api-provider-libvirt/pkg/apis"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	ovirtproviderapi "github.com/openshift/cluster-api-provider-ovirt/pkg/apis"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/baremetal"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	ibmcloudapi "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis"
	ibmcloudprovider "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
)

// Arbiter generates the machines for the `arbiter` machine pool.
type Arbiter struct {
	UserDataFile       *asset.File
	MachineConfigFiles []*asset.File
	MachineFiles       []*asset.File
	IPClaimFiles       []*asset.File
	IPAddrFiles        []*asset.File

	// SecretFiles is used by the baremetal platform to register the
	// credential information for communicating with management
	// controllers on hosts.
	SecretFiles []*asset.File

	// NetworkConfigSecretFiles is used by the baremetal platform to
	// store the networking configuration per host
	NetworkConfigSecretFiles []*asset.File

	// HostFiles is the list of baremetal hosts provided in the
	// installer configuration.
	HostFiles []*asset.File
}

const (

	// arbiterMachineFileName is the format string for constucting the
	// arbiter Machine filenames.
	arbiterMachineFileName = "99_openshift-cluster-api_arbiter-machines-%s.yaml"

	// arbiterUserDataFileName is the filename used for the arbiter
	// user-data secret.
	arbiterUserDataFileName = "99_openshift-cluster-api_arbiter-user-data-secret.yaml"

	arbiterHostFileName                = "99_openshift-cluster-api_arbiter_hosts-%s.yaml"
	arbiterSecretFileName              = "99_openshift-cluster-api_arbiter_host-bmc-secrets-%s.yaml"            // #nosec G101
	arbiterNetworkConfigSecretFileName = "99_openshift-cluster-api_arbiter_host-network-config-secrets-%s.yaml" // #nosec G101
)

var (
	arbiterMachineFileNamePattern   = fmt.Sprintf(arbiterMachineFileName, "*")
	arbiterIPClaimFileNamePattern   = fmt.Sprintf(ipClaimFileName, "*arbiter*")
	arbiterIPAddressFileNamePattern = fmt.Sprintf(ipAddressFileName, "*arbiter*")

	_ asset.WritableAsset = (*Arbiter)(nil)
)

// Name returns a human friendly name for the Arbiter Asset.
func (m *Arbiter) Name() string {
	return "Arbiter Machines"
}

// Dependencies returns all of the dependencies directly needed by the
// Arbiter asset.
func (m *Arbiter) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		&machine.Arbiter{},
	}
}

// Generate generates the Arbiter asset.
//
//nolint:gocyclo
func (m *Arbiter) Generate(ctx context.Context, dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	mign := &machine.Arbiter{}
	dependencies.Get(clusterID, installConfig, rhcosImage, mign)

	ic := installConfig.Config

	if ic.Arbiter == nil {
		return nil
	}
	if ic.Platform.Name() != baremetaltypes.Name {
		return fmt.Errorf("only BareMetal platform is supported for Arbiter deployments")
	}

	pool := *ic.Arbiter
	var err error
	machines := []machinev1beta1.Machine{}
	var ipClaims []ipamv1.IPAddressClaim
	var ipAddrs []ipamv1.IPAddress

	mpool := defaultBareMetalMachinePoolPlatform()
	mpool.Set(ic.Platform.BareMetal.DefaultMachinePlatform)
	mpool.Set(pool.Platform.BareMetal)
	pool.Platform.BareMetal = &mpool

	// Use managed user data secret, since we always have up to date images
	// available in the cluster
	arbiterUserDataSecretName := "arbiter-user-data-managed" // #nosec G101
	enabledCaps := installConfig.Config.GetEnabledCapabilities()
	if enabledCaps.Has(configv1.ClusterVersionCapabilityMachineAPI) {
		machines, err = baremetal.Machines(clusterID.InfraID, ic, &pool, "arbiter", arbiterUserDataSecretName)
		if err != nil {
			return fmt.Errorf("failed to create arbiter machine objects: %w", err)
		}

		hostSettings, err := baremetal.ArbiterHosts(ic, machines, arbiterUserDataSecretName)
		if err != nil {
			return fmt.Errorf("failed to assemble host data: %w", err)
		}

		hosts, err := createHostAssetFiles(hostSettings.Hosts, arbiterHostFileName)
		if err != nil {
			return err
		}
		m.HostFiles = append(m.HostFiles, hosts...)

		secrets, err := createSecretAssetFiles(hostSettings.Secrets, arbiterSecretFileName)
		if err != nil {
			return err
		}
		m.SecretFiles = append(m.SecretFiles, secrets...)

		networkSecrets, err := createSecretAssetFiles(hostSettings.NetworkConfigSecrets, arbiterNetworkConfigSecretFileName)
		if err != nil {
			return err
		}
		m.NetworkConfigSecretFiles = append(m.NetworkConfigSecretFiles, networkSecrets...)
	}

	data, err := UserDataSecret(arbiterUserDataSecretName, mign.File.Data)
	if err != nil {
		return fmt.Errorf("failed to create user-data secret for arbiter machines: %w", err)
	}

	m.UserDataFile = &asset.File{
		Filename: filepath.Join(directory, arbiterUserDataFileName),
		Data:     data,
	}

	machineConfigs := []*mcfgv1.MachineConfig{}
	if pool.Hyperthreading == types.HyperthreadingDisabled {
		ignHT, err := machineconfig.ForHyperthreadingDisabled("arbiter")
		if err != nil {
			return fmt.Errorf("failed to create ignition for hyperthreading disabled for arbiter machines: %w", err)
		}
		machineConfigs = append(machineConfigs, ignHT)
	}
	if ic.SSHKey != "" {
		ignSSH, err := machineconfig.ForAuthorizedKeys(ic.SSHKey, "arbiter")
		if err != nil {
			return fmt.Errorf("failed to create ignition for authorized SSH keys for arbiter machines: %w", err)
		}
		machineConfigs = append(machineConfigs, ignSSH)
	}
	if ic.FIPS {
		ignFIPS, err := machineconfig.ForFIPSEnabled("arbiter")
		if err != nil {
			return fmt.Errorf("failed to create ignition for FIPS enabled for arbiter machines: %w", err)
		}
		machineConfigs = append(machineConfigs, ignFIPS)
	}

	// The maximum number of networks supported on ServiceNetwork is two, one IPv4 and one IPv6 network.
	// The cluster-network-operator handles the validation of this field.
	// Reference: https://github.com/openshift/cluster-network-operator/blob/fc3e0e25b4cfa43e14122bdcdd6d7f2585017d75/pkg/network/cluster_config.go#L45-L52
	if ic.Networking != nil && len(ic.Networking.ServiceNetwork) == 2 {
		// Only configure kernel args for dual-stack clusters.
		ignIPv6, err := machineconfig.ForDualStackAddresses("arbiter")
		if err != nil {
			return fmt.Errorf("failed to create ignition to configure IPv6 for arbiter machines: %w", err)
		}
		machineConfigs = append(machineConfigs, ignIPv6)
	}

	m.MachineConfigFiles, err = machineconfig.Manifests(machineConfigs, "arbiter", directory)
	if err != nil {
		return fmt.Errorf("failed to create MachineConfig manifests for arbiter machines: %w", err)
	}

	m.MachineFiles = make([]*asset.File, len(machines))

	m.IPClaimFiles = make([]*asset.File, len(ipClaims))
	for i, claim := range ipClaims {
		data, err := yaml.Marshal(claim)
		if err != nil {
			return fmt.Errorf("unable to marshal ip claim %v: %w", claim.Name, err)
		}

		m.IPClaimFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(ipClaimFileName, claim.Name)),
			Data:     data,
		}
	}

	m.IPAddrFiles = make([]*asset.File, len(ipAddrs))
	for i, address := range ipAddrs {
		data, err := yaml.Marshal(address)
		if err != nil {
			return fmt.Errorf("unable to marshal ip claim %v: %w", address.Name, err)
		}

		m.IPAddrFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(ipAddressFileName, address.Name)),
			Data:     data,
		}
	}

	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(machines))))
	for i, machine := range machines {
		data, err := yaml.Marshal(machine)
		if err != nil {
			return fmt.Errorf("marshal arbiter %d: %w", i, err)
		}

		padded := fmt.Sprintf(padFormat, i)
		m.MachineFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(arbiterMachineFileName, padded)),
			Data:     data,
		}
	}
	return nil
}

// Files returns the files generated by the asset.
func (m *Arbiter) Files() []*asset.File {
	files := make([]*asset.File, 0, 1+len(m.MachineConfigFiles)+len(m.MachineFiles))
	if m.UserDataFile != nil {
		files = append(files, m.UserDataFile)
	}
	files = append(files, m.MachineConfigFiles...)
	// Hosts refer to secrets, so place the secrets before the hosts
	// to avoid unnecessary reconciliation errors.
	files = append(files, m.SecretFiles...)
	files = append(files, m.NetworkConfigSecretFiles...)
	// Machines are linked to hosts via the machineRef, so we create
	// the hosts first to ensure if the operator starts trying to
	// reconcile a machine it can pick up the related host.
	files = append(files, m.HostFiles...)
	files = append(files, m.MachineFiles...)
	files = append(files, m.IPClaimFiles...)
	files = append(files, m.IPAddrFiles...)
	return files
}

// Load reads the asset files from disk.
func (m *Arbiter) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(filepath.Join(directory, arbiterUserDataFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	m.UserDataFile = file

	m.MachineConfigFiles, err = machineconfig.Load(f, "arbiter", directory)
	if err != nil {
		return true, err
	}

	var fileList []*asset.File

	fileList, err = f.FetchByPattern(filepath.Join(directory, secretFileNamePattern))
	if err != nil {
		return true, err
	}
	m.SecretFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, networkConfigSecretFileNamePattern))
	if err != nil {
		return true, err
	}
	m.NetworkConfigSecretFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, hostFileNamePattern))
	if err != nil {
		return true, err
	}
	m.HostFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, arbiterMachineFileNamePattern))
	if err != nil {
		return true, err
	}
	m.MachineFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, arbiterIPClaimFileNamePattern))
	if err != nil {
		return true, err
	}
	m.IPClaimFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, arbiterIPAddressFileNamePattern))
	if err != nil {
		return true, err
	}
	m.IPAddrFiles = fileList

	return true, nil
}

// Machines returns arbiter Machine manifest structures.
func (m *Arbiter) Machines() ([]machinev1beta1.Machine, error) {
	scheme := runtime.NewScheme()
	utilruntime.Must(baremetalapi.AddToScheme(scheme))
	utilruntime.Must(ibmcloudapi.AddToScheme(scheme))
	utilruntime.Must(libvirtapi.AddToScheme(scheme))
	utilruntime.Must(ovirtproviderapi.AddToScheme(scheme))
	utilruntime.Must(machinev1beta1.AddToScheme(scheme))
	utilruntime.Must(machinev1.Install(scheme))

	scheme.AddKnownTypes(machinev1alpha1.GroupVersion,
		&machinev1alpha1.OpenstackProviderSpec{},
	)
	scheme.AddKnownTypes(machinev1beta1.SchemeGroupVersion,
		&machinev1beta1.AWSMachineProviderConfig{},
		&machinev1beta1.VSphereMachineProviderSpec{},
		&machinev1beta1.AzureMachineProviderSpec{},
		&machinev1beta1.GCPMachineProviderSpec{},
	)
	scheme.AddKnownTypes(machinev1.GroupVersion,
		&machinev1.NutanixMachineProviderConfig{},
		&machinev1.PowerVSMachineProviderConfig{},
	)

	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		machinev1.GroupVersion,
		baremetalprovider.SchemeGroupVersion,
		ibmcloudprovider.SchemeGroupVersion,
		libvirtprovider.SchemeGroupVersion,
		machinev1alpha1.GroupVersion,
		machinev1beta1.SchemeGroupVersion,
		ovirtprovider.SchemeGroupVersion,
	)

	machines := []machinev1beta1.Machine{}
	for i, file := range m.MachineFiles {
		machine := &machinev1beta1.Machine{}
		err := yaml.Unmarshal(file.Data, &machine)
		if err != nil {
			return machines, fmt.Errorf("unmarshal arbiter %d, %w", i, err)
		}

		obj, _, err := decoder.Decode(machine.Spec.ProviderSpec.Value.Raw, nil, nil)
		if err != nil {
			return machines, fmt.Errorf("unmarshal arbiter %d: %w", i, err)
		}

		machine.Spec.ProviderSpec.Value = &runtime.RawExtension{Object: obj}
		machines = append(machines, *machine)
	}

	return machines, nil
}
