package ovirt

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
)

var diskBaseSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_domain_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the storage domain to use for disk creation.",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
	"format": {
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"Format for the disk. One of: `%s`",
			strings.Join(ovirtclient.ImageFormatValues().Strings(), "`, `"),
		),
		ValidateDiagFunc: validateFormat,
		ForceNew:         true,
	},
	"alias": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Human-readable alias for the disk.",
	},
	"sparse": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Description: "Use sparse provisioning for disk.",
	},
	"total_size": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Size of the actual image size on the disk in bytes.",
	},
	"status": {
		Type:     schema.TypeString,
		Computed: true,
		Description: fmt.Sprintf(
			"Status of the disk. One of: `%s`.",
			strings.Join(ovirtclient.VMStatusValues().Strings(), "`, `"),
		),
	},
}

var diskSchema = schemaMerge(
	diskBaseSchema, map[string]*schema.Schema{
		"size": {
			Type:             schema.TypeInt,
			Required:         true,
			Description:      "Disk size in bytes.",
			ValidateDiagFunc: validateDiskSize,
			ForceNew:         true,
		},
	},
)

func (p *provider) diskResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.diskCreate,
		ReadContext:   p.diskRead,
		UpdateContext: p.diskUpdate,
		DeleteContext: p.diskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: p.diskImport,
		},
		Schema:      diskSchema,
		Description: "The ovirt_disk resource creates disks in oVirt.",
	}
}

func (p *provider) diskCreate(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	var err error

	storageDomainID := data.Get("storage_domain_id").(string)
	format := data.Get("format").(string)
	size := data.Get("size").(int)

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
	// GetOkExists is necessary here due to GetOk check for default values (for sparse=false, ok would be false, too)
	// see: https://github.com/hashicorp/terraform/pull/15723
	//nolint:staticcheck
	if sparse, ok := data.GetOkExists("sparse"); ok {
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

	disk, err := client.CreateDisk(
		ovirtclient.StorageDomainID(storageDomainID),
		ovirtclient.ImageFormat(format),
		uint64(size),
		params,
	)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create disk.",
				Detail:   err.Error(),
			},
		}
	}

	return diskResourceUpdate(disk, data)
}

func diskResourceUpdate(disk ovirtclient.Disk, data *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	data.SetId(string(disk.ID()))
	diags = setResourceField(data, "alias", disk.Alias(), diags)
	diags = setResourceField(data, "format", string(disk.Format()), diags)
	diags = setResourceField(data, "size", disk.ProvisionedSize(), diags)
	diags = setResourceField(data, "sparse", disk.Sparse(), diags)
	diags = setResourceField(data, "total_size", disk.TotalSize(), diags)
	diags = setResourceField(data, "status", disk.Status(), diags)

	desiredStorageDomainID := ovirtclient.StorageDomainID(data.Get("storage_domain_id").(string))
	foundStorageDomain := false
	for _, storageDomainID := range disk.StorageDomainIDs() {
		if desiredStorageDomainID == storageDomainID {
			foundStorageDomain = true
		}
	}
	if foundStorageDomain {
		diags = setResourceField(data, "storage_domain_id", desiredStorageDomainID, diags)
	} else {
		diags = setResourceField(data, "storage_domain_id", "", diags)
	}

	return diags
}

func (p *provider) diskRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	disk, err := client.GetDisk(ovirtclient.DiskID(data.Id()))
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
			return nil
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to fetch disk.",
				Detail:   err.Error(),
			},
		}
	}
	return diskResourceUpdate(disk, data)
}

func (p *provider) diskUpdate(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	params := ovirtclient.UpdateDiskParams()
	var err error
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
	disk, err := client.UpdateDisk(ovirtclient.DiskID(data.Id()), params)
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
			return nil
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to update disk.",
				Detail:   err.Error(),
			},
		}
	}
	return diskResourceUpdate(disk, data)
}

func (p *provider) diskDelete(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	if err := client.RemoveDisk(ovirtclient.DiskID(data.Id())); err != nil {
		if isNotFound(err) {
			data.SetId("")
			return nil
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to remove disk.",
				Detail:   err.Error(),
			},
		}
	}
	data.SetId("")
	return nil
}

func (p *provider) diskImport(ctx context.Context, data *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData,
	error,
) {
	client := p.client.WithContext(ctx)
	disk, err := client.GetDisk(ovirtclient.DiskID(data.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to import disk %s (%w)", data.Id(), err)
	}
	diags := diskResourceUpdate(disk, data)
	if err := diagsToError(diags); err != nil {
		return nil, fmt.Errorf("failed to import disk (%w)", err)
	}
	return []*schema.ResourceData{
		data,
	}, nil
}
