package aws

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/sirupsen/logrus"
)

// GetBaseDomain returns a base domain chosen from among the account's
// public routes.
func GetBaseDomain(ctx context.Context, client *route53.Client) (string, error) {
	logrus.Debugf("listing AWS hosted zones")

	publicZoneMap := map[string]struct{}{}
	exists := struct{}{}

	paginator := route53.NewListHostedZonesPaginator(client, &route53.ListHostedZonesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to list hosted zones: %w", err)
		}

		for _, zone := range page.HostedZones {
			if zone.Config != nil && !zone.Config.PrivateZone {
				publicZoneMap[strings.TrimSuffix(aws.ToString(zone.Name), ".")] = exists
			}
		}
	}

	publicZones := make([]string, 0, len(publicZoneMap))
	for name := range publicZoneMap {
		publicZones = append(publicZones, name)
	}
	sort.Strings(publicZones)
	if len(publicZones) == 0 {
		return "", errors.New("no public Route 53 hosted zones found")
	}

	var domain string
	if err := survey.AskOne(
		&survey.Select{
			Message: "Base Domain",
			Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new public Route53 hosted zone and rerun the installer.",
			Options: publicZones,
		},
		&domain,
		survey.WithValidator(func(ans interface{}) error {
			choice := ans.(core.OptionAnswer).Value
			i := sort.SearchStrings(publicZones, choice)
			if i == len(publicZones) || publicZones[i] != choice {
				return fmt.Errorf("invalid base domain %q", choice)
			}
			return nil
		}),
	); err != nil {
		return "", fmt.Errorf("failed UserInput: %w", err)
	}

	return domain, nil
}

// GetPublicZone returns a public route53 zone that matches the name.
func GetPublicZone(ctx context.Context, client *route53.Client, name string) (*route53types.HostedZone, error) {
	paginator := route53.NewListHostedZonesPaginator(client, &route53.ListHostedZonesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list hosted zones: %w", err)
		}

		for _, zone := range page.HostedZones {
			if zone.Config != nil && !zone.Config.PrivateZone && strings.TrimSuffix(aws.ToString(zone.Name), ".") == strings.TrimSuffix(name, ".") {
				return &zone, nil
			}
		}
	}

	return nil, fmt.Errorf("no public route53 zone found matching name %q", name)
}
