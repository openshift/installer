/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package extensions

import (
	"context"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// Importer is an optional interface that can be implemented by resource extensions to customize the import process.
// This extension is invoked during 'asoctl import' operations, after retrieving the resource from Azure but before
// writing it to Kubernetes. It allows resources to skip import for system-managed resources, read-only configurations,
// or resources that only have default values.
// Implement this extension when:
// - System-managed or auto-created resources should be excluded from import
// - Default or empty configurations don't need management
// - Resources need validation before allowing import
type Importer interface {
	// Import allows interception of the import process to skip or modify resources being imported.
	// ctx is the current asynchronous context
	// rsrc is the resource being imported.
	// owner is an optional owner for the resource.
	// next is a function to call to do the actual import.
	// Returns an ImportResult indicating success or skip with a reason, and an error if import fails.
	Import(
		ctx context.Context,
		rsrc genruntime.ImportableResource,
		owner *genruntime.ResourceReference,
		next ImporterFunc,
	) (ImportResult, error)
}

// ImportResult is the result of doing an import.
type ImportResult struct {
	because string
}

// ImporterFunc is the signature of the function that does the actual import.
type ImporterFunc func(
	ctx context.Context,
	resource genruntime.ImportableResource,
	owner *genruntime.ResourceReference,
) (ImportResult, error)

// ImportSucceeded creates a new ImportResult with a resource that was imported successfully.
func ImportSucceeded() ImportResult {
	return ImportResult{}
}

// ImportSkipped creates a new ImportResult for a resource that was not imported.
func ImportSkipped(because string) ImportResult {
	return ImportResult{
		because: because,
	}
}

// Skipped returns a reason and true if the import was skipped, empty string and false otherwise.
func (r ImportResult) Skipped() (string, bool) {
	return r.because, r.because != ""
}
