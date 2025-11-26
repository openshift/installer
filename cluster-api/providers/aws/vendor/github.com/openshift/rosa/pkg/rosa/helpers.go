package rosa

import (
	"os"

	"github.com/spf13/cobra"
)

const hostedCpFlagName = "hosted-cp"

func HostedClusterOnlyFlag(r *Runtime, cmd *cobra.Command, flagName string) {
	if cmd.Flag(hostedCpFlagName) == nil || (cmd.Flag(hostedCpFlagName) != nil && !cmd.Flag(hostedCpFlagName).Changed) {
		isFlagSet := cmd.Flags().Changed(flagName)
		if isFlagSet {
			r.Reporter.Errorf("Setting the `%s` flag is only supported for Hosted Control Plane clusters",
				flagName)
			os.Exit(1)
		}
	}
}
