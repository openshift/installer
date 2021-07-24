package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const floatingIPTypeName = "floating ip"

// listFloatingIPs lists floating IPs in the vpc
func (o *ClusterUninstaller) listFloatingIPs() (cloudResources, error) {
	o.Logger.Debugf("Listing floating IPs")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListFloatingIpsOptions()
	resources, _, err := o.vpcSvc.ListFloatingIpsWithContext(ctx, options)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to list floating IPs")
	}

	result := []cloudResource{}
	for _, floatingIPs := range resources.FloatingIps {
		if strings.Contains(*floatingIPs.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *floatingIPs.ID,
				name:     *floatingIPs.Name,
				status:   *floatingIPs.Status,
				typeName: floatingIPTypeName,
				id:       *floatingIPs.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteFloatingIP(item cloudResource) error {
	if item.status == vpcv1.FloatingIPStatusDeletingConst {
		o.Logger.Debugf("Waiting for floating IP %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting floating IP %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeleteFloatingIPOptions(item.id)
	details, err := o.vpcSvc.DeleteFloatingIPWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted floating IP %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete floating IP %s", item.name)
	}

	return nil
}

// destroyFloatingIPs removes all floating IP resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyFloatingIPs() error {
	found, err := o.listFloatingIPs()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(floatingIPTypeName, found.list())

	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted floating IP %q", item.name)
			continue
		}
		err = o.deleteFloatingIP(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(floatingIPTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
