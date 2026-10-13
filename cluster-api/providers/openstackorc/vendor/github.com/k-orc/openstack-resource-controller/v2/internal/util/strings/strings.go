/*
Copyright 2024 The ORC Authors.

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

package strings

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const ORCK8SPrefix = "openstack.k-orc.cloud"

type SSATransactionID string

const (
	// Field owner of the object finalizer.
	SSATransactionFinalizer SSATransactionID = "finalizer"
	SSATransactionStatus    SSATransactionID = "status"
)

func getSSAFieldOwnerString(controllerName string) string {
	return ORCK8SPrefix + "/" + controllerName + "controller"
}

// GetSSAFieldOwner returns a field owner string for a controller without a
// transaction identifier. It is intended to be used when setting fields owned by
// a controller to an object it does not control, e.g. a Finalizer.
//
// The returned string is of the form:
//
//	openstack.k-orc.cloud/<controllername>controller
func GetSSAFieldOwner(controllerName string) client.FieldOwner {
	return client.FieldOwner(getSSAFieldOwnerString(controllerName))
}

// GetSSAFieldOwnerWithTxn returns a field owner string for a specific named SSA
// transaction. It is intended to be used when setting fields owned by a
// controller to objects it controls.
//
// The returned string is of the form:
//
//	openstack.k-orc.cloud/<controllername>controller/<txn>
func GetSSAFieldOwnerWithTxn(controllerName string, txn SSATransactionID) client.FieldOwner {
	return client.FieldOwner(getSSAFieldOwnerString(controllerName) + "/" + string(txn))
}

// GetFinalizerName returns the finalizer to be used for the given actuator
//
// The returned string is of the form:
//
//	openstack.k-orc.cloud/<controllername>
func GetFinalizerName(controllerName string) string {
	return ORCK8SPrefix + "/" + controllerName
}
