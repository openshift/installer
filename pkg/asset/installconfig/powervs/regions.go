package powervs

import (
	"fmt"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/pkg/errors"
)

func knownRegions() map[string]string {

	regions := make(map[string]string)

	for _, region := range rhcos.PowerVSRegions {
		regions[region.Name] = region.Description
	}
	return regions
}

// IsKnownRegion return true is a specified region is Known to the installer.
// A known region is subset of AWS regions and the regions where RHEL CoreOS images are published.
func IsKnownRegion(region string) bool {
	if _, ok := knownRegions()[region]; ok {
		return true
	}
	return false
}

func knownZones(region string) []string {
	if _, ok := rhcos.PowerVSRegions[region]; ok {
		return rhcos.PowerVSRegions[region].Zones
	}
	return []string{}
}

// IsKnownZone return true is a specified zone is known to the installer.
func IsKnownZone(region string, zone string) bool {
	zones := knownZones(region)
	for _, z := range zones {
		if z == zone {
			return true
		}
	}
	return false
}

// GetRegion prompts the user to select a region and returns that region
func GetRegion() (string, error) {
	regions := knownRegions()

	longRegions := make([]string, 0, len(regions))
	shortRegions := make([]string, 0, len(regions))
	for id, location := range regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	var regionTransform survey.Transformer = func(ans interface{}) interface{} {
		switch v := ans.(type) {
		case core.OptionAnswer:
			return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
		case string:
			return strings.SplitN(v, " ", 2)[0]
		}
		return ""
	}

	var region string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The Power VS region to be used for installation.",
				Options: longRegions,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := regionTransform(ans).(core.OptionAnswer).Value
				i := sort.SearchStrings(shortRegions, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return errors.Errorf("Invalid region %q", choice)
				}
				return nil
			}),
			Transform: regionTransform,
		},
	}, &region)
	if err != nil {
		return "", err
	}

	return region, nil
}

// GetZone prompts the user for a zone given a zone
func GetZone(region string) (string, error) {
	zones := knownZones(region)
	defaultZone := zones[0]

	var zoneTransform survey.Transformer = func(ans interface{}) interface{} {
		switch v := ans.(type) {
		case core.OptionAnswer:
			return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
		case string:
			return strings.SplitN(v, " ", 2)[0]
		}
		return ""
	}

	var zone string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Zone",
				Help:    "The Power VS zone within the region to be used for installation.",
				Default: fmt.Sprintf("%s", defaultZone),
				Options: zones,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := zoneTransform(ans).(core.OptionAnswer).Value
				i := sort.SearchStrings(zones, choice)
				if i == len(zones) || zones[i] != choice {
					return errors.Errorf("Invalid zone %q", choice)
				}
				return nil
			}),
			Transform: zoneTransform,
		},
	}, &zone)
	if err != nil {
		return "", err
	}
	return zone, err
}
