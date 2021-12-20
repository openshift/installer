package ibm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIBMContainerVpcWorkerVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMContainerVpcWorkerVolumeAttachmentRead,

		Schema: map[string]*schema.Schema{
			"volume_attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The volume attachment ID",
			},

			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster name or ID",
			},

			"worker": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Worker node ID",
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the resource group.",
			},

			"volume": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume ID",
			},

			"volume_attachment_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume attachment name",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume attachment status",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of volume",
			},
		},
	}
}

func dataSourceIBMContainerVpcWorkerVolumeAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wpClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	workersAPI := wpClient.Workers()
	target, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterNameorID := d.Get("cluster").(string)
	volumeAttachmentID := d.Get("volume_attachment_id").(string)
	workerID := d.Get("worker").(string)

	volume, err := workersAPI.GetStorageAttachment(clusterNameorID, workerID, volumeAttachmentID, target)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("volume_attachment_name", volume.Name)
	d.Set("status", volume.Status)
	d.Set("volume_type", volume.Type)
	d.SetId(fmt.Sprintf("%s/%s/%s", clusterNameorID, workerID, volumeAttachmentID))
	return nil
}
