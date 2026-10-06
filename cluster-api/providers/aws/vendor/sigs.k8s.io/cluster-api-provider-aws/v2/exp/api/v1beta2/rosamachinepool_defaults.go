/*
Copyright 2026 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

// Default satisfies the defaulting webhook interface.
func (r *ROSAMachinePool) Default() {
	if r.Spec.NodeDrainGracePeriod == nil {
		r.Spec.NodeDrainGracePeriod = &metav1.Duration{}
	}

	if r.Spec.UpdateConfig == nil {
		r.Spec.UpdateConfig = &RosaUpdateConfig{}
	}
	if r.Spec.UpdateConfig.RollingUpdate == nil {
		r.Spec.UpdateConfig.RollingUpdate = &RollingUpdate{
			MaxUnavailable: ptr.To(intstr.FromInt32(0)),
			MaxSurge:       ptr.To(intstr.FromInt32(1)),
		}
	}
}
