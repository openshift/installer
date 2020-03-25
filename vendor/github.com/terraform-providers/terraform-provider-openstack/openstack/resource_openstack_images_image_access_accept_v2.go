package openstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/members"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceImagesImageAccessAcceptV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceImagesImageAccessAcceptV2Create,
		Read:   resourceImagesImageAccessAcceptV2Read,
		Update: resourceImagesImageAccessAcceptV2Update,
		Delete: resourceImagesImageAccessAcceptV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceImagesImageAccessAcceptV2Import,
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
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceImagesImageAccessAcceptV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID := d.Get("image_id").(string)
	memberID := d.Get("member_id").(string)
	status := d.Get("status").(string)

	if memberID == "" {
		memberID, err = resourceImagesImageAccessV2DetectMemberID(imageClient, imageID)
		if err != nil {
			return err
		}
	}

	// accept status on the consumer side
	opts := members.UpdateOpts{
		Status: status,
	}
	_, err = members.Update(imageClient, imageID, memberID, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error setting a member status to the %q image share for the %q member: %s", imageID, memberID, err)
	}

	id := fmt.Sprintf("%s/%s", imageID, memberID)
	d.SetId(id)

	return resourceImagesImageAccessAcceptV2Read(d, meta)
}

func resourceImagesImageAccessAcceptV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return err
	}

	member, err := members.Get(imageClient, imageID, memberID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving the openstack_images_image_access_accept_v2")
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

func resourceImagesImageAccessAcceptV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return err
	}

	status := d.Get("status").(string)

	opts := members.UpdateOpts{
		Status: status,
	}
	_, err = members.Update(imageClient, imageID, memberID, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error updateing the %q image with the %q member: %s", imageID, memberID, err)
	}

	return resourceImagesImageAccessAcceptV2Read(d, meta)
}

func resourceImagesImageAccessAcceptV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Rejecting Image membership %s", d.Id())
	// reject status on the consumer side
	opts := members.UpdateOpts{
		Status: "rejected",
	}
	if err := members.Update(imageClient, imageID, memberID, opts).Err; err != nil {
		return CheckDeleted(d, err, "Error rejecting the openstack_images_image_access_accept_v2")
	}

	return nil
}

func resourceImagesImageAccessAcceptV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)

	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID := parts[0]
	memberID := ""
	if len(parts) > 1 {
		memberID = parts[1]
	} else {
		memberID, err = resourceImagesImageAccessV2DetectMemberID(imageClient, imageID)
		if err != nil {
			return nil, err
		}
	}

	id := fmt.Sprintf("%s/%s", imageID, memberID)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
