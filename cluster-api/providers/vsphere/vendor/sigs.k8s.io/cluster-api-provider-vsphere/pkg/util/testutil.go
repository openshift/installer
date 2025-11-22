/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	netopv1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	ncpv1 "github.com/vmware-tanzu/vm-operator/external/ncp/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	bootstrapv1 "sigs.k8s.io/cluster-api/api/bootstrap/kubeadm/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	topologyv1 "sigs.k8s.io/cluster-api-provider-vsphere/internal/apis/topology/v1alpha1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
)

const (
	clusterKind          = "Cluster"
	infraClusterKind     = "VSphereCluster"
	machineKind          = "Machine"
	infraMachineKind     = "VSphereMachine"
	clusterNameLabelName = "cluster.x-k8s.io/cluster-name"
)

func CreateCluster(clusterName string) *clusterv1.Cluster {
	return &clusterv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterv1.GroupVersion.String(),
			Kind:       clusterKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterName,
			Namespace: corev1.NamespaceDefault,
		},
		Spec: clusterv1.ClusterSpec{
			InfrastructureRef: clusterv1.ContractVersionedObjectReference{
				APIGroup: vmwarev1.GroupVersion.Group,
				Kind:     infraClusterKind,
				Name:     clusterName,
			},
		},
	}
}

func CreateVSphereCluster(clusterName string) *vmwarev1.VSphereCluster {
	return &vmwarev1.VSphereCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: vmwarev1.GroupVersion.String(),
			Kind:       infraClusterKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterName,
			Namespace: corev1.NamespaceDefault,
		},
	}
}

func CreateMachine(machineName, clusterName, k8sVersion string, controlPlaneLabel bool) *clusterv1.Machine {
	machine := &clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterv1.GroupVersion.String(),
			Kind:       machineKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineName,
			Namespace: corev1.NamespaceDefault,
			Labels: map[string]string{
				clusterNameLabelName: clusterName,
			},
		},
		Spec: clusterv1.MachineSpec{
			Version: k8sVersion,
			Bootstrap: clusterv1.Bootstrap{
				ConfigRef: clusterv1.ContractVersionedObjectReference{
					APIGroup: bootstrapv1.GroupVersion.Group,
					Name:     machineName,
				},
			},
			InfrastructureRef: clusterv1.ContractVersionedObjectReference{
				APIGroup: vmwarev1.GroupVersion.Group,
				Kind:     infraMachineKind,
				Name:     machineName,
			},
		},
	}
	if controlPlaneLabel {
		labels := machine.GetLabels()
		labels[clusterv1.MachineControlPlaneLabel] = ""
		machine.SetLabels(labels)
	}
	return machine
}

func CreateVSphereMachine(machineName, clusterName, className, imageName, storageClass string, controlPlaneLabel bool) *vmwarev1.VSphereMachine {
	vsphereMachine := &vmwarev1.VSphereMachine{
		TypeMeta: metav1.TypeMeta{
			APIVersion: vmwarev1.GroupVersion.String(),
			Kind:       infraMachineKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineName,
			Namespace: corev1.NamespaceDefault,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: clusterName,
			},
		},
		Spec: vmwarev1.VSphereMachineSpec{
			ClassName:          className,
			ImageName:          imageName,
			StorageClass:       storageClass,
			MinHardwareVersion: "vmx-17",
		},
	}
	if controlPlaneLabel {
		labels := vsphereMachine.GetLabels()
		labels[clusterv1.MachineControlPlaneLabel] = ""
		vsphereMachine.SetLabels(labels)
	}
	return vsphereMachine
}

func createScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = vmwarev1.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	_ = topologyv1.AddToScheme(scheme)
	_ = vmoprv1.AddToScheme(scheme)
	_ = netopv1.AddToScheme(scheme)
	_ = ncpv1.AddToScheme(scheme)
	return scheme
}

func CreateClusterContext(cluster *clusterv1.Cluster, vsphereCluster *vmwarev1.VSphereCluster) (*vmware.ClusterContext, *capvcontext.ControllerManagerContext) {
	scheme := createScheme()

	// Build the cluster context.
	return &vmware.ClusterContext{
			Cluster:        cluster,
			VSphereCluster: vsphereCluster,
		}, &capvcontext.ControllerManagerContext{
			Logger: klog.Background().WithName("controller-manager-logger"),
			Scheme: scheme,
			Client: fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(
				&vmoprv1.VirtualMachineService{},
				&vmoprv1.VirtualMachine{},
			).Build(),
		}
}

func CreateMachineContext(clusterContext *vmware.ClusterContext, machine *clusterv1.Machine,
	vsphereMachine *vmwarev1.VSphereMachine) *vmware.SupervisorMachineContext {
	// Build the machine context.
	return &vmware.SupervisorMachineContext{
		BaseMachineContext: &capvcontext.BaseMachineContext{
			Machine: machine,
			Cluster: clusterContext.Cluster,
		},
		VSphereCluster: clusterContext.VSphereCluster,
		VSphereMachine: vsphereMachine,
	}
}
