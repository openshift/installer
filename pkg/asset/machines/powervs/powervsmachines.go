// Package powervs generates Machine objects for powerVS.
package powervs

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(clusterID string, ic *types.InstallConfig, pool *types.MachinePool, role string) ([]*asset.RuntimeFile, error) {
	if poolPlatform := pool.Platform.Name(); poolPlatform != powervs.Name {
		return nil, fmt.Errorf("non-Power VS machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.PowerVS

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	var (
		result         []*asset.RuntimeFile
		image          string
		service        capibm.IBMPowerVSResourceReference
		name           string
		powerVSMachine *capibm.IBMPowerVSMachine
		dataSecret     string
		machine        *capi.Machine
	)

	// Note: This will be created later
	image = fmt.Sprintf("rhcos-%s", clusterID)

	if ic.PowerVS.ServiceInstanceGUID == "" {
		serviceName := fmt.Sprintf("%s-power-iaas", clusterID)

		service = capibm.IBMPowerVSResourceReference{
			Name: &serviceName,
		}
	} else {
		service = capibm.IBMPowerVSResourceReference{
			ID: &ic.PowerVS.ServiceInstanceGUID,
		}
	}

	for idx := int64(0); idx < total; idx++ {
		name = fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx)

		powerVSMachine = GenerateMachine(ic, service, mpool, name, image)

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", powerVSMachine.Name)},
			Object: powerVSMachine,
		})

		dataSecret = fmt.Sprintf("%s-%s", clusterID, "master")
		machine = GenerateCAPIMachine(clusterID, powerVSMachine.Name, dataSecret)

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", machine.Name)},
			Object: machine,
		})
	}

	name = fmt.Sprintf("%s-bootstrap", clusterID)
	powerVSMachine = GenerateMachine(ic, service, mpool, name, image)
	powerVSMachine.Labels["install.openshift.io/bootstrap"] = ""

	result = append(result, &asset.RuntimeFile{
		File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", powerVSMachine.Name)},
		Object: powerVSMachine,
	})

	dataSecret = fmt.Sprintf("%s-%s", clusterID, "bootstrap")
	machine = GenerateCAPIMachine(clusterID, powerVSMachine.Name, dataSecret)
	machine.Labels["install.openshift.io/bootstrap"] = ""

	result = append(result, &asset.RuntimeFile{
		File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", machine.Name)},
		Object: machine,
	})

	return result, nil
}

// GenerateMachine creates a capibm.IBMPowerVSMachine struct.
func GenerateMachine(ic *types.InstallConfig, service capibm.IBMPowerVSResourceReference, mpool *powervs.MachinePool, name string, image string) *capibm.IBMPowerVSMachine {
	return &capibm.IBMPowerVSMachine{
		TypeMeta: metav1.TypeMeta{
			APIVersion: capibm.GroupVersion.String(),
			Kind:       "IBMPowerVSMachine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: capiutils.Namespace,
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capibm.IBMPowerVSMachineSpec{
			ServiceInstanceID: ic.PowerVS.ServiceInstanceGUID,
			ServiceInstance:   &service,
			SSHKey:            "",
			ImageRef: &v1.LocalObjectReference{
				Name: image,
			},
			SystemType:    mpool.SysType,
			ProcessorType: capibm.PowerVSProcessorType(mpool.ProcType),
			Processors:    mpool.Processors,
			MemoryGiB:     mpool.MemoryGiB,
		},
	}
}

// GenerateCAPIMachine creates a capi.Machine struct.
func GenerateCAPIMachine(clusterID string, name string, dataSecret string) *capi.Machine {
	return &capi.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Machine",
			APIVersion: capi.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capi.MachineSpec{
			ClusterName: clusterID,
			Bootstrap: capi.Bootstrap{
				DataSecretName: ptr.To(dataSecret),
			},
			InfrastructureRef: v1.ObjectReference{
				APIVersion: capibm.GroupVersion.String(),
				Kind:       "IBMPowerVSMachine",
				Name:       name,
			},
		},
	}
}
