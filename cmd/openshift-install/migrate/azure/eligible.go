package azure

import (
	"fmt"

	"github.com/spf13/cobra"

	azmigrate "github.com/openshift/installer/pkg/migrate/azure"
	"github.com/openshift/installer/pkg/types/azure"
)

// NewMigrateAzurePrivateDNSEligibleCmd adds the eligble command to openshift-install
func NewMigrateAzurePrivateDNSEligibleCmd() *cobra.Command {
	var cloudName string

	cmd := &cobra.Command{
		Use:   "azure-privatedns-eligible",
		Short: "Show legacy Azure zones that are eligible to be migrated",
		Long:  "This will show legacy Azure private zones that can be migrated to new private zones.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return azmigrate.Eligible(azure.CloudEnvironment(cloudName))
		},
	}

	cmd.Flags().StringVar(
		&cloudName,
		"cloud-name",
		string(azure.PublicCloud),
		fmt.Sprintf("cloud environment name, one of: %s, %s, %s, %s", azure.PublicCloud, azure.USGovernmentCloud, azure.ChinaCloud, azure.GermanCloud),
	)

	return cmd
}
