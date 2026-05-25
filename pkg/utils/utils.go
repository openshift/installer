package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capi "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563

	"github.com/openshift/api/features"
	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
)

// SetMachineOSStreamLabels adds the OS image stream label to a Machine if the OSStreams
// feature gate is enabled.
func SetMachineOSStreamLabels[T metav1.Object](obj T, ic *types.InstallConfig) {
	if ic == nil || !ic.Enabled(features.FeatureGateOSStreams) {
		return
	}
	labels := obj.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[types.OSStreamLabelKey] = string(ic.OSImageStream)
	obj.SetLabels(labels)
}

// SetMachineSetOSStreamLabels adds the OS image stream label to a MachineSet's metadata
// and Spec.Template if the OSStreams feature gate is enabled.
func SetMachineSetOSStreamLabels(machineSet *machineapi.MachineSet, ic *types.InstallConfig) {
	if ic == nil || !ic.Enabled(features.FeatureGateOSStreams) {
		return
	}
	// Set the metadata labels
	labels := machineSet.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[types.OSStreamLabelKey] = string(ic.OSImageStream)
	machineSet.SetLabels(labels)
	// Set the Spec.Template labels
	if machineSet.Spec.Template.Labels == nil {
		machineSet.Spec.Template.Labels = make(map[string]string)
	}
	machineSet.Spec.Template.Labels[types.OSStreamLabelKey] = string(ic.OSImageStream)
}

// SetCAPIMachineSetOSStreamLabels adds the OS image stream label to a CAPI MachineSet's
// metadata and Spec.Template if the OSStreams feature gate is enabled.
func SetCAPIMachineSetOSStreamLabels(machineSet *capi.MachineSet, ic *types.InstallConfig) {
	if ic == nil || !ic.Enabled(features.FeatureGateOSStreams) {
		return
	}
	labels := machineSet.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[types.OSStreamLabelKey] = string(ic.OSImageStream)
	machineSet.SetLabels(labels)
	if machineSet.Spec.Template.ObjectMeta.Labels == nil {
		machineSet.Spec.Template.ObjectMeta.Labels = make(map[string]string)
	}
	machineSet.Spec.Template.ObjectMeta.Labels[types.OSStreamLabelKey] = string(ic.OSImageStream)
}

// SetCPMSOSStreamLabels adds the OS image stream label to a ControlPlaneMachineSet's
// metadata and Spec.Template if the OSStreams feature gate is enabled.
func SetCPMSOSStreamLabels(cpms *machinev1.ControlPlaneMachineSet, ic *types.InstallConfig) {
	if ic == nil || !ic.Enabled(features.FeatureGateOSStreams) {
		return
	}
	// Set the metadata labels
	labels := cpms.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[types.OSStreamLabelKey] = string(ic.OSImageStream)
	cpms.SetLabels(labels)
	// Set the Spec.Template labels
	if cpms.Spec.Template.OpenShiftMachineV1Beta1Machine == nil {
		cpms.Spec.Template.OpenShiftMachineV1Beta1Machine = &machinev1.OpenShiftMachineV1Beta1MachineTemplate{}
	}
	if cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.ObjectMeta.Labels == nil {
		cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.ObjectMeta.Labels = make(map[string]string)
	}
	cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.ObjectMeta.Labels[types.OSStreamLabelKey] = string(ic.OSImageStream)
}
