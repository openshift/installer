package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const publicGatewayTypeName = "public gateway"

// listPublicGateways lists public gateways in the vpc
func (o *ClusterUninstaller) listPublicGateways() (cloudResources, error) {
	o.Logger.Debugf("Listing public gateways")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListPublicGatewaysOptions()
	resources, _, err := o.vpcSvc.ListPublicGatewaysWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list public gateways")
	}

	result := []cloudResource{}
	for _, publicGateway := range resources.PublicGateways {
		if strings.Contains(*publicGateway.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *publicGateway.ID,
				name:     *publicGateway.Name,
				status:   *publicGateway.Status,
				typeName: publicGatewayTypeName,
				id:       *publicGateway.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deletePublicGateway(item cloudResource) error {
	if item.status == vpcv1.PublicGatewayStatusDeletingConst {
		o.Logger.Debugf("Waiting for public gateway %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting public gateway %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeletePublicGatewayOptions(item.id)
	details, err := o.vpcSvc.DeletePublicGatewayWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted public gateway %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete public gateway %s", item.name)
	}

	return nil
}

// destroyPublicGateways removes all public gateway resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyPublicGateways() error {
	if len(o.UserProvidedSubnets) > 0 {
		o.Logger.Info("Skipping deletion of public gateways in case of user-provided subnets")
		return nil
	}

	found, err := o.listPublicGateways()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(publicGatewayTypeName, found.list())
	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted public gateway %q", item.name)
			continue
		}

		err := o.deletePublicGateway(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(publicGatewayTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
