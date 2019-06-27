package aws

import (
	"net/http"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Zone represents an AWS route53 DNS Zone
type Zone struct {
	ID   string
	Name string
}

// IsForbidden returns true if and only if the input error is an HTTP
// 403 error from the AWS API.
func IsForbidden(err error) bool {
	requestError, ok := err.(awserr.RequestFailure)
	return ok && requestError.StatusCode() == http.StatusForbidden
}

// GetBaseDomain returns a base domain zone chosen from among the
// account's public routes. The zone struct contains the zone ID
// for the chosen base domain, to distinguish between multiple
// public zones that share the same base domain.
func GetBaseDomain() (zone *Zone, err error) {
	session, err := GetSession()
	if err != nil {
		return nil, err
	}

	logrus.Debugf("listing AWS hosted zones")
	var publicZoneName string
	var publicZoneID string
	client := route53.New(session)
	publicZoneMap := map[string][]string{}
	if err := client.ListHostedZonesPages(
		&route53.ListHostedZonesInput{},
		func(resp *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
			for _, zone := range resp.HostedZones {
				if zone.Config != nil && !aws.BoolValue(zone.Config.PrivateZone) {
					zoneName := strings.TrimSuffix(aws.StringValue(zone.Name), ".")
					publicZoneID = strings.TrimPrefix(aws.StringValue(zone.Id), "/hostedzone/")
					publicZoneMap[zoneName] = append(publicZoneMap[zoneName], publicZoneID)
				}
			}
			return !lastPage
		},
	); err != nil {
		return nil, errors.Wrap(err, "list hosted zones")
	}

	publicZones := make([]string, 0, len(publicZoneMap))
	for name := range publicZoneMap {
		publicZones = append(publicZones, name)
	}
	sort.Strings(publicZones)
	if len(publicZones) == 0 {
		return nil, errors.New("no public Route 53 hosted zones found")
	}

	if err := survey.AskOne(&survey.Select{
		Message: "Base Domain",
		Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new public Route53 hosted zone and rerun the installer.",
		Options: publicZones,
	}, &publicZoneName, func(ans interface{}) error {
		choice := ans.(string)
		i := sort.SearchStrings(publicZones, choice)
		if i == len(publicZones) || publicZones[i] != choice {
			return errors.Errorf("invalid base domain %q", choice)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed UserInput for base domain")
	}

	zoneList, ok := publicZoneMap[publicZoneName]
	if !ok {
		return nil, errors.Errorf("invalid base domain %q", publicZoneName)
	}

	zoneListLen := len(zoneList)
	if zoneListLen == 0 {
		return nil, errors.Wrap(err, "there must be at least one public zones")
	}
	if zoneListLen == 1 {
		return &Zone{
			ID:   zoneList[0],
			Name: publicZoneName,
		}, nil
	}

	err = survey.AskOne(&survey.Select{
		Message: "Zone ID",
		Help:    "The zone ID for the base domain of the cluster",
		Options: zoneList,
	}, &publicZoneID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed UserInput for zone ID")
	}

	return &Zone{
		ID:   publicZoneID,
		Name: publicZoneName,
	}, nil
}

// GetPublicZone returns a public route53 zone that matches the name.
func GetPublicZone(name string) (*route53.HostedZone, error) {
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

	session, err := GetSession()
	if err != nil {
		return nil, errors.Wrap(err, "getting AWS session")
	}
	client := route53.New(session)
	if err := client.ListHostedZonesPages(&route53.ListHostedZonesInput{}, f); err != nil {
		return nil, errors.Wrap(err, "listing hosted zones")
	}
	if res == nil {
		return nil, errors.Errorf("No public route53 zone found matching name %q", name)
	}
	return res, nil
}
