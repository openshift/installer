package azure

import (
	"github.com/spf13/cobra"

	azmigrate "github.com/openshift/installer/pkg/migrate/azure"
)

func runMigrateAzurePrivateDNSEligibleCmd(cmd *cobra.Command, args []string) error {
	return azmigrate.Eligible()
}

// NewMigrateAzurePrivateDNSEligibleCmd adds the eligble command to openshift-install
func NewMigrateAzurePrivateDNSEligibleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "azure-privatedns-eligible",
		Short: "Show legacy Azure zones that are eligible to be migrated",
		Long:  "This will show legacy Azure private zones that can be migrated to new private zones.",
		RunE:  runMigrateAzurePrivateDNSEligibleCmd,
	}

	return cmd
}
