/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

// LocatableResource represents a resource with a location.
type LocatableResource interface {
	Location() string
}
