package gcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

// GetPublicZone returns a DNS managed zone from the provided project which matches the baseDomain
// If multiple zones match the basedomain, it uses the last public zone in the list as provided by the GCP API.
func GetPublicZone(ctx context.Context, project, baseDomain string) (managedZone *dns.ManagedZone, err error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	ssn, err := GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	svc, err := dns.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}

	if !strings.HasSuffix(baseDomain, ".") {
		baseDomain = fmt.Sprintf("%s.", baseDomain)
	}
	req := svc.ManagedZones.List(project).DnsName(baseDomain).Context(ctx)

	var res *dns.ManagedZone
	if err := req.Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for idx, v := range page.ManagedZones {
			if v.Visibility == "public" {
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
