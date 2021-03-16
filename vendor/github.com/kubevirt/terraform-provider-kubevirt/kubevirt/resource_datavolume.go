package kubevirt

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/client"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/datavolume"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils/patch"
	"k8s.io/apimachinery/pkg/api/errors"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func resourceKubevirtDataVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubevirtDataVolumeCreate,
		Read:   resourceKubevirtDataVolumeRead,
		Update: resourceKubevirtDataVolumeUpdate,
		Delete: resourceKubevirtDataVolumeDelete,
		Exists: resourceKubevirtDataVolumeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: datavolume.DataVolumeFields(),
	}
}

func resourceKubevirtDataVolumeCreate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	dv, err := datavolume.FromResourceData(resourceData)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating new data volume: %#v", dv)
	if err := cli.CreateDataVolume(dv); err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new data volume: %#v", dv)
	if err := datavolume.ToResourceData(*dv, resourceData); err != nil {
		return err
	}
	resourceData.SetId(utils.BuildId(dv.ObjectMeta))

	// Wait for data volume instance's status phase to be succeeded:
	name := dv.ObjectMeta.Name
	namespace := dv.ObjectMeta.Namespace

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Creating"},
		Target:  []string{"Succeeded"},
		Timeout: resourceData.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			var err error
			dv, err = cli.GetDataVolume(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Printf("[DEBUG] data volume %s is not created yet", name)
					return dv, "Creating", nil
				}
				return dv, "", err
			}

			switch dv.Status.Phase {
			case cdiv1.Succeeded, cdiv1.WaitForFirstConsumer:
				return dv, "Succeeded", nil
			case cdiv1.Failed:
				return dv, "", fmt.Errorf("data volume failed to be created, finished with phase=\"failed\"")
			}

			log.Printf("[DEBUG] data volume %s is being created", name)
			return dv, "Creating", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}
	return datavolume.ToResourceData(*dv, resourceData)
}

func resourceKubevirtDataVolumeRead(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading data volume %s", name)

	dv, err := cli.GetDataVolume(namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received data volume: %#v", dv)

	return datavolume.ToResourceData(*dv, resourceData)
}

func resourceKubevirtDataVolumeUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	ops := datavolume.AppendPatchOps("", "", resourceData, make([]patch.PatchOperation, 0, 0))
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating data volume: %s", ops)
	out := &cdiv1.DataVolume{}
	if err := cli.UpdateDataVolume(namespace, name, out, data); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted updated data volume: %#v", out)

	return resourceKubevirtDataVolumeRead(resourceData, meta)
}

func resourceKubevirtDataVolumeDelete(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting data volume: %#v", name)
	if err := cli.DeleteDataVolume(namespace, name); err != nil {
		return err
	}

	// Wait for data volume instance to be removed:
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Timeout: resourceData.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			dv, err := cli.GetDataVolume(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					return nil, "", nil
				}
				return dv, "", err
			}

			log.Printf("[DEBUG] data volume %s is being deleted", dv.GetName())
			return dv, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("[INFO] data volume %s deleted", name)

	resourceData.SetId("")
	return nil
}

func resourceKubevirtDataVolumeExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking data volume %s", name)
	if _, err := cli.GetDataVolume(namespace, name); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return true, err
	}
	return true, nil
}
