package nutanix

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceNutanixFoundationImage resource creates image in foundation vm for all types of hypervisor & nos images
// Note: source is the path to file in setup where this terraform file runs
func resourceNutanixFoundationImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixFoundationImageCreate,
		ReadContext:   resourceNutanixFoundationImageRead,
		DeleteContext: resourceNutanixFoundationImageDelete,
		Schema: map[string]*schema.Schema{
			"installer_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"md5sum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"in_whitelist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// resourceNutanixFoundationImageCreate creates a image as per installer type, filename and source path
func resourceNutanixFoundationImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// create connection
	conn := meta.(*Client).FoundationClientAPI

	fileName, ok := d.GetOk("filename")
	if !ok {
		return diag.Errorf(getRequiredErrorMessage("filename"))
	}
	installerType, ok := d.GetOk("installer_type")
	if !ok {
		return diag.Errorf(getRequiredErrorMessage("installer_type"))
	}
	source, ok := d.GetOk("source")
	if !ok {
		return diag.Errorf(getRequiredErrorMessage("source"))
	}

	resp, err := conn.FileManagement.UploadImage(ctx, installerType.(string), fileName.(string), source.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("md5sum", resp.Md5Sum)
	d.Set("name", resp.Name)
	d.Set("in_whitelist", resp.InWhitelist)
	d.SetId(fmt.Sprintf("%s_%s", installerType, fileName))
	return nil
}

// resourceNutanixFoundationImageRead will skip as there is no way to read all images
func resourceNutanixFoundationImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// resourceNutanixFoundationImageDelete deletes the existing image
func resourceNutanixFoundationImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// create foundation client
	conn := meta.(*Client).FoundationClientAPI
	installerType := d.Get("installer_type").(string)
	fileName := d.Get("filename").(string)

	// delete the existing image as per the state details
	log.Printf("[Debug] Destroying the image with the name %s", fileName)
	if err := conn.FileManagement.DeleteImage(ctx, installerType, fileName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
