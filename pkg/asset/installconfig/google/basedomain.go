package google

import (
	"sort"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// GetBaseDomain returns a base domain chosen from among the account's
// public routes.
func GetBaseDomain() (string, error) {

	publicZones := make([]string, 0)
	sort.Strings(publicZones)
	if len(publicZones) == 0 {
		return "", errors.New("no public hosted zones found")
	}

	var domain string
	if err := survey.AskOne(&survey.Select{
		Message: "Base Domain",
		Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new public hosted zone and rerun the installer.",
		Options: publicZones,
	}, &domain, func(ans interface{}) error {
		choice := ans.(string)
		i := sort.SearchStrings(publicZones, choice)
		if i == len(publicZones) || publicZones[i] != choice {
			return errors.Errorf("invalid base domain %q", choice)
		}
		return nil
	}); err != nil {
		return "", errors.Wrap(err, "failed UserInput for base domain")
	}

	return domain, nil
}
