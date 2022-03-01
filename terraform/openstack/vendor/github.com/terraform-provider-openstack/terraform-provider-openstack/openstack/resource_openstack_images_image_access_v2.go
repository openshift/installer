package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/members"
)

func resourceImagesImageAccessV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImagesImageAccessV2Create,
		ReadContext:   resourceImagesImageAccessV2Read,
		UpdateContext: resourceImagesImageAccessV2Update,
		DeleteContext: resourceImagesImageAccessV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"member_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"accepted", "rejected", "pending",
				}, false),
			},

			// Computed-only
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImagesImageAccessV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID := d.Get("image_id").(string)
	memberID := d.Get("member_id").(string)

	// create share request on the producer side
	_, err = members.Create(imageClient, imageID, memberID).Extract()
	if err != nil {
		return diag.Errorf("Error sharing a %q image with the %q member: %s", imageID, memberID, err)
	}

	id := fmt.Sprintf("%s/%s", imageID, memberID)
	d.SetId(id)

	if v, ok := d.GetOkExists("status"); ok {
		d.Partial(true)

		opts := members.UpdateOpts{
			Status: v.(string),
		}
		_, err = members.Update(imageClient, imageID, memberID, opts).Extract()
		if err != nil {
			return diag.Errorf("Error updating the %q image with the %q member: %s", imageID, memberID, err)
		}

		d.Partial(false)
	}

	return resourceImagesImageAccessV2Read(ctx, d, meta)
}

func resourceImagesImageAccessV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	member, err := members.Get(imageClient, imageID, memberID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving the openstack_images_image_access_v2"))
	}

	log.Printf("[DEBUG] Retrieved Image member %s: %#v", d.Id(), member)

	d.Set("status", member.Status)
	d.Set("image_id", member.ImageID)
	d.Set("member_id", member.MemberID)
	// Computed
	d.Set("schema", member.Schema)
	d.Set("created_at", member.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", member.UpdatedAt.Format(time.RFC3339))
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceImagesImageAccessV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	status := d.Get("status").(string)

	opts := members.UpdateOpts{
		Status: status,
	}
	_, err = members.Update(imageClient, imageID, memberID, opts).Extract()
	if err != nil {
		return diag.Errorf("Error updating the %q image with the %q member: %s", imageID, memberID, err)
	}

	return resourceImagesImageAccessV2Read(ctx, d, meta)
}

func resourceImagesImageAccessV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting Image member %s", d.Id())

	if err := members.Delete(imageClient, imageID, memberID).Err; err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting the openstack_images_image_access_v2"))
	}

	return nil
}
