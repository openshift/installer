// Package init contains the list of backends that can be initialized and
// basic helper functions for initializing those backends.
package init

import (
	"sync"

	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform-plugin-sdk/tfdiags"
	"github.com/zclconf/go-cty/cty"

	backendLocal "github.com/hashicorp/terraform/backend/local"
)

// backends is the list of available backends. This is a global variable
// because backends are currently hardcoded into Terraform and can't be
// modified without recompilation.
//
// To read an available backend, use the Backend function. This ensures
// safe concurrent read access to the list of built-in backends.
//
// Backends are hardcoded into Terraform because the API for backends uses
// complex structures and supporting that over the plugin system is currently
// prohibitively difficult. For those wanting to implement a custom backend,
// they can do so with recompilation.
var backends map[string]backend.InitFn
var backendsLock sync.Mutex

// Init initializes the backends map with all our hardcoded backends.
func Init(services *disco.Disco) {
	backendsLock.Lock()
	defer backendsLock.Unlock()

	backends = map[string]backend.InitFn{
		// Enhanced backends.
		"local": func() backend.Backend { return backendLocal.New() },
	}
}

// Backend returns the initialization factory for the given backend, or
// nil if none exists.
func Backend(name string) backend.InitFn {
	backendsLock.Lock()
	defer backendsLock.Unlock()
	return backends[name]
}

// Set sets a new backend in the list of backends. If f is nil then the
// backend will be removed from the map. If this backend already exists
// then it will be overwritten.
//
// This method sets this backend globally and care should be taken to do
// this only before Terraform is executing to prevent odd behavior of backends
// changing mid-execution.
func Set(name string, f backend.InitFn) {
	backendsLock.Lock()
	defer backendsLock.Unlock()

	if f == nil {
		delete(backends, name)
		return
	}

	backends[name] = f
}

// deprecatedBackendShim is used to wrap a backend and inject a deprecation
// warning into the Validate method.
type deprecatedBackendShim struct {
	backend.Backend
	Message string
}

// PrepareConfig delegates to the wrapped backend to validate its config
// and then appends shim's deprecation warning.
func (b deprecatedBackendShim) PrepareConfig(obj cty.Value) (cty.Value, tfdiags.Diagnostics) {
	newObj, diags := b.Backend.PrepareConfig(obj)
	return newObj, diags.Append(tfdiags.SimpleWarning(b.Message))
}

// DeprecateBackend can be used to wrap a backend to retrun a deprecation
// warning during validation.
func deprecateBackend(b backend.Backend, message string) backend.Backend {
	// Since a Backend wrapped by deprecatedBackendShim can no longer be
	// asserted as an Enhanced or Local backend, disallow those types here
	// entirely.  If something other than a basic backend.Backend needs to be
	// deprecated, we can add that functionality to schema.Backend or the
	// backend itself.
	if _, ok := b.(backend.Enhanced); ok {
		panic("cannot use DeprecateBackend on an Enhanced Backend")
	}

	if _, ok := b.(backend.Local); ok {
		panic("cannot use DeprecateBackend on a Local Backend")
	}

	return deprecatedBackendShim{
		Backend: b,
		Message: message,
	}
}
