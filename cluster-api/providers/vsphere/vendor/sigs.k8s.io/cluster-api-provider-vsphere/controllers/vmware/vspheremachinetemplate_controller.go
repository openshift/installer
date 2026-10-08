/*
Copyright 2024 The Kubernetes Authors.

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

package vmware

import (
	"context"

	pkgerrors "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	capicontrollerutil "sigs.k8s.io/cluster-api/util/controller"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
	conversionclient "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/client"
)

// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vmoperator.vmware.com,resources=virtualmachineclasses,verbs=get;list;watch
// +kubebuilder:rbac:groups=vmoperator.vmware.com,resources=clustervirtualmachineimages,verbs=get;list;watch

// AddVSphereMachineTemplateControllerToManager adds the machine template controller to the provided
// manager.
func AddVSphereMachineTemplateControllerToManager(ctx context.Context, controllerManagerContext *capvcontext.ControllerManagerContext, mgr manager.Manager, options controller.Options) error {
	r := &vSphereMachineTemplateReconciler{
		Client: controllerManagerContext.Client,
	}
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "vspheremachinetemplate")

	// NOTE: use vm-operator native types for watches (the reconciler uses the internal hub version).
	vmClass, err := conversionclient.WatchObject(r.Client, &vmoprvhub.VirtualMachineClass{})
	if err != nil {
		return pkgerrors.Wrap(err, "failed to create watch object for VirtualMachineClass")
	}

	return capicontrollerutil.NewControllerManagedBy(mgr, predicateLog).
		For(&vmwarev1.VSphereMachineTemplate{}).
		WithOptions(options).
		Watches(
			vmClass,
			handler.EnqueueRequestsFromMapFunc(r.enqueueVirtualMachineClassToVSphereMachineTemplateRequests),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), predicateLog, controllerManagerContext.WatchFilterValue)).
		Complete(ctx, r)
}

type vSphereMachineTemplateReconciler struct {
	Client client.Client
}

func (r *vSphereMachineTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	// Fetch VSphereMachineTemplate object
	vSphereMachineTemplate := &vmwarev1.VSphereMachineTemplate{}
	if err := r.Client.Get(ctx, req.NamespacedName, vSphereMachineTemplate); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// If ClassName is not set, there is nothing to do.
	if vSphereMachineTemplate.Spec.Template.Spec.ClassName == "" {
		return reconcile.Result{}, nil
	}

	// Fetch the VirtualMachineClass
	vmClass := &vmoprvhub.VirtualMachineClass{}
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: vSphereMachineTemplate.Spec.Template.Spec.ClassName}, vmClass); err != nil {
		return reconcile.Result{}, pkgerrors.Wrapf(err, "failed to get VirtualMachineClass %q for VSphereMachineTemplate", vSphereMachineTemplate.Spec.Template.Spec.ClassName)
	}

	patchHelper, err := patch.NewHelper(vSphereMachineTemplate, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}

	if vSphereMachineTemplate.Status.Capacity == nil {
		vSphereMachineTemplate.Status.Capacity = corev1.ResourceList{}
	}
	if vmClass.Spec.Hardware.Cpus > 0 {
		vSphereMachineTemplate.Status.Capacity[vmwarev1.VSphereResourceCPU] = *resource.NewQuantity(vmClass.Spec.Hardware.Cpus, resource.DecimalSI)
	}
	if !vmClass.Spec.Hardware.Memory.IsZero() {
		vSphereMachineTemplate.Status.Capacity[vmwarev1.VSphereResourceMemory] = vmClass.Spec.Hardware.Memory
	}

	// retrieve the os and arch info from the ClusterVirtualMachineImage
	os, arch, err := getOSAndArchFromClusterVirtualMachineImage(ctx, r.Client, vSphereMachineTemplate.Spec.Template.Spec.ImageName)
	if err != nil {
		return reconcile.Result{}, err
	}

	if validOS := normalizeOperatingSystem(os); validOS != "" {
		vSphereMachineTemplate.Status.NodeInfo.OperatingSystem = validOS
	}
	if validArch := normalizeArchitecture(arch); validArch != "" {
		vSphereMachineTemplate.Status.NodeInfo.Architecture = validArch
	}

	return reconcile.Result{}, patchHelper.Patch(ctx, vSphereMachineTemplate)
}

// normalizeOperatingSystem converts the OS string from CVMI to a valid OperatingSystem constant.
// Returns empty string if the value is not recognized.
func normalizeOperatingSystem(os string) vmwarev1.OperatingSystem {
	switch os {
	case "linux":
		return vmwarev1.OperatingSystemLinux
	case "windows":
		return vmwarev1.OperatingSystemWindows
	default:
		return ""
	}
}

// normalizeArchitecture converts the architecture string from CVMI to a valid Architecture constant.
// Returns empty string if the value is not recognized.
func normalizeArchitecture(arch string) vmwarev1.Architecture {
	switch arch {
	case "amd64":
		return vmwarev1.ArchitectureAmd64
	case "arm64":
		return vmwarev1.ArchitectureArm64
	case "s390x":
		return vmwarev1.ArchitectureS390x
	case "ppc64le":
		return vmwarev1.ArchitecturePpc64le
	default:
		return ""
	}
}

func getOSAndArchFromClusterVirtualMachineImage(ctx context.Context, c client.Client, imageName string) (string, string, error) {
	if imageName == "" {
		return "", "", nil
	}
	// Try to fetch the ClusterVirtualMachineImage with the given name
	cvmi := &vmoprvhub.ClusterVirtualMachineImage{}
	if err := c.Get(ctx, client.ObjectKey{Name: imageName}, cvmi); err != nil {
		if apierrors.IsNotFound(err) {
			return "", "", nil
		}
		return "", "", pkgerrors.Wrapf(err, "failed to get ClusterVirtualMachineImage %q", imageName)
	}

	// Extract OS type and architecture from vmwareSystemProperties
	var osType, osArch string
	for _, prop := range cvmi.Status.VMwareSystemProperties {
		switch prop.Key {
		case vmwarev1.VMwareSystemOSTypePropertyKey:
			osType = prop.Value
		case vmwarev1.VMwareSystemOSArchPropertyKey:
			osArch = prop.Value
		}
	}

	return osType, osArch, nil
}

// enqueueVirtualMachineClassToVSphereMachineTemplateRequests returns a list of VSphereMachineTemplate reconcile requests
// that use a specific VirtualMachineClass.
func (r *vSphereMachineTemplateReconciler) enqueueVirtualMachineClassToVSphereMachineTemplateRequests(ctx context.Context, virtualMachineClass client.Object) []reconcile.Request {
	requests := []reconcile.Request{}

	vSphereMachineTemplates := &vmwarev1.VSphereMachineTemplateList{}
	if err := r.Client.List(ctx, vSphereMachineTemplates, client.InNamespace(virtualMachineClass.GetNamespace())); err != nil {
		return nil
	}

	for _, vSphereMachineTemplate := range vSphereMachineTemplates.Items {
		if vSphereMachineTemplate.Spec.Template.Spec.ClassName != virtualMachineClass.GetName() {
			continue
		}

		requests = append(requests, reconcile.Request{
			NamespacedName: client.ObjectKey{Namespace: vSphereMachineTemplate.Namespace, Name: vSphereMachineTemplate.Name},
		})
	}

	return requests
}
