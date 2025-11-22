// Package aws generates Machine objects for aws.
package aws

import (
	"bytes"
	"context"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/vincent-petithory/dataurl"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// MachineInput defines the inputs needed to generate a machine asset.
type MachineInput struct {
	Role           string
	Pool           *types.MachinePool
	Subnets        aws.SubnetsByZone
	Tags           capa.Tags
	PublicIP       bool
	PublicIpv4Pool string
	Ignition       *capa.Ignition
}

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(clusterID string, in *MachineInput) ([]*asset.RuntimeFile, error) {
	if poolPlatform := in.Pool.Platform.Name(); poolPlatform != awstypes.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	mpool := in.Pool.Platform.AWS

	total := int64(1)
	if in.Pool.Replicas != nil {
		total = *in.Pool.Replicas
	}

	imds := capa.HTTPTokensStateOptional
	if mpool.EC2Metadata.Authentication == "Required" {
		imds = capa.HTTPTokensStateRequired
	}

	instanceProfile := in.Pool.Platform.AWS.IAMProfile
	if len(instanceProfile) == 0 {
		instanceProfile = fmt.Sprintf("%s-master-profile", clusterID)
	}

	var result []*asset.RuntimeFile

	for idx := int64(0); idx < total; idx++ {
		subnet := &capa.AWSResourceReference{}
		zone := mpool.Zones[int(idx)%len(mpool.Zones)]

		// BYO VPC deployments when subnet IDs are set on install-config.yaml
		if len(in.Subnets) > 0 {
			subnetMeta, ok := in.Subnets[zone]
			if len(in.Subnets) > 0 && !ok {
				return nil, fmt.Errorf("no subnet for zone %s", zone)
			}
			if subnetMeta.ID == "" {
				return nil, fmt.Errorf("invalid subnet ID for zone %s", zone)
			}
			subnet.ID = ptr.To(subnetMeta.ID)
		} else {
			subnetInternetScope := "private"
			if in.PublicIP {
				subnetInternetScope = "public"
			}
			subnet.Filters = []capa.Filter{
				{
					Name:   "tag:Name",
					Values: []string{fmt.Sprintf("%s-subnet-%s-%s", clusterID, subnetInternetScope, zone)},
				},
			}
		}

		awsMachine := &capa.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%d", clusterID, in.Pool.Name, idx),
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capa.AWSMachineSpec{
				Ignition:             in.Ignition,
				UncompressedUserData: ptr.To(true),
				InstanceType:         mpool.InstanceType,
				AMI:                  capa.AMIReference{ID: ptr.To(mpool.AMIID)},
				SSHKeyName:           ptr.To(""),
				IAMInstanceProfile:   instanceProfile,
				Subnet:               subnet,
				PublicIP:             ptr.To(in.PublicIP),
				AdditionalTags:       in.Tags,
				RootVolume: &capa.Volume{
					Size:          int64(mpool.EC2RootVolume.Size),
					Type:          capa.VolumeType(mpool.EC2RootVolume.Type),
					IOPS:          int64(mpool.EC2RootVolume.IOPS),
					Throughput:    ptr.To(mpool.EC2RootVolume.Throughput),
					Encrypted:     ptr.To(true),
					EncryptionKey: mpool.KMSKeyARN,
				},
				InstanceMetadataOptions: &capa.InstanceMetadataOptions{
					HTTPTokens:   imds,
					HTTPEndpoint: capa.InstanceMetadataEndpointStateEnabled,
				},
			},
		}
		awsMachine.SetGroupVersionKind(capa.GroupVersion.WithKind("AWSMachine"))

		if in.Role == "bootstrap" {
			awsMachine.Name = capiutils.GenerateBoostrapMachineName(clusterID)
			awsMachine.Labels["install.openshift.io/bootstrap"] = ""

			// Enable BYO Public IPv4 Pool when defined on install-config.yaml.
			if len(in.PublicIpv4Pool) > 0 {
				awsMachine.Spec.ElasticIPPool = &capa.ElasticIPPool{
					PublicIpv4Pool:              ptr.To(in.PublicIpv4Pool),
					PublicIpv4PoolFallBackOrder: ptr.To(capa.PublicIpv4PoolFallbackOrderAmazonPool),
				}
			}
		}

		// Handle additional security groups.
		for _, sg := range mpool.AdditionalSecurityGroupIDs {
			awsMachine.Spec.AdditionalSecurityGroups = append(
				awsMachine.Spec.AdditionalSecurityGroups,
				capa.AWSResourceReference{ID: ptr.To(sg)},
			)
		}

		if mpool.CPUOptions != nil {
			cpuOptions := capa.CPUOptions{}

			if mpool.CPUOptions.ConfidentialCompute != nil {
				cpuOptions.ConfidentialCompute = capa.AWSConfidentialComputePolicy(*mpool.CPUOptions.ConfidentialCompute)
			}

			awsMachine.Spec.CPUOptions = cpuOptions
		}

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", awsMachine.Name)},
			Object: awsMachine,
		})

		machine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: awsMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterID,
				Bootstrap: capi.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, in.Role)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: capa.GroupVersion.String(),
					Kind:       "AWSMachine",
					Name:       awsMachine.Name,
				},
			},
		}
		machine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", machine.Name)},
			Object: machine,
		})
	}
	return result, nil
}

// CapaTagsFromUserTags converts a map of user tags to a map of capa.Tags.
func CapaTagsFromUserTags(clusterID string, usertags map[string]string) (capa.Tags, error) {
	tags := capa.Tags{}
	tags[fmt.Sprintf("kubernetes.io/cluster/%s", clusterID)] = "owned"

	forbiddenTags := sets.New[string]()
	for key := range tags {
		forbiddenTags.Insert(key)
	}

	userTagKeys := sets.New[string]()
	for key := range usertags {
		userTagKeys.Insert(key)
	}

	if clobberedTags := userTagKeys.Intersection(forbiddenTags); clobberedTags.Len() > 0 {
		return nil, fmt.Errorf("user tag keys %v are not allowed", sets.List(clobberedTags))
	}

	for _, k := range sets.List(userTagKeys) {
		tags[k] = usertags[k]
	}
	return tags, nil
}

// CapaIgnitionWithCertBundleAndProxy generates CAPA ignition config with cert and proxy information.
func CapaIgnitionWithCertBundleAndProxy(userCA string, proxy *types.Proxy) (*capa.Ignition, error) {
	carefs, err := parseCertificateBundle([]byte(userCA))
	if err != nil {
		return nil, err
	}
	return &capa.Ignition{
		Version: "3.2",
		TLS: &capa.IgnitionTLS{
			CASources: carefs,
		},
		Proxy: capaIgnitionProxy(proxy),
	}, nil
}

// TODO: try to share this code with ignition.bootstrap package?
// parseCertificateBundle loads each certificate in the bundle to the CAPA
// carrier type, ignoring any invisible character before, after and in between
// certificates.
func parseCertificateBundle(userCA []byte) ([]capa.IgnitionCASource, error) {
	userCA = bytes.TrimSpace(userCA)

	var carefs []capa.IgnitionCASource
	for len(userCA) > 0 {
		var block *pem.Block
		block, userCA = pem.Decode(userCA)
		if block == nil {
			return nil, fmt.Errorf("unable to parse certificate, please check the certificates")
		}

		carefs = append(carefs, capa.IgnitionCASource(dataurl.EncodeBytes(pem.EncodeToMemory(block))))

		userCA = bytes.TrimSpace(userCA)
	}

	return carefs, nil
}

func capaIgnitionProxy(proxy *types.Proxy) *capa.IgnitionProxy {
	capaProxy := &capa.IgnitionProxy{}
	if proxy == nil {
		return capaProxy
	}
	if httpProxy := proxy.HTTPProxy; httpProxy != "" {
		capaProxy.HTTPProxy = &httpProxy
	}
	if httpsProxy := proxy.HTTPSProxy; httpsProxy != "" {
		capaProxy.HTTPSProxy = &httpsProxy
	}
	capaProxy.NoProxy = make([]capa.IgnitionNoProxy, 0, len(proxy.NoProxy))
	if noProxy := proxy.NoProxy; noProxy != "" {
		noProxySplit := strings.Split(noProxy, ",")
		for _, p := range noProxySplit {
			capaProxy.NoProxy = append(capaProxy.NoProxy, capa.IgnitionNoProxy(p))
		}
	}
	return capaProxy
}

// MachineSubnetsByZones returns a map of subnets by zones for a subnet role type
// to place corresponding nodes in. Type must be one of: ClusterNode, EdgeNode and BootstrapNode.
func MachineSubnetsByZones(ctx context.Context, ic *installconfig.InstallConfig, roleType awstypes.SubnetRoleType) (aws.SubnetsByZone, error) {
	machineSubnets := make(aws.SubnetsByZone)

	// If managed subnets, leave empty to use capa subnet filters.
	if ic.Config.AWS == nil || len(ic.Config.AWS.VPC.Subnets) == 0 {
		return machineSubnets, nil
	}

	subnetIDsByRole := make(map[awstypes.SubnetRoleType][]string)
	for _, subnet := range ic.Config.AWS.VPC.Subnets {
		for _, role := range subnet.Roles {
			subnetIDsByRole[role.Type] = append(subnetIDsByRole[role.Type], string(subnet.ID))
		}
	}

	// BYO-subnet install case and subnet roles are specified.
	if len(subnetIDsByRole) > 0 {
		for _, subnetID := range subnetIDsByRole[roleType] {
			subnetMeta, err := ic.AWS.SubnetByID(ctx, subnetID)
			if err != nil {
				return machineSubnets, err
			}
			machineSubnets[subnetMeta.Zone.Name] = subnetMeta
		}
		return machineSubnets, nil
	}

	// BYO-subnet install case and subnet roles are not specified.
	var subnetMeta aws.Subnets
	var err error
	switch roleType {
	case awstypes.ClusterNodeSubnetRole:
		// fetch private subnets to master nodes.
		subnetMeta, err = ic.AWS.PrivateSubnets(ctx)
	case awstypes.EdgeNodeSubnetRole:
		// fetch edge subnets to edge compute nodes.
		subnetMeta, err = ic.AWS.EdgeSubnets(ctx)
	case awstypes.BootstrapNodeSubnetRole:
		// fetch public subnets for bootstrap, when exists (external cluster).
		// Otherwise, use private.
		if ic.Config.Publish == types.ExternalPublishingStrategy {
			subnetMeta, err = ic.AWS.PublicSubnets(ctx)
		} else {
			subnetMeta, err = ic.AWS.PrivateSubnets(ctx)
		}
	}

	if err != nil {
		return machineSubnets, err
	}

	for _, subnet := range subnetMeta {
		machineSubnets[subnet.Zone.Name] = subnet
	}

	return machineSubnets, nil
}
