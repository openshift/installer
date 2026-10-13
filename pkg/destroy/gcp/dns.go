package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	dns "google.golang.org/api/dns/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset/installconfig/gcp"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

type dnsZone struct {
	name    string
	domain  string
	project string
	labels  map[string]string
}

func (o *ClusterUninstaller) listDNSZones(ctx context.Context) (private *dnsZone, public []dnsZone, err error) {
	o.Logger.Debugf("Listing DNS Zones")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	projects := []string{o.ProjectID}
	if o.NetworkProjectID != "" {
		projects = append(projects, o.NetworkProjectID)
	}
	if o.PrivateZoneProjectID != "" {
		projects = append(projects, o.PrivateZoneProjectID)
	}

	for _, project := range projects {
		req := o.dnsSvc.ManagedZones.List(project).Fields("managedZones(name,dnsName,visibility,labels),nextPageToken")
		err = req.Pages(ctx, func(response *dns.ManagedZonesListResponse) error {
			for _, zone := range response.ManagedZones {
				switch zone.Visibility {
				case "private":
					// if a dns private managed zone is provided make sure the correct zone is found and
					// added to the list for destroy.
					if o.PrivateZoneProjectID != "" && project != o.PrivateZoneProjectID {
						// The private zone should exist in the private zone project if one was provided.
						continue
					}
					if o.isClusterResource(zone.Name) || (o.PrivateZoneDomain != "" && o.PrivateZoneDomain == zone.DnsName) ||
						o.isManagedResource(zone.Labels) {
						private = &dnsZone{name: zone.Name, domain: zone.DnsName, project: project, labels: zone.Labels}
					}
				default:
					if project == o.ProjectID {
						public = append(public, dnsZone{name: zone.Name, domain: zone.DnsName, project: project, labels: zone.Labels})
					}
				}
			}
			return nil
		})
		if err != nil {
			err = errors.Wrapf(err, "failed to fetch dns zones")
		}
	}
	return
}

// removeSharedTag will remove the shared tag associated with this cluster from the DNS Managed Zone.
func (o *ClusterUninstaller) removeSharedTag(ctx context.Context, svc *dns.Service, clusterID, baseDomain, project, zoneName string, isPublic bool) error {
	params := gcptypes.DNSZoneParams{
		Project:    project,
		Name:       zoneName,
		BaseDomain: baseDomain,
		IsPublic:   isPublic,
	}
	zone, err := gcp.GetDNSZoneFromParams(ctx, svc, params)
	if err != nil {
		return err
	}

	if zone == nil {
		return fmt.Errorf("failed to find matching DNS zone for %s in project %s", zoneName, project)
	}

	if zone.Labels == nil {
		return nil
	}

	// Remove the shared tag for this cluster from the DNS Managed Zone.
	// First, make sure that the tag value is shared and not owned or anything else.
	labelKey := fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, clusterID)
	if value, ok := zone.Labels[labelKey]; ok {
		if value == sharedLabelValue {
			delete(zone.Labels, labelKey)
		}
	}

	return gcp.UpdateDNSManagedZone(ctx, svc, project, zoneName, zone)
}

func (o *ClusterUninstaller) deleteDNSZone(ctx context.Context, name string, labels map[string]string) error {
	// The destroy code is executed in a loop. On the first pass it is possible that the
	// zone contained the "shared" label and the label is removed. The zone should NOT
	// be deleted in this case, so the process is skipped.
	if !o.isManagedResource(labels) && !o.isClusterResource(name) {
		o.Logger.Debugf("DNS Zone %s is not owned or shared by the cluster, skipping deletion", name)
		return nil
	}

	if o.isSharedResource(labels) {
		o.Logger.Warnf("Skipping deletion of DNS Zone %s, not created by installer", name)
		// currently, this function only runs for the private managed zone, but
		// let's add a check just to be sure.
		err := o.removeSharedTag(ctx, o.dnsSvc, o.ClusterID, o.PrivateZoneDomain, o.PrivateZoneProjectID, name, false)
		if err != nil {
			return fmt.Errorf("failed to remove shared tag from DNS Zone %s: %w", name, err)
		}
		return nil
	}

	o.Logger.Debugf("Deleting DNS zones %s", name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	err := o.dnsSvc.ManagedZones.Delete(o.PrivateZoneProjectID, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete DNS zone %s", name)
	}
	o.Logger.Infof("Deleted DNS zone %s", name)
	return nil
}

func (o *ClusterUninstaller) listDNSZoneRecordSets(ctx context.Context, zone *dnsZone) ([]*dns.ResourceRecordSet, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	req := o.dnsSvc.ResourceRecordSets.List(zone.project, zone.name)
	result := []*dns.ResourceRecordSet{}
	err := req.Pages(ctx, func(response *dns.ResourceRecordSetsListResponse) error {
		result = append(result, response.Rrsets...)
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch resource record sets for zone: %s", zone.name)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteDNSZoneRecordSets(ctx context.Context, zone *dnsZone, recordSets []*dns.ResourceRecordSet) error {
	change := &dns.Change{}
	for _, rr := range recordSets {
		if (rr.Type == "NS" || rr.Type == "SOA") && strings.TrimRight(rr.Name, ".") == strings.TrimRight(zone.domain, ".") {
			continue
		}
		change.Deletions = append(change.Deletions, rr)
	}
	if len(change.Deletions) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	o.Logger.Debugf("Deleting %d recordset(s) in zone %s", len(change.Deletions), zone.name)
	change, err := o.dnsSvc.Changes.Create(zone.project, zone.name, change).ClientOperationId(o.requestID("recordsets", zone.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("recordsets", zone.name)
		return errors.Wrapf(err, "failed to delete DNS zone %s recordsets", zone.name)
	}
	if change != nil && isErrorStatus(int64(change.ServerResponse.HTTPStatusCode)) {
		o.resetRequestID("recordsets", zone.name)
		return errors.Errorf("failed to delete DNS zone %s recordsets with code: %d", zone.name, change.ServerResponse.HTTPStatusCode)
	}
	o.resetRequestID("recordsets", zone.name)
	o.Logger.Infof("Deleted %d recordset(s) in zone %s", len(change.Deletions), zone.name)
	return nil
}

func possibleZoneParents(dnsDomain string) []string {
	result := []string{}
	parts := strings.Split(dnsDomain, ".")
	for len(parts) > 0 {
		result = append(result, strings.Join(parts, "."))
		parts = parts[1:]
	}
	return result
}

func getParentDNSZones(dnsDomain string, publicZones []dnsZone, logger logrus.FieldLogger) []*dnsZone {
	parentZones := []*dnsZone{}

	possibleParents := possibleZoneParents(dnsDomain)
	for _, parentDomain := range possibleParents {
		for _, zone := range publicZones {
			if zone.domain == parentDomain {
				logger.Debugf("Found parent dns zone: %s", zone.name)
				parentZones = append(parentZones, &dnsZone{name: zone.name, domain: parentDomain, project: zone.project})
			}
		}
	}
	return parentZones
}

// getMatchingRecordSets finds all recordsets in the parent list that match recordsets in the child list
// using the record type and record name as keys.
func (o *ClusterUninstaller) getMatchingRecordSets(parentRecords, childRecords []*dns.ResourceRecordSet) []*dns.ResourceRecordSet {
	matchingRecordSets := []*dns.ResourceRecordSet{}
	recordKey := func(r *dns.ResourceRecordSet) string {
		return fmt.Sprintf("%s %s", r.Type, r.Name)
	}
	childKeys := sets.NewString()
	for _, record := range childRecords {
		childKeys.Insert(recordKey(record))
	}
	for _, record := range parentRecords {
		if childKeys.Has(recordKey(record)) {
			matchingRecordSets = append(matchingRecordSets, record)
		}
	}
	return matchingRecordSets
}

// destroyDNS deletes DNS resources associated with the cluster. It first finds
// the private DNS zone that belongs to the cluster by looking for a zone prefixed
// with the cluster's infra ID. It then finds a public zone that is the parent of
// the cluster's private zone by searching for zones with a matching domain (in
// order from specific to general). If/when a parent DNS zone is found, the records
// from the private zone are matched to records in the parent zone (by using type
// and name for each record). Matching records are removed from the public zone.
// Finally all records are removed from the private zone and the private zone is removed.
func (o *ClusterUninstaller) destroyDNS(ctx context.Context) error {
	privateZone, publicZones, err := o.listDNSZones(ctx)
	if err != nil {
		return err
	}
	if privateZone == nil {
		o.Logger.Debugf("Private DNS zone not found")
		return nil
	}

	zoneRecordSets, err := o.listDNSZoneRecordSets(ctx, privateZone)
	if err != nil {
		return err
	}

	parentZones := getParentDNSZones(privateZone.domain, publicZones, o.Logger)
	for _, parentZone := range parentZones {
		parentRecordSets, err := o.listDNSZoneRecordSets(ctx, parentZone)
		if err != nil {
			return err
		}
		matchingRecordSets := o.getMatchingRecordSets(parentRecordSets, zoneRecordSets)
		if len(matchingRecordSets) > 0 {
			err = o.deleteDNSZoneRecordSets(ctx, parentZone, matchingRecordSets)
			if err != nil {
				return err
			}
		}
	}
	err = o.deleteDNSZoneRecordSets(ctx, privateZone, zoneRecordSets)
	if err != nil {
		return err
	}
	err = o.deleteDNSZone(ctx, privateZone.name, privateZone.labels)
	if err != nil {
		return err
	}
	return nil
}
