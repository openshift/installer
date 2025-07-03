package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func GetCoreOS() (string, error) {
	logrus.Info("Getting CoreOS stream data")
	streamData, err := rhcos.FetchRawCoreOSStream(context.Background())
	if err != nil {
		return "", err
	}
	return string(streamData), nil
}

func GetExampleInstallConfig(platform, pullSecret, baseDomain, clusterName, sshKey string) (string, error) {

	logrus.Info("in getInstallConfigResource")
	logrus.Infof("platform: %s", platform)

	// Create a new InstallConfig
	installConfig := &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName,
		},
		BaseDomain: baseDomain,
		PullSecret: pullSecret,
		SSHKey:     sshKey,
		Publish:    types.ExternalPublishingStrategy,
		Networking: &types.Networking{
			NetworkType:    "OVNKubernetes",
			MachineNetwork: []types.MachineNetworkEntry{{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")}},
			ClusterNetwork: []types.ClusterNetworkEntry{{CIDR: *ipnet.MustParseCIDR("10.128.0.0/14"), HostPrefix: 23}},
			ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("172.30.0.0/16")},
		},
		ControlPlane: &types.MachinePool{
			Name:           "master",
			Replicas:       pointer.Int64Ptr(3),
			Hyperthreading: types.HyperthreadingEnabled,
			Architecture:   types.ArchitectureAMD64,
		},
		Compute: []types.MachinePool{
			{
				Name:           "worker",
				Replicas:       pointer.Int64Ptr(3),
				Hyperthreading: types.HyperthreadingEnabled,
				Architecture:   types.ArchitectureAMD64,
			},
		},
	}

	// Set platform-specific configuration
	switch platform {
	case aws.Name:
		installConfig.Platform = types.Platform{
			AWS: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "subnet-12345678", Roles: []aws.SubnetRole{{Type: aws.ClusterNodeSubnetRole}}},
						{ID: "subnet-87654321", Roles: []aws.SubnetRole{{Type: aws.ControlPlaneExternalLBSubnetRole}}},
					},
				},
			},
		}
		// Set AWS-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			AWS: &aws.MachinePool{
				InstanceType: "m5.xlarge",
				EC2RootVolume: aws.EC2RootVolume{
					Size: 120,
					Type: "gp3",
				},
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			AWS: &aws.MachinePool{
				InstanceType: "m5.large",
				EC2RootVolume: aws.EC2RootVolume{
					Size: 120,
					Type: "gp3",
				},
			},
		}

	case azure.Name:
		installConfig.Platform = types.Platform{
			Azure: &azure.Platform{
				Region: "eastus",
			},
		}
		// Set Azure-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			Azure: &azure.MachinePool{
				InstanceType: "Standard_D4s_v3",
				OSDisk: azure.OSDisk{
					DiskSizeGB: 120,
					DiskType:   "Premium_LRS",
				},
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			Azure: &azure.MachinePool{
				InstanceType: "Standard_D2s_v3",
				OSDisk: azure.OSDisk{
					DiskSizeGB: 120,
					DiskType:   "Premium_LRS",
				},
			},
		}

	case gcp.Name:
		installConfig.Platform = types.Platform{
			GCP: &gcp.Platform{
				ProjectID: "example-project",
				Region:    "us-central1",
			},
		}
		// Set GCP-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			GCP: &gcp.MachinePool{
				InstanceType: "n1-standard-4",
				OSDisk: gcp.OSDisk{
					DiskSizeGB: 120,
					DiskType:   "pd-ssd",
				},
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			GCP: &gcp.MachinePool{
				InstanceType: "n1-standard-2",
				OSDisk: gcp.OSDisk{
					DiskSizeGB: 120,
					DiskType:   "pd-ssd",
				},
			},
		}

	case vsphere.Name:
		installConfig.Platform = types.Platform{
			VSphere: &vsphere.Platform{
				VCenters: []vsphere.VCenter{
					{
						Server:      "vcenter.example.com",
						Port:        443,
						Username:    "administrator@vsphere.local",
						Password:    "password",
						Datacenters: []string{"datacenter-1"},
					},
				},
				FailureDomains: []vsphere.FailureDomain{
					{
						Name:   "failure-domain-1",
						Region: "region-1",
						Zone:   "zone-1",
						Server: "vcenter.example.com",
						Topology: vsphere.Topology{
							Datacenter:     "datacenter-1",
							ComputeCluster: "cluster-1",
							Networks:       []string{"network-1"},
							Datastore:      "datastore-1",
						},
					},
				},
				APIVIPs:     []string{"10.0.0.10"},
				IngressVIPs: []string{"10.0.0.11"},
			},
		}
		// Set vSphere-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			VSphere: &vsphere.MachinePool{
				NumCPUs:           4,
				NumCoresPerSocket: 4,
				MemoryMiB:         16384,
				OSDisk: vsphere.OSDisk{
					DiskSizeGB: 120,
				},
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			VSphere: &vsphere.MachinePool{
				NumCPUs:           2,
				NumCoresPerSocket: 2,
				MemoryMiB:         8192,
				OSDisk: vsphere.OSDisk{
					DiskSizeGB: 120,
				},
			},
		}

	case baremetal.Name:
		installConfig.Platform = types.Platform{
			BareMetal: &baremetal.Platform{
				LibvirtURI:              "qemu:///system",
				ClusterProvisioningIP:   "172.22.0.3",
				BootstrapProvisioningIP: "172.22.0.2",
				ExternalBridge:          "baremetal",
				ProvisioningBridge:      "provisioning",
				Hosts: []*baremetal.Host{
					{
						Name:            "host-0",
						BMC:             baremetal.BMC{Address: "ipmi://192.168.111.1:6230", Username: "admin", Password: "password"},
						BootMACAddress:  "52:54:00:82:68:77",
						HardwareProfile: "default",
					},
				},
			},
		}
		// Set baremetal-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			BareMetal: &baremetal.MachinePool{},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			BareMetal: &baremetal.MachinePool{},
		}

	case openstack.Name:
		installConfig.Platform = types.Platform{
			OpenStack: &openstack.Platform{
				Cloud:           "openstack",
				ExternalNetwork: "external",
			},
		}
		// Set OpenStack-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			OpenStack: &openstack.MachinePool{
				FlavorName: "m1.xlarge",
				Zones:      []string{"zone-1"},
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			OpenStack: &openstack.MachinePool{
				FlavorName: "m1.large",
				Zones:      []string{"zone-1"},
			},
		}

	case ibmcloud.Name:
		installConfig.Platform = types.Platform{
			IBMCloud: &ibmcloud.Platform{
				Region: "us-south",
			},
		}
		// Set IBM Cloud-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			IBMCloud: &ibmcloud.MachinePool{
				InstanceType: "bx2-4x16",
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			IBMCloud: &ibmcloud.MachinePool{
				InstanceType: "bx2-2x8",
			},
		}

	case powervs.Name:
		installConfig.Platform = types.Platform{
			PowerVS: &powervs.Platform{
				Region:              "us-south",
				Zone:                "us-south-1",
				ServiceInstanceGUID: "00000000-0000-0000-0000-000000000000",
			},
		}
		// Set PowerVS-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			PowerVS: &powervs.MachinePool{
				SysType:    "s922",
				MemoryGiB:  32,
				Processors: intstr.FromInt(1),
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			PowerVS: &powervs.MachinePool{
				SysType:    "s922",
				MemoryGiB:  16,
				Processors: intstr.FromInt(1),
			},
		}

	case nutanix.Name:
		installConfig.Platform = types.Platform{
			Nutanix: &nutanix.Platform{
				PrismCentral: nutanix.PrismCentral{
					Endpoint: nutanix.PrismEndpoint{
						Address: "prism.example.com",
						Port:    9440,
					},
					Username: "admin",
					Password: "password",
				},
				PrismElements: []nutanix.PrismElement{
					{
						Endpoint: nutanix.PrismEndpoint{
							Address: "pe.example.com",
							Port:    9440,
						},
						UUID: "00000000-0000-0000-0000-000000000000",
						Name: "pe-1",
					},
				},
			},
		}
		// Set Nutanix-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			Nutanix: &nutanix.MachinePool{
				NumCPUs:           4,
				NumCoresPerSocket: 1,
				MemoryMiB:         16384,
				OSDisk: nutanix.OSDisk{
					DiskSizeGiB: 120,
				},
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			Nutanix: &nutanix.MachinePool{
				NumCPUs:           2,
				NumCoresPerSocket: 1,
				MemoryMiB:         8192,
				OSDisk: nutanix.OSDisk{
					DiskSizeGiB: 120,
				},
			},
		}

	case ovirt.Name:
		installConfig.Platform = types.Platform{
			Ovirt: &ovirt.Platform{
				ClusterID:       "00000000-0000-0000-0000-000000000000",
				StorageDomainID: "00000000-0000-0000-0000-000000000000",
				NetworkName:     "ovirtmgmt",
			},
		}
		// Set oVirt-specific machine pool
		installConfig.ControlPlane.Platform = types.MachinePoolPlatform{
			Ovirt: &ovirt.MachinePool{
				CPU: &ovirt.CPU{
					Cores:   4,
					Sockets: 1,
					Threads: 1,
				},
				MemoryMB: 16384,
				OSDisk: &ovirt.Disk{
					SizeGB: 120,
				},
				VMType: ovirt.VMTypeServer,
			},
		}
		installConfig.Compute[0].Platform = types.MachinePoolPlatform{
			Ovirt: &ovirt.MachinePool{
				CPU: &ovirt.CPU{
					Cores:   2,
					Sockets: 1,
					Threads: 1,
				},
				MemoryMB: 8192,
				OSDisk: &ovirt.Disk{
					SizeGB: 120,
				},
				VMType: ovirt.VMTypeServer,
			},
		}

	case none.Name:
		installConfig.Platform = types.Platform{
			None: &none.Platform{},
		}
		// No platform-specific machine pool for none platform

	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}

	// Remove deprecated fields by creating a clean copy
	cleanConfig := &types.InstallConfig{
		TypeMeta:                    installConfig.TypeMeta,
		ObjectMeta:                  installConfig.ObjectMeta,
		AdditionalTrustBundle:       installConfig.AdditionalTrustBundle,
		AdditionalTrustBundlePolicy: installConfig.AdditionalTrustBundlePolicy,
		SSHKey:                      installConfig.SSHKey,
		BaseDomain:                  installConfig.BaseDomain,
		Networking: &types.Networking{
			NetworkType:         installConfig.Networking.NetworkType,
			MachineNetwork:      installConfig.Networking.MachineNetwork,
			ClusterNetwork:      installConfig.Networking.ClusterNetwork,
			ServiceNetwork:      installConfig.Networking.ServiceNetwork,
			ClusterNetworkMTU:   installConfig.Networking.ClusterNetworkMTU,
			OVNKubernetesConfig: installConfig.Networking.OVNKubernetesConfig,
		},
		ControlPlane:               installConfig.ControlPlane,
		Arbiter:                    installConfig.Arbiter,
		Compute:                    installConfig.Compute,
		Platform:                   installConfig.Platform,
		PullSecret:                 installConfig.PullSecret,
		Proxy:                      installConfig.Proxy,
		ImageDigestSources:         installConfig.ImageDigestSources,
		Publish:                    installConfig.Publish,
		OperatorPublishingStrategy: installConfig.OperatorPublishingStrategy,
		FIPS:                       installConfig.FIPS,
		CPUPartitioning:            installConfig.CPUPartitioning,
		CredentialsMode:            installConfig.CredentialsMode,
		BootstrapInPlace:           installConfig.BootstrapInPlace,
		Capabilities:               installConfig.Capabilities,
		FeatureSet:                 installConfig.FeatureSet,
		FeatureGates:               installConfig.FeatureGates,
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(cleanConfig, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal install config to JSON: %w", err)
	}

	return string(jsonData), nil
}

// MergeInstallConfigPlatform merges an install config with a platform struct at InstallConfig.Platform
// Inputs:
// - installConfigStr: string containing install config as either JSON or YAML
// - platformStr: string containing platform struct as either JSON or YAML
// - outputFormat: "json" or "yaml" to specify output format (defaults to "yaml" if empty or unsupported)
// Outputs:
// - string containing the merged install config in the specified format
func MergeInstallConfigPlatform(installConfigStr, platformStr, outputFormat string) (string, error) {
	logrus.Info("Merging install config with platform configuration")

	// Parse install config
	var installConfig types.InstallConfig
	if err := parseConfig(installConfigStr, &installConfig); err != nil {
		return "", fmt.Errorf("failed to parse install config: %w", err)
	}

	// Parse platform
	var platform types.Platform
	if err := parseConfig(platformStr, &platform); err != nil {
		return "", fmt.Errorf("failed to parse platform: %w", err)
	}

	// Merge platform into install config
	installConfig.Platform = platform

	// Default to yaml if outputFormat is empty or unsupported
	if outputFormat == "" {
		outputFormat = "yaml"
	}

	// Output in requested format
	switch strings.ToLower(outputFormat) {
	case "json":
		jsonData, err := json.MarshalIndent(installConfig, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal install config to JSON: %w", err)
		}
		return string(jsonData), nil
	case "yaml":
		yamlData, err := yaml.Marshal(installConfig)
		if err != nil {
			return "", fmt.Errorf("failed to marshal install config to YAML: %w", err)
		}
		return string(yamlData), nil
	default:
		// Default to YAML for unsupported formats
		logrus.Warnf("Unsupported output format: %s, defaulting to YAML", outputFormat)
		yamlData, err := yaml.Marshal(installConfig)
		if err != nil {
			return "", fmt.Errorf("failed to marshal install config to YAML: %w", err)
		}
		return string(yamlData), nil
	}
}

// parseConfig attempts to parse a string as either JSON or YAML into the target struct
func parseConfig(configStr string, target interface{}) error {
	// Try JSON first
	if err := json.Unmarshal([]byte(configStr), target); err == nil {
		return nil
	}

	// Try YAML if JSON fails
	if err := yaml.Unmarshal([]byte(configStr), target); err != nil {
		return fmt.Errorf("failed to parse as JSON or YAML: %w", err)
	}

	return nil
}
