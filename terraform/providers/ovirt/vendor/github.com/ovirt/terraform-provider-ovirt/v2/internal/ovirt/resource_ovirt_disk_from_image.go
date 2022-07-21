package ovirt

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

var diskFromImageSchema = schemaMerge(
	diskBaseSchema, map[string]*schema.Schema{
		"source_file": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      "Path to the local file to upload as the disk image.",
			ValidateDiagFunc: validateLocalFile,
		},
		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Disk size in bytes.",
		},
	},
)

func (p *provider) diskFromImageResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.diskFromImageCreate,
		ReadContext:   p.diskRead,
		UpdateContext: p.diskUpdate,
		DeleteContext: p.diskDelete,
		Schema:        diskFromImageSchema,
		Description:   "The ovirt_disk_from_image resource creates disks in oVirt from a local image file.",
	}
}

func (p *provider) diskFromImageCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	var err error

	storageDomainID := data.Get("storage_domain_id").(string)
	format := data.Get("format").(string)

	params := ovirtclient.CreateDiskParams()
	if alias, ok := data.GetOk("alias"); ok {
		params, err = params.WithAlias(alias.(string))
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid alias value.",
					Detail:   err.Error(),
				},
			}
		}
	}
	if sparse, ok := data.GetOk("sparse"); ok {
		params, err = params.WithSparse(sparse.(bool))
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid sparse value.",
					Detail:   err.Error(),
				},
			}
		}
	}
	sourceFileName := data.Get("source_file").(string)
	sourceFile, err := filepath.Abs(sourceFileName)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to find absolute path for %s", sourceFileName),
				Detail:   err.Error(),
			},
		}
	}
	// We actually want to include the file here, so this is not gosec-relevant.
	fh, err := os.Open(sourceFile) //nolint:gosec
	if err != nil {
		return errorToDiags(fmt.Sprintf("opening file %s", sourceFile), err)
	}
	stat, err := fh.Stat()
	if err != nil {
		return errorToDiags(fmt.Sprintf("opening file %s", sourceFile), err)
	}
	upload, err := client.UploadToNewDisk(
		ovirtclient.StorageDomainID(storageDomainID),
		ovirtclient.ImageFormat(format),
		uint64(stat.Size()),
		params,
		fh,
	)
	var disk ovirtclient.Disk
	if upload != nil {
		disk = upload.Disk()
	}
	if err != nil {
		diags := diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create disk.",
				Detail:   err.Error(),
			},
		}
		if disk != nil {
			if err := disk.Remove(); err != nil && !ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
				diags = append(
					diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Failed to remove created disk %s", disk.ID()),
						Detail:   err.Error(),
					},
				)
			}
		}
		return diags
	}
	return diskResourceUpdate(disk, data)
}
