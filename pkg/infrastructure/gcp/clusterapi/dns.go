package clusterapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/dns/v1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types"
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

	zone, err := client.GetDNSZone(cctx, ic.Config.GCP.ProjectID, ic.Config.ClusterDomain(), isPublic)
	if err != nil {
		return "", fmt.Errorf("failed to get dns zone name: %w", err)
	}

	if zone != nil {
		return zone.Name, nil
	}

	return "", nil
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

	privateZoneName := fmt.Sprintf("%s-private-zone", clusterID)
	potentialPrivateZoneName, err := getDNSZoneName(ctx, ic, false)
	if err != nil {
		return nil, err
	}
	if potentialPrivateZoneName != "" {
		privateZoneName = potentialPrivateZoneName
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
		potentialPublicZoneName, err := getDNSZoneName(ctx, ic, true)
		if err != nil {
			return nil, err
		}

		apiRecord := recordSet{
			projectID: ic.Config.GCP.ProjectID,
			zoneName:  potentialPublicZoneName,
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
	// TODO: use the opts for the service to restrict scopes see google.golang.org/api/option.WithScopes
	dnsService, err := dns.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create the gcp dns service: %w", err)
	}

	records, err := createRecordSets(ctx, ic, clusterID, apiIP, apiIntIP)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
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
func createPrivateManagedZone(ctx context.Context, ic *installconfig.InstallConfig, clusterID string) error {
	// TODO: use the opts for the service to restrict scopes see google.golang.org/api/option.WithScopes
	dnsService, err := dns.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create the gcp dns service: %w", err)
	}

	labels := mergeLabels(ic, clusterID)

	managedZone := &dns.ManagedZone{
		Name:        fmt.Sprintf("%s-private-zone", clusterID),
		Description: "Created By OpenShift Installer",
		DnsName:     fmt.Sprintf("%s.", ic.Config.ClusterDomain()),
		Visibility:  "private",
		Labels:      labels,
		PrivateVisibilityConfig: &dns.ManagedZonePrivateVisibilityConfig{
			Networks: []*dns.ManagedZonePrivateVisibilityConfigNetwork{
				{
					NetworkUrl: ic.Config.GCP.Network,
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
