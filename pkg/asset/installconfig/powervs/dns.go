package powervs

import (
	"context"
	"fmt"
	"sort"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"

	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// Zone represents a DNS Zone
type Zone struct {
	Name            string
	InstanceCRN     string
	ResourceGroupID string
	Publish         types.PublishingStrategy
}

// GetDNSZone returns a DNS Zone chosen by survey.
func GetDNSZone() (*Zone, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var options []string
	var optionToZoneMap = make(map[string]*Zone, 10)
	isInternal := ""
	strategies := []types.PublishingStrategy{types.ExternalPublishingStrategy, types.InternalPublishingStrategy}
	for _, s := range strategies {
		zones, err := client.GetDNSZones(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve base domains: %w", err)
		}

		for _, zone := range zones {
			if s == types.InternalPublishingStrategy {
				isInternal = " (Internal)"
			}
			option := fmt.Sprintf("%s%s", zone.Name, isInternal)
			optionToZoneMap[option] = &Zone{
				Name:            zone.Name,
				InstanceCRN:     zone.InstanceCRN,
				ResourceGroupID: zone.ResourceGroupID,
				Publish:         s,
			}
			options = append(options, option)
		}
	}
	sort.Strings(options)

	var zoneChoice string
	if err := survey.AskOne(&survey.Select{
		Message: "Base Domain",
		Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see your intended base-domain listed, create a new public hosted zone and rerun the installer.",
		Options: options,
	},
		&zoneChoice,
		survey.WithValidator(func(ans interface{}) error {
			choice := ans.(core.OptionAnswer).Value
			i := sort.SearchStrings(options, choice)
			if i == len(options) || options[i] != choice {
				return fmt.Errorf("invalid base domain %q", choice)
			}
			return nil
		}),
	); err != nil {
		return nil, fmt.Errorf("failed UserInput: %w", err)
	}

	return optionToZoneMap[zoneChoice], nil
}

// GetVPC returns a VPC name chosen by survey from those available in the resource group.
func GetVPC(resourceGroup string, region string) (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Resolve the VPC region from the PowerVS region.
	vpcRegion, err := powervstypes.VPCRegionForPowerVSRegion(region)
	if err != nil {
		return "", fmt.Errorf("could not determine VPC region for PowerVS region %q: %w", region, err)
	}

	// Resolve the resource group name to its ID.
	resourceGroups, err := client.ListResourceGroups(ctx)
	if err != nil {
		return "", fmt.Errorf("could not list resource groups: %w", err)
	}
	resourceGroupID := ""
	for _, rg := range resourceGroups.Resources {
		if rg.Name != nil && *rg.Name == resourceGroup {
			resourceGroupID = *rg.ID
			break
		}
	}
	if resourceGroupID == "" {
		return "", fmt.Errorf("resource group %q not found", resourceGroup)
	}

	vpcs, err := client.GetVPCsInResourceGroup(ctx, resourceGroupID, vpcRegion)
	if err != nil {
		return "", fmt.Errorf("could not retrieve VPCs: %w", err)
	}
	if len(vpcs) == 0 {
		return "", fmt.Errorf("no VPCs found in resource group %q (VPC region %q)", resourceGroup, vpcRegion)
	}

	vpcNames := make([]string, 0, len(vpcs))
	for _, vpc := range vpcs {
		if vpc.Name != nil {
			vpcNames = append(vpcNames, *vpc.Name)
		}
	}
	sort.Strings(vpcNames)

	var vpcChoice string
	if err := survey.AskOne(&survey.Select{
		Message: "VPC",
		Help:    "The VPC to use for the internal cluster. Must be in the same resource group as the PowerVS workspace.",
		Options: vpcNames,
	},
		&vpcChoice,
		survey.WithValidator(func(ans interface{}) error {
			choice := ans.(core.OptionAnswer).Value
			i := sort.SearchStrings(vpcNames, choice)
			if i == len(vpcNames) || vpcNames[i] != choice {
				return fmt.Errorf("invalid VPC %q", choice)
			}
			return nil
		}),
	); err != nil {
		return "", fmt.Errorf("failed UserInput: %w", err)
	}

	return vpcChoice, nil
}
