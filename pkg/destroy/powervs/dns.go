package powervs

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/networking-go-sdk/resourcerecordsv1"
	"github.com/pkg/errors"
)

const (
	cisDNSRecordTypeName = "cis dns record"
	ibmDNSRecordTypeName = "ibm dns record"
)

// listDNSRecords lists DNS records for the cluster.
func (o *ClusterUninstaller) listDNSRecords() (cloudResources, error) {
	o.Logger.Debugf("Listing DNS records")

	ctx, _ := o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listLoadBalancers: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	var (
		foundOne       = false
		perPage  int64 = 20
		page     int64 = 1
		moreData bool  = true
	)

	dnsRecordsOptions := o.dnsRecordsSvc.NewListAllDnsRecordsOptions()
	dnsRecordsOptions.PerPage = &perPage
	dnsRecordsOptions.Page = &page

	result := []cloudResource{}

	for moreData {
		dnsResources, detailedResponse, err := o.dnsRecordsSvc.ListAllDnsRecordsWithContext(ctx, dnsRecordsOptions)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list DNS records: %v and the response is: %s", err, detailedResponse)
		}

		for _, record := range dnsResources.Result {
			// Match all of the cluster's DNS records
			exp := fmt.Sprintf(`.*\Q%s.%s\E$`, o.ClusterName, o.BaseDomain)
			nameMatches, _ := regexp.Match(exp, []byte(*record.Name))
			contentMatches, _ := regexp.Match(exp, []byte(*record.Content))
			if nameMatches || contentMatches {
				foundOne = true
				o.Logger.Debugf("listDNSRecords: FOUND: %v, %v", *record.ID, *record.Name)
				result = append(result, cloudResource{
					key:      *record.ID,
					name:     *record.Name,
					status:   "",
					typeName: cisDNSRecordTypeName,
					id:       *record.ID,
				})
			}
		}

		o.Logger.Debugf("listDNSRecords: PerPage = %v, Page = %v, Count = %v", *dnsResources.ResultInfo.PerPage, *dnsResources.ResultInfo.Page, *dnsResources.ResultInfo.Count)

		moreData = *dnsResources.ResultInfo.PerPage == *dnsResources.ResultInfo.Count
		o.Logger.Debugf("listDNSRecords: moreData = %v", moreData)

		page++
	}
	if !foundOne {
		o.Logger.Debugf("listDNSRecords: NO matching DNS against: %s", o.InfraID)
		for moreData {
			dnsResources, detailedResponse, err := o.dnsRecordsSvc.ListAllDnsRecordsWithContext(ctx, dnsRecordsOptions)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to list DNS records: %v and the response is: %s", err, detailedResponse)
			}
			for _, record := range dnsResources.Result {
				o.Logger.Debugf("listDNSRecords: FOUND: DNS: %v, %v", *record.ID, *record.Name)
			}
			moreData = *dnsResources.ResultInfo.PerPage == *dnsResources.ResultInfo.Count
			page++
		}
	}

	return cloudResources{}.insert(result...), nil
}

// listDNSRecords lists DNS records for the cluster.
func (o *ClusterUninstaller) listResourceRecords() (cloudResources, error) {
	o.Logger.Debugf("Listing DNS resource records")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listLoadBalancers: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	result := []cloudResource{}

	dnsCRN, err := crn.Parse(o.DNSInstanceCRN)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse DNSInstanceCRN")
	}
	records, _, err := o.resourceRecordsSvc.ListResourceRecords(&resourcerecordsv1.ListResourceRecordsOptions{
		InstanceID: &dnsCRN.ServiceInstance,
		DnszoneID:  &o.dnsZoneID,
	})
	for _, record := range records.ResourceRecords {
		// Match all of the cluster's DNS records
		exp := fmt.Sprintf(`.*\Q%s.%s\E$`, o.ClusterName, o.BaseDomain)
		nameMatches, _ := regexp.Match(exp, []byte(*record.Name))
		if nameMatches {
			o.Logger.Debugf("listResourceRecords: FOUND: %v, %v", *record.ID, *record.Name)
			result = append(result, cloudResource{
				key:      *record.ID,
				name:     *record.Name,
				status:   "",
				typeName: ibmDNSRecordTypeName,
				id:       *record.ID,
			})
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve DNS records")
	}
	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) destroyDNSRecord(item cloudResource) error {
	var (
		response *core.DetailedResponse
		err      error
	)

	ctx, _ := o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deleteFloatingIP: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	switch item.typeName {
	case cisDNSRecordTypeName:
		getOptions := o.dnsRecordsSvc.NewGetDnsRecordOptions(item.id)
		_, response, err = o.dnsRecordsSvc.GetDnsRecordWithContext(ctx, getOptions)

		if err != nil && response != nil && response.StatusCode == http.StatusNotFound {
			// The resource is gone
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted DNS record %q", item.name)
			return nil
		}
		if err != nil && response != nil && response.StatusCode == http.StatusInternalServerError {
			o.Logger.Infof("destroyDNSRecord: internal server error")
			return nil
		}

		o.Logger.Debugf("Deleting DNS record %q", item.name)

		deleteOptions := o.dnsRecordsSvc.NewDeleteDnsRecordOptions(item.id)
		_, _, err = o.dnsRecordsSvc.DeleteDnsRecordWithContext(ctx, deleteOptions)
		if err != nil {
			return errors.Wrapf(err, "failed to delete DNS record %s", item.name)
		}
	case ibmDNSRecordTypeName:
		if err != nil {
			return errors.Wrapf(err, "failed to delete DNS record %s", item.name)
		}
		dnsCRN, err := crn.Parse(o.DNSInstanceCRN)
		if err != nil {
			return errors.Wrap(err, "Failed to parse DNSInstanceCRN")
		}
		getOptions := o.resourceRecordsSvc.NewGetResourceRecordOptions(dnsCRN.ServiceInstance, o.dnsZoneID, item.id)
		_, response, err = o.resourceRecordsSvc.GetResourceRecord(getOptions)

		if err != nil && response != nil && response.StatusCode == http.StatusNotFound {
			// The resource is gone
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted DNS record %q", item.name)
			return nil
		}
		if err != nil && response != nil && response.StatusCode == http.StatusInternalServerError {
			o.Logger.Infof("destroyDNSRecord: internal server error")
			return nil
		}

		o.Logger.Debugf("Deleting DNS record %q", item.name)

		deleteOptions := o.resourceRecordsSvc.NewDeleteResourceRecordOptions(dnsCRN.ServiceInstance, o.dnsZoneID, item.id)
		_, err = o.resourceRecordsSvc.DeleteResourceRecord(deleteOptions)
		if err != nil {
			return errors.Wrapf(err, "failed to delete DNS record %s", item.name)
		}
	}
	return nil
}

// destroyDNSRecords removes all DNS record resources that have a name containing
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyDNSRecords() error {
	var (
		err   error
		items []cloudResource
		found cloudResources
	)
	if o.dnsRecordsSvc != nil {
		found, err = o.listDNSRecords()
		if err != nil {
			return err
		}
		items = o.insertPendingItems(cisDNSRecordTypeName, found.list())
	}

	if o.resourceRecordsSvc != nil {
		found, err = o.listResourceRecords()
		if err != nil {
			return err
		}
		items = o.insertPendingItems(ibmDNSRecordTypeName, found.list())
	}

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyDNSRecords: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted DNS record %q", item.name)
				continue
			}
			err = o.destroyDNSRecord(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(cisDNSRecordTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(cisDNSRecordTypeName); len(items) > 0 {
		return errors.Errorf("destroyDNSRecords: %d undeleted items pending", len(items))
	}
	return nil
}
