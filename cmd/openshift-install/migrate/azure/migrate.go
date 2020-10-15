package azure

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	azmigrate "github.com/openshift/installer/pkg/migrate/azure"
	"github.com/openshift/installer/pkg/types/azure"
)

var (
	azureMigrateOpts struct {
		cloudName         string
		zone              string
		resourceGroup     string
		virtualNetwork    string
		vnetResourceGroup string
		link              bool
	}
)

func runMigrateAzurePrivateDNSMigrateCmd(cmd *cobra.Command, args []string) error {
	switch azure.CloudEnvironment(azureMigrateOpts.cloudName) {
	case azure.PublicCloud, azure.USGovernmentCloud, azure.ChinaCloud, azure.GermanCloud:
	default:
		return errors.Errorf("cloud-name must be one of %s, %s, %s, %s", azure.PublicCloud, azure.USGovernmentCloud, azure.ChinaCloud, azure.GermanCloud)
	}
	if azureMigrateOpts.zone == "" {
		return errors.New("zone is a required argument")
	}
	if azureMigrateOpts.resourceGroup == "" {
		return errors.New("resource-group is a required argument")
	}
	if azureMigrateOpts.link && azureMigrateOpts.virtualNetwork == "" {
		return errors.New("link requires virtual-network to be set")
	}
	if azureMigrateOpts.virtualNetwork != "" && azureMigrateOpts.vnetResourceGroup == "" {
		return errors.New("virtual-network requires virtual-network-resource-group to be set")
	}

	return azmigrate.Migrate(
		azure.CloudEnvironment(azureMigrateOpts.cloudName),
		azureMigrateOpts.resourceGroup,
		azureMigrateOpts.zone,
		azureMigrateOpts.virtualNetwork,
		azureMigrateOpts.vnetResourceGroup,
		azureMigrateOpts.link,
	)
}

// NewMigrateAzurePrivateDNSMigrateCmd adds the migrate command to openshift-install
func NewMigrateAzurePrivateDNSMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "azure-privatedns",
		Short: "Migrate a legacy Azure zone",
		Long:  "This will migrate a legacy Azure private zone to a new style private zone.",
		RunE:  runMigrateAzurePrivateDNSMigrateCmd,
	}

	cmd.PersistentFlags().StringVar(
		&azureMigrateOpts.cloudName,
		"cloud-name",
		string(azure.PublicCloud),
		fmt.Sprintf("cloud environment name, one of: %s, %s, %s, %s", azure.PublicCloud, azure.USGovernmentCloud, azure.ChinaCloud, azure.GermanCloud),
	)
	cmd.PersistentFlags().StringVar(&azureMigrateOpts.zone, "zone", "", "The zone to migrate")
	cmd.PersistentFlags().StringVar(&azureMigrateOpts.resourceGroup, "resource-group", "", "The resource group of the zone")
	cmd.PersistentFlags().StringVar(&azureMigrateOpts.virtualNetwork, "virtual-network", "", "The virtual network to create the private zone in")
	cmd.PersistentFlags().StringVar(&azureMigrateOpts.vnetResourceGroup, "virtual-network-resource-group", "", "The resource group the virtual network is in")
	cmd.PersistentFlags().BoolVar(&azureMigrateOpts.link, "link", false, "Link the newly created private zone to the virtual network")

	return cmd
}
