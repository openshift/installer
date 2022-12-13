package ibmcloud

import (
	"fmt"
	"net/http"

	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/pkg/errors"
)

const cosTypeName = "cos instance"

// Resource ID collected via following command using IBM Cloud CLI:
// $ ibmcloud catalog service cloud-object-storage --output json | jq -r '.[].id' .
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

func (o *ClusterUninstaller) findReclaimedCOSInstance(item cloudResource) (*resourcecontrollerv2.ResourceInstance, *resourcecontrollerv2.Reclamation) {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	reclamationOptions := o.controllerSvc.NewListReclamationsOptions()
	reclamation, response, err := o.controllerSvc.ListReclamationsWithContext(ctx, reclamationOptions)
	if err != nil {
		o.Logger.Debugf("Failed listing reclamations: %v, with response: %v", err, response)
		return nil, nil
	}

	o.Logger.Debugf("Checking reclamations for COS Instance: %s", item.name)
	for _, reclamation := range reclamation.Resources {
		getOptions := o.controllerSvc.NewGetResourceInstanceOptions(*reclamation.ResourceInstanceID)
		cosInstance, _, err := o.controllerSvc.GetResourceInstanceWithContext(ctx, getOptions)
		if err != nil {
			o.Logger.Debugf("Failed checking if reclamation is for COS Instance %s: %v", item.name, err)
			return nil, nil
		}

		if *cosInstance.Name == item.name {
			o.Logger.Debugf("Found reclamation for COS Instance %s - %s", item.name, *reclamation.ID)
			return cosInstance, &reclamation
		}
	}

	return nil, nil
}

func (o *ClusterUninstaller) deleteCOSInstance(item cloudResource) error {
	o.Logger.Debugf("Deleting COS instance %s", item.name)

	cosInstance, _ := o.findReclaimedCOSInstance(item)
	if cosInstance != nil {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted COS Instance %s", item.name)
		return nil
	}

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	getOptions := o.controllerSvc.NewGetResourceInstanceOptions(item.id)
	_, getResponse, err := o.controllerSvc.GetResourceInstanceWithContext(ctx, getOptions)
	if err != nil {
		if getResponse != nil && getResponse.StatusCode == http.StatusNotFound {
			// The resource is gone
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted COS Instance %s", item.name)
			return nil
		}

		return errors.Wrapf(err, "Failed to delete COS Instance %s", item.name)
	}

	deleteOptions := o.controllerSvc.NewDeleteResourceInstanceOptions(item.id)
	deleteOptions.SetRecursive(true)
	deleteResponse, err := o.controllerSvc.DeleteResourceInstanceWithContext(ctx, deleteOptions)
	if err != nil && deleteResponse != nil && deleteResponse.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete COS Instance %s", item.name)
	}

	cosInstance, reclamation := o.findReclaimedCOSInstance(item)
	if cosInstance != nil {
		o.Logger.Infof("Reclaiming COS Instance %s: Reclamation: %s", item.name, *reclamation.ID)
		reclamationOptions := o.controllerSvc.NewRunReclamationActionOptions(*reclamation.ID, "reclaim")
		_, _, err := o.controllerSvc.RunReclamationActionWithContext(ctx, reclamationOptions)
		if err != nil {
			return errors.Wrapf(err, "Failed to reclaim COS Instance: %s", item.name)
		}
		o.Logger.Infof("Deleted COS Instance Reclamation: %s - %s", item.name, *reclamation.ID)
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted COS Instance %s", item.name)

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
			o.Logger.Infof("Deleted COS instance %s", item.name)
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
			o.cosInstanceID = instance.id
			return instance.id, nil
		}
	}
	return "", errors.Errorf("COS instance not found")
}
