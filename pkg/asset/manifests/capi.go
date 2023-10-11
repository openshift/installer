package manifests

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	awsmachines "github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/openshiftinstall"
	"github.com/openshift/installer/pkg/asset/rhcos"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

const (
	infraManifestDir    = "infra"
	capiGuestsNamespace = "openshift-cluster-api-guests"
)

var (
	_ asset.WritableAsset = (*ClusterAPI)(nil)
)

// ClusterAPI generates manifests for target cluster
// creation using CAPI.
type ClusterAPI struct {
	FileList  []*asset.File
	Manifests []Manifest `json:"-"`
}

type Manifest struct {
	Object   client.Object
	filename string
}

// Name returns a human friendly name for the operator
func (c *ClusterAPI) Name() string {
	return "ClusterAPI Manifests"
}

// Dependencies returns all of the dependencies directly needed by the
// ClusterAPI asset
func (c *ClusterAPI) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&installconfig.ClusterID{},
		&openshiftinstall.Config{},
		&FeatureGate{},
		new(rhcos.Image),
	}
}

// Generate generates the respective operator config.yml files
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
			InfrastructureRef: &v1.ObjectReference{},
		},
	}

	region := installConfig.Config.Platform.AWS.Region

	tags, err := awsmachines.CapaTagsFromUserTags(clusterID.InfraID, installConfig.Config.AWS.UserTags)
	if err != nil {
		return errors.Wrap(err, "error in user-provided tags")
	}

	switch platform {
	case awstypes.Name:

		// Not sure if this is the best place to create IAM roles.
		if err := awsmachines.PutIAMRoles(clusterID.InfraID, installConfig); err != nil {
			return errors.Wrap(err, "failed to create IAM roles")
		}

		awsCluster := &capa.AWSCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterID.InfraID,
				Namespace: capiGuestsNamespace,
			},
			Spec: capa.AWSClusterSpec{
				AdditionalTags: tags,
				Region:         installConfig.Config.Platform.AWS.Region,
				NetworkSpec: capa.NetworkSpec{
					VPC: capa.VPCSpec{
						AvailabilityZoneUsageLimit: pointer.Int(6),
						AvailabilityZoneSelection:  &capa.AZSelectionSchemeOrdered,
					},
					Subnets: capa.Subnets{
						{
							ID:               clusterID.InfraID + "-private-" + region + "a",
							AvailabilityZone: region + "a",
							CidrBlock:        "10.0.0.0/19",
						},
						{
							ID:               clusterID.InfraID + "-private-" + region + "b",
							AvailabilityZone: region + "b",
							CidrBlock:        "10.0.32.0/19",
						},
						{
							ID:               clusterID.InfraID + "-private-" + region + "c",
							AvailabilityZone: region + "c",
							CidrBlock:        "10.0.64.0/19",
						},
						{
							ID:               clusterID.InfraID + "-public-" + region + "a",
							IsPublic:         true,
							AvailabilityZone: region + "a",
							CidrBlock:        "10.0.128.0/19",
						},
						{
							ID:               clusterID.InfraID + "-public-" + region + "b",
							IsPublic:         true,
							AvailabilityZone: region + "b",
							CidrBlock:        "10.0.160.0/19",
						},
						{
							ID:               clusterID.InfraID + "-public-" + region + "c",
							IsPublic:         true,
							AvailabilityZone: region + "c",
							CidrBlock:        "10.0.192.0/19",
						},
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
								Description: "Port 6441-6442 (TCP)",
								Protocol:    capa.SecurityGroupProtocolTCP,
								FromPort:    6441,
								ToPort:      6442,
							},
							{
								Description: "Port 6081 (UDP)",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    6081,
								ToPort:      6081,
							},
							{
								Description: "Port 500 (UDP)",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    500,
								ToPort:      500,
							},
							{
								Description: "Port 4789 (UDP)",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    4789,
								ToPort:      4789,
							},
							{
								Description: "Port 4500 (UDP)",
								Protocol:    capa.SecurityGroupProtocolUDP,
								FromPort:    4500,
								ToPort:      4500,
							},
							{
								Description: "Port 10257 (TCP)",
								Protocol:    capa.SecurityGroupProtocolTCP,
								FromPort:    10257,
								ToPort:      10257,
							},
							{
								Description: "Port 10259 (TCP)",
								Protocol:    capa.SecurityGroupProtocolTCP,
								FromPort:    10259,
								ToPort:      10259,
							},
							{
								Description: "Port 22 (TCP)",
								Protocol:    capa.SecurityGroupProtocolTCP,
								FromPort:    22,
								ToPort:      22,
							},
							{
								Description: "ESP",
								Protocol:    capa.SecurityGroupProtocolESP,
								FromPort:    -1,
								ToPort:      -1,
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
							Description:              "Kubelet traffic from nodes",
							Protocol:                 capa.SecurityGroupProtocolTCP,
							FromPort:                 10250,
							ToPort:                   10250,
							SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane"},
						},
						{
							Description:              "Service node ports (TCP)",
							Protocol:                 capa.SecurityGroupProtocolTCP,
							FromPort:                 30000,
							ToPort:                   32767,
							SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane"},
						},
						{
							Description:              "Service node ports (UDP)",
							Protocol:                 capa.SecurityGroupProtocolUDP,
							FromPort:                 30000,
							ToPort:                   32767,
							SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane"},
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
					Name:             pointer.String(clusterID.InfraID + "-ext"),
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
		awsClusterFn := "01_aws-cluster.yaml"
		c.Manifests = append(c.Manifests, Manifest{awsCluster, awsClusterFn})

		cluster.Spec.InfrastructureRef.APIVersion = "infrastructure.cluster.x-k8s.io/v1beta2"
		cluster.Spec.InfrastructureRef.Kind = "AWSCluster"
		cluster.Spec.InfrastructureRef.Name = clusterID.InfraID

		id := &capa.AWSClusterControllerIdentity{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
				Kind:       "AWSClusterControllerIdentity",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
			Spec: capa.AWSClusterControllerIdentitySpec{
				AWSClusterIdentitySpec: capa.AWSClusterIdentitySpec{
					AllowedNamespaces: &capa.AllowedNamespaces{
						// TODO: The godoc for this field indicates:
						// An nil or empty list indicates that AWSClusters cannot use the identity from any namespace.
						// Our internal notes say:
						// https://github.com/openshift-cloud-team/cluster-api-installer-poc/blob/main/templates/00_aws-cluster-controller-identity-default.yaml
						// allowedNamespaces: {}  # matches all namespaces
						// Check if this is a discrepency.
					},
				},
			},
		}
		idFn := "00_aws-cluster-controller-identity-default.yaml"
		c.Manifests = append(c.Manifests, Manifest{id, idFn})
	}

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
