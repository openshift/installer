package powervs

import (
	"context"
	"fmt"
	"sort"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// Zone represents a DNS Zone
type Zone struct {
	Name            string
	CISInstanceCRN  string
	ResourceGroupID string
}

// GetDNSZone returns a DNS Zone chosen by survey.
func GetDNSZone() (*Zone, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Default to offering only external DNS entries, as IBM Cloud does in their provider
	// TODO(mjturek): Offer private zones (IBM DNS) when deploying a private cluster
	publicZones, err := client.GetDNSZones(ctx, types.ExternalPublishingStrategy)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve base domains")
	}

	var options []string
	var optionToZoneMap = make(map[string]*Zone, len(publicZones))
	for _, zone := range publicZones {
		option := fmt.Sprintf("%s (%s)", zone.Name, zone.InstanceName)
		optionToZoneMap[option] = &Zone{
			Name:            zone.Name,
			CISInstanceCRN:  zone.InstanceCRN,
			ResourceGroupID: zone.ResourceGroupID,
		}
		options = append(options, option)
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
			if i == len(publicZones) || options[i] != choice {
				return errors.Errorf("invalid base domain %q", choice)
			}
			return nil
		}),
	); err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	return optionToZoneMap[zoneChoice], nil
}
