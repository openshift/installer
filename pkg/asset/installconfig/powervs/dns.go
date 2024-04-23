package powervs

import (
	"context"
	"fmt"
	"sort"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"

	"github.com/openshift/installer/pkg/types"
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
