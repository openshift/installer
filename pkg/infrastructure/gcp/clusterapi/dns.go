package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/dns/v1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/types"
)

var (
	errNotFound = errors.New("not found")
)

func getDNSZoneName(ctx context.Context, ic *installconfig.InstallConfig, clusterID string, isPublic bool) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	client, err := gcpic.NewClient(ctx, ic.Config.GCP.ServiceEndpoints)
	if err != nil {
		return "", fmt.Errorf("failed to create new client: %w", err)
	}

	cctx, ccancel := context.WithTimeout(ctx, time.Minute*1)
	defer ccancel()

	params, err := manifests.GetGCPPrivateZoneInfo(ctx, client, ic, clusterID)
	if err != nil {
		return "", fmt.Errorf("failed to get private zone info: %w", err)
	}

	domain := params.BaseDomain
	project := params.Project
	if isPublic {
		project = ic.Config.GCP.ProjectID
		domain = ic.Config.BaseDomain
	}

	zone, err := client.GetDNSZone(cctx, project, domain, isPublic)
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
func createRecordSets(ctx context.Context, client *gcpic.Client, ic *installconfig.InstallConfig, clusterID, apiIP, apiIntIP string) ([]recordSet, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	privateZoneParams, err := manifests.GetGCPPrivateZoneInfo(ctx, client, ic, clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get private zone info for record creation: %w", err)
	}

	records := []recordSet{
		{
			// api_internal
			projectID: privateZoneParams.Project,
			zoneName:  privateZoneParams.Name,
			record: &dns.ResourceRecordSet{
				Name:    fmt.Sprintf("api-int.%s.", ic.Config.ClusterDomain()),
				Type:    "A",
				Ttl:     60,
				Rrdatas: []string{apiIntIP},
			},
		},
		{
			// api_external_internal_zone
			projectID: privateZoneParams.Project,
			zoneName:  privateZoneParams.Name,
			record: &dns.ResourceRecordSet{
				Name:    fmt.Sprintf("api.%s.", ic.Config.ClusterDomain()),
				Type:    "A",
				Ttl:     60,
				Rrdatas: []string{apiIntIP},
			},
		},
	}

	if ic.Config.Publish == types.ExternalPublishingStrategy {
		existingPublicZoneName, err := getDNSZoneName(ctx, ic, clusterID, true)
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
func createDNSRecords(ctx context.Context, client *gcpic.Client, ic *installconfig.InstallConfig, clusterID, apiIP, apiIntIP string) error {
	// TODO: use the opts for the service to restrict scopes see google.golang.org/api/option.WithScopes
	dnsService, err := gcpic.GetDNSService(ctx, ic.Config.GCP.ServiceEndpoints)
	if err != nil {
		return fmt.Errorf("failed to create the gcp dns service: %w", err)
	}

	records, err := createRecordSets(ctx, client, ic, clusterID, apiIP, apiIntIP)
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
	client, err := gcpic.NewClient(ctx, ic.Config.GCP.ServiceEndpoints)
	if err != nil {
		return err
	}

	params, err := manifests.GetGCPPrivateZoneInfo(ctx, client, ic, clusterID)
	if err != nil {
		return err
	}

	if ic.Config.GCP.NetworkProjectID != "" {
		if !params.InstallerCreated {
			logrus.Debugf("found private zone %s, skipping creation of private zone", params.Name)
			// The private zone already exists, so we need to add the shared label to the zone.
			labels := mergeLabels(ic, clusterID, sharedLabelValue)
			if err := client.UpdateDNSPrivateZoneLabels(ctx, ic.Config.ClusterDomain(), params.Project, params.Name, labels); err != nil {
				return fmt.Errorf("failed to update dns private zone labels: %w", err)
			}
			return nil
		}
	}
	logrus.Debugf("creating private zone %s", params.Name)

	// TODO: use the opts for the service to restrict scopes see google.golang.org/api/option.WithScopes
	dnsService, err := gcpic.GetDNSService(ctx, ic.Config.GCP.ServiceEndpoints)
	if err != nil {
		return fmt.Errorf("failed to create the gcp dns service: %w", err)
	}

	managedZone := &dns.ManagedZone{
		Name:        params.Name,
		Description: resourceDescription,
		DnsName:     fmt.Sprintf("%s.", ic.Config.ClusterDomain()),
		Visibility:  "private",
		Labels:      mergeLabels(ic, clusterID, ownedLabelValue),
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

	if _, err = dnsService.ManagedZones.Create(params.Project, managedZone).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to create private managed zone: %w", err)
	}

	return nil
}
