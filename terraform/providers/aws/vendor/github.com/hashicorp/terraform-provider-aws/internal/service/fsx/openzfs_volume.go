package fsx

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_fsx_openzfs_volume", name="OpenZFS Volume")
// @Tags(identifierAttribute="arn")
func ResourceOpenzfsVolume() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceOpenzfsVolumeCreate,
		ReadWithoutTimeout:   resourceOpenzfsVolumeRead,
		UpdateWithoutTimeout: resourceOpenzfsVolumeUpdate,
		DeleteWithoutTimeout: resourceOpenzfsVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"copy_tags_to_snapshots": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"data_compression_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validation.StringInSlice(fsx.OpenZFSDataCompressionType_Values(), false),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 203),
			},
			"nfs_exports": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_configurations": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 25,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"clients": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 128),
											validation.StringMatch(regexp.MustCompile(`^[ -~]{1,128}$`), "must be either IP Address or CIDR"),
										),
									},
									"options": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										MaxItems: 20,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringLenBetween(1, 128),
										},
									},
								},
							},
						},
					},
				},
			},
			"origin_snapshot": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"copy_strategy": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(fsx.OpenZFSCopyStrategy_Values(), false),
						},
						"snapshot_arn": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(8, 512),
								validation.StringMatch(regexp.MustCompile(`^arn:.*`), "must specify the full ARN of the snapshot"),
							),
						},
					},
				},
			},
			"parent_volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(23, 23),
					validation.StringMatch(regexp.MustCompile(`^(fsvol-[0-9a-f]{17,})$`), "must specify a filesystem id i.e. fs-12345678"),
				),
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"record_size_kib": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      128,
				ValidateFunc: validation.IntInSlice([]int{4, 8, 16, 32, 64, 128, 256, 512, 1024}),
			},
			"storage_capacity_quota_gib": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 2147483647),
			},
			"storage_capacity_reservation_gib": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 2147483647),
			},
			"user_and_group_quotas": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 2147483647),
						},
						"storage_capacity_quota_gib": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 2147483647),
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(fsx.OpenZFSQuotaType_Values(), false),
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"volume_type": {
				Type:         schema.TypeString,
				Default:      fsx.VolumeTypeOpenzfs,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(fsx.VolumeType_Values(), false),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceOpenzfsVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	input := &fsx.CreateVolumeInput{
		ClientRequestToken: aws.String(id.UniqueId()),
		Name:               aws.String(d.Get("name").(string)),
		VolumeType:         aws.String(d.Get("volume_type").(string)),
		OpenZFSConfiguration: &fsx.CreateOpenZFSVolumeConfiguration{
			ParentVolumeId: aws.String(d.Get("parent_volume_id").(string)),
		},
		Tags: GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("copy_tags_to_snapshots"); ok {
		input.OpenZFSConfiguration.CopyTagsToSnapshots = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("data_compression_type"); ok {
		input.OpenZFSConfiguration.DataCompressionType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("nfs_exports"); ok {
		input.OpenZFSConfiguration.NfsExports = expandOpenzfsVolumeNFSExports(v.([]interface{}))
	}

	if v, ok := d.GetOk("read_only"); ok {
		input.OpenZFSConfiguration.ReadOnly = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("record_size_kib"); ok {
		input.OpenZFSConfiguration.RecordSizeKiB = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("storage_capacity_quota_gib"); ok {
		input.OpenZFSConfiguration.StorageCapacityQuotaGiB = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("storage_capacity_reservation_gib"); ok {
		input.OpenZFSConfiguration.StorageCapacityReservationGiB = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("user_and_group_quotas"); ok {
		input.OpenZFSConfiguration.UserAndGroupQuotas = expandOpenzfsVolumeUserAndGroupQuotas(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("origin_snapshot"); ok {
		input.OpenZFSConfiguration.OriginSnapshot = expandOpenzfsCreateVolumeOriginSnapshot(v.([]interface{}))

		log.Printf("[DEBUG] Creating FSx OpenZFS Volume: %s", input)
		result, err := conn.CreateVolumeWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "creating FSx OpenZFS Volume from snapshot: %s", err)
		}

		d.SetId(aws.StringValue(result.Volume.VolumeId))
	} else {
		log.Printf("[DEBUG] Creating FSx OpenZFS Volume: %s", input)
		result, err := conn.CreateVolumeWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "creating FSx OpenZFS Volume: %s", err)
		}

		d.SetId(aws.StringValue(result.Volume.VolumeId))
	}

	if _, err := waitVolumeCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for FSx OpenZFS Volume(%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceOpenzfsVolumeRead(ctx, d, meta)...)
}

func resourceOpenzfsVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	volume, err := FindVolumeByID(ctx, conn, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] FSx OpenZFS volume (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading FSx OpenZFS Volume (%s): %s", d.Id(), err)
	}

	openzfsConfig := volume.OpenZFSConfiguration

	if volume.OntapConfiguration != nil {
		return sdkdiag.AppendErrorf(diags, "expected FSx OpeZFS Volume, found FSx ONTAP Volume: %s", d.Id())
	}

	if openzfsConfig == nil {
		return sdkdiag.AppendErrorf(diags, "describing FSx OpenZFS Volume (%s): empty Openzfs configuration", d.Id())
	}

	d.Set("arn", volume.ResourceARN)
	d.Set("copy_tags_to_snapshots", openzfsConfig.CopyTagsToSnapshots)
	d.Set("data_compression_type", openzfsConfig.DataCompressionType)
	d.Set("name", volume.Name)
	d.Set("parent_volume_id", openzfsConfig.ParentVolumeId)
	d.Set("read_only", openzfsConfig.ReadOnly)
	d.Set("record_size_kib", openzfsConfig.RecordSizeKiB)
	d.Set("storage_capacity_quota_gib", openzfsConfig.StorageCapacityQuotaGiB)
	d.Set("storage_capacity_reservation_gib", openzfsConfig.StorageCapacityReservationGiB)
	d.Set("volume_type", volume.VolumeType)

	if err := d.Set("origin_snapshot", flattenOpenzfsVolumeOriginSnapshot(openzfsConfig.OriginSnapshot)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting nfs_exports: %s", err)
	}

	if err := d.Set("nfs_exports", flattenOpenzfsVolumeNFSExports(openzfsConfig.NfsExports)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting nfs_exports: %s", err)
	}

	if err := d.Set("user_and_group_quotas", flattenOpenzfsVolumeUserAndGroupQuotas(openzfsConfig.UserAndGroupQuotas)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting user_and_group_quotas: %s", err)
	}

	return diags
}

func resourceOpenzfsVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	if d.HasChangesExcept("tags_all", "tags") {
		input := &fsx.UpdateVolumeInput{
			ClientRequestToken:   aws.String(id.UniqueId()),
			VolumeId:             aws.String(d.Id()),
			OpenZFSConfiguration: &fsx.UpdateOpenZFSVolumeConfiguration{},
		}

		if d.HasChange("data_compression_type") {
			input.OpenZFSConfiguration.DataCompressionType = aws.String(d.Get("data_compression_type").(string))
		}

		if d.HasChange("name") {
			input.Name = aws.String(d.Get("name").(string))
		}

		if d.HasChange("nfs_exports") {
			input.OpenZFSConfiguration.NfsExports = expandOpenzfsVolumeNFSExports(d.Get("nfs_exports").([]interface{}))
		}

		if d.HasChange("read_only") {
			input.OpenZFSConfiguration.ReadOnly = aws.Bool(d.Get("read_only").(bool))
		}

		if d.HasChange("record_size_kib") {
			input.OpenZFSConfiguration.RecordSizeKiB = aws.Int64(int64(d.Get("record_size_kib").(int)))
		}

		if d.HasChange("storage_capacity_quota_gib") {
			input.OpenZFSConfiguration.StorageCapacityQuotaGiB = aws.Int64(int64(d.Get("storage_capacity_quota_gib").(int)))
		}

		if d.HasChange("storage_capacity_reservation_gib") {
			input.OpenZFSConfiguration.StorageCapacityReservationGiB = aws.Int64(int64(d.Get("storage_capacity_reservation_gib").(int)))
		}

		if d.HasChange("user_and_group_quotas") {
			input.OpenZFSConfiguration.UserAndGroupQuotas = expandOpenzfsVolumeUserAndGroupQuotas(d.Get("user_and_group_quotas").(*schema.Set).List())
		}

		_, err := conn.UpdateVolumeWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating FSx OpenZFS Volume (%s): %s", d.Id(), err)
		}

		if _, err := waitVolumeUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for FSx OpenZFS Volume (%s) update: %s", d.Id(), err)
		}
	}

	return append(diags, resourceOpenzfsVolumeRead(ctx, d, meta)...)
}

func resourceOpenzfsVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	log.Printf("[DEBUG] Deleting FSx OpenZFS Volume: %s", d.Id())
	_, err := conn.DeleteVolumeWithContext(ctx, &fsx.DeleteVolumeInput{
		VolumeId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, fsx.ErrCodeVolumeNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting FSx OpenZFS Volume (%s): %s", d.Id(), err)
	}

	if _, err := waitVolumeDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for FSx OpenZFS Volume (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func expandOpenzfsVolumeUserAndGroupQuotas(cfg []interface{}) []*fsx.OpenZFSUserOrGroupQuota {
	quotas := []*fsx.OpenZFSUserOrGroupQuota{}

	for _, quota := range cfg {
		expandedQuota := expandOpenzfsVolumeUserAndGroupQuota(quota.(map[string]interface{}))
		if expandedQuota != nil {
			quotas = append(quotas, expandedQuota)
		}
	}

	return quotas
}

func expandOpenzfsVolumeUserAndGroupQuota(conf map[string]interface{}) *fsx.OpenZFSUserOrGroupQuota {
	if len(conf) < 1 {
		return nil
	}

	out := fsx.OpenZFSUserOrGroupQuota{}

	if v, ok := conf["id"].(int); ok {
		out.Id = aws.Int64(int64(v))
	}

	if v, ok := conf["storage_capacity_quota_gib"].(int); ok {
		out.StorageCapacityQuotaGiB = aws.Int64(int64(v))
	}

	if v, ok := conf["type"].(string); ok {
		out.Type = aws.String(v)
	}

	return &out
}

func expandOpenzfsVolumeNFSExports(cfg []interface{}) []*fsx.OpenZFSNfsExport {
	exports := []*fsx.OpenZFSNfsExport{}

	for _, export := range cfg {
		expandedExport := expandOpenzfsVolumeNFSExport(export.(map[string]interface{}))
		if expandedExport != nil {
			exports = append(exports, expandedExport)
		}
	}

	return exports
}

func expandOpenzfsVolumeNFSExport(cfg map[string]interface{}) *fsx.OpenZFSNfsExport {
	out := fsx.OpenZFSNfsExport{}

	if v, ok := cfg["client_configurations"]; ok {
		out.ClientConfigurations = expandOpenzfsVolumeClinetConfigurations(v.(*schema.Set).List())
	}

	return &out
}

func expandOpenzfsVolumeClinetConfigurations(cfg []interface{}) []*fsx.OpenZFSClientConfiguration {
	configurations := []*fsx.OpenZFSClientConfiguration{}

	for _, configuration := range cfg {
		expandedConfiguration := expandOpenzfsVolumeClientConfiguration(configuration.(map[string]interface{}))
		if expandedConfiguration != nil {
			configurations = append(configurations, expandedConfiguration)
		}
	}

	return configurations
}

func expandOpenzfsVolumeClientConfiguration(conf map[string]interface{}) *fsx.OpenZFSClientConfiguration {
	out := fsx.OpenZFSClientConfiguration{}

	if v, ok := conf["clients"].(string); ok && len(v) > 0 {
		out.Clients = aws.String(v)
	}

	if v, ok := conf["options"].([]interface{}); ok {
		out.Options = flex.ExpandStringList(v)
	}

	return &out
}

func expandOpenzfsCreateVolumeOriginSnapshot(cfg []interface{}) *fsx.CreateOpenZFSOriginSnapshotConfiguration {
	if len(cfg) < 1 {
		return nil
	}

	conf := cfg[0].(map[string]interface{})

	out := fsx.CreateOpenZFSOriginSnapshotConfiguration{}

	if v, ok := conf["copy_strategy"].(string); ok {
		out.CopyStrategy = aws.String(v)
	}

	if v, ok := conf["snapshot_arn"].(string); ok {
		out.SnapshotARN = aws.String(v)
	}

	return &out
}

func flattenOpenzfsVolumeNFSExports(rs []*fsx.OpenZFSNfsExport) []map[string]interface{} {
	exports := make([]map[string]interface{}, 0)

	for _, export := range rs {
		if export != nil {
			cfg := make(map[string]interface{})
			cfg["client_configurations"] = flattenOpenzfsVolumeClientConfigurations(export.ClientConfigurations)
			exports = append(exports, cfg)
		}
	}

	if len(exports) > 0 {
		return exports
	}

	return nil
}

func flattenOpenzfsVolumeClientConfigurations(rs []*fsx.OpenZFSClientConfiguration) []map[string]interface{} {
	configurations := make([]map[string]interface{}, 0)

	for _, configuration := range rs {
		if configuration != nil {
			cfg := make(map[string]interface{})
			cfg["clients"] = aws.StringValue(configuration.Clients)
			cfg["options"] = flex.FlattenStringList(configuration.Options)
			configurations = append(configurations, cfg)
		}
	}

	if len(configurations) > 0 {
		return configurations
	}

	return nil
}

func flattenOpenzfsVolumeUserAndGroupQuotas(rs []*fsx.OpenZFSUserOrGroupQuota) []map[string]interface{} {
	quotas := make([]map[string]interface{}, 0)

	for _, quota := range rs {
		if quota != nil {
			cfg := make(map[string]interface{})
			cfg["id"] = aws.Int64Value(quota.Id)
			cfg["storage_capacity_quota_gib"] = aws.Int64Value(quota.StorageCapacityQuotaGiB)
			cfg["type"] = aws.StringValue(quota.Type)
			quotas = append(quotas, cfg)
		}
	}

	if len(quotas) > 0 {
		return quotas
	}

	return nil
}

func flattenOpenzfsVolumeOriginSnapshot(rs *fsx.OpenZFSOriginSnapshotConfiguration) []interface{} {
	if rs == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})
	if rs.CopyStrategy != nil {
		m["copy_strategy"] = aws.StringValue(rs.CopyStrategy)
	}
	if rs.SnapshotARN != nil {
		m["snapshot_arn"] = aws.StringValue(rs.SnapshotARN)
	}

	return []interface{}{m}
}
