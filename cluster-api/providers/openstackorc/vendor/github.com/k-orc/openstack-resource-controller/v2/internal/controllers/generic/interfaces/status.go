/*
Copyright 2025 The ORC Authors.

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

package interfaces

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	applyconfigv1 "k8s.io/client-go/applyconfigurations/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
)

// ORCApplyConfig is an interface implemented by any apply configuration for an
// ORC API object. Specifically its WithStatus method is constrained to return
// an ORCStatusApplyConfig.
type ORCApplyConfig[objectApplyPT any, statusApplyPT ORCStatusApplyConfig[statusApplyPT]] interface {
	WithUID(types.UID) objectApplyPT
	WithStatus(statusApplyPT) objectApplyPT
}

// ORCStatusApplyConfig is an interface implemented by the status of any apply
// configuration for an ORC API object. It has Conditions and an ID field.
type ORCStatusApplyConfig[statusApplyPT any] interface {
	WithConditions(...*applyconfigv1.ConditionApplyConfiguration) statusApplyPT
	WithID(id string) statusApplyPT
}

// ResourceStatusWriter defines methods for writing an ORC object status
type ResourceStatusWriter[objectPT orcv1alpha1.ObjectWithConditions, osResourcePT any, objectApplyPT ORCApplyConfig[objectApplyPT, statusApplyPT], statusApplyPT ORCStatusApplyConfig[statusApplyPT]] interface {
	// GetApplyConfig returns an ORCApplyConfig for this object for use in an
	// SSA transaction, initialised with a name and a namespace.
	GetApplyConfig(name, namespace string) objectApplyPT

	// ResourceAvailableStatus returns what the status of the Available
	// condition should be set to based on the observed state of the given
	// orcObject and osResource.
	ResourceAvailableStatus(orcObject objectPT, osResource osResourcePT) (metav1.ConditionStatus, progress.ReconcileStatus)

	// ApplyResourceStatus writes status.resource to the given status apply
	// configuration based on the given osResource
	ApplyResourceStatus(log logr.Logger, osResource osResourcePT, statusApply statusApplyPT)
}
