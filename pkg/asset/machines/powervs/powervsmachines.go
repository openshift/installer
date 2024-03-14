// Package powervs generates Machine objects for powerVS.
package powervs

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(clusterID string, pool *types.MachinePool, role string) ([]*asset.RuntimeFile, error) {
	if poolPlatform := pool.Platform.Name(); poolPlatform != powervs.Name {
		return nil, fmt.Errorf("non-Power VS machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.PowerVS

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	var (
		result []*asset.RuntimeFile
		image  string
	)

	// Note: This will be created later
	image = fmt.Sprintf("rhcos-%v", clusterID)

	for idx := int64(0); idx < total; idx++ {
		powervsMachine := &capibm.IBMPowerVSMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capibm.IBMPowerVSMachineSpec{
				SSHKey: "",
				ImageRef: &v1.LocalObjectReference{
					Name: image,
				},
				SystemType:    mpool.SysType,
				ProcessorType: capibm.PowerVSProcessorType(mpool.ProcType),
				Processors:    mpool.Processors,
				MemoryGiB:     mpool.MemoryGiB,
			},
		}
		powervsMachine.SetGroupVersionKind(capibm.GroupVersion.WithKind("IBMPowerVSMachine"))

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", powervsMachine.Name)},
			Object: powervsMachine,
		})

		machine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: powervsMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterID,
				Bootstrap: capi.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, role)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
					Kind:       "IBMPowerVSMachine",
					Name:       powervsMachine.Name,
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
