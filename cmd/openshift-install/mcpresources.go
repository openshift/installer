package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/mcpserver"
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

// ok, today resources are a no-go. there are very few mcp clients that actually support
// this and the ones that work, don't. Maybe later...
// https://modelcontextprotocol.io/clients

func ResourceTemplates() []mcpserver.ServerResourceTemplate {
	return []mcpserver.ServerResourceTemplate{
		{
			/*
				Handler: func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
					logrus.Info("in Handler")
					result, err := mcpserver.ProcessResourceResults(getInstallConfigResource(ctx, req))
					return result.Contents, err
				},

				ResourceTemplate: mcp.NewResourceTemplate("install://{platform}/config",
					"install-config",
					mcp.WithTemplateMIMEType("application/json"),
					mcp.WithTemplateDescription("provides an example install-config by platform")),

			*/
		},
	}
}
func Resources() []server.ServerResource {

	return []server.ServerResource{
		{
			/*
				Resource: mcp.NewResource(
					"config://{platform}",
					"install-config",
					mcp.WithResourceDescription("provides an example install-config by platform"),
					mcp.WithMIMEType("application/json"),
				),
				Handler: func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
					logrus.Info("in Handler")
					result, err := mcpserver.ProcessResourceResults(getInstallConfigResource(ctx, req))
					return result.Contents, err
				},
			*/
		},
	}
}

func getInstallConfigResource(ctx context.Context, req mcp.ReadResourceRequest) (string, error) {
	logrus.Info("in getInstallConfigResource")
	platform := req.Params.Arguments["platform"]
	logrus.Infof("platform: %s", platform)

	// Create a new InstallConfig
	installConfig := &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "example-cluster",
		},
		BaseDomain: "example.com",
		PullSecret: "{\"auths\":{\"fake-registry.example.com\":{\"auth\":\"ZmFrZS1hdXRoLXRva2Vu\"}}}",
		SSHKey:     "ssh-rsa AAAA...",
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
						Server:   "vcenter.example.com",
						Username: "administrator@vsphere.local",
						Password: "password",
					},
				},
				FailureDomains: []vsphere.FailureDomain{
					{
						Name:   "failure-domain-1",
						Region: "region-1",
						Zone:   "zone-1",
						Topology: vsphere.Topology{
							Datacenter:     "datacenter-1",
							ComputeCluster: "cluster-1",
							Networks:       []string{"network-1"},
							Datastore:      "datastore-1",
						},
					},
				},
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
