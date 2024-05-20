package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/snapshots"
)

func dataSourceSharedFilesystemSnapshotV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedFilesystemSnapshotV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"share_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"share_proto": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"share_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSharedFilesystemSnapshotV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem sfsClient: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	listOpts := snapshots.ListOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ProjectID:   d.Get("project_id").(string),
		Status:      d.Get("status").(string),
	}

	allPages, err := snapshots.ListDetail(sfsClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query snapshots: %s", err)
	}

	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve snapshots: %s", err)
	}

	if len(allSnapshots) < 1 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var share snapshots.Snapshot
	if len(allSnapshots) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allSnapshots)
		return diag.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}
	share = allSnapshots[0]

	dataSourceSharedFilesystemSnapshotV2Attributes(d, &share, GetRegion(d, config))

	return nil
}

func dataSourceSharedFilesystemSnapshotV2Attributes(d *schema.ResourceData, snapshot *snapshots.Snapshot, region string) {
	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("region", region)
	d.Set("project_id", snapshot.ProjectID)
	d.Set("description", snapshot.Description)
	d.Set("size", snapshot.Size)
	d.Set("status", snapshot.Status)
	d.Set("share_id", snapshot.ShareID)
	d.Set("share_proto", snapshot.ShareProto)
	d.Set("share_size", snapshot.ShareSize)
}
