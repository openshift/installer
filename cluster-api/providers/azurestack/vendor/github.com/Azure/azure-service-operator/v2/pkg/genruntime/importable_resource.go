/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

// ImportableResource is implemented by any resource that can be imported into the operator
type ImportableResource interface {
	// InitializeSpec initializes the Spec of the resource from the provided Status.
	InitializeSpec(status ConvertibleStatus) error
}

// ImportableARMResource represents an ARM based resource that can be imported into the operator
type ImportableARMResource interface {
	ImportableResource
	ARMMetaObject
}
