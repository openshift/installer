package powervs

import (
	"fmt"
	"os"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Platform collects powervs-specific configuration.
func Platform() (*powervs.Platform, error) {
	regions := knownRegions()

	// TODO(cklokman): This section came from aws and transforms the response from knownRegions
	//                 into long and short regions to prompt the user for region select this section
	//                 need need to be different based on powervs's implementation of knownRegions
	//

	longRegions := make([]string, 0, len(regions))
	shortRegions := make([]string, 0, len(regions))
	for id, location := range regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}

	var regionTransform survey.Transformer = func(ans interface{}) interface{} {
		switch v := ans.(type) {
		case core.OptionAnswer:
			return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
		case string:
			return strings.SplitN(v, " ", 2)[0]
		}
		return ""
	}

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	var region string

	sessionRegion := ssn.Session.Region
	if sessionRegion != "" {
		if IsKnownRegion(sessionRegion) {
			region = sessionRegion
		} else {
			logrus.Warnf("Unrecognized Power VS region %s, ignoring IC_REGION", sessionRegion)
		}
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)
	if region == "" {
		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Select{
					Message: "Region",
					Help:    "The Power VS region to be used for installation.",
					// Default: fmt.Sprintf("%s (%s)", defaultRegion, regions[defaultRegion]),
					Options: longRegions,
				},
				Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
					choice := regionTransform(ans).(core.OptionAnswer).Value
					i := sort.SearchStrings(shortRegions, choice)
					if i == len(shortRegions) || shortRegions[i] != choice {
						return errors.Errorf("invalid region %q", choice)
					}
					return nil
				}),
				Transform: regionTransform,
			},
		}, &region)
		if err != nil {
			return nil, err
		}
	}

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
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Zone",
				Help:    "The powervs zone within the region to be used for installation.",
				Default: fmt.Sprintf("%s", defaultZone),
				Options: zones,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := zoneTransform(ans).(core.OptionAnswer).Value
				i := sort.SearchStrings(zones, choice)
				if i == len(zones) || zones[i] != choice {
					return errors.Errorf("invalid zone %q", choice)
				}
				return nil
			}),
			Transform: zoneTransform,
		},
	}, &zone)
	if err != nil {
		return nil, err
	}

	var p powervs.Platform
	if osOverride := os.Getenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); len(osOverride) != 0 {
		p.BootstrapOSImage = osOverride
		p.ClusterOSImage = osOverride
	}

	p.Region = region
	p.Zone = zone
	p.APIKey = ssn.Creds.APIKey
	p.UserID = ssn.Creds.UserID

	return &p, nil
}
