package alibabacloud

import (
	"sort"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// GetBaseDomain returns a base domain chosen from among the account's domains.
func GetBaseDomain() (string, error) {
	client, err := NewClient(defaultRegion)
	if err != nil {
		return "", err
	}

	logrus.Debugf("listing Alibaba Cloud domains")
	resp, err := client.ListDNSDomain()
	if err != nil {
		return "", err
	}

	domains := []string{}
	for _, domain := range resp.Domains.Domain {
		domains = append(domains, domain.DomainName)
	}

	sort.Strings(domains)
	if len(domains) == 0 {
		return "", errors.New("no domain found")
	}

	var basedomain string
	if err := survey.AskOne(
		&survey.Select{
			Message: "Base Domain",
			Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new domain and rerun the installer.",
			Options: domains,
		},
		&basedomain,
		survey.WithValidator(func(ans interface{}) error {
			choice := ans.(core.OptionAnswer).Value
			i := sort.SearchStrings(domains, choice)
			if i == len(domains) || domains[i] != choice {
				return errors.Errorf("invalid base domain %q", choice)
			}
			return nil
		}),
	); err != nil {
		return "", errors.Wrap(err, "failed UserInput")
	}

	return basedomain, nil
}
