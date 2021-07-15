package azurestack

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/compute/mgmt/compute"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagedDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagedDiskCreateUpdate,
		Read:   resourceArmManagedDiskRead,
		Update: resourceArmManagedDiskCreateUpdate,
		Delete: resourceArmManagedDiskDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"storage_account_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.PremiumLRS),
					string(compute.StandardLRS),
				}, false),
			},

			"create_option": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Copy),
					string(compute.Empty),
					string(compute.FromImage),
					string(compute.Import),
				}, false),
			},

			"source_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"source_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"image_reference_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Windows),
					string(compute.Linux),
				}, false),
			},

			"disk_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateDiskSizeGB,
			},

			"tags": tagsSchema(),
		},
	}
}

func validateDiskSizeGB(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 32767 {
		errors = append(errors, fmt.Errorf(
			"The `disk_size_gb` can only be between 0 and 32767"))
	}
	return warnings, errors
}

func resourceArmManagedDiskCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).diskClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	location := azureStackNormalizeLocation(d.Get("location").(string))
	storageAccountType := d.Get("storage_account_type").(string)
	osType := d.Get("os_type").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	createDisk := compute.Disk{
		Name:     &name,
		Location: &location,
		DiskProperties: &compute.DiskProperties{
			OsType: compute.OperatingSystemTypes(osType),
		},
		Sku: &compute.DiskSku{
			Name: compute.StorageAccountTypes(storageAccountType),
		},
		Tags: *expandedTags,
	}

	if v := d.Get("disk_size_gb"); v != 0 {
		diskSize := int32(v.(int))
		createDisk.DiskProperties.DiskSizeGB = &diskSize
	}

	createOption := d.Get("create_option").(string)
	createDisk.CreationData = &compute.CreationData{
		CreateOption: compute.DiskCreateOption(createOption),
	}

	if strings.EqualFold(createOption, string(compute.Import)) {
		if sourceUri := d.Get("source_uri").(string); sourceUri != "" {
			createDisk.CreationData.SourceURI = &sourceUri
		} else {
			return fmt.Errorf("[ERROR] source_uri must be specified when create_option is `%s`", compute.Import)
		}
	} else if strings.EqualFold(createOption, string(compute.Copy)) {
		if sourceResourceId := d.Get("source_resource_id").(string); sourceResourceId != "" {
			createDisk.CreationData.SourceResourceID = &sourceResourceId
		} else {
			return fmt.Errorf("[ERROR] source_resource_id must be specified when create_option is `%s`", compute.Copy)
		}
	} else if strings.EqualFold(createOption, string(compute.FromImage)) {
		if imageReferenceId := d.Get("image_reference_id").(string); imageReferenceId != "" {
			createDisk.CreationData.ImageReference = &compute.ImageDiskReference{
				ID: utils.String(imageReferenceId),
			}
		} else {
			return fmt.Errorf("[ERROR] image_reference_id must be specified when create_option is `%s`", compute.FromImage)
		}
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, createDisk)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("[ERROR] Cannot read Managed Disk %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmManagedDiskRead(d, meta)
}

func resourceArmManagedDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).diskClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["disks"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Azure Managed Disk %s (resource group %s): %s", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureStackNormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("storage_account_type", string(sku.Name))
	}

	if props := resp.DiskProperties; props != nil {
		if diskSize := props.DiskSizeGB; diskSize != nil {
			d.Set("disk_size_gb", *diskSize)
		}
		if osType := props.OsType; osType != "" {
			d.Set("os_type", string(osType))
		}
	}

	if resp.CreationData != nil {
		flattenAzureRmManagedDiskCreationData(d, resp.CreationData)
	}

	flattenAndSetTags(d, &resp.Tags)
	return nil
}

func resourceArmManagedDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).diskClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["disks"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return err
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return err
		}
	}

	return nil
}

func flattenAzureRmManagedDiskCreationData(d *schema.ResourceData, creationData *compute.CreationData) {
	d.Set("create_option", string(creationData.CreateOption))
	if ref := creationData.ImageReference; ref != nil {
		d.Set("image_reference_id", *ref.ID)
	}
	if id := creationData.SourceResourceID; id != nil {
		d.Set("source_resource_id", *id)
	}
	if creationData.SourceURI != nil {
		d.Set("source_uri", *creationData.SourceURI)
	}
}
