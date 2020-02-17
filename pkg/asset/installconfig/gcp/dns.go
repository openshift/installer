package gcp

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Zone represents a GCP cloud DNS zone
type Zone struct {
	ID      string
	DNSName string
}

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

// GetPublicZoneByID returns a DNS managed zone from the provided project which
// matches the public zone id.
func GetPublicZoneByID(ctx context.Context, project, publicZoneID string) (*dns.ManagedZone, error) {
	client, err := NewClient(context.TODO())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	dnsZone, err := client.GetPublicDNSZoneByID(ctx, project, publicZoneID)
	if err != nil {
		return nil, err
	}
	return dnsZone, nil
}

// GetBaseDomain returns a base domain chosen from among the project's public DNS zones.
func GetBaseDomain(project string) (*Zone, error) {
	client, err := NewClient(context.TODO())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	publicDomains, err := client.GetPublicDomains(ctx, project)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve base domains")
	}

	publicZones := make([]string, 0, len(publicDomains))
	for name, ids := range publicDomains {
		for _, id := range ids {
			publicZones = append(publicZones, fmt.Sprintf("%s (%s)", name, id))
		}
	}

	if len(publicZones) == 0 {
		return nil, errors.New("no domain names found in project")
	}
	sort.Strings(publicZones)

	var publicZoneNameChoice string
	if err := survey.AskOne(&survey.Select{
		Message: "Base Domain",
		Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new public hosted zone and rerun the installer.",
		Options: publicZones,
	}, &publicZoneNameChoice, func(ans interface{}) error {
		choice := ans.(string)
		i := sort.SearchStrings(publicZones, choice)
		if i == len(publicZones) || publicZones[i] != choice {
			return errors.Errorf("invalid base domain %q", choice)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed UserInput for base domain")
	}

	parts := strings.Split(publicZoneNameChoice, " ")
	publicZoneName := parts[0]
	publicZoneID := parts[1][1 : len(parts[1])-1]

	return &Zone{
		ID:      publicZoneID,
		DNSName: publicZoneName,
	}, nil
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
