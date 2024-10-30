package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"time"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	volumeAttaching                = "attaching"
	volumeAttached                 = "attached"
	volumeAttachSuccessStatus      = "active"
	volumeAttachProgressStatus     = "in progress"
	volumeAttachProvisioningStatus = "provisioning"
	volumeAttachInactiveStatus     = "inactive"
	volumeAttachFailStatus         = "failed"
	volumeAttachRemovedStatus      = "removed"
	volumeAttachReclamation        = "pending_reclamation"
)

func ResourceIBMContainerVpcWorkerVolumeAttachment() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceIBMContainerVpcWorkerVolumeAttachmentCreate,
		ReadContext:   resourceIBMContainerVpcWorkerVolumeAttachmentRead,
		DeleteContext: resourceIBMContainerVpcWorkerVolumeAttachmentDelete,
		Exists:        resourceIBMContainerVpcWorkerVolumeAttachmentExists,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"volume": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC Volume ID",
			},

			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster name or ID",
				ValidateFunc: validate.InvokeValidator(
					"ibm_container_storage_attachment",
					"cluster"),
			},

			"worker": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "worker node ID",
			},

			"volume_attachment_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume attachment name",
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the resource group.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume attachment status",
			},

			"volume_attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume attachment ID",
			},

			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of volume",
			},
		},
	}
}
func ResourceIBMContainerVpcWorkerVolumeAttachmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMContainerVpcWorkerVolumeAttachmentValidator := validate.ResourceValidator{ResourceName: "ibm_container_storage_attachment", Schema: validateSchema}
	return &iBMContainerVpcWorkerVolumeAttachmentValidator
}

func resourceIBMContainerVpcWorkerVolumeAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	workersAPI := wpClient.Workers()

	volumeID := d.Get("volume").(string)
	clusterNameorID := d.Get("cluster").(string)
	workerID := d.Get("worker").(string)

	target, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return diag.FromErr(err)
	}
	attachVolumeRequest := v2.VolumeRequest{
		Cluster:  clusterNameorID,
		VolumeID: volumeID,
		Worker:   workerID,
	}

	volumeattached, err := workersAPI.CreateStorageAttachment(attachVolumeRequest, target)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", clusterNameorID, workerID, volumeattached.Id))
	_, attachErr := waitforVolumetoAttach(d, meta)
	if attachErr != nil {
		return diag.FromErr(attachErr)
	}
	return resourceIBMContainerVpcWorkerVolumeAttachmentRead(context, d, meta)
}

func resourceIBMContainerVpcWorkerVolumeAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	workersAPI := wpClient.Workers()
	target, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	clusterNameorID := parts[0]
	workerID := parts[1]
	volumeAttachmentID := parts[2]

	volumeAttachment, err := workersAPI.GetStorageAttachment(clusterNameorID, workerID, volumeAttachmentID, target)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("cluster", clusterNameorID)
	d.Set("worker", workerID)
	d.Set("volume", volumeAttachment.Volume.Id)
	d.Set("volume_attachment_id", volumeAttachmentID)
	d.Set("volume_attachment_name", volumeAttachment.Name)
	d.Set("status", volumeAttachment.Status)
	d.Set("volume_type", volumeAttachment.Type)
	return nil
}

func resourceIBMContainerVpcWorkerVolumeAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	workersAPI := wpClient.Workers()
	target, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	clusterNameorID := parts[0]
	workerID := parts[1]
	volumeAttachmentID := parts[2]

	detachVolumeRequest := v2.VolumeRequest{
		Cluster:            clusterNameorID,
		VolumeAttachmentID: volumeAttachmentID,
		Worker:             workerID,
	}

	response, deleteErr := workersAPI.DeleteStorageAttachment(detachVolumeRequest, target)
	if deleteErr != nil && !strings.Contains(deleteErr.Error(), "EmptyResponseBody") {
		if response != "Ok" && strings.Contains(response, "Not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Failed to delete the volume attachment: %s", deleteErr))
	}

	_, err = waitForStorageAttachmentDelete(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for storage attachment (%s) to be deleted: %s", d.Id(), err))
	}

	d.SetId("")
	return nil
}

func resourceIBMContainerVpcWorkerVolumeAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}

	workersAPI := wpClient.Workers()
	target, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return false, err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	clusterNameorID := parts[0]
	workerID := parts[1]
	volumeAttachmentID := parts[2]

	_, getErr := workersAPI.GetStorageAttachment(clusterNameorID, workerID, volumeAttachmentID, target)
	if getErr != nil {
		if apiErr, ok := getErr.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting storage attachement: %s", getErr)
	}
	return true, nil
}

func waitforVolumetoAttach(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}

	workersAPI := wpClient.Workers()

	target, trgetErr := getVpcClusterTargetHeader(d)
	if trgetErr != nil {
		return nil, trgetErr
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil, err
	}
	clusterNameorID := parts[0]
	workerID := parts[1]
	volumeAttachmentID := parts[2]

	createStateConf := &resource.StateChangeConf{
		Pending: []string{volumeAttaching},
		Target:  []string{volumeAttached},
		Refresh: func() (interface{}, string, error) {
			volume, err := workersAPI.GetStorageAttachment(clusterNameorID, workerID, volumeAttachmentID, target)

			if err != nil {
				return volume, volumeAttaching, err
			}

			if volume.Status == volumeAttached {
				return volume, volumeAttached, nil
			}
			return volume, volumeAttaching, nil

		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	return createStateConf.WaitForState()
}

func waitForStorageAttachmentDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}

	workersAPI := wpClient.Workers()

	target, trgetErr := getVpcClusterTargetHeader(d)
	if trgetErr != nil {
		return nil, trgetErr
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil, err
	}
	clusterNameorID := parts[0]
	workerID := parts[1]
	volumeAttachmentID := parts[2]

	stateConf := &resource.StateChangeConf{
		Pending: []string{"inprogress"},
		Target:  []string{"removed"},
		Refresh: func() (interface{}, string, error) {
			volume, getErr := workersAPI.GetStorageAttachment(clusterNameorID, workerID, volumeAttachmentID, target)
			if getErr != nil {
				if apiErr, ok := getErr.(bmxerror.RequestFailure); ok {
					if apiErr.StatusCode() == 404 {
						return volume, "removed", nil
					}
					return nil, "", fmt.Errorf("[ERROR] Reading volume attach %s failed with resp code: %s, err: %v", d.Id(), volume, getErr)
				}
				return nil, "", fmt.Errorf("[ERROR] Reading volume attach failed: %s", getErr)
			}
			return volume, "inprogress", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForState()
}
