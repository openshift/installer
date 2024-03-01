package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types"
)

var (
	errNotFound = errors.New("not found")
)

func getDNSZoneName(ctx context.Context, ic *installconfig.InstallConfig, isPublic bool) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	client, err := gcpic.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create new client: %w", err)
	}

	cctx, ccancel := context.WithTimeout(ctx, time.Minute*1)
	defer ccancel()

	domain := ic.Config.ClusterDomain()
	if isPublic {
		domain = ic.Config.BaseDomain
	}

	zone, err := client.GetDNSZone(cctx, ic.Config.GCP.ProjectID, domain, isPublic)
	if err != nil {
		return "", fmt.Errorf("failed to get dns zone name: %w", err)
	}

	if zone != nil {
		return zone.Name, nil
	}

	return "", errNotFound
}

type recordSet struct {
	projectID string
	zoneName  string
	record    *dns.ResourceRecordSet
}

// createRecordSets will create a list of records that will be created during the install.
func createRecordSets(ctx context.Context, ic *installconfig.InstallConfig, clusterID, apiIP, apiIntIP string) ([]recordSet, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	// A shared VPC install allows a user to preconfigure a private zone. If there is a private zone found,
	// use the existing private zone otherwise default to the installer private zone pattern below.
	privateZoneName, err := getDNSZoneName(ctx, ic, false)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, fmt.Errorf("failed to find private zone: %w", err)
		}
		privateZoneName = fmt.Sprintf("%s-private-zone", clusterID)
	}

	records := []recordSet{
		{
			// api_internal
			projectID: ic.Config.GCP.ProjectID,
			zoneName:  privateZoneName,
			record: &dns.ResourceRecordSet{
				Name:    fmt.Sprintf("api-int.%s.", ic.Config.ClusterDomain()),
				Type:    "A",
				Ttl:     60,
				Rrdatas: []string{apiIntIP},
			},
		},
		{
			// api_external_internal_zone
			projectID: ic.Config.GCP.ProjectID,
			zoneName:  privateZoneName,
			record: &dns.ResourceRecordSet{
				Name:    fmt.Sprintf("api.%s.", ic.Config.ClusterDomain()),
				Type:    "A",
				Ttl:     60,
				Rrdatas: []string{apiIntIP},
			},
		},
	}

	if ic.Config.Publish == types.ExternalPublishingStrategy {
		existingPublicZoneName, err := getDNSZoneName(ctx, ic, true)
		if err != nil {
			return nil, fmt.Errorf("failed to find a public zone: %w", err)
		}

		apiRecord := recordSet{
			projectID: ic.Config.GCP.ProjectID,
			zoneName:  existingPublicZoneName,
			record: &dns.ResourceRecordSet{
				Name:    fmt.Sprintf("api.%s.", ic.Config.ClusterDomain()),
				Type:    "A",
				Ttl:     60,
				Rrdatas: []string{apiIP},
			},
		}
		records = append(records, apiRecord)
	}

	return records, nil
}

// createDNSRecords will get the list of records to be created and execute their creation through the gcp dns api.
func createDNSRecords(ctx context.Context, ic *installconfig.InstallConfig, clusterID, apiIP, apiIntIP string) error {
	ssn, err := gcpic.GetSession(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	// TODO: use the opts for the service to restrict scopes see google.golang.org/api/option.WithScopes
	dnsService, err := dns.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return fmt.Errorf("failed to create the gcp dns service: %w", err)
	}

	records, err := createRecordSets(ctx, ic, clusterID, apiIP, apiIntIP)
	if err != nil {
		return err
	}

	// 1 minute timeout for each record
	ctx, cancel := context.WithTimeout(ctx, time.Minute*time.Duration(len(records)))
	defer cancel()

	for _, record := range records {
		if _, err := dnsService.ResourceRecordSets.Create(record.projectID, record.zoneName, record.record).Context(ctx).Do(); err != nil {
			return fmt.Errorf("failed to create record set %s: %w", record.record.Name, err)
		}
	}

	return nil
}

// createPrivateManagedZone will create a private managed zone in the GCP project specified in the install config. The
// private managed zone should only be created when one is not specified in the install config.
func createPrivateManagedZone(ctx context.Context, ic *installconfig.InstallConfig, clusterID, network string) error {
	// TODO: use the opts for the service to restrict scopes see google.golang.org/api/option.WithScopes
	ssn, err := gcpic.GetSession(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	dnsService, err := dns.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return fmt.Errorf("failed to create the gcp dns service: %w", err)
	}

	managedZone := &dns.ManagedZone{
		Name:        fmt.Sprintf("%s-private-zone", clusterID),
		Description: resourceDescription,
		DnsName:     fmt.Sprintf("%s.", ic.Config.ClusterDomain()),
		Visibility:  "private",
		Labels:      mergeLabels(ic, clusterID),
		PrivateVisibilityConfig: &dns.ManagedZonePrivateVisibilityConfig{
			Networks: []*dns.ManagedZonePrivateVisibilityConfigNetwork{
				{
					NetworkUrl: network,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	if _, err = dnsService.ManagedZones.Create(ic.Config.GCP.ProjectID, managedZone).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to create private managed zone: %w", err)
	}

	return nil
}
