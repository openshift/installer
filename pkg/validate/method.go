package validate

import "os"

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
