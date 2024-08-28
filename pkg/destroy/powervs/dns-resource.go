package powervs

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/resourcerecordsv1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	ibmDNSRecordTypeName = "ibm dns record"
)

// listResourceRecords lists DNS Resource records for the cluster.
func (o *ClusterUninstaller) listResourceRecords() (cloudResources, error) {
	o.Logger.Debugf("Listing DNS resource records")

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listResourceRecords: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	result := []cloudResource{}

	dnsCRN, err := crn.Parse(o.DNSInstanceCRN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
	}
	records, _, err := o.resourceRecordsSvc.ListResourceRecords(&resourcerecordsv1.ListResourceRecordsOptions{
		InstanceID: &dnsCRN.ServiceInstance,
		DnszoneID:  &o.dnsZoneID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list resource records: %w", err)
	}

	dnsMatcher, err := regexp.Compile(fmt.Sprintf(`.*\Q%s.%s\E$`, o.ClusterName, o.BaseDomain))
	if err != nil {
		return nil, fmt.Errorf("failed to build DNS records matcher: %w", err)
	}

	for _, record := range records.ResourceRecords {
		// Match all of the cluster's DNS records
		nameMatches := dnsMatcher.Match([]byte(*record.Name))
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
		return nil, fmt.Errorf("could not retrieve DNS records: %w", err)
	}
	return cloudResources{}.insert(result...), nil
}

// destroyResourceRecord destroys a Resource Record.
func (o *ClusterUninstaller) destroyResourceRecord(item cloudResource) error {
	var (
		response *core.DetailedResponse
		err      error
	)

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyResourceRecord: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	if err != nil {
		return fmt.Errorf("failed to delete DNS Resource record %s: %w", item.name, err)
	}
	dnsCRN, err := crn.Parse(o.DNSInstanceCRN)
	if err != nil {
		return fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
	}
	getOptions := o.resourceRecordsSvc.NewGetResourceRecordOptions(dnsCRN.ServiceInstance, o.dnsZoneID, item.id)
	_, response, err = o.resourceRecordsSvc.GetResourceRecord(getOptions)

	if err != nil && response != nil && response.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted DNS Resource Record %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == http.StatusInternalServerError {
		o.Logger.Infof("destroyResourceRecord: internal server error")
		return nil
	}

	deleteOptions := o.resourceRecordsSvc.NewDeleteResourceRecordOptions(dnsCRN.ServiceInstance, o.dnsZoneID, item.id)

	_, err = o.resourceRecordsSvc.DeleteResourceRecord(deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete DNS Resource record %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted DNS Resource Record %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyResourceRecords removes all DNS record resources that have a name containing
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyResourceRecords() error {
	if o.resourceRecordsSvc == nil {
		// Install config didn't specify using these resources
		return nil
	}

	firstPassList, err := o.listResourceRecords()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(ibmDNSRecordTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyResourceRecords: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.destroyResourceRecord(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyResourceRecords: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(ibmDNSRecordTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyResourceRecords: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyResourceRecords: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listResourceRecords()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyResourceRecords: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyResourceRecords: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
