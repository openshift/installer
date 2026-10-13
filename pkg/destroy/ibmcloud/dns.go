package ibmcloud

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
)

const dnsRecordTypeName = "dns record"

// listDNSRecords lists DNS records for the cluster for CIS or DNS Service
func (o *ClusterUninstaller) listDNSRecords() (cloudResources, error) {
	if len(o.CISInstanceCRN) > 0 {
		return o.listCISDNSRecords()
	}
	return o.listDNSSvcDNSRecords()
}

// listCISDNSRecords lists CIS DNS records for the cluster
func (o *ClusterUninstaller) listCISDNSRecords() (cloudResources, error) {
	o.Logger.Debug("Listing CIS DNS records")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	result := []cloudResource{}
	moreData := true
	for moreData {
		options := o.dnsRecordsSvc.NewListAllDnsRecordsOptions()
		resources, _, err := o.dnsRecordsSvc.ListAllDnsRecordsWithContext(ctx, options)

		if err != nil {
			return nil, errors.Wrapf(err, "Failed to list DNS records")
		}

		for _, record := range resources.Result {
			// Match all of the cluster's DNS records
			exp := fmt.Sprintf(`.*\Q.%s.%s\E$`, o.ClusterName, o.BaseDomain)
			nameMatches, _ := regexp.Match(exp, []byte(*record.Name))
			contentMatches, _ := regexp.Match(exp, []byte(*record.Content))
			if nameMatches || contentMatches {
				result = append(result, cloudResource{
					key:      *record.ID,
					name:     *record.Name,
					status:   "",
					typeName: dnsRecordTypeName,
					id:       *record.ID,
				})
			}
		}

		resultInfo := *resources.ResultInfo
		moreData = (*resultInfo.PerPage * *resultInfo.Page) < *resultInfo.Count
	}

	return cloudResources{}.insert(result...), nil
}

// listDNSSvcDNSRecords lists DNS Services DNS records for the cluster
func (o *ClusterUninstaller) listDNSSvcDNSRecords() (cloudResources, error) {
	o.Logger.Debug("Listing DNS Services DNS records")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	result := []cloudResource{}
	resourceRecordsRemaining := true
	viewedResourceRecords := int64(0)
	for resourceRecordsRemaining {
		options := o.dnsServicesSvc.NewListResourceRecordsOptions(o.DNSInstanceID, o.zoneID)
		options = options.SetOffset(viewedResourceRecords)
		resources, _, err := o.dnsServicesSvc.ListResourceRecordsWithContext(ctx, options)

		if err != nil {
			return nil, errors.Wrap(err, "Failed to list DNS records")
		}

		for _, record := range resources.ResourceRecords {
			// Match all of the cluster's DNS records
			exp := fmt.Sprintf(`.*\Q.%s.%s\E$`, o.ClusterName, o.BaseDomain)
			if nameMatches, _ := regexp.Match(exp, []byte(*record.Name)); nameMatches {
				result = append(result, cloudResource{
					key:      *record.ID,
					name:     *record.Name,
					status:   "",
					typeName: dnsRecordTypeName,
					id:       *record.ID,
				})
			}
		}

		viewedResourceRecords += *resources.Limit
		resourceRecordsRemaining = viewedResourceRecords < *resources.TotalCount
	}

	return cloudResources{}.insert(result...), nil
}

// deleteDNSRecord deletes a CIS or DNS Services DNS Record
func (o *ClusterUninstaller) deleteDNSRecord(item cloudResource) error {
	if len(o.CISInstanceCRN) > 0 {
		return o.deleteCISDNSRecord(item)
	}
	return o.deleteDNSSvcDNSRecord(item)
}

// deleteCISDNSRecord deletes a CIS DNS Record
func (o *ClusterUninstaller) deleteCISDNSRecord(item cloudResource) error {
	o.Logger.Debugf("Deleting DNS record %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.dnsRecordsSvc.NewDeleteDnsRecordOptions(item.id)
	_, details, err := o.dnsRecordsSvc.DeleteDnsRecordWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted DNS record %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete DNS record %s", item.name)
	}

	return nil
}

// deleteDNSSvcDNSRecord deletes a DNS Services DNS Record
func (o *ClusterUninstaller) deleteDNSSvcDNSRecord(item cloudResource) error {
	o.Logger.Debugf("Deleting DNS record %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.dnsServicesSvc.NewDeleteResourceRecordOptions(o.DNSInstanceID, o.zoneID, item.id)
	details, err := o.dnsServicesSvc.DeleteResourceRecordWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted DNS record %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete DNS record %s", item.name)
	}
	return nil
}

// destroyDNSRecords removes all DNS record resources that have a name containing
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyDNSRecords() error {
	// If neither CIS CRN or DNS Services ID is not set, skip DNS records cleanup
	if len(o.CISInstanceCRN) == 0 && len(o.DNSInstanceID) == 0 {
		o.Logger.Info("Skipping deletion of DNS Records, no CIS CRN or DNS Instance ID found")
		return nil
	}

	found, err := o.listDNSRecords()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(dnsRecordTypeName, found.list())

	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted DNS record %q", item.name)
			continue
		}
		err = o.deleteDNSRecord(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(dnsRecordTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
