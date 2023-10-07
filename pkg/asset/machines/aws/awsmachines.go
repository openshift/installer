// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	capiGuestsNamespace = "openshift-cluster-api-guests"
)

// Machines returns a list of machines for a machinepool.
func AWSMachines(clusterID string, region string, subnets map[string]string, pool *types.MachinePool, role string, userTags map[string]string) ([]client.Object, error) {
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.AWS

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	// tags, err := tagsFromUserTags(clusterID, userTags)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to create machineapi.TagSpecifications from UserTags")
	// }

	var result []client.Object

	for idx := int64(0); idx < total; idx++ {
		//zone := mpool.Zones[int(idx)%len(mpool.Zones)]
		// subnet, ok := subnets[zone]
		// if len(subnets) > 0 && !ok {
		// 	return nil, errors.Errorf("no subnet for zone %s", zone)
		// }

		awsMachine := &capa.AWSMachine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
				Kind:       "AWSMachine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: capiGuestsNamespace,
				Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capa.AWSMachineSpec{
				//failureDomain?
				Ignition:             &capa.Ignition{Version: "3.2"},
				UncompressedUserData: pointer.Bool(true),
				InstanceType:         mpool.InstanceType,
				AMI:                  capa.AMIReference{ID: pointer.String(mpool.AMIID)},
				SSHKeyName:           pointer.String(""),
				// IAMInstanceProfile:   fmt.Sprintf("%s-master-profile", clusterID),
				//Subnet: ?
				//AdditionalTags: tags,
			},
			// RootVolume: capa.Volume{
			// 	Size:      int64(mpool.EC2RootVolume.Size),
			// 	Type:      mpool.EC2RootVolume.Type,
			// 	IOPS:      int64(mpool.EC2RootVolume.IOPS),
			// 	Encrypted: true, // is this configurable? Use KMS?
			// },
		}

		machine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: capiGuestsNamespace,
				Name:      awsMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterID,
				Bootstrap: capi.Bootstrap{
					DataSecretName: pointer.String(fmt.Sprintf("%s-%s", clusterID, role)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
					Kind:       "AWSMachine",
					Name:       awsMachine.Name,
				},
			},
		}

		result = append(result, awsMachine, machine)

	}
	return result, nil
}
