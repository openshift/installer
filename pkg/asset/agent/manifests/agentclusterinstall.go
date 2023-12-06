package manifests

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	operv1 "github.com/openshift/api/operator/v1"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/defaults"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/none"
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
		&agent.OptionalInstallConfig{},
	}
}

// Generate generates the AgentClusterInstall manifest.
func (a *AgentClusterInstall) Generate(dependencies asset.Parents) error {
	installConfig := &agent.OptionalInstallConfig{}
	dependencies.Get(installConfig)

	if installConfig.Config != nil {
		var numberOfWorkers int = 0
		for _, compute := range installConfig.Config.Compute {
			numberOfWorkers = numberOfWorkers + int(*compute.Replicas)
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
			ObjectMeta: metav1.ObjectMeta{
				Name:      getAgentClusterInstallName(installConfig),
				Namespace: getObjectMetaNamespace(installConfig),
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

		if installConfig.Config.Platform.Name() == none.Name || installConfig.Config.Platform.Name() == external.Name {
			logrus.Debugf("Setting UserManagedNetworking to true for %s platform", installConfig.Config.Platform.Name())
			agentClusterInstall.Spec.Networking.UserManagedNetworking = swag.Bool(true)
		}

		icOverridden := false
		icOverrides := agentClusterInstallInstallConfigOverrides{}
		if installConfig.Config.FIPS {
			icOverridden = true
			icOverrides.FIPS = installConfig.Config.FIPS
		}

		if installConfig.Config.Proxy != nil {
			agentClusterInstall.Spec.Proxy = (*hiveext.Proxy)(getProxy(installConfig))
		}

		if installConfig.Config.Platform.BareMetal != nil {
			agentClusterInstall.Spec.APIVIPs = installConfig.Config.Platform.BareMetal.APIVIPs
			agentClusterInstall.Spec.IngressVIPs = installConfig.Config.Platform.BareMetal.IngressVIPs
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
	}

	// Set the default value for userManagedNetworking, as would be done by the
	// mutating webhook in ZTP.
	if agentClusterInstall.Spec.Networking.UserManagedNetworking == nil {
		switch agentClusterInstall.Spec.PlatformType {
		case hiveext.NonePlatformType, hiveext.ExternalPlatformType:
			logrus.Debugf("Setting UserManagedNetworking to true for %s platform", agentClusterInstall.Spec.PlatformType)
			agentClusterInstall.Spec.Networking.UserManagedNetworking = swag.Bool(true)
		}
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

	if a.Config.Spec.Networking.NetworkType == string(operv1.NetworkTypeOpenShiftSDN) {
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
