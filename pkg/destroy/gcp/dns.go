package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	dns "google.golang.org/api/dns/v1"
)

type dnsZone struct {
	name   string
	domain string
}

func (o *ClusterUninstaller) listDNSZones() (private *dnsZone, public []dnsZone, err error) {
	o.Logger.Debugf("Listing DNS Zones")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.dnsSvc.ManagedZones.List(o.ProjectID).Fields("managedZones(name,dnsName,visibility),nextPageToken")
	err = req.Pages(ctx, func(response *dns.ManagedZonesListResponse) error {
		for _, zone := range response.ManagedZones {
			switch zone.Visibility {
			case "private":
				if o.isClusterResource(zone.Name) {
					o.Logger.Debugf("Found cluster private dns zone: %s\n", zone.Name)
					private = &dnsZone{name: zone.Name, domain: zone.DnsName}
				}
			default:
				public = append(public, dnsZone{name: zone.Name, domain: zone.DnsName})
			}
		}
		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to fetch dns zones")
	}
	return
}

func (o *ClusterUninstaller) deleteDNSZone(name string) error {
	o.Logger.Debugf("Deleting DNS zones %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	err := o.dnsSvc.ManagedZones.Delete(o.ProjectID, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete DNS zone %s", name)
	}
	o.Logger.Infof("Deleted DNS zone %s", name)
	return nil
}

func (o *ClusterUninstaller) listDNSZoneRecordSets(dnsZoneName string) ([]*dns.ResourceRecordSet, error) {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.dnsSvc.ResourceRecordSets.List(o.ProjectID, dnsZoneName)
	result := []*dns.ResourceRecordSet{}
	err := req.Pages(ctx, func(response *dns.ResourceRecordSetsListResponse) error {
		result = append(result, response.Rrsets...)
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch resource record sets for zone: %s", dnsZoneName)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteDNSZoneRecordSets(zoneName string, zoneDomain string, recordSets []*dns.ResourceRecordSet) error {
	change := &dns.Change{}
	for _, rr := range recordSets {
		if (rr.Type == "NS" || rr.Type == "SOA") && strings.TrimRight(rr.Name, ".") == strings.TrimRight(zoneDomain, ".") {
			continue
		}
		change.Deletions = append(change.Deletions, rr)
	}
	if len(change.Deletions) == 0 {
		return nil
	}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting %d recordset(s) in zone %s", len(change.Deletions), zoneName)
	change, err := o.dnsSvc.Changes.Create(o.ProjectID, zoneName, change).ClientOperationId(o.requestID("recordsets", zoneName)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("recordsets", zoneName)
		return errors.Wrapf(err, "failed to delete DNS zone %s recordsets", zoneName)
	}
	if change != nil && isErrorStatus(int64(change.ServerResponse.HTTPStatusCode)) {
		o.resetRequestID("recordsets", zoneName)
		return errors.Errorf("failed to delete DNS zone %s recordsets with code: %d", zoneName, change.ServerResponse.HTTPStatusCode)
	}
	o.resetRequestID("recordsets", zoneName)
	o.Logger.Infof("Deleted %d recordset(s) in zone %s", len(change.Deletions), zoneName)
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

func getParentDNSZone(dnsDomain string, publicZones []dnsZone, logger logrus.FieldLogger) *dnsZone {
	possibleParents := possibleZoneParents(dnsDomain)
	for _, parentDomain := range possibleParents {
		for _, zone := range publicZones {
			if zone.domain == parentDomain {
				logger.Debugf("Found parent dns zone: %s", zone.name)
				return &dnsZone{name: zone.name, domain: parentDomain}
			}
		}
	}
	return nil
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
func (o *ClusterUninstaller) destroyDNS() error {
	privateZone, publicZones, err := o.listDNSZones()
	if err != nil {
		return err
	}
	if privateZone == nil {
		o.Logger.Debugf("Private DNS zone not found")
		return nil
	}

	zoneRecordSets, err := o.listDNSZoneRecordSets(privateZone.name)
	if err != nil {
		return err
	}

	parentZone := getParentDNSZone(privateZone.domain, publicZones, o.Logger)
	if parentZone != nil {
		parentRecordSets, err := o.listDNSZoneRecordSets(parentZone.name)
		if err != nil {
			return err
		}
		matchingRecordSets := o.getMatchingRecordSets(parentRecordSets, zoneRecordSets)
		err = o.deleteDNSZoneRecordSets(parentZone.name, parentZone.domain, matchingRecordSets)
		if err != nil {
			return err
		}
	}
	err = o.deleteDNSZoneRecordSets(privateZone.name, privateZone.domain, zoneRecordSets)
	if err != nil {
		return err
	}
	err = o.deleteDNSZone(privateZone.name)
	if err != nil {
		return err
	}
	return nil
}
