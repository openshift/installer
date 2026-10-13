package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/api/features"
	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

// GetOSImageStream returns the OS image stream value from the install config, defaulting to the
// value of `DefaultOSImageStream` in as defined in pkg/rhcos/stream.go if not specified.
func GetOSImageStream(ic *types.InstallConfig) string {
	if ic.OSImageStream == "" {
		return string(rhcos.DefaultOSImageStream)
	}
	return string(ic.OSImageStream)
}

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
	labels[types.OSStreamLabelKey] = GetOSImageStream(ic)
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
	labels[types.OSStreamLabelKey] = GetOSImageStream(ic)
	machineSet.SetLabels(labels)
	// Set the Spec.Template labels
	if machineSet.Spec.Template.Labels == nil {
		machineSet.Spec.Template.Labels = make(map[string]string)
	}
	machineSet.Spec.Template.Labels[types.OSStreamLabelKey] = GetOSImageStream(ic)
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
	labels[types.OSStreamLabelKey] = GetOSImageStream(ic)
	cpms.SetLabels(labels)
	// Set the Spec.Template labels
	if cpms.Spec.Template.OpenShiftMachineV1Beta1Machine == nil {
		cpms.Spec.Template.OpenShiftMachineV1Beta1Machine = &machinev1.OpenShiftMachineV1Beta1MachineTemplate{}
	}
	if cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.ObjectMeta.Labels == nil {
		cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.ObjectMeta.Labels = make(map[string]string)
	}
	cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.ObjectMeta.Labels[types.OSStreamLabelKey] = GetOSImageStream(ic)
}
