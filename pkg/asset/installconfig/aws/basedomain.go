package aws

import (
	"net/http"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
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
	session, err := getSession()
	if err != nil {
		return "", err
	}

	client := route53.New(session)
	publicZoneMap := map[string]struct{}{}
	exists := struct{}{}
	input := route53.ListHostedZonesInput{}
	for i := 0; true; i++ {
		logrus.Debugf("listing AWS hosted zones (page %d)", i)
		response, err := client.ListHostedZones(&input)
		if err != nil {
			return "", errors.Wrap(err, "list hosted zones")
		}

		for _, zone := range response.HostedZones {
			if !*zone.Config.PrivateZone {
				publicZoneMap[strings.TrimSuffix(*zone.Name, ".")] = exists
			}
		}

		if !*response.IsTruncated {
			break
		}

		input.Marker = response.NextMarker
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
