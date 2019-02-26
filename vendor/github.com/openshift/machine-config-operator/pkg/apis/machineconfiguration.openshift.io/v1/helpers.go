package v1

import (
	"sort"

	ignv2_2 "github.com/coreos/ignition/config/v2_2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MergeMachineConfigs combines multiple machineconfig objects into one object.
// It sorts all the configs in increasing order of their name.
// It uses the Ign config from first object as base and appends all the rest.
// It only uses the OSImageURL from first object and ignores it from rest.
func MergeMachineConfigs(configs []*MachineConfig) *MachineConfig {
	if len(configs) == 0 {
		return nil
	}
	sort.Slice(configs, func(i, j int) bool { return configs[i].Name < configs[j].Name })

	outOSImageURL := configs[0].Spec.OSImageURL
	outIgn := configs[0].Spec.Config
	for idx := 1; idx < len(configs); idx++ {
		outIgn = ignv2_2.Append(outIgn, configs[idx].Spec.Config)
	}

	return &MachineConfig{
		Spec: MachineConfigSpec{
			OSImageURL: outOSImageURL,
			Config:     outIgn,
		},
	}
}

// NewMachineConfigPoolCondition creates a new MachineConfigPool condition.
func NewMachineConfigPoolCondition(condType MachineConfigPoolConditionType, status corev1.ConditionStatus, reason, message string) *MachineConfigPoolCondition {
	return &MachineConfigPoolCondition{
		Type:               condType,
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

// GetMachineConfigPoolCondition returns the condition with the provided type.
func GetMachineConfigPoolCondition(status MachineConfigPoolStatus, condType MachineConfigPoolConditionType) *MachineConfigPoolCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

// SetMachineConfigPoolCondition updates the MachineConfigPool to include the provided condition. If the condition that
// we are about to add already exists and has the same status and reason then we are not going to update.
func SetMachineConfigPoolCondition(status *MachineConfigPoolStatus, condition MachineConfigPoolCondition) {
	currentCond := GetMachineConfigPoolCondition(*status, condition.Type)
	if currentCond != nil && currentCond.Status == condition.Status && currentCond.Reason == condition.Reason {
		return
	}
	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == condition.Status {
		condition.LastTransitionTime = currentCond.LastTransitionTime
	}
	newConditions := filterOutCondition(status.Conditions, condition.Type)
	status.Conditions = append(newConditions, condition)
}

// RemoveMachineConfigPoolCondition removes the MachineConfigPool condition with the provided type.
func RemoveMachineConfigPoolCondition(status *MachineConfigPoolStatus, condType MachineConfigPoolConditionType) {
	status.Conditions = filterOutCondition(status.Conditions, condType)
}

// filterOutCondition returns a new slice of MachineConfigPool conditions without conditions with the provided type.
func filterOutCondition(conditions []MachineConfigPoolCondition, condType MachineConfigPoolConditionType) []MachineConfigPoolCondition {
	var newConditions []MachineConfigPoolCondition
	for _, c := range conditions {
		if c.Type == condType {
			continue
		}
		newConditions = append(newConditions, c)
	}
	return newConditions
}
