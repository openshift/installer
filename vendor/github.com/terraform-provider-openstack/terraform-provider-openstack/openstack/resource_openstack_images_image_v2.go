package openstack

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imageimport"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceImagesImageV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceImagesImageV2Create,
		Read:   resourceImagesImageV2Read,
		Update: resourceImagesImageV2Update,
		Delete: resourceImagesImageV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: resourceImagesImageV2UpdateComputedAttributes,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"container_format": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ami", "ari", "aki", "bare", "ovf", "ova",
				}, false),
			},

			"disk_format": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ami", "ari", "aki", "vhd", "vmdk", "raw", "qcow2", "vdi", "iso",
				}, false),
			},

			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"image_cache_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  fmt.Sprintf("%s/.terraform/image_cache", os.Getenv("HOME")),
			},

			"image_source_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"local_file_path"},
			},

			"local_file_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"image_source_url", "web_download"},
			},

			"min_disk_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      0,
			},

			"min_ram_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      0,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"verify_checksum": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"web_download"},
			},

			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				ValidateFunc: validation.StringInSlice([]string{
					"public", "private", "shared", "community",
				}, false),
				Default: "private",
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"web_download": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"local_file_path", "verify_checksum"},
			},

			// Computed-only
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"update_at": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use updated_at instead",
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImagesImageV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	protected := d.Get("protected").(bool)
	visibility := resourceImagesImageV2VisibilityFromString(d.Get("visibility").(string))

	properties := d.Get("properties").(map[string]interface{})
	imageProperties := resourceImagesImageV2ExpandProperties(properties)

	createOpts := &images.CreateOpts{
		Name:            d.Get("name").(string),
		ContainerFormat: d.Get("container_format").(string),
		DiskFormat:      d.Get("disk_format").(string),
		MinDisk:         d.Get("min_disk_gb").(int),
		MinRAM:          d.Get("min_ram_mb").(int),
		Protected:       &protected,
		Visibility:      &visibility,
		Properties:      imageProperties,
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		createOpts.Tags = resourceImagesImageV2BuildTags(tags)
	}

	d.Partial(true)

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newImg, err := images.Create(imageClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Image: %s", err)
	}

	d.SetId(newImg.ID)

	var fileChecksum string
	useWebDownload := d.Get("web_download").(bool)
	if !useWebDownload {
		// variable declaration
		var err error
		var imgFilePath string
		var fileSize int64
		var imgFile *os.File

		// downloading/getting image file props
		imgFilePath, err = resourceImagesImageV2File(d)
		if err != nil {
			return fmt.Errorf("Error opening file for Image: %s", err)
		}
		fileSize, fileChecksum, err = resourceImagesImageV2FileProps(imgFilePath)
		if err != nil {
			return fmt.Errorf("Error getting file props: %s", err)
		}

		// upload
		imgFile, err = os.Open(imgFilePath)
		if err != nil {
			return fmt.Errorf("Error opening file %q: %s", imgFilePath, err)
		}
		defer imgFile.Close()
		log.Printf("[WARN] Uploading image %s (%d bytes). This can be pretty long.", d.Id(), fileSize)

		res := imagedata.Upload(imageClient, d.Id(), imgFile)
		if res.Err != nil {
			return fmt.Errorf("Error while uploading file %q: %s", imgFilePath, res.Err)
		}
	} else {
		// import
		imgURL := d.Get("image_source_url").(string)

		importOpts := &imageimport.CreateOpts{
			Name: imageimport.WebDownloadMethod,
			URI:  imgURL,
		}

		log.Printf("[DEBUG] Import Options: %#v", importOpts)
		res := imageimport.Create(imageClient, d.Id(), importOpts)
		if res.Err != nil {
			return fmt.Errorf("Error while importing url %q: %s", imgURL, res.Err)
		}
	}

	//wait for active
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(images.ImageStatusQueued), string(images.ImageStatusSaving), string(images.ImageStatusImporting)},
		Target:     []string{string(images.ImageStatusActive)},
		Refresh:    resourceImagesImageV2RefreshFunc(imageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Image: %s", err)
	}

	img, err := images.Get(imageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "image")
	}

	if v, ok := d.GetOkExists("verify_checksum"); !useWebDownload && (!ok || (ok && v.(bool))) {
		if img.Checksum != fileChecksum {
			return fmt.Errorf("Error wrong checksum: got %q, expected %q", img.Checksum, fileChecksum)
		}
	}

	d.Partial(false)

	return resourceImagesImageV2Read(d, meta)
}

func resourceImagesImageV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	img, err := images.Get(imageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "image")
	}

	log.Printf("[DEBUG] Retrieved Image %s: %#v", d.Id(), img)

	d.Set("owner", img.Owner)
	d.Set("status", img.Status)
	d.Set("file", img.File)
	d.Set("schema", img.Schema)
	d.Set("checksum", img.Checksum)
	d.Set("size_bytes", img.SizeBytes)
	d.Set("metadata", img.Metadata)
	d.Set("created_at", img.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", img.UpdatedAt.Format(time.RFC3339))
	d.Set("container_format", img.ContainerFormat)
	d.Set("disk_format", img.DiskFormat)
	d.Set("min_disk_gb", img.MinDiskGigabytes)
	d.Set("min_ram_mb", img.MinRAMMegabytes)
	d.Set("file", img.File)
	d.Set("name", img.Name)
	d.Set("protected", img.Protected)
	d.Set("size_bytes", img.SizeBytes)
	d.Set("tags", img.Tags)
	d.Set("visibility", img.Visibility)
	d.Set("region", GetRegion(d, config))

	// Deprecated
	d.Set("update_at", img.UpdatedAt.Format(time.RFC3339))

	properties := resourceImagesImageV2ExpandProperties(img.Properties)
	if err := d.Set("properties", properties); err != nil {
		log.Printf("[WARN] unable to set properties for image %s: %s", img.ID, err)
	}

	return nil
}

func resourceImagesImageV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	updateOpts := make(images.UpdateOpts, 0)

	if d.HasChange("visibility") {
		visibility := resourceImagesImageV2VisibilityFromString(d.Get("visibility").(string))
		v := images.UpdateVisibility{Visibility: visibility}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("name") {
		v := images.ReplaceImageName{NewName: d.Get("name").(string)}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("tags") {
		tags := d.Get("tags").(*schema.Set).List()
		v := images.ReplaceImageTags{
			NewTags: resourceImagesImageV2BuildTags(tags),
		}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("properties") {
		o, n := d.GetChange("properties")
		oldProperties := resourceImagesImageV2ExpandProperties(o.(map[string]interface{}))
		newProperties := resourceImagesImageV2ExpandProperties(n.(map[string]interface{}))

		// Check for new and changed properties
		for newKey, newValue := range newProperties {
			var changed bool

			oldValue, found := oldProperties[newKey]
			if found && (newValue != oldValue) {
				changed = true
			}

			// os_ keys are provided by the OpenStack Image service.
			// These are read-only properties that cannot be modified.
			// Ignore them here and let CustomizeDiff handle them.
			if strings.HasPrefix(newKey, "os_") {
				found = true
				changed = false
			}

			// direct_url is provided by some storage drivers.
			// This is a read-only property that cannot be modified.
			// Ignore it here and let CustomizeDiff handle it.
			if newKey == "direct_url" {
				found = true
				changed = false
			}

			if !found {
				v := images.UpdateImageProperty{
					Op:    images.AddOp,
					Name:  newKey,
					Value: newValue,
				}

				updateOpts = append(updateOpts, v)
			}

			if found && changed {
				v := images.UpdateImageProperty{
					Op:    images.ReplaceOp,
					Name:  newKey,
					Value: newValue,
				}

				updateOpts = append(updateOpts, v)
			}
		}

		// Check for removed properties
		for oldKey := range oldProperties {
			_, found := newProperties[oldKey]

			if !found {
				v := images.UpdateImageProperty{
					Op:   images.RemoveOp,
					Name: oldKey,
				}

				updateOpts = append(updateOpts, v)
			}
		}
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = images.Update(imageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating image: %s", err)
	}

	return resourceImagesImageV2Read(d, meta)
}

func resourceImagesImageV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	log.Printf("[DEBUG] Deleting Image %s", d.Id())
	if err := images.Delete(imageClient, d.Id()).Err; err != nil {
		return fmt.Errorf("Error deleting Image: %s", err)
	}

	d.SetId("")
	return nil
}
