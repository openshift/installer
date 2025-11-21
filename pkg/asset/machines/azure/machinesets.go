package azure

import (
	"fmt"
	"slices"
	"sort"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"

	clusterapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, ic *installconfig.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string, capabilities map[string]string, useImageGallery bool, subnetZones []string, session *icazure.Session) ([]*clusterapi.MachineSet, error) {
	config := ic.Config
	if configPlatform := config.Platform.Name(); configPlatform != azure.Name {
		return nil, fmt.Errorf("non-azure configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != azure.Name {
		return nil, fmt.Errorf("non-azure machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Azure
	mpool := pool.Platform.Azure

	if len(mpool.Zones) == 0 {
		// if no azs are given we set to []string{""} for convenience over later operations.
		// It means no-zoned for the machine API
		mpool.Zones = []string{""}
	}
	azs := mpool.Zones

	total := int64(0)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	networkResourceGroup, virtualNetworkName, subnets, err := getNetworkInfo(platform, clusterID, role, subnetZones)
	if err != nil {
		return nil, fmt.Errorf("failed to get subnets for role %s : %w", role, err)
	}

	sort.Strings(subnets)
	numOfAZs := int64(len(azs))
	sort.Strings(azs)
	subnetIndex := -1
	var machinesets []*clusterapi.MachineSet

	if config.Azure.OutboundType == azure.NATGatewayMultiZoneOutboundType {
		return getMultiZoneMachineSets(multiZoneMachineSetInput{
			networkResourceGroup: networkResourceGroup,
			virtualNetworkName:   virtualNetworkName,
			platform:             platform,
			mpool:                mpool,
			osImage:              osImage,
			userDataSecret:       userDataSecret,
			clusterID:            clusterID,
			role:                 role,
			capabilities:         capabilities,
			useImageGallery:      useImageGallery,
			session:              session,
			subnetSpec:           config.Azure.Subnets,
			replicas:             total,
			ic:                   ic,
			azs:                  azs,
			pool:                 pool,
		})
	}
	for idx, az := range azs {
		replicas := int32(total / numOfAZs)
		if int64(idx) < total%numOfAZs {
			replicas++
		}
		subnetIndex = (subnetIndex + 1) % len(subnets)
		provider, err := provider(platform, mpool, osImage, userDataSecret, clusterID, role, &idx, capabilities, useImageGallery, session, networkResourceGroup, virtualNetworkName, subnets[subnetIndex])
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		name := fmt.Sprintf("%s-%s-%s%s", clusterID, pool.Name, platform.Region, az)
		mset := &clusterapi.MachineSet{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "MachineSet",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      name,
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clusterID,
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
				},
			},
			Spec: clusterapi.MachineSetSpec{
				Replicas: &replicas,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"machine.openshift.io/cluster-api-machineset": name,
						"machine.openshift.io/cluster-api-cluster":    clusterID,
					},
				},
				Template: clusterapi.MachineTemplateSpec{
					ObjectMeta: clusterapi.ObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-machineset":   name,
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
						},
					},
					Spec: clusterapi.MachineSpec{
						ProviderSpec: clusterapi.ProviderSpec{
							Value: &runtime.RawExtension{Object: provider},
						},
						// we don't need to set Versions, because we control those via cluster operators.
					},
				},
			},
		}
		machinesets = append(machinesets, mset)
	}
	return machinesets, nil
}

type multiZoneMachineSetInput struct {
	networkResourceGroup string
	platform             *azure.Platform
	mpool                *azure.MachinePool
	osImage              string
	userDataSecret       string
	clusterID            string
	role                 string
	capabilities         map[string]string
	useImageGallery      bool
	session              *icazure.Session
	virtualNetworkName   string
	subnetSpec           []azure.SubnetSpec
	replicas             int64
	ic                   *installconfig.InstallConfig
	azs                  []string
	pool                 *types.MachinePool
}

func getMultiZoneMachineSets(in multiZoneMachineSetInput) ([]*clusterapi.MachineSet, error) {
	// Deep copy metadata map.
	zoneSubnetmap := map[string][]string{}
	subnetCount := 0
	// Filter for the zones the user provided for compute nodes.
	for key, value := range in.ic.Azure.ZonesSubnetMap {
		if slices.Contains(in.azs, key) {
			zoneSubnetmap[key] = sets.NewString(value...).List()
			subnetCount += len(value)
		}
	}
	machineSets := []*clusterapi.MachineSet{}
	replicasToCreate := int32(in.replicas)
	// Calculate the replicas per machine set.
	// This just first finds the nearest multiple of subnet count
	// then distributes the remainder across the machine sets one by one.
	// If there are 3 subnets and 8 replicas, first we would
	// set 8/3 = 2 replicas for each subnet (2,2,2) and distribute the
	// remaining machines (2) evenly to have (3,3,2).
	replicaPerSet := max(replicasToCreate/int32(subnetCount), 1)
	remainder := replicasToCreate % int32(subnetCount)
	if replicasToCreate < int32(subnetCount) {
		remainder = 0
	}
	numAZUsed := map[string]int{}
	for _, az := range in.azs {
		numAZUsed[az] = 0
	}

	// Iterate till we used up all the replicas mentioned.
	// Iterate through the zones provided and find a subnet to use.
	for replicasToCreate != 0 && len(zoneSubnetmap) != 0 {
		for idx, az := range in.azs {
			if replicaPerSet == 0 || len(zoneSubnetmap) == 0 {
				break
			}
			if _, ok := zoneSubnetmap[az]; !ok {
				continue
			}
			subnet := zoneSubnetmap[az][0]
			if len(zoneSubnetmap[az]) == 1 {
				delete(zoneSubnetmap, az)
			} else {
				zoneSubnetmap[az] = zoneSubnetmap[az][1:]
			}
			currentReplica := replicaPerSet
			if remainder != 0 {
				currentReplica++
				remainder--
			}
			provider, err := provider(in.platform, in.mpool, in.osImage, in.userDataSecret, in.clusterID, in.role, &idx, in.capabilities, in.useImageGallery, in.session, in.networkResourceGroup, in.virtualNetworkName, subnet)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create provider")
			}
			name := fmt.Sprintf("%s-%s-%s%s-%d", in.clusterID, in.pool.Name, in.platform.Region, az, numAZUsed[az])
			if numAZUsed[az] == 0 {
				name = fmt.Sprintf("%s-%s-%s%s", in.clusterID, in.pool.Name, in.platform.Region, az)
			}
			numAZUsed[az]++
			mset := &clusterapi.MachineSet{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "machine.openshift.io/v1beta1",
					Kind:       "MachineSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "openshift-machine-api",
					Name:      name,
					Labels: map[string]string{
						"machine.openshift.io/cluster-api-cluster":      in.clusterID,
						"machine.openshift.io/cluster-api-machine-role": in.role,
						"machine.openshift.io/cluster-api-machine-type": in.role,
					},
				},
				Spec: clusterapi.MachineSetSpec{
					Replicas: &currentReplica,
					Selector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"machine.openshift.io/cluster-api-machineset": name,
							"machine.openshift.io/cluster-api-cluster":    in.clusterID,
						},
					},
					Template: clusterapi.MachineTemplateSpec{
						ObjectMeta: clusterapi.ObjectMeta{
							Labels: map[string]string{
								"machine.openshift.io/cluster-api-machineset":   name,
								"machine.openshift.io/cluster-api-cluster":      in.clusterID,
								"machine.openshift.io/cluster-api-machine-role": in.role,
								"machine.openshift.io/cluster-api-machine-type": in.role,
							},
						},
						Spec: clusterapi.MachineSpec{
							ProviderSpec: clusterapi.ProviderSpec{
								Value: &runtime.RawExtension{Object: provider},
							},
							// we don't need to set Versions, because we control those via cluster operators.
						},
					},
				},
			}
			machineSets = append(machineSets, mset)
			replicasToCreate -= currentReplica
		}
	}
	return machineSets, nil
}
