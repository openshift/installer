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

package v1beta2

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

func validateLifecycleHooks(hooks []AWSLifecycleHook) field.ErrorList {
	var allErrs field.ErrorList

	for _, hook := range hooks {
		if hook.Name == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.lifecycleHooks.name"), "Name is required"))
		}
		if hook.NotificationTargetARN != nil && hook.RoleARN == nil {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.lifecycleHooks.roleARN"), "RoleARN is required if NotificationTargetARN is provided"))
		}
		if hook.RoleARN != nil && hook.NotificationTargetARN == nil {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.lifecycleHooks.notificationTargetARN"), "NotificationTargetARN is required if RoleARN is provided"))
		}
		if hook.LifecycleTransition != LifecycleHookTransitionInstanceLaunching && hook.LifecycleTransition != LifecycleHookTransitionInstanceTerminating {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.lifecycleTransition"), hook.LifecycleTransition, fmt.Sprintf("LifecycleTransition must be either %q or %q", LifecycleHookTransitionInstanceLaunching, LifecycleHookTransitionInstanceTerminating)))
		}
		if hook.DefaultResult != nil && (*hook.DefaultResult != LifecycleHookDefaultResultContinue && *hook.DefaultResult != LifecycleHookDefaultResultAbandon) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.defaultResult"), *hook.DefaultResult, fmt.Sprintf("DefaultResult must be either %s or %s", LifecycleHookDefaultResultContinue, LifecycleHookDefaultResultAbandon)))
		}
		if hook.HeartbeatTimeout != nil && (hook.HeartbeatTimeout.Seconds() < float64(30) || hook.HeartbeatTimeout.Seconds() > float64(172800)) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.heartbeatTimeout"), *hook.HeartbeatTimeout, "HeartbeatTimeout must be between 30 and 172800 seconds"))
		}
	}

	return allErrs
}
