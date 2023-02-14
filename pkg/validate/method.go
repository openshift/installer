package validate

import "os"

// IsAgentBasedInstallation determines whether we are using the 'agent'
// subcommand to install
func IsAgentBasedInstallation() bool {
	if len(os.Args) > 1 {
		for _, arg := range os.Args {
			if arg == "agent" {
				return true
			}
		}
	}
	return false
}
