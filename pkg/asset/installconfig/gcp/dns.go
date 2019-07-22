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
	"google.golang.org/api/option"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getDNSService(ctx context.Context) (*dns.Service, error) {
	ssn, err := GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	svc, err := dns.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}
	return svc, nil
}

// GetPublicZone returns a DNS managed zone from the provided project which matches the baseDomain
// If multiple zones match the basedomain, it uses the last public zone in the list as provided by the GCP API.
func GetPublicZone(ctx context.Context, project, baseDomain string) (*dns.ManagedZone, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	svc, err := getDNSService(ctx)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(baseDomain, ".") {
		baseDomain = fmt.Sprintf("%s.", baseDomain)
	}
	req := svc.ManagedZones.List(project).DnsName(baseDomain).Context(ctx)

	var res *dns.ManagedZone
	if err := req.Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for idx, v := range page.ManagedZones {
			if v.Visibility != "private" {
				res = page.ManagedZones[idx]
			}
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to list DNS Zones")
	}
	if res == nil {
		return nil, errors.New("no matching public DNS Zone found")
	}
	return res, nil
}

// GetBaseDomain returns a base domain chosen from among the project's public DNS zones.
func GetBaseDomain(project string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	svc, err := getDNSService(ctx)
	if err != nil {
		return "", err
	}

	var publicZones []string
	req := svc.ManagedZones.List(project).Context(ctx)
	if err := req.Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for _, v := range page.ManagedZones {
			if v.Visibility != "private" {
				publicZones = append(publicZones, strings.TrimSuffix(v.DnsName, "."))
			}
		}
		return nil
	}); err != nil {
		return "", err
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
