package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

var (
	errNotFound = errors.New("not found")
)

func getDNSZoneName(ctx context.Context, ic *installconfig.InstallConfig, clusterID string, isPublic bool) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	client, err := gcpic.NewClient(ctx, ic.Config.GCP.Endpoint)
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
	opts := []option.ClientOption{option.WithScopes(dns.CloudPlatformScope)}
	endpoint := ic.Config.GCP.Endpoint
	if gcptypes.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcpic.CreateEndpointOption(endpoint.Name, gcpic.ServiceNameGCPDNS))
	}
	dnsService, err := gcpic.GetDNSService(ctx, opts...)
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

// createPSCRecords will create and publish the records for the private service connect managed zone.
func createPSCRecords(ctx context.Context, ic *installconfig.InstallConfig, clusterID string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*2) // 1 minute for each record below
	defer cancel()

	computeOpts := []option.ClientOption{option.WithScopes(compute.CloudPlatformScope)}
	dnsOpts := []option.ClientOption{option.WithScopes(dns.CloudPlatformScope)}
	pscEndpoint := ic.Config.GCP.Endpoint
	if gcptypes.ShouldUseEndpointForInstaller(pscEndpoint) {
		computeOpts = append(computeOpts, gcpic.CreateEndpointOption(pscEndpoint.Name, gcpic.ServiceNameGCPCompute))
		dnsOpts = append(dnsOpts, gcpic.CreateEndpointOption(pscEndpoint.Name, gcpic.ServiceNameGCPDNS))
	}

	computeService, err := gcpic.GetComputeService(ctx, computeOpts...)
	if err != nil {
		return fmt.Errorf("failed to create Compute service: %w", err)
	}

	endpoint, err := gcpic.GetPrivateServiceConnectEndpoint(computeService, ic.Config.GCP.ProjectID, ic.Config.GCP.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to find private service connect endpoint: %w", err)
	}
	if endpoint == nil {
		return fmt.Errorf("private service connect endpoint %q not found in project %s", ic.Config.GCP.Endpoint.Name, ic.Config.GCP.ProjectID)
	}

	records := []recordSet{
		{
			// "A" record to forward all googleapis.com traffic to the Private Service Connect Endpoint IP Address.
			projectID: ic.Config.GCP.ProjectID,
			zoneName:  fmt.Sprintf("%s-psc-zone", clusterID),
			record: &dns.ResourceRecordSet{
				Name:    "googleapis.com.",
				Type:    "A",
				Ttl:     60,
				Rrdatas: []string{endpoint.IPAddress},
			},
		},
		{
			// A CNAME record that will associate all *.googleapis.com traffic with googleapis.com
			projectID: ic.Config.GCP.ProjectID,
			zoneName:  fmt.Sprintf("%s-psc-zone", clusterID),
			record: &dns.ResourceRecordSet{
				Name:    "*.googleapis.com.",
				Type:    "CNAME",
				Ttl:     60,
				Rrdatas: []string{"googleapis.com."},
			},
		},
	}

	dnsService, err := gcpic.GetDNSService(ctx, dnsOpts...)
	if err != nil {
		return fmt.Errorf("failed to create the gcp dns service: %w", err)
	}

	for _, record := range records {
		if _, err := dnsService.ResourceRecordSets.Create(record.projectID, record.zoneName, record.record).Context(ctx).Do(); err != nil {
			return fmt.Errorf("failed to create private service connect record set %s: %w", record.record.Name, err)
		}
	}

	return nil
}

// createPrivateServiceConnectZone creates a new private zone when a private service connect endpoint is provided.
// The zone is used to forward all traffic intended for the default googleapis endpoints to the IP address
// associated with the Private Service Connect Endpoint.
func createPrivateServiceConnectZone(ctx context.Context, ic *installconfig.InstallConfig, clusterID, network string) error {
	opts := []option.ClientOption{option.WithScopes(dns.CloudPlatformScope)}
	pscEndpoint := ic.Config.GCP.Endpoint
	if gcptypes.ShouldUseEndpointForInstaller(pscEndpoint) {
		opts = append(opts, gcpic.CreateEndpointOption(pscEndpoint.Name, gcpic.ServiceNameGCPDNS))
	}
	dnsService, err := gcpic.GetDNSService(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create the gcp dns service: %w", err)
	}

	managedZone := &dns.ManagedZone{
		Name:        fmt.Sprintf("%s-psc-zone", clusterID),
		Description: "Private Service Connect Private Zone Created By OpenShift Installer",
		DnsName:     "googleapis.com.",
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

	if _, err = dnsService.ManagedZones.Create(ic.Config.GCP.ProjectID, managedZone).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to create private service connect managed zone: %w", err)
	}
	return nil
}

// createPrivateManagedZone will create a private managed zone in the GCP project specified in the install config. The
// private managed zone should only be created when one is not specified in the install config.
func createPrivateManagedZone(ctx context.Context, ic *installconfig.InstallConfig, clusterID, network string) error {
	client, err := gcpic.NewClient(ctx, ic.Config.GCP.Endpoint)
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

	opts := []option.ClientOption{option.WithScopes(dns.CloudPlatformScope)}
	pscEndpoint := ic.Config.GCP.Endpoint
	if gcptypes.ShouldUseEndpointForInstaller(pscEndpoint) {
		opts = append(opts, gcpic.CreateEndpointOption(ic.Config.GCP.Endpoint.Name, gcpic.ServiceNameGCPDNS))
	}
	dnsService, err := gcpic.GetDNSService(ctx, opts...)
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
