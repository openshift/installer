package ibmcloud

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
)

const dnsRecordTypeName = "dns record"

// listDNSRecords lists DNS records for the cluster
func (o *ClusterUninstaller) listDNSRecords() (cloudResources, error) {
	o.Logger.Debugf("Listing DNS records")
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
			exp := fmt.Sprintf(`.*\Q%s.%s\E$`, o.ClusterName, o.BaseDomain)
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

func (o *ClusterUninstaller) deleteDNSRecord(item cloudResource) error {
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

// destroyDNSRecords removes all DNS record resources that have a name containing
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyDNSRecords() error {
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
