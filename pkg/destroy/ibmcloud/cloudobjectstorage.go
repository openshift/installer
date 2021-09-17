package ibmcloud

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const cosTypeName = "cos instance"
const cosResourceID = "dff97f5c-bc5e-4455-b470-411c3edbe49c"

// listCOSInstances lists COS service instances
func (o *ClusterUninstaller) listCOSInstances() (cloudResources, error) {
	o.Logger.Debugf("Listing COS instances")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	resourceGroupID, err := o.ResourceGroupID()
	if err != nil {
		return nil, err
	}

	options := o.controllerSvc.NewListResourceInstancesOptions()
	options.SetResourceGroupID(resourceGroupID)
	options.SetResourceID(cosResourceID)
	options.SetType("service_instance")
	resources, _, err := o.controllerSvc.ListResourceInstancesWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list COS instances")
	}

	result := []cloudResource{}
	for _, instance := range resources.Resources {
		// Match the COS instances created by both the installer and the
		// cluster-image-registry-operator.
		if fmt.Sprintf("%s-cos", o.InfraID) == *instance.Name ||
			fmt.Sprintf("%s-image-registry", o.InfraID) == *instance.Name {
			result = append(result, cloudResource{
				key:      *instance.ID,
				name:     *instance.Name,
				status:   *instance.State,
				typeName: cosTypeName,
				id:       *instance.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteCOSInstance(item cloudResource) error {
	o.Logger.Debugf("Deleting COS instance %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.controllerSvc.NewDeleteResourceInstanceOptions(item.id)
	options.SetRecursive(true)
	details, err := o.controllerSvc.DeleteResourceInstanceWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted COS instance %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete COS instance %s", item.name)
	}

	return nil
}

// destroyCOSInstances removes the COS service instance resources that have a
// name prefixed with the cluster's infra ID.
func (o *ClusterUninstaller) destroyCOSInstances() error {
	found, err := o.listCOSInstances()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(cosTypeName, found.list())

	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted COS instance %q", item.name)
			continue
		}
		err = o.deleteCOSInstance(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(cosTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}

// COSInstanceID returns the ID of the Cloud Object Storage service instance
// created by the installer during installation.
func (o *ClusterUninstaller) COSInstanceID() (string, error) {
	if o.cosInstanceID != "" {
		return o.cosInstanceID, nil
	}
	cosInstances, err := o.listCOSInstances()
	if err != nil {
		return "", err
	}
	instanceList := cosInstances.list()
	if len(instanceList) == 0 {
		return "", errors.Errorf("COS instance not found")
	}

	// Locate the installer's COS instance by name.
	for _, instance := range instanceList {
		if instance.name == fmt.Sprintf("%s-cos", o.InfraID) {
                        split := strings.Split(instance.id, ":")
                        if len(split) >= 8 { //CRN string
                                o.cosInstanceID = split[7]
                        } else {
                                o.cosInstanceID = instance.id
                        }
                        return o.cosInstanceID, nil
		}
	}
	return "", errors.Errorf("COS instance not found")
}
