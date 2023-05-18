package manifests

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

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

type agentClusterInstallPlatform struct {
	// BareMetal is the configuration used when installing on bare metal.
	// +optional
	BareMetal *agentClusterInstallOnPremPlatform `json:"baremetal,omitempty"`
	// VSphere is the configuration used when installing on vSphere.
	// +optional
	VSphere *agentClusterInstallOnPremPlatform `json:"vsphere,omitempty"`
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
	// AdditionalTrustBundle must be set here when mirroring not configured
	AdditionalTrustBundle string `json:"additionalTrustBundle,omitempty"`
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
			if len(installConfig.Config.Platform.BareMetal.APIVIPs) > 1 {
				icOverridden = true
				icOverrides.Platform = &agentClusterInstallPlatform{
					BareMetal: &agentClusterInstallOnPremPlatform{
						APIVIPs:     installConfig.Config.Platform.BareMetal.APIVIPs,
						IngressVIPs: installConfig.Config.Platform.BareMetal.IngressVIPs,
					},
				}
			}
			agentClusterInstall.Spec.APIVIP = installConfig.Config.Platform.BareMetal.APIVIPs[0]
			agentClusterInstall.Spec.IngressVIP = installConfig.Config.Platform.BareMetal.IngressVIPs[0]
		} else if installConfig.Config.Platform.VSphere != nil {
			if len(installConfig.Config.Platform.VSphere.APIVIPs) > 1 {
				icOverridden = true
				icOverrides.Platform = &agentClusterInstallPlatform{
					VSphere: &agentClusterInstallOnPremPlatform{
						APIVIPs:     installConfig.Config.Platform.VSphere.APIVIPs,
						IngressVIPs: installConfig.Config.Platform.VSphere.IngressVIPs,
					},
				}
			}
			agentClusterInstall.Spec.APIVIP = installConfig.Config.Platform.VSphere.APIVIPs[0]
			agentClusterInstall.Spec.IngressVIP = installConfig.Config.Platform.VSphere.IngressVIPs[0]
		}

		setNetworkType(agentClusterInstall, installConfig.Config, "NetworkType is not specified in InstallConfig.")

		if installConfig.Config.Capabilities != nil {
			icOverrides.Capabilities = installConfig.Config.Capabilities
			icOverridden = true
		}

		if installConfig.Config.AdditionalTrustBundle != "" {
			// Add trust bundle to the config overrides to be included in installed image
			// TODO: when MGMT-11520 adds support for AdditionalTrustBundle as part of the InfraEnv CRD
			// then it must be set in the infraEnv manifest instead of below
			icOverrides.AdditionalTrustBundle = installConfig.Config.AdditionalTrustBundle
			icOverridden = true
		}
		if icOverridden {
			overrides, err := json.Marshal(icOverrides)
			if err != nil {
				return errors.Wrap(err, "failed to marshal AgentClusterInstall installConfigOverrides")
			}
			agentClusterInstall.SetAnnotations(map[string]string{
				installConfigOverrides: fmt.Sprintf("%s", overrides),
			})
		}

		a.Config = agentClusterInstall

		agentClusterInstallData, err := yaml.Marshal(agentClusterInstall)
		if err != nil {
			return errors.Wrap(err, "failed to marshal agent installer AgentClusterInstall")
		}

		a.File = &asset.File{
			Filename: agentClusterInstallFilename,
			Data:     agentClusterInstallData,
		}
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

	a.File = agentClusterInstallFile

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
	case none.Name:
		agentClusterInstall.Spec.PlatformType = hiveext.NonePlatformType
	case vsphere.Name:
		agentClusterInstall.Spec.PlatformType = hiveext.VSpherePlatformType
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

	return nil
}

// Sets the default network type to OVNKubernetes if it is unspecified in the
// AgentClusterInstall or InstallConfig
func setNetworkType(aci *hiveext.AgentClusterInstall, installConfig *types.InstallConfig,
	warningMessage string) {

	if aci.Spec.Networking.NetworkType != "" {
		return
	}

	if installConfig != nil && installConfig.Networking != nil &&
		installConfig.Networking.NetworkType != "" {
		aci.Spec.Networking.NetworkType = installConfig.NetworkType
		return
	}

	defaults.SetInstallConfigDefaults(installConfig)
	logrus.Infof("%s Defaulting NetworkType to %s.", warningMessage, installConfig.NetworkType)
	aci.Spec.Networking.NetworkType = installConfig.NetworkType
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

	fieldPath := field.NewPath("spec", "platformType")

	if a.Config.Spec.PlatformType != "" && !agent.IsSupportedPlatform(a.Config.Spec.PlatformType) {
		allErrs = append(allErrs, field.NotSupported(fieldPath, a.Config.Spec.PlatformType, agent.SupportedHivePlatforms()))
	}
	return allErrs
}
