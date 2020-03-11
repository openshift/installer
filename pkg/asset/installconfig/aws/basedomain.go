package aws

import (
	"net/http"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// IsForbidden returns true if and only if the input error is an HTTP
// 403 error from the AWS API.
func IsForbidden(err error) bool {
	requestError, ok := err.(awserr.RequestFailure)
	return ok && requestError.StatusCode() == http.StatusForbidden
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
		return "", errors.Wrap(err, "list hosted zones")
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
	if err := survey.AskOne(&survey.Select{
		Message: "Base Domain",
		Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new public Route53 hosted zone and rerun the installer.",
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
		return nil, errors.Wrap(err, "listing hosted zones")
	}
	if res == nil {
		return nil, errors.Errorf("No public route53 zone found matching name %q", name)
	}
	return res, nil
}
