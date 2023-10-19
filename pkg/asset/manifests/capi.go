package manifests

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/manifests/internal/cidr"
	"github.com/openshift/installer/pkg/asset/openshiftinstall"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/ipnet"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

const (
	infraManifestDir    = "infra"
	capiGuestsNamespace = "openshift-cluster-api-guests"
)

var _ asset.WritableAsset = (*ClusterAPI)(nil)

// ClusterAPI generates manifests for target cluster
// creation using CAPI.
type ClusterAPI struct {
	FileList  []*asset.File
	Manifests []Manifest `json:"-"`
}

// Manifest is a wrapper for a CAPI manifest.
type Manifest struct {
	Object   client.Object
	filename string
}

// Name returns a human friendly name for the operator.
func (c *ClusterAPI) Name() string {
	return "ClusterAPI Manifests"
}

// Dependencies returns all of the dependencies directly needed by the
// ClusterAPI asset.
func (c *ClusterAPI) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&installconfig.ClusterID{},
		&openshiftinstall.Config{},
		&FeatureGate{},
		new(rhcos.Image),
	}
}

// Generate generates the respective operator config.yml files.
func (c *ClusterAPI) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	openshiftInstall := &openshiftinstall.Config{}
	featureGate := &FeatureGate{}
	rhcosImage := new(rhcos.Image)
	dependencies.Get(installConfig, clusterID, openshiftInstall, featureGate, rhcosImage)

	c.FileList = []*asset.File{}
	c.Manifests = []Manifest{}

	platform := installConfig.Config.Platform.Name()

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiGuestsNamespace,
		},
		Spec: clusterv1.ClusterSpec{
			InfrastructureRef: &v1.ObjectReference{
				Namespace: capiGuestsNamespace,
				Name:      clusterID.InfraID,
			},
		},
	}

	// Retrieve the AZs available and generate the subnets private
	// and public.
	mainCIDR := ipnet.MustParseCIDR("10.0.0.0/16")
	if len(installConfig.Config.MachineNetwork) > 0 {
		mainCIDR = &installConfig.Config.MachineNetwork[0].CIDR
	}

	switch platform {
	case awstypes.Name:
		// Not sure if this is the best place to create IAM roles.
		if err := aws.PutIAMRoles(clusterID.InfraID, installConfig); err != nil {
			return errors.Wrap(err, "failed to create IAM roles")
		}

		zones, err := installConfig.AWS.AvailabilityZones(context.TODO())
		if err != nil {
			return errors.Wrap(err, "failed to get availability zones")
		}

		awsCluster := &capa.AWSCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterID.InfraID,
				Namespace: capiGuestsNamespace,
			},
			Spec: capa.AWSClusterSpec{
				Region: installConfig.Config.AWS.Region,
				NetworkSpec: capa.NetworkSpec{
					VPC: capa.VPCSpec{
						CidrBlock:                  mainCIDR.String(),
						AvailabilityZoneUsageLimit: ptr.To(len(zones)),
						AvailabilityZoneSelection:  &capa.AZSelectionSchemeOrdered,
					},
					CNI: &capa.CNISpec{
						CNIIngressRules: capa.CNIIngressRules{
							{
								Description: "ICMP",
								Protocol:    capa.SecurityGroupProtocolICMP,
								FromPort:    -1,
								ToPort:      -1,
							},
							{
								Description: "Port 22 (TCP)",
								Protocol:    capa.SecurityGroupProtocolTCP,
								FromPort:    22,
								ToPort:      22,
							},
							{
								Description: "Port 4789 (UDP) for VXLAN",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    4789,
								ToPort:      4789,
							},
							{
								Description: "Port 6081 (UDP) for geneve",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    6081,
								ToPort:      6081,
							},
							{
								Description: "Port 500 (UDP) for IKE",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    500,
								ToPort:      500,
							},
							{
								Description: "Port 4500 (UDP) for IKE NAT",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    4500,
								ToPort:      4500,
							},
							{
								Description: "ESP",
								Protocol:    capa.SecurityGroupProtocolESP,
								FromPort:    -1,
								ToPort:      -1,
							},
							{
								Description: "Port 6441-6442 (TCP) for ovndb",
								Protocol:    capa.SecurityGroupProtocolTCP,
								FromPort:    6441,
								ToPort:      6442,
							},
							{
								Description: "Port 9000-9999 for node ports (TCP)",
								Protocol:    capa.SecurityGroupProtocolTCP,
								FromPort:    9000,
								ToPort:      9999,
							},
							{
								Description: "Port 9000-9999 for node ports (UDP)",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    9000,
								ToPort:      9999,
							},
							{
								Description: "Service node ports (TCP)",
								Protocol:    capa.SecurityGroupProtocolTCP,
								FromPort:    30000,
								ToPort:      32767,
							},
							{
								Description: "Service node ports (UDP)",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    30000,
								ToPort:      32767,
							},
						},
					},
					AdditionalControlPlaneIngressRules: []capa.IngressRule{
						{
							Description:              "MCS traffic from cluster network",
							Protocol:                 capa.SecurityGroupProtocolTCP,
							FromPort:                 22623,
							ToPort:                   22623,
							SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane"},
						},
						{
							Description:              "controller-manager",
							Protocol:                 capa.SecurityGroupProtocolTCP,
							FromPort:                 10257,
							ToPort:                   10257,
							SourceSecurityGroupRoles: []capa.SecurityGroupRole{"controlplane", "node"},
						},
						{
							Description:              "kube-scheduler",
							Protocol:                 capa.SecurityGroupProtocolTCP,
							FromPort:                 10259,
							ToPort:                   10259,
							SourceSecurityGroupRoles: []capa.SecurityGroupRole{"controlplane", "node"},
						},
						{
							Description: "SSH everyone",
							Protocol:    capa.SecurityGroupProtocolTCP,
							FromPort:    22,
							ToPort:      22,
							CidrBlocks:  []string{"0.0.0.0/0"},
						},
					},
				},
				S3Bucket: &capa.S3Bucket{
					Name:                 fmt.Sprintf("openshift-bootstrap-data-%s", clusterID.InfraID),
					PresignedURLDuration: &metav1.Duration{Duration: 1 * time.Hour},
				},
				ControlPlaneLoadBalancer: &capa.AWSLoadBalancerSpec{
					Name:             ptr.To(clusterID.InfraID + "-ext"),
					LoadBalancerType: capa.LoadBalancerTypeNLB,
					Scheme:           &capa.ELBSchemeInternetFacing,
					AdditionalListeners: []*capa.AdditionalListenerSpec{
						{
							Port:     22623,
							Protocol: capa.ELBProtocolTCP,
						},
					},
				},
			},
		}

		// If the install config has subnets, use them.
		if len(installConfig.AWS.Subnets) > 0 {
			privateSubnets, err := installConfig.AWS.PrivateSubnets(context.TODO())
			if err != nil {
				return errors.Wrap(err, "failed to get private subnets")
			}
			for _, subnet := range privateSubnets {
				awsCluster.Spec.NetworkSpec.Subnets = append(awsCluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
					ID:               subnet.ID,
					CidrBlock:        subnet.CIDR,
					AvailabilityZone: subnet.Zone.Name,
					IsPublic:         subnet.Public,
				})
			}
			publicSubnets, err := installConfig.AWS.PublicSubnets(context.TODO())
			if err != nil {
				return errors.Wrap(err, "failed to get public subnets")
			}

			for _, subnet := range publicSubnets {
				awsCluster.Spec.NetworkSpec.Subnets = append(awsCluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
					ID:               subnet.ID,
					CidrBlock:        subnet.CIDR,
					AvailabilityZone: subnet.Zone.Name,
					IsPublic:         subnet.Public,
				})
			}

			vpc, err := installConfig.AWS.VPC(context.TODO())
			if err != nil {
				return errors.Wrap(err, "failed to get VPC")
			}
			awsCluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
				ID: vpc,
			}
		}

		awsClusterFn := "01_aws-cluster.yaml"
		c.Manifests = append(c.Manifests, Manifest{awsCluster, awsClusterFn})

		cluster.Spec.InfrastructureRef.APIVersion = "infrastructure.cluster.x-k8s.io/v1beta2"
		cluster.Spec.InfrastructureRef.Kind = "AWSCluster"

		id := &capa.AWSClusterControllerIdentity{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
			Spec: capa.AWSClusterControllerIdentitySpec{
				AWSClusterIdentitySpec: capa.AWSClusterIdentitySpec{
					AllowedNamespaces: &capa.AllowedNamespaces{}, // Allow all namespaces.
				},
			},
		}
		idFn := "00_aws-cluster-controller-identity-default.yaml"
		c.Manifests = append(c.Manifests, Manifest{id, idFn})
	case "azure":
		session, err := installConfig.Azure.Session()
		if err != nil {
			return errors.Wrap(err, "failed to create Azure session")
		}

		subnets, err := cidr.SplitIntoSubnetsIPv4(mainCIDR.String(), 2)
		if err != nil {
			return errors.Wrap(err, "failed to split CIDR into subnets")
		}

		// CAPZ expects the capz-system to be created.
		azureNamespace := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "capz-system"}}
		azureNamespaceFn := "00_azure-namespace.yaml"
		c.Manifests = append(c.Manifests, Manifest{azureNamespace, azureNamespaceFn})

		azureCluster := &capz.AzureCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterID.InfraID,
				Namespace: capiGuestsNamespace,
			},
			Spec: capz.AzureClusterSpec{
				ResourceGroup: clusterID.InfraID,
				AzureClusterClassSpec: capz.AzureClusterClassSpec{
					SubscriptionID:   session.Credentials.SubscriptionID,
					Location:         installConfig.Config.Azure.Region,
					AzureEnvironment: string(installConfig.Azure.CloudName),
					IdentityRef: &v1.ObjectReference{
						APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
						Kind:       "AzureClusterIdentity",
						Name:       clusterID.InfraID,
					},
				},
				NetworkSpec: capz.NetworkSpec{
					Vnet: capz.VnetSpec{
						ID: installConfig.Config.Azure.VirtualNetwork,
						VnetClassSpec: capz.VnetClassSpec{
							CIDRBlocks: []string{
								mainCIDR.String(),
							},
						},
					},
					Subnets: capz.Subnets{
						{
							SubnetClassSpec: capz.SubnetClassSpec{
								Name: "control-plane-subnet",
								Role: capz.SubnetControlPlane,
								CIDRBlocks: []string{
									subnets[0].String(),
								},
							},
						},
						{
							SubnetClassSpec: capz.SubnetClassSpec{
								Name: "worker-subnet",
								Role: capz.SubnetNode,
								CIDRBlocks: []string{
									subnets[1].String(),
								},
							},
						},
					},
				},
			},
		}

		azureClusterFn := "01_azure-cluster.yaml"
		c.Manifests = append(c.Manifests, Manifest{azureCluster, azureClusterFn})

		cluster.Spec.InfrastructureRef.APIVersion = "infrastructure.cluster.x-k8s.io/v1beta1"
		cluster.Spec.InfrastructureRef.Kind = "AzureCluster"

		azureClientSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterID.InfraID + "-azure-client-secret",
				Namespace: capiGuestsNamespace,
			},
			StringData: map[string]string{
				"clientSecret": session.Credentials.ClientSecret,
			},
		}
		azureClientSecretFn := "00_azure-client-secret.yaml"
		c.Manifests = append(c.Manifests, Manifest{azureClientSecret, azureClientSecretFn})

		id := &capz.AzureClusterIdentity{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterID.InfraID,
			},
			Spec: capz.AzureClusterIdentitySpec{
				Type:              capz.ManualServicePrincipal,
				AllowedNamespaces: &capz.AllowedNamespaces{}, // Allow all namespaces.
				ClientID:          session.Credentials.ClientID,
				ClientSecret: v1.SecretReference{
					Name:      azureClientSecret.Name,
					Namespace: azureClientSecret.Namespace,
				},
				TenantID: session.Credentials.TenantID,
			},
		}
		idFn := "00_aws-cluster-controller-identity-default.yaml"
		c.Manifests = append(c.Manifests, Manifest{id, idFn})
	default:
		return nil
	}

	// Create the infrastructure manifest.
	clusterFn := "01-capi-cluster.yaml"
	c.Manifests = append(c.Manifests, Manifest{cluster, clusterFn})

	for _, m := range c.Manifests {
		objData, err := yaml.Marshal(m.Object)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create infrastructure manifest %s from InstallConfig", m.filename)
			return errors.Wrapf(err, errMsg)
		}

		c.FileList = append(c.FileList, &asset.File{
			Filename: filepath.Join(infraManifestDir, m.filename),
			Data:     objData,
		})
	}

	asset.SortFiles(c.FileList)
	return nil
}

// Files returns the files generated by the asset.
func (c *ClusterAPI) Files() []*asset.File {
	return c.FileList
}

// Load returns the openshift asset from disk.
func (c *ClusterAPI) Load(f asset.FileFetcher) (bool, error) {
	// yamlFileList, err := f.FetchByPattern(filepath.Join(infraManifestDir, "*.yaml"))
	// if err != nil {
	// 	return false, errors.Wrap(err, "failed to load *.yaml files")
	// }
	// ymlFileList, err := f.FetchByPattern(filepath.Join(infraManifestDir, "*.yml"))
	// if err != nil {
	// 	return false, errors.Wrap(err, "failed to load *.yml files")
	// }
	// jsonFileList, err := f.FetchByPattern(filepath.Join(infraManifestDir, "*.json"))
	// if err != nil {
	// 	return false, errors.Wrap(err, "failed to load *.json files")
	// }
	// fileList := append(yamlFileList, ymlFileList...)
	// fileList = append(fileList, jsonFileList...)

	// for _, file := range fileList {
	// 	c.FileList = append(c.FileList, file)
	// }

	// asset.SortFiles(c.FileList)
	// return len(c.FileList) > 0, nil
	return false, nil
}
