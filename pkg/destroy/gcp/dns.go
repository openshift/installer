package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	dns "google.golang.org/api/dns/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type dnsZone struct {
	name    string
	domain  string
	project string
}

func (o *ClusterUninstaller) listDNSZones(ctx context.Context) (private *dnsZone, public []dnsZone, err error) {
	o.Logger.Debugf("Listing DNS Zones")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	projects := []string{o.ProjectID}
	if o.NetworkProjectID != "" {
		projects = append(projects, o.NetworkProjectID)
	}

	for _, project := range projects {
		req := o.dnsSvc.ManagedZones.List(project).Fields("managedZones(name,dnsName,visibility),nextPageToken")
		err = req.Pages(ctx, func(response *dns.ManagedZonesListResponse) error {
			for _, zone := range response.ManagedZones {
				switch zone.Visibility {
				case "private":
					if o.isClusterResource(zone.Name) || (o.PrivateZoneDomain != "" && o.PrivateZoneDomain == zone.DnsName) {
						private = &dnsZone{name: zone.Name, domain: zone.DnsName, project: project}
					}
				default:
					public = append(public, dnsZone{name: zone.Name, domain: zone.DnsName, project: project})
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

func (o *ClusterUninstaller) deleteDNSZone(ctx context.Context, name string) error {
	if !o.isClusterResource(name) {
		o.Logger.Warnf("Skipping deletion of DNS Zone %s, not created by installer", name)
		return nil
	}

	o.Logger.Debugf("Deleting DNS zones %s", name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	err := o.dnsSvc.ManagedZones.Delete(o.ProjectID, name).Context(ctx).Do()
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
	err = o.deleteDNSZone(ctx, privateZone.name)
	if err != nil {
		return err
	}
	return nil
}
