package ibmcloud

import (
	"fmt"

	machineapi "github.com/openshift/api/machine/v1beta1"
	ibmcloudprovider "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != ibmcloud.Name {
		return nil, fmt.Errorf("non-IBMCloud configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != ibmcloud.Name {
		return nil, fmt.Errorf("non-IBMCloud machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.IBMCloud
	mpool := pool.Platform.IBMCloud
	azs := mpool.Zones

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		azIndex := int(idx) % len(azs)
		provider, err := provider(clusterID, platform, mpool, azIndex, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clusterID,
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
				},
			},
			Spec: machineapi.MachineSpec{
				ProviderSpec: machineapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
			},
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterID string,
	platform *ibmcloud.Platform,
	mpool *ibmcloud.MachinePool,
	azIdx int,
	role string,
	userDataSecret string,
) (*ibmcloudprovider.IBMCloudMachineProviderSpec, error) {
	az := mpool.Zones[azIdx]

	var vpc string
	if platform.VPC != "" {
		vpc = platform.VPC
	} else {
		vpc = fmt.Sprintf("%s-vpc", clusterID)
	}

	var resourceGroup string
	if platform.ResourceGroupName != "" {
		resourceGroup = platform.ResourceGroupName
	} else {
		resourceGroup = clusterID
	}

	subnet, err := getSubnetName(clusterID, role, az)
	if err != nil {
		return nil, err
	}

	securityGroups, err := getSecurityGroupNames(clusterID, role)
	if err != nil {
		return nil, err
	}

	var dedicatedHost string
	if len(mpool.DedicatedHosts) == len(mpool.Zones) {
		if mpool.DedicatedHosts[azIdx].Name != "" {
			dedicatedHost = mpool.DedicatedHosts[azIdx].Name
		} else {
			dedicatedHost, err = getDedicatedHostNameForZone(clusterID, role, az)
			if err != nil {
				return nil, err
			}
		}
	}

	return &ibmcloudprovider.IBMCloudMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ibmcloudproviderconfig.openshift.io/v1beta1",
			Kind:       "IBMCloudMachineProviderSpec",
		},
		VPC:           vpc,
		DedicatedHost: dedicatedHost,
		Tags:          []ibmcloudprovider.TagSpecs{},
		Image:         fmt.Sprintf("%s-rhcos", clusterID),
		Profile:       mpool.InstanceType,
		Region:        platform.Region,
		ResourceGroup: resourceGroup,
		Zone:          az,
		PrimaryNetworkInterface: ibmcloudprovider.NetworkInterface{
			Subnet:         subnet,
			SecurityGroups: securityGroups,
		},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "ibmcloud-credentials"},
		// TODO: IBM: Boot volume encryption key
	}, nil
}

func getDedicatedHostNameForZone(clusterID string, role string, zone string) (string, error) {
	switch role {
	case "master":
		return fmt.Sprintf("%s-dhost-control-plane-%s", clusterID, zone), nil
	case "worker":
		return fmt.Sprintf("%s-dhost-compute-%s", clusterID, zone), nil
	default:
		return "", fmt.Errorf("invalid machine role %v", role)
	}
}

func getSubnetName(clusterID string, role string, zone string) (string, error) {
	switch role {
	case "master":
		return fmt.Sprintf("%s-subnet-control-plane-%s", clusterID, zone), nil
	case "worker":
		return fmt.Sprintf("%s-subnet-compute-%s", clusterID, zone), nil
	default:
		return "", fmt.Errorf("invalid machine role %v", role)
	}
}

func getSecurityGroupNames(clusterID string, role string) ([]string, error) {
	switch role {
	case "master":
		return []string{
			fmt.Sprintf("%s-security-group-cluster-wide", clusterID),
			fmt.Sprintf("%s-security-group-openshift-network", clusterID),
			fmt.Sprintf("%s-security-group-control-plane", clusterID),
			fmt.Sprintf("%s-security-group-control-plane-internal", clusterID),
		}, nil
	case "worker":
		return []string{
			fmt.Sprintf("%s-security-group-cluster-wide", clusterID),
			fmt.Sprintf("%s-security-group-openshift-network", clusterID),
		}, nil
	default:
		return nil, fmt.Errorf("invalid machine role %v", role)
	}
}
