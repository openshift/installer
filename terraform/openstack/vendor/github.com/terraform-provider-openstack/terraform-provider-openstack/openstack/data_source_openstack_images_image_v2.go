package openstack

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func dataSourceImagesImageV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagesImageV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(images.ImageVisibilityPublic),
					string(images.ImageVisibilityPrivate),
					string(images.ImageVisibilityShared),
					string(images.ImageVisibilityCommunity),
				}, false),
			},

			"member_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(images.ImageMemberStatusAccepted),
					string(images.ImageMemberStatusPending),
					string(images.ImageMemberStatusRejected),
					string(images.ImageMemberStatusAll),
				}, false),
			},

			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"size_min": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"size_max": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "name",
			},

			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "asc",
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, false),
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"hidden": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},

			// Computed values
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"min_disk_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"min_ram_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

// dataSourceImagesImageV2Read performs the image lookup.
func dataSourceImagesImageV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	visibility := resourceImagesImageV2VisibilityFromString(d.Get("visibility").(string))
	memberStatus := resourceImagesImageV2MemberStatusFromString(d.Get("member_status").(string))

	var tags []string
	tag := d.Get("tag").(string)
	if tag != "" {
		tags = append(tags, tag)
	}

	listOpts := images.ListOpts{
		Name:         d.Get("name").(string),
		Visibility:   visibility,
		Hidden:       d.Get("hidden").(bool),
		Owner:        d.Get("owner").(string),
		Status:       images.ImageStatusActive,
		SizeMin:      int64(d.Get("size_min").(int)),
		SizeMax:      int64(d.Get("size_max").(int)),
		SortKey:      d.Get("sort_key").(string),
		SortDir:      d.Get("sort_direction").(string),
		Tags:         tags,
		MemberStatus: memberStatus,
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var image images.Image
	allPages, err := images.List(imageClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query images: %s", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve images: %s", err)
	}

	properties := resourceImagesImageV2ExpandProperties(
		d.Get("properties").(map[string]interface{}))

	if len(allImages) > 1 {
		allImages = imagesFilterByProperties(allImages, properties)

		log.Printf("[DEBUG] Image list filtered by properties: %#v", properties)
	}

	if len(allImages) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allImages) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] Multiple results found and `most_recent` is set to: %t", recent)
		if recent {
			image = mostRecentImage(allImages)
		} else {
			log.Printf("[DEBUG] Multiple results found: %#v", allImages)
			return diag.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		image = allImages[0]
	}

	log.Printf("[DEBUG] Single Image found: %s", image.ID)

	log.Printf("[DEBUG] openstack_images_image details: %#v", image)

	d.SetId(image.ID)
	d.Set("name", image.Name)
	d.Set("region", GetRegion(d, config))
	d.Set("tags", image.Tags)
	d.Set("container_format", image.ContainerFormat)
	d.Set("disk_format", image.DiskFormat)
	d.Set("min_disk_gb", image.MinDiskGigabytes)
	d.Set("min_ram_mb", image.MinRAMMegabytes)
	d.Set("owner", image.Owner)
	d.Set("protected", image.Protected)
	d.Set("visibility", image.Visibility)
	d.Set("checksum", image.Checksum)
	d.Set("size_bytes", image.SizeBytes)
	d.Set("metadata", image.Metadata)
	d.Set("created_at", image.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", image.UpdatedAt.Format(time.RFC3339))
	d.Set("file", image.File)
	d.Set("schema", image.Schema)
	d.Set("hidden", image.Hidden)

	return nil
}

type imageSort []images.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	itime := a[i].CreatedAt
	jtime := a[j].CreatedAt
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []images.Image) images.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}
