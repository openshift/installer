package aws

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
)

// IsForbidden returns true if and only if the input error is an HTTP
// 403 error from the AWS API.
func IsForbidden(err error) bool {
	var requestError awserr.RequestFailure
	return errors.As(err, &requestError) && requestError.StatusCode() == http.StatusForbidden
}

// GetBaseDomain returns a base domain chosen from among the account's
// public routes.
func GetBaseDomain() (string, error) {
	session, err := GetSession()
	if err != nil {
		return "", err
	}

	logrus.Debugf("listing AWS hosted zones")
	client := route53.New(session)
	publicZoneMap := map[string]struct{}{}
	exists := struct{}{}
	if err := client.ListHostedZonesPages(
		&route53.ListHostedZonesInput{},
		func(resp *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
			for _, zone := range resp.HostedZones {
				if zone.Config != nil && !aws.BoolValue(zone.Config.PrivateZone) {
					publicZoneMap[strings.TrimSuffix(*zone.Name, ".")] = exists
				}
			}
			return !lastPage
		},
	); err != nil {
		return "", fmt.Errorf("list hosted zones: %w", err)
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
func GetPublicZone(sess *session.Session, name string) (*route53.HostedZone, error) {
	var res *route53.HostedZone
	f := func(resp *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
		for idx, zone := range resp.HostedZones {
			if zone.Config != nil && !aws.BoolValue(zone.Config.PrivateZone) && strings.TrimSuffix(aws.StringValue(zone.Name), ".") == strings.TrimSuffix(name, ".") {
				res = resp.HostedZones[idx]
				return false
			}
		}
		return !lastPage
	}

	client := route53.New(sess)
	if err := client.ListHostedZonesPages(&route53.ListHostedZonesInput{}, f); err != nil {
		return nil, fmt.Errorf("listing hosted zones: %w", err)
	}
	if res == nil {
		return nil, fmt.Errorf("no public route53 zone found matching name %q", name)
	}
	return res, nil
}
