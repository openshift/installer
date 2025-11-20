// go mod won't pull in code that isn't depended upon, but we have some code we don't depend on from code that must be included
// for our build to work.
package dependencymagnet

import (
	// this gives us clear dependency control of our generator, easy replaces for development, ease of vendored inspection, and fully local builds.
	_ "github.com/openshift-online/ocm-api-model/model/dependencymagnet"
)
