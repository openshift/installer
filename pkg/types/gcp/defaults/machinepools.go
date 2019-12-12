package defaults

import (
	"github.com/openshift/installer/pkg/types"
)

// SetMachinePoolDefaults sets the defaults for the machine pool.
func SetMachinePoolDefaults(p *types.MachinePool) {
	if p.OSEncryption == "" {
		p.OSEncryption = types.OSEncryptionPolicyTPM2
	}
}
