package main

import (
	"github.com/spf13/cobra"

	azure "github.com/openshift/installer/cmd/openshift-install/migrate/azure"
)

func newMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Do a migration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	migrateCmd.AddCommand(azure.NewMigrateAzurePrivateDNSEligibleCmd())
	migrateCmd.AddCommand(azure.NewMigrateAzurePrivateDNSMigrateCmd())

	return migrateCmd
}
