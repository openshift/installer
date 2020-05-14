package gcp

import (
	"context"
	"sort"
	"time"

	"github.com/pkg/errors"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// GetPublicZone returns a DNS managed zone from the provided project which matches the baseDomain
// If multiple zones match the basedomain, it uses the last public zone in the list as provided by the GCP API.
func GetPublicZone(ctx context.Context, project, baseDomain string) (*dns.ManagedZone, error) {
	client, err := NewClient(context.TODO())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	dnsZone, err := client.GetPublicDNSZone(ctx, project, baseDomain)
	if err != nil {
		return nil, err
	}
	return dnsZone, nil
}

// GetBaseDomain returns a base domain chosen from among the project's public DNS zones.
func GetBaseDomain(project string) (string, error) {
	client, err := NewClient(context.TODO())
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	publicZones, err := client.GetPublicDomains(ctx, project)
	if err != nil {
		return "", errors.Wrap(err, "could not retrieve base domains")
	}
	if len(publicZones) == 0 {
		return "", errors.New("no domain names found in project")
	}
	sort.Strings(publicZones)

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

// IsForbidden checks whether a response from the GPC API was forbidden,
// indicating that a given service account cannot access the specified project.
func IsForbidden(err error) bool {
	gErr, ok := err.(*googleapi.Error)
	return ok && gErr.Code == 403
}

// IsThrottled checks whether a response from the GPC API returns Too Many Requests
func IsThrottled(err error) bool {
	gErr, ok := err.(*googleapi.Error)
	return ok && gErr.Code == 429
}
