package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
)

func dataSourceBlockStorageQuotasetV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBlockStorageQuotasetV3Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volumes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"snapshots": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"gigabytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"per_volume_gigabytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"backups": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"backup_gigabytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"groups": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"volume_type_quota": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceBlockStorageQuotasetV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	blockStorageClient, err := config.BlockStorageV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	projectID := d.Get("project_id").(string)

	q, err := quotasets.Get(blockStorageClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_blockstorage_quotaset_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_quotaset_v3 %s: %#v", d.Id(), q)

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)
	d.Set("project_id", projectID)
	d.Set("region", region)
	d.Set("volumes", q.Volumes)
	d.Set("snapshots", q.Snapshots)
	d.Set("gigabytes", q.Gigabytes)
	d.Set("per_volume_gigabytes", q.PerVolumeGigabytes)
	d.Set("backups", q.Backups)
	d.Set("backup_gigabytes", q.BackupGigabytes)
	d.Set("groups", q.Groups)

	volumeTypeQuota, _ := blockStorageQuotasetVolTypeQuotaToStr(q.Extra)
	if err := d.Set("volume_type_quota", volumeTypeQuota); err != nil {
		log.Printf("[WARN] Unable to set openstack_blockstorage_quotaset_v3 %s volume_type_quotas: %s", d.Id(), err)
	}

	return nil
}
