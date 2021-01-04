package externalprovider

import (
	"github.com/openshift/installer/pkg/types/ovirt"
)

// Name is a const for convenience for external provider names.
type Name string

const (
	// NameOvirt is the name of the ovirt provider.
	NameOvirt = Name(ovirt.Name)
)
