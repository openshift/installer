package manifests

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	operv1 "github.com/openshift/api/operator/v1"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/defaults"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const (
	installConfigOverrides = aiv1beta1.Group + "/install-config-overrides"
)

var (
	agentClusterInstallFilename = filepath.Join(clusterManifestDir, "agent-cluster-install.yaml")
)

// AgentClusterInstall generates the agent-cluster-install.yaml file.
type AgentClusterInstall struct {
	File   *asset.File
	Config *hiveext.AgentClusterInstall
}

type agentClusterInstallOnPremPlatform struct {
	// APIVIPs contains the VIP(s) to use for internal API communication. In
	// dual stack clusters it contains an IPv4 and IPv6 address, otherwise only
	// one VIP
	APIVIPs []string `json:"apiVIPs,omitempty"`

	// IngressVIPs contains the VIP(s) to use for ingress traffic. In dual stack
	// clusters it contains an IPv4 and IPv6 address, otherwise only one VIP
	IngressVIPs []string `json:"ingressVIPs,omitempty"`

	// Host, including BMC, configuration.
	Hosts []baremetal.Host `json:"hosts,omitempty"`

	// ClusterProvisioningIP is the IP on the dedicated provisioning network.
	ClusterProvisioningIP string `json:"clusterProvisioningIP,omitempty"`

	// ProvisioningNetwork is used to indicate if we will have a provisioning network, and how it will be managed.
	ProvisioningNetwork baremetal.ProvisioningNetwork `json:"provisioningNetwork,omitempty"`

	// ProvisioningNetworkInterface is the name of the network interface on a control plane
	// baremetal host that is connected to the provisioning network.
	ProvisioningNetworkInterface string `json:"provisioningNetworkInterface,omitempty"`

	// ProvisioningNetworkCIDR defines the network to use for provisioning.
	ProvisioningNetworkCIDR *ipnet.IPNet `json:"provisioningNetworkCIDR,omitempty"`

	// ProvisioningDHCPRange is used to provide DHCP services to hosts
	// for provisioning.
	ProvisioningDHCPRange string `json:"provisioningDHCPRange,omitempty"`
}

type agentClusterInstallOnPremExternalPlatform struct {
	// PlatformName holds the arbitrary string representing the infrastructure provider name, expected to be set at the installation time.
	PlatformName string `json:"platformName,omitempty"`
	// CloudControllerManager when set to external, this property will enable an external cloud provider.
	CloudControllerManager external.CloudControllerManager `json:"cloudControllerManager,omitempty"`
}

type agentClusterInstallPlatform struct {
	// BareMetal is the configuration used when installing on bare metal.
	// +optional
	BareMetal *agentClusterInstallOnPremPlatform `json:"baremetal,omitempty"`
	// VSphere is the configuration used when installing on vSphere.
	// +optional
	VSphere *vsphere.Platform `json:"vsphere,omitempty"`
	// External is the configuration used when installing on external cloud provider.
	// +optional
	External *agentClusterInstallOnPremExternalPlatform `json:"external,omitempty"`
	// Nutanix is the configuration used when installing on nutanix platform.
	// +optional
	Nutanix *nutanix.Platform `json:"nutanix,omitempty"`
}

// Used to generate InstallConfig overrides for Assisted Service to apply
type agentClusterInstallInstallConfigOverrides struct {
	// FIPS configures https://www.nist.gov/itl/fips-general-information
	//
	// +kubebuilder:default=false
	// +optional
	FIPS bool `json:"fips,omitempty"`
	// Platform is the configuration for the specific platform upon which to
	// perform the installation.
	Platform *agentClusterInstallPlatform `json:"platform,omitempty"`
	// Capabilities selects the managed set of optional, core cluster components.
	Capabilities *types.Capabilities `json:"capabilities,omitempty"`
	// Allow override of network type
	Networking *types.Networking `json:"networking,omitempty"`
	// Allow override of CPUPartitioning
	CPUPartitioning types.CPUPartitioningMode `json:"cpuPartitioningMode,omitempty"`
	// Allow override of AdditionalTrustBundlePolicy
	AdditionalTrustBundlePolicy types.PolicyType `json:"additionalTrustBundlePolicy,omitempty"`
	// Allow override of FeatureSet
	FeatureSet configv1.FeatureSet `json:"featureSet,omitempty"`
	// Allow override of FeatureGates
	FeatureGates []string `json:"featureGates,omitempty"`
}

var _ asset.WritableAsset = (*AgentClusterInstall)(nil)

// Name returns a human friendly name for the asset.
func (*AgentClusterInstall) Name() string {
	return "AgentClusterInstall Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*AgentClusterInstall) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&agent.OptionalInstallConfig{},
		&agentconfig.AgentHosts{},
		&agentconfig.AgentConfig{},
	}
}

// Generate generates the AgentClusterInstall manifest.
//
//nolint:gocyclo
func (a *AgentClusterInstall) Generate(_ context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	installConfig := &agent.OptionalInstallConfig{}
	agentHosts := &agentconfig.AgentHosts{}
	agentConfig := &agentconfig.AgentConfig{}
	dependencies.Get(agentWorkflow, agentHosts, installConfig, agentConfig)

	// This manifest is not required for AddNodes workflow
	if agentWorkflow.Workflow == workflow.AgentWorkflowTypeAddNodes {
		// Add empty file to keep config ISO loader happy
		a.File = &asset.File{
			Filename: agentClusterInstallFilename,
		}
		return nil
	}

	if installConfig.Config != nil {
		var numberOfWorkers int = 0
		for _, compute := range installConfig.Config.Compute {
			numberOfWorkers = numberOfWorkers + int(*compute.Replicas)
		}

		numberOfArbiters := 0
		if installConfig.Config.IsArbiterEnabled() {
			numberOfArbiters = int(*installConfig.Config.Arbiter.Replicas)
		}

		clusterNetwork := []hiveext.ClusterNetworkEntry{}
		for _, cn := range installConfig.Config.Networking.ClusterNetwork {
			entry := hiveext.ClusterNetworkEntry{
				CIDR:       cn.CIDR.String(),
				HostPrefix: cn.HostPrefix,
			}
			clusterNetwork = append(clusterNetwork, entry)
		}

		serviceNetwork := []string{}
		for _, sn := range installConfig.Config.Networking.ServiceNetwork {
			serviceNetwork = append(serviceNetwork, sn.String())
		}

		machineNetwork := []hiveext.MachineNetworkEntry{}
		for _, mn := range installConfig.Config.Networking.MachineNetwork {
			entry := hiveext.MachineNetworkEntry{
				CIDR: mn.CIDR.String(),
			}
			machineNetwork = append(machineNetwork, entry)
		}

		agentClusterInstall := &hiveext.AgentClusterInstall{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AgentClusterInstall",
				APIVersion: hiveext.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      getAgentClusterInstallName(installConfig),
				Namespace: installConfig.ClusterNamespace(),
			},
			Spec: hiveext.AgentClusterInstallSpec{
				ImageSetRef: &hivev1.ClusterImageSetReference{
					Name: getClusterImageSetReferenceName(),
				},
				ClusterDeploymentRef: corev1.LocalObjectReference{
					Name: getClusterDeploymentName(installConfig),
				},
				Networking: hiveext.Networking{
					MachineNetwork: machineNetwork,
					ClusterNetwork: clusterNetwork,
					ServiceNetwork: serviceNetwork,
				},
				SSHPublicKey: strings.Trim(installConfig.Config.SSHKey, "|\n\t"),
				ProvisionRequirements: hiveext.ProvisionRequirements{
					ControlPlaneAgents: int(*installConfig.Config.ControlPlane.Replicas),
					ArbiterAgents:      numberOfArbiters,
					WorkerAgents:       numberOfWorkers,
				},
				PlatformType: agent.HivePlatformType(installConfig.Config.Platform),
			},
		}

		if agentClusterInstall.Spec.PlatformType == hiveext.ExternalPlatformType {
			agentClusterInstall.Spec.ExternalPlatformSpec = &hiveext.ExternalPlatformSpec{
				PlatformName: installConfig.Config.Platform.External.PlatformName,
			}
		}
		if installConfig.Config.Platform.Name() == external.Name && installConfig.Config.Platform.External.PlatformName == agent.ExternalPlatformNameOci {
			agentClusterInstall.Spec.ExternalPlatformSpec.CloudControllerManager = external.CloudControllerManagerTypeExternal
		}

		agentClusterInstall.Spec.Networking.UserManagedNetworking = agent.GetUserManagedNetworkingByPlatformType(agent.HivePlatformType(installConfig.Config.Platform))

		icOverridden := false
		icOverrides := agentClusterInstallInstallConfigOverrides{}
		if installConfig.Config.FIPS {
			icOverridden = true
			icOverrides.FIPS = installConfig.Config.FIPS
		}

		if len(installConfig.Config.FeatureSet) > 0 {
			icOverridden = true
			icOverrides.FeatureSet = installConfig.Config.FeatureSet
		}

		if len(installConfig.Config.FeatureGates) > 0 {
			icOverridden = true
			icOverrides.FeatureGates = installConfig.Config.FeatureGates
		}

		if installConfig.Config.Proxy != nil {
			rendezvousIP := ""
			if agentConfig.Config != nil {
				rendezvousIP = agentConfig.Config.RendezvousIP
			}

			agentClusterInstall.Spec.Proxy = (*hiveext.Proxy)(getProxy(installConfig.Config.Proxy, &installConfig.Config.Networking.MachineNetwork, rendezvousIP))
		}

		if installConfig.Config.Platform.BareMetal != nil {
			baremetalPlatform := agentClusterInstallOnPremPlatform{}
			bmIcOverridden := false

			for _, host := range agentHosts.Hosts {
				// Override if BMC values are not the same as default
				if host.BMC.Username != "" || host.BMC.Password != "" || host.BMC.Address != "" {
					bmhost := baremetal.Host{
						Name: host.Hostname,
						Role: host.Role,
					}
					if len(host.Interfaces) > 0 {
						// Boot MAC address is stored as first interface
						bmhost.BootMACAddress = host.Interfaces[0].MacAddress
					} else {
						logrus.Infof("Could not obtain baremetal BootMacAddress for %s", installConfig.Config.Platform.Name())
					}
					bmIcOverridden = true
					bmhost.BMC = host.BMC
					baremetalPlatform.Hosts = append(baremetalPlatform.Hosts, bmhost)

					// Set provisioning network configuration
					baremetalPlatform.ClusterProvisioningIP = installConfig.Config.Platform.BareMetal.ClusterProvisioningIP
					baremetalPlatform.ProvisioningNetwork = installConfig.Config.Platform.BareMetal.ProvisioningNetwork
					baremetalPlatform.ProvisioningNetworkInterface = installConfig.Config.Platform.BareMetal.ProvisioningNetworkInterface
					baremetalPlatform.ProvisioningNetworkCIDR = installConfig.Config.Platform.BareMetal.ProvisioningNetworkCIDR
					baremetalPlatform.ProvisioningDHCPRange = installConfig.Config.Platform.BareMetal.ProvisioningDHCPRange
				}
			}
			if bmIcOverridden {
				icOverridden = true
				icOverrides.Platform = &agentClusterInstallPlatform{}
				icOverrides.Platform = &agentClusterInstallPlatform{
					BareMetal: &baremetalPlatform,
				}
			}

			agentClusterInstall.Spec.APIVIPs = installConfig.Config.Platform.BareMetal.APIVIPs
			agentClusterInstall.Spec.IngressVIPs = installConfig.Config.Platform.BareMetal.IngressVIPs
			agentClusterInstall.Spec.APIVIP = installConfig.Config.Platform.BareMetal.APIVIPs[0]
			agentClusterInstall.Spec.IngressVIP = installConfig.Config.Platform.BareMetal.IngressVIPs[0]

			// Copy LoadBalancer configuration to allow UserManaged load balancer with same API/Ingress VIPs
			if installConfig.Config.Platform.BareMetal.LoadBalancer != nil {
				agentClusterInstall.Spec.LoadBalancer = &hiveext.LoadBalancer{
					Type: convertLoadBalancerType(installConfig.Config.Platform.BareMetal.LoadBalancer.Type),
				}
			}
		} else if installConfig.Config.Platform.VSphere != nil {
			vspherePlatform := vsphere.Platform{}
			if len(installConfig.Config.Platform.VSphere.APIVIPs) > 1 {
				icOverridden = true
				vspherePlatform.APIVIPs = installConfig.Config.Platform.VSphere.APIVIPs
				vspherePlatform.IngressVIPs = installConfig.Config.Platform.VSphere.IngressVIPs
			}
			hasCredentials := false
			if len(installConfig.Config.Platform.VSphere.VCenters) > 0 {
				for _, vcenter := range installConfig.Config.Platform.VSphere.VCenters {
					if agent.VCenterCredentialsAreProvided(vcenter) {
						icOverridden = true
						hasCredentials = true
						vspherePlatform.VCenters = append(vspherePlatform.VCenters, vcenter)
					}
				}
			}
			if hasCredentials && len(installConfig.Config.Platform.VSphere.FailureDomains) > 0 {
				icOverridden = true
				vspherePlatform.FailureDomains = append(vspherePlatform.FailureDomains, installConfig.Config.VSphere.FailureDomains...)
			}
			if icOverridden {
				icOverrides.Platform = &agentClusterInstallPlatform{
					VSphere: &vspherePlatform,
				}
			}
			agentClusterInstall.Spec.APIVIPs = installConfig.Config.Platform.VSphere.APIVIPs
			agentClusterInstall.Spec.IngressVIPs = installConfig.Config.Platform.VSphere.IngressVIPs
		} else if installConfig.Config.Platform.Nutanix != nil {
			icNutanixPlatformBytes, err := json.Marshal(*installConfig.Config.Platform.Nutanix)
			if err != nil {
				logrus.Errorf("failed to marshal installConfig.platform.nutanix: %v", err)
			}
			nutanixPlatform := nutanix.Platform{}
			err = json.Unmarshal(icNutanixPlatformBytes, &nutanixPlatform)
			if err != nil {
				logrus.Errorf("failed to unmarshal installConfig.platform.nutanix: %v", err)
			}

			// Skip the below agent installer not supported fields
			nutanixPlatform.ClusterOSImage = ""
			nutanixPlatform.PreloadedOSImageName = ""
			nutanixPlatform.DefaultMachinePlatform = nil
			nutanixPlatform.LoadBalancer = nil
			nutanixPlatform.FailureDomains = nil
			nutanixPlatform.PrismAPICallTimeout = nil

			icOverridden = true
			icOverrides.Platform = &agentClusterInstallPlatform{
				Nutanix: &nutanixPlatform,
			}
			agentClusterInstall.Spec.APIVIPs = installConfig.Config.Platform.Nutanix.APIVIPs
			agentClusterInstall.Spec.IngressVIPs = installConfig.Config.Platform.Nutanix.IngressVIPs
		} else if installConfig.Config.Platform.External != nil {
			icOverridden = true
			icOverrides.Platform = &agentClusterInstallPlatform{
				External: &agentClusterInstallOnPremExternalPlatform{
					PlatformName:           installConfig.Config.External.PlatformName,
					CloudControllerManager: installConfig.Config.External.CloudControllerManager,
				},
			}
		}

		networkOverridden := setNetworkType(agentClusterInstall, installConfig.Config, "NetworkType is not specified in InstallConfig.")
		if networkOverridden {
			icOverridden = true
			icOverrides.Networking = installConfig.Config.Networking
		}

		if installConfig.Config.Capabilities != nil {
			icOverrides.Capabilities = installConfig.Config.Capabilities
			icOverridden = true
		}

		if installConfig.Config.CPUPartitioning != "" {
			icOverridden = true
			icOverrides.CPUPartitioning = installConfig.Config.CPUPartitioning
		}

		if installConfig.Config.AdditionalTrustBundlePolicy != "" && installConfig.Config.AdditionalTrustBundlePolicy != types.PolicyProxyOnly {
			icOverridden = true
			icOverrides.AdditionalTrustBundlePolicy = installConfig.Config.AdditionalTrustBundlePolicy
		}

		if icOverridden {
			overrides, err := json.Marshal(icOverrides)
			if err != nil {
				return errors.Wrap(err, "failed to marshal AgentClusterInstall installConfigOverrides")
			}
			agentClusterInstall.SetAnnotations(map[string]string{
				installConfigOverrides: string(overrides),
			})
		}

		a.Config = agentClusterInstall

	}
	return a.finish()
}

// Files returns the files generated by the asset.
func (a *AgentClusterInstall) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return []*asset.File{}
}

// Load returns agentclusterinstall asset from the disk.
func (a *AgentClusterInstall) Load(f asset.FileFetcher) (bool, error) {

	agentClusterInstallFile, err := f.FetchByName(agentClusterInstallFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", agentClusterInstallFilename))
	}

	agentClusterInstall := &hiveext.AgentClusterInstall{}
	if err := yaml.UnmarshalStrict(agentClusterInstallFile.Data, agentClusterInstall); err != nil {
		err = errors.Wrapf(err, "failed to unmarshal %s", agentClusterInstallFilename)
		return false, err
	}

	setNetworkType(agentClusterInstall, &types.InstallConfig{}, "NetworkType is not specified in AgentClusterInstall.")

	// Due to OCPBUGS-7495 we previously required lowercase platform names here,
	// even though that is incorrect. Rewrite to the correct mixed case names
	// for backward compatibility.
	switch string(agentClusterInstall.Spec.PlatformType) {
	case baremetal.Name:
		agentClusterInstall.Spec.PlatformType = hiveext.BareMetalPlatformType
	case external.Name:
		agentClusterInstall.Spec.PlatformType = hiveext.ExternalPlatformType
	case none.Name:
		agentClusterInstall.Spec.PlatformType = hiveext.NonePlatformType
	case vsphere.Name:
		agentClusterInstall.Spec.PlatformType = hiveext.VSpherePlatformType
	case nutanix.Name:
		agentClusterInstall.Spec.PlatformType = hiveext.NutanixPlatformType
	}

	// Set the default value for userManagedNetworking, as would be done by the
	// mutating webhook in ZTP.
	if agentClusterInstall.Spec.Networking.UserManagedNetworking == nil {
		agentClusterInstall.Spec.Networking.UserManagedNetworking = agent.GetUserManagedNetworkingByPlatformType(agentClusterInstall.Spec.PlatformType)
	}

	a.Config = agentClusterInstall

	if err = a.finish(); err != nil {
		return false, err
	}
	return true, nil
}

func (a *AgentClusterInstall) finish() error {

	if a.Config == nil {
		return errors.New("missing configuration or manifest file")
	}

	if err := a.validateIPAddressAndNetworkType().ToAggregate(); err != nil {
		return errors.Wrapf(err, "invalid NetworkType configured")
	}

	if err := a.validateSupportedPlatforms().ToAggregate(); err != nil {
		return errors.Wrapf(err, "invalid PlatformType configured")
	}

	if err := a.validateDiskEncryption().ToAggregate(); err != nil {
		return errors.Wrapf(err, "invalid DiskEncryption configured")
	}

	agentClusterInstallData, err := yaml.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal agent installer AgentClusterInstall")
	}

	a.File = &asset.File{
		Filename: agentClusterInstallFilename,
		Data:     agentClusterInstallData,
	}
	return nil
}

// Sets the default network type to OVNKubernetes if it is unspecified in the
// AgentClusterInstall or InstallConfig.
func setNetworkType(aci *hiveext.AgentClusterInstall, installConfig *types.InstallConfig,
	warningMessage string) bool {
	if aci.Spec.Networking.NetworkType != "" {
		return false
	}

	if installConfig != nil && installConfig.Networking != nil &&
		installConfig.Networking.NetworkType != "" {
		if installConfig.Networking.NetworkType == string(operv1.NetworkTypeOVNKubernetes) || installConfig.Networking.NetworkType == string(operv1.NetworkTypeOpenShiftSDN) {
			aci.Spec.Networking.NetworkType = installConfig.NetworkType
			return false
		}

		// Set OVNKubernetes in AgentClusterInstall and return true to indicate InstallConfigOverride should be used
		aci.Spec.Networking.NetworkType = string(operv1.NetworkTypeOVNKubernetes)
		return true
	}

	defaults.SetInstallConfigDefaults(installConfig)
	logrus.Infof("%s Defaulting NetworkType to %s.", warningMessage, installConfig.NetworkType)
	aci.Spec.Networking.NetworkType = installConfig.NetworkType
	return false
}

func isIPv6(ipAddress net.IP) bool {
	// Using To16() on IPv4 addresses does not return nil so it cannot be used to determine if
	// IP addresses are IPv6. Instead we are checking if the address is IPv6 by using To4().
	// Same as https://github.com/openshift/installer/blob/6eca978b89fc0be17f70fc8a28fa20aab1316843/pkg/types/validation/installconfig.go#L193
	ip := ipAddress.To4()
	return ip == nil
}

func (a *AgentClusterInstall) validateIPAddressAndNetworkType() field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("spec", "networking", "networkType")
	clusterNetworkPath := field.NewPath("spec", "networking", "clusterNetwork")
	serviceNetworkPath := field.NewPath("spec", "networking", "serviceNetwork")

	switch a.Config.Spec.Networking.NetworkType {
	case string(operv1.NetworkTypeOpenShiftSDN):
		hasIPv6 := false
		for _, cn := range a.Config.Spec.Networking.ClusterNetwork {
			ipNet, errCIDR := ipnet.ParseCIDR(cn.CIDR)
			if errCIDR != nil {
				allErrs = append(allErrs, field.Required(clusterNetworkPath, "error parsing the clusterNetwork CIDR"))
				continue
			}
			if isIPv6(ipNet.IP) {
				hasIPv6 = true
			}
		}
		if hasIPv6 {
			allErrs = append(allErrs, field.Required(fieldPath,
				fmt.Sprintf("clusterNetwork CIDR is IPv6 and is not compatible with networkType %s",
					operv1.NetworkTypeOpenShiftSDN)))
		}

		hasIPv6 = false
		for _, cidr := range a.Config.Spec.Networking.ServiceNetwork {
			ipNet, errCIDR := ipnet.ParseCIDR(cidr)
			if errCIDR != nil {
				allErrs = append(allErrs, field.Required(serviceNetworkPath, "error parsing the clusterNetwork CIDR"))
				continue
			}
			if isIPv6(ipNet.IP) {
				hasIPv6 = true
			}
		}
		if hasIPv6 {
			allErrs = append(allErrs, field.Required(fieldPath,
				fmt.Sprintf("serviceNetwork CIDR is IPv6 and is not compatible with networkType %s",
					operv1.NetworkTypeOpenShiftSDN)))
		}
	case string(operv1.NetworkTypeOVNKubernetes):
		for i, cn := range a.Config.Spec.Networking.ClusterNetwork {
			path := clusterNetworkPath.Index(i)
			ipNet, errCIDR := ipnet.ParseCIDR(cn.CIDR)
			if errCIDR != nil {
				allErrs = append(allErrs, field.Required(path.Child("cidr"), "error parsing the clusterNetwork CIDR"))
				continue
			}
			cnOnes, cnBits := ipNet.Mask.Size()
			maxHostPrefix := int32(cnBits) - 7
			if cn.HostPrefix > maxHostPrefix {
				allErrs = append(allErrs, field.Invalid(path.Child("hostPrefix"), cn.HostPrefix, fmt.Sprintf("must be at most %d", maxHostPrefix)))
			}

			numHosts := a.Config.Spec.ProvisionRequirements.ControlPlaneAgents + a.Config.Spec.ProvisionRequirements.WorkerAgents
			var minPrefixDiff int32
			for (1 << minPrefixDiff) < numHosts {
				minPrefixDiff++
			}
			if (cn.HostPrefix - int32(cnOnes)) < minPrefixDiff {
				allErrs = append(allErrs, field.Invalid(path, cn.CIDR, fmt.Sprintf("prefix length %d not large enough to accommodate %d hosts with hostPrefix length %d", cnOnes, numHosts, cn.HostPrefix)))
			}
		}
	}

	return allErrs
}

func (a *AgentClusterInstall) validateSupportedPlatforms() field.ErrorList {
	var allErrs field.ErrorList

	if a.Config.Spec.PlatformType != "" && !agent.IsSupportedPlatform(a.Config.Spec.PlatformType) {
		fieldPath := field.NewPath("spec", "platformType")
		allErrs = append(allErrs, field.NotSupported(fieldPath, a.Config.Spec.PlatformType, agent.SupportedHivePlatforms()))
	}

	switch a.Config.Spec.PlatformType {
	case hiveext.NonePlatformType, hiveext.ExternalPlatformType:
		if a.Config.Spec.Networking.UserManagedNetworking != nil && !*a.Config.Spec.Networking.UserManagedNetworking {
			fieldPath := field.NewPath("spec", "networking", "userManagedNetworking")
			allErrs = append(allErrs, field.Forbidden(fieldPath,
				fmt.Sprintf("%s platform requires user-managed networking",
					a.Config.Spec.PlatformType)))
		}
	}
	return allErrs
}

// FIPSEnabled returns whether FIPS is enabled in the cluster configuration.
func (a *AgentClusterInstall) FIPSEnabled() bool {
	icOverrides := agentClusterInstallInstallConfigOverrides{}
	if err := json.Unmarshal([]byte(a.Config.Annotations[installConfigOverrides]), &icOverrides); err == nil {
		return icOverrides.FIPS
	}
	return false
}

// GetExternalPlatformName returns the platform name for the external platform.
func (a *AgentClusterInstall) GetExternalPlatformName() string {
	if a.Config != nil && a.Config.Spec.ExternalPlatformSpec != nil {
		return a.Config.Spec.ExternalPlatformSpec.PlatformName
	}
	return ""
}

func (a *AgentClusterInstall) validateDiskEncryption() field.ErrorList {
	var allErrs field.ErrorList
	supportedEnableOn := []string{models.DiskEncryptionEnableOnNone, models.DiskEncryptionEnableOnAll, models.DiskEncryptionEnableOnMasters, models.DiskEncryptionEnableOnWorkers}
	supportedMode := []string{models.DiskEncryptionModeTpmv2, models.DiskEncryptionModeTang}

	if a.Config.Spec.DiskEncryption != nil {
		if !slices.Contains(supportedEnableOn, swag.StringValue(a.Config.Spec.DiskEncryption.EnableOn)) {
			fieldPath := field.NewPath("spec", "diskEncryption", "enableOn")
			allErrs = append(allErrs, field.NotSupported(fieldPath, a.Config.Spec.DiskEncryption.EnableOn, supportedEnableOn))
		}

		if !slices.Contains(supportedMode, swag.StringValue(a.Config.Spec.DiskEncryption.Mode)) {
			fieldPath := field.NewPath("spec", "diskEncryption", "mode")
			allErrs = append(allErrs, field.NotSupported(fieldPath, a.Config.Spec.DiskEncryption.Mode, supportedMode))
		}
	}
	return allErrs
}

// convertLoadBalancerType converts the configv1 PlatformLoadBalancerType to hiveext LoadBalancerType.
func convertLoadBalancerType(lbType configv1.PlatformLoadBalancerType) hiveext.LoadBalancerType {
	switch lbType {
	case configv1.LoadBalancerTypeUserManaged:
		return hiveext.LoadBalancerTypeUserManaged
	case configv1.LoadBalancerTypeOpenShiftManagedDefault:
		return hiveext.LoadBalancerTypeClusterManaged
	default:
		// Default to ClusterManaged if type is empty or unknown
		return hiveext.LoadBalancerTypeClusterManaged
	}
}
