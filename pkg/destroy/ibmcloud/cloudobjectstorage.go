package ibmcloud

import (
	"fmt"
	"net/http"

	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/pkg/errors"
)

const (
	cosTypeName = "cos instance"
	// reclamationReclaim will delete the resource, reclaim it.
	reclamationReclaim = "reclaim"
)

// Resource ID collected via following command using IBM Cloud CLI:
// $ ibmcloud catalog service cloud-object-storage --output json | jq -r '.[].id' .
const cosResourceID = "dff97f5c-bc5e-4455-b470-411c3edbe49c"

// findCOSInstanceReclamation Checks current reclamations for one that matches the COS instance name.
func (o *ClusterUninstaller) findCOSInstanceReclamation(instance cloudResource) (*resourcecontrollerv2.Reclamation, error) {
	o.Logger.Debugf("Searching for COS reclamations for instance %s", instance.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	reclamationOptions := o.controllerSvc.NewListReclamationsOptions()
	resources, _, err := o.controllerSvc.ListReclamationsWithContext(ctx, reclamationOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed listing reclamations for instance %s", instance.name)
	}

	o.Logger.Debugf("Checking reclamations that match instance %s", instance.name)
	for _, reclamation := range resources.Resources {
		getOptions := o.controllerSvc.NewGetResourceInstanceOptions(*reclamation.ResourceInstanceID)
		cosInstance, _, err := o.controllerSvc.GetResourceInstanceWithContext(ctx, getOptions)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed checking reclamation %s", *reclamation.ResourceInstanceID)
		}
		if *cosInstance.Name == instance.name {
			o.Logger.Debugf("Found COS instance reclamation %s - %s", instance.name, *reclamation.ID)
			return &reclamation, nil
		}
	}

	return nil, nil
}

// reclaimCOSInstanceReclamation reclaims (deletes) a reclamation from a deleted COS instance.
func (o *ClusterUninstaller) reclaimCOSInstanceReclamation(reclamationID string) error {
	o.Logger.Debugf("Reclaming COS instance reclamation %s", reclamationID)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.controllerSvc.NewRunReclamationActionOptions(reclamationID, reclamationReclaim)
	_, response, err := o.controllerSvc.RunReclamationActionWithContext(ctx, options)
	if err != nil {
		// If reclaim attempt failed because the reclamation doesn't exist (404) don't return an error
		if response != nil && response.StatusCode == http.StatusNotFound {
			o.Logger.Debugf("Reclamation not found, it has likely already been reclaimed %s", reclamationID)
			return nil
		}
		return errors.Wrapf(err, "Failed to reclaim COS instance reclamation %s", reclamationID)
	}

	o.Logger.Infof("Reclaimed %s", reclamationID)
	return nil
}

// listCOSInstances lists COS service instances.
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
	o.Logger.Debugf("Deleting COS instance %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.controllerSvc.NewDeleteResourceInstanceOptions(item.id)
	options.SetRecursive(true)
	details, err := o.controllerSvc.DeleteResourceInstanceWithContext(ctx, options)
	if err != nil {
		// If the deletion attempt failed because the COS instance doesn't exist (404) don't return an error
		if details != nil && details.StatusCode == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "Failed to delete COS Instance %s", item.name)
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
			// Check if a COS reclamation was created for the deleted instance
			reclamation, err := o.findCOSInstanceReclamation(item)
			if err != nil {
				return err
			}
			if reclamation != nil {
				err = o.reclaimCOSInstanceReclamation(*reclamation.ID)
				if err != nil {
					return err
				}
				continue
			}
			o.Logger.Debugf("No reclamations found for COS instance %s", item.name)
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
