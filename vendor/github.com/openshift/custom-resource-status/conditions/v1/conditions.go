package v1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetStatusCondition sets the corresponding condition in conditions to newCondition.
// The return value indicates if this resulted in any changes *other than* LastHeartbeatTime.
func SetStatusCondition(conditions *[]Condition, newCondition Condition) bool {
	if conditions == nil {
		conditions = &[]Condition{}
	}
	existingCondition := FindStatusCondition(*conditions, newCondition.Type)
	if existingCondition == nil {
		newCondition.LastTransitionTime = metav1.NewTime(time.Now())
		newCondition.LastHeartbeatTime = metav1.NewTime(time.Now())
		*conditions = append(*conditions, newCondition)
		return true
	}

	changed := updateCondition(existingCondition, newCondition)
	existingCondition.LastHeartbeatTime = metav1.NewTime(time.Now())
	return changed
}

// SetStatusConditionNoHearbeat sets the corresponding condition in conditions to newCondition
// without setting lastHeartbeatTime.
// The return value indicates if this resulted in any changes.
func SetStatusConditionNoHeartbeat(conditions *[]Condition, newCondition Condition) bool {
	if conditions == nil {
		conditions = &[]Condition{}
	}
	existingCondition := FindStatusCondition(*conditions, newCondition.Type)
	if existingCondition == nil {
		newCondition.LastTransitionTime = metav1.NewTime(time.Now())
		*conditions = append(*conditions, newCondition)
		return true
	}

	return updateCondition(existingCondition, newCondition)
}

// RemoveStatusCondition removes the corresponding conditionType from conditions.
func RemoveStatusCondition(conditions *[]Condition, conditionType ConditionType) {
	if conditions == nil {
		return
	}
	newConditions := []Condition{}
	for _, condition := range *conditions {
		if condition.Type != conditionType {
			newConditions = append(newConditions, condition)
		}
	}

	*conditions = newConditions
}

func updateCondition(existingCondition *Condition, newCondition Condition) bool {
	changed := false
	if existingCondition.Status != newCondition.Status {
		changed = true
		existingCondition.Status = newCondition.Status
		existingCondition.LastTransitionTime = metav1.NewTime(time.Now())
	}

	if existingCondition.Reason != newCondition.Reason {
		changed = true
		existingCondition.Reason = newCondition.Reason
	}
	if existingCondition.Message != newCondition.Message {
		changed = true
		existingCondition.Message = newCondition.Message
	}
	return changed
}

// FindStatusCondition finds the conditionType in conditions.
func FindStatusCondition(conditions []Condition, conditionType ConditionType) *Condition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return &conditions[i]
		}
	}

	return nil
}

// IsStatusConditionTrue returns true when the conditionType is present and set to `corev1.ConditionTrue`
func IsStatusConditionTrue(conditions []Condition, conditionType ConditionType) bool {
	return IsStatusConditionPresentAndEqual(conditions, conditionType, corev1.ConditionTrue)
}

// IsStatusConditionFalse returns true when the conditionType is present and set to `corev1.ConditionFalse`
func IsStatusConditionFalse(conditions []Condition, conditionType ConditionType) bool {
	return IsStatusConditionPresentAndEqual(conditions, conditionType, corev1.ConditionFalse)
}

// IsStatusConditionUnknown returns true when the conditionType is present and set to `corev1.ConditionUnknown`
func IsStatusConditionUnknown(conditions []Condition, conditionType ConditionType) bool {
	return IsStatusConditionPresentAndEqual(conditions, conditionType, corev1.ConditionUnknown)
}

// IsStatusConditionPresentAndEqual returns true when conditionType is present and equal to status.
func IsStatusConditionPresentAndEqual(conditions []Condition, conditionType ConditionType, status corev1.ConditionStatus) bool {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return condition.Status == status
		}
	}
	return false
}
