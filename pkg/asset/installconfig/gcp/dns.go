package gcp

import (
	"golang.org/x/net/context"
	"google.golang.org/api/dns/v1"
)

// GetPublicZone returns a DNS managed zone from the provided project which matches the baseDomain
func GetPublicZone(baseDomain, project string) (managedZone *dns.ManagedZone, err error) {
	ctx := context.Background()
	//TODO: Determine if this will be proper/acceptable way to invoke services with credentials
	dnsService, err := dns.NewService(ctx)
	var res dns.ManagedZone
	req := dnsService.ManagedZones.List(project)

	//TODO: How are we selecting managed zones? This returns the last match.
	// Check for public/private?
	if err := req.Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for _, managedZone := range page.ManagedZones {
			//fmt.Printf("%#v\n", managedZone)
			res = *managedZone
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &res, nil
}
