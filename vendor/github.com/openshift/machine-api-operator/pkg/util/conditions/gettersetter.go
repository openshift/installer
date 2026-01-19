package conditions

import (
	"fmt"
	"sort"
	"time"

	machinev1 "github.com/openshift/api/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type GetterSetter interface {
	runtime.Object
	metav1.Object

	// GetConditions returns the list of conditions for a machine API object.
	GetConditions() []machinev1.Condition

	// SetConditions sets the list of conditions for a machine API object.
	SetConditions([]machinev1.Condition)
}

// Get returns the condition with the given type, if the condition does not exists,
// it returns nil.
func Get(from interface{}, t machinev1.ConditionType) *machinev1.Condition {
	obj := getWrapperObject(from)
	conditions := obj.GetConditions()
	if conditions == nil {
		return nil
	}

	for _, condition := range conditions {
		if condition.Type == t {
			return &condition
		}
	}
	return nil
}

// Set sets the given condition.
//
// NOTE: If a condition already exists, the LastTransitionTime is updated only if a change is detected
// in any of the following fields: Status, Reason, Severity and Message.
func Set(to interface{}, condition *machinev1.Condition) {
	if to == nil || condition == nil {
		return
	}

	obj := getWrapperObject(to)

	// Check if the new conditions already exists, and change it only if there is a status
	// transition (otherwise we should preserve the current last transition time)-
	conditions := obj.GetConditions()
	exists := false
	for i := range conditions {
		existingCondition := conditions[i]
		if existingCondition.Type == condition.Type {
			exists = true
			if !hasSameState(&existingCondition, condition) {
				condition.LastTransitionTime = metav1.NewTime(time.Now().UTC().Truncate(time.Second))
				conditions[i] = *condition
				break
			}
			condition.LastTransitionTime = existingCondition.LastTransitionTime
			break
		}
	}

	// If the condition does not exist, add it, setting the transition time only if not already set
	if !exists {
		if condition.LastTransitionTime.IsZero() {
			condition.LastTransitionTime = metav1.NewTime(time.Now().UTC().Truncate(time.Second))
		}
		conditions = append(conditions, *condition)
	}

	// Sorts conditions for convenience of the consumer, i.e. kubectl.
	sort.Slice(conditions, func(i, j int) bool {
		return lexicographicLess(&conditions[i], &conditions[j])
	})

	obj.SetConditions(conditions)
}

// TrueCondition returns a condition with Status=True and the given type.
func TrueCondition(t machinev1.ConditionType) *machinev1.Condition {
	return &machinev1.Condition{
		Type:   t,
		Status: corev1.ConditionTrue,
	}
}

// TrueConditionWithReason returns a condition with Status=True and the given type.
func TrueConditionWithReason(t machinev1.ConditionType, reason string, messageFormat string, messageArgs ...interface{}) *machinev1.Condition {
	return &machinev1.Condition{
		Type:    t,
		Status:  corev1.ConditionTrue,
		Reason:  reason,
		Message: fmt.Sprintf(messageFormat, messageArgs...),
	}
}

// FalseCondition returns a condition with Status=False and the given type.
func FalseCondition(t machinev1.ConditionType, reason string, severity machinev1.ConditionSeverity, messageFormat string, messageArgs ...interface{}) *machinev1.Condition {
	return &machinev1.Condition{
		Type:     t,
		Status:   corev1.ConditionFalse,
		Reason:   reason,
		Severity: severity,
		Message:  fmt.Sprintf(messageFormat, messageArgs...),
	}
}

// UnknownCondition returns a condition with Status=Unknown and the given type.
func UnknownCondition(t machinev1.ConditionType, reason string, messageFormat string, messageArgs ...interface{}) *machinev1.Condition {
	return &machinev1.Condition{
		Type:    t,
		Status:  corev1.ConditionUnknown,
		Reason:  reason,
		Message: fmt.Sprintf(messageFormat, messageArgs...),
	}
}

// MarkTrue sets Status=True for the condition with the given type.
func MarkTrue(to interface{}, t machinev1.ConditionType) {
	Set(to, TrueCondition(t))
}

// MarkFalse sets Status=False for the condition with the given type.
func MarkFalse(to interface{}, t machinev1.ConditionType, reason string, severity machinev1.ConditionSeverity, messageFormat string, messageArgs ...interface{}) {
	Set(to, FalseCondition(t, reason, severity, messageFormat, messageArgs...))
}

// IsTrue is true if the condition with the given type is True, otherwise it return false
// if the condition is not True or if the condition does not exist (is nil).
func IsTrue(from interface{}, t machinev1.ConditionType) bool {
	if c := Get(from, t); c != nil {
		return c.Status == corev1.ConditionTrue
	}
	return false
}

// IsFalse is true if the condition with the given type is False, otherwise it return false
// if the condition is not False or if the condition does not exist (is nil).
func IsFalse(from interface{}, t machinev1.ConditionType) bool {
	if c := Get(from, t); c != nil {
		return c.Status == corev1.ConditionFalse
	}
	return false
}

// IsEquivalentTo returns true if condition a is equivalent to condition b,
// by checking for equality of the following fields: Type, Status, Reason, Severity and Message (it excludes LastTransitionTime).
func IsEquivalentTo(a, b *machinev1.Condition) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	}

	return hasSameState(a, b)
}

// lexicographicLess returns true if a condition is less than another with regards to the
// to order of conditions designed for convenience of the consumer, i.e. kubectl.
func lexicographicLess(i, j *machinev1.Condition) bool {
	return i.Type < j.Type
}

// hasSameState returns true if a condition has the same state of another; state is defined
// by the union of following fields: Type, Status, Reason, Severity and Message (it excludes LastTransitionTime).
func hasSameState(i, j *machinev1.Condition) bool {
	return i.Type == j.Type &&
		i.Status == j.Status &&
		i.Reason == j.Reason &&
		i.Severity == j.Severity &&
		i.Message == j.Message
}

func getWrapperObject(from interface{}) GetterSetter {
	switch obj := from.(type) {
	case *machinev1.Machine:
		return &MachineWrapper{obj}
	case *machinev1.MachineSet:
		return &MachineSetWrapper{obj}
	case *machinev1.MachineHealthCheck:
		return &MachineHealthCheckWrapper{obj}
	default:
		panic("type is not supported as conditions getter or setter")
	}
}
