package fsx

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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

// @SDKResource("aws_fsx_lustre_file_system", name="Lustre File System")
// @Tags(identifierAttribute="arn")
func ResourceLustreFileSystem() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLustreFileSystemCreate,
		ReadWithoutTimeout:   resourceLustreFileSystemRead,
		UpdateWithoutTimeout: resourceLustreFileSystemUpdate,
		DeleteWithoutTimeout: resourceLustreFileSystemDelete,
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
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 900),
					validation.StringMatch(regexp.MustCompile(`^s3://`), "must begin with s3://"),
				),
			},
			"import_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 900),
					validation.StringMatch(regexp.MustCompile(`^s3://`), "must begin with s3://"),
				),
			},
			"imported_file_chunk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 512000),
			},
			"mount_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_ids": {
				// As explained in https://docs.aws.amazon.com/fsx/latest/LustreGuide/mounting-on-premises.html, the first
				// network_interface_id is the primary one, so ordering matters. Use TypeList instead of TypeSet to preserve it.
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 50,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"storage_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1200),
			},
			"subnet_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weekly_maintenance_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(7, 7),
					validation.StringMatch(regexp.MustCompile(`^[1-7]:([01]\d|2[0-3]):?([0-5]\d)$`), "must be in the format d:HH:MM"),
				),
			},
			"deployment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      fsx.LustreDeploymentTypeScratch1,
				ValidateFunc: validation.StringInSlice(fsx.LustreDeploymentType_Values(), false),
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"per_unit_storage_throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					12,
					40,
					50,
					100,
					125,
					200,
					250,
					500,
					1000,
				}),
			},
			"automatic_backup_retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 90),
			},
			"daily_automatic_backup_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(5, 5),
					validation.StringMatch(regexp.MustCompile(`^([01]\d|2[0-3]):?([0-5]\d)$`), "must be in the format HH:MM"),
				),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      fsx.StorageTypeSsd,
				ValidateFunc: validation.StringInSlice(fsx.StorageType_Values(), false),
			},
			"drive_cache_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(fsx.DriveCacheType_Values(), false),
			},
			"auto_import_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(fsx.AutoImportPolicyType_Values(), false),
			},
			"copy_tags_to_backups": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"data_compression_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(fsx.DataCompressionType_Values(), false),
				Default:      fsx.DataCompressionTypeNone,
			},
			"file_system_type_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 20),
					validation.StringMatch(regexp.MustCompile(`^[0-9].[0-9]+$`), "must be in format x.y"),
				),
			},
			"log_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidARN,
							StateFunc:    logStateFunc,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return strings.HasPrefix(old, fmt.Sprintf("%s:", new))
							},
						},
						"level": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(fsx.LustreAccessAuditLogLevel_Values(), false),
						},
					},
				},
			},
			"root_squash_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"no_squash_nids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([0-9\[\]\-]*\.){3}([0-9\[\]\-]*)@tcp$`), "must be in the standard Lustre NID foramt"),
							},
						},
						"root_squash": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([0-9]{1,10}):([0-9]{1,10})$`), "must be in the format UID:GID"),
						},
					},
				},
			},
		},

		CustomizeDiff: customdiff.Sequence(
			verify.SetTagsDiff,
			resourceLustreFileSystemSchemaCustomizeDiff,
		),
	}
}

func resourceLustreFileSystemSchemaCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// we want to force a new resource if the new storage capacity is less than the old one
	if d.HasChange("storage_capacity") {
		o, n := d.GetChange("storage_capacity")
		if n.(int) < o.(int) || d.Get("deployment_type").(string) == fsx.LustreDeploymentTypeScratch1 {
			if err := d.ForceNew("storage_capacity"); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceLustreFileSystemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	input := &fsx.CreateFileSystemInput{
		ClientRequestToken: aws.String(id.UniqueId()),
		FileSystemType:     aws.String(fsx.FileSystemTypeLustre),
		StorageCapacity:    aws.Int64(int64(d.Get("storage_capacity").(int))),
		StorageType:        aws.String(d.Get("storage_type").(string)),
		SubnetIds:          flex.ExpandStringList(d.Get("subnet_ids").([]interface{})),
		LustreConfiguration: &fsx.CreateFileSystemLustreConfiguration{
			DeploymentType: aws.String(d.Get("deployment_type").(string)),
		},
		Tags: GetTagsIn(ctx),
	}

	backupInput := &fsx.CreateFileSystemFromBackupInput{
		ClientRequestToken: aws.String(id.UniqueId()),
		StorageType:        aws.String(d.Get("storage_type").(string)),
		SubnetIds:          flex.ExpandStringList(d.Get("subnet_ids").([]interface{})),
		LustreConfiguration: &fsx.CreateFileSystemLustreConfiguration{
			DeploymentType: aws.String(d.Get("deployment_type").(string)),
		},
		Tags: GetTagsIn(ctx),
	}

	//Applicable only for TypePersistent1 and TypePersistent2
	if v, ok := d.GetOk("kms_key_id"); ok {
		input.KmsKeyId = aws.String(v.(string))
		backupInput.KmsKeyId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("automatic_backup_retention_days"); ok {
		input.LustreConfiguration.AutomaticBackupRetentionDays = aws.Int64(int64(v.(int)))
		backupInput.LustreConfiguration.AutomaticBackupRetentionDays = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("daily_automatic_backup_start_time"); ok {
		input.LustreConfiguration.DailyAutomaticBackupStartTime = aws.String(v.(string))
		backupInput.LustreConfiguration.DailyAutomaticBackupStartTime = aws.String(v.(string))
	}

	if v, ok := d.GetOk("export_path"); ok {
		input.LustreConfiguration.ExportPath = aws.String(v.(string))
		backupInput.LustreConfiguration.ExportPath = aws.String(v.(string))
	}

	if v, ok := d.GetOk("import_path"); ok {
		input.LustreConfiguration.ImportPath = aws.String(v.(string))
		backupInput.LustreConfiguration.ImportPath = aws.String(v.(string))
	}

	if v, ok := d.GetOk("imported_file_chunk_size"); ok {
		input.LustreConfiguration.ImportedFileChunkSize = aws.Int64(int64(v.(int)))
		backupInput.LustreConfiguration.ImportedFileChunkSize = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		input.SecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
		backupInput.SecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("weekly_maintenance_start_time"); ok {
		input.LustreConfiguration.WeeklyMaintenanceStartTime = aws.String(v.(string))
		backupInput.LustreConfiguration.WeeklyMaintenanceStartTime = aws.String(v.(string))
	}

	if v, ok := d.GetOk("per_unit_storage_throughput"); ok {
		input.LustreConfiguration.PerUnitStorageThroughput = aws.Int64(int64(v.(int)))
		backupInput.LustreConfiguration.PerUnitStorageThroughput = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("drive_cache_type"); ok {
		input.LustreConfiguration.DriveCacheType = aws.String(v.(string))
		backupInput.LustreConfiguration.DriveCacheType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("auto_import_policy"); ok {
		input.LustreConfiguration.AutoImportPolicy = aws.String(v.(string))
		backupInput.LustreConfiguration.AutoImportPolicy = aws.String(v.(string))
	}

	if v, ok := d.GetOk("copy_tags_to_backups"); ok {
		input.LustreConfiguration.CopyTagsToBackups = aws.Bool(v.(bool))
		backupInput.LustreConfiguration.CopyTagsToBackups = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("data_compression_type"); ok {
		input.LustreConfiguration.DataCompressionType = aws.String(v.(string))
		backupInput.LustreConfiguration.DataCompressionType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("file_system_type_version"); ok {
		input.FileSystemTypeVersion = aws.String(v.(string))
		backupInput.FileSystemTypeVersion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("log_configuration"); ok && len(v.([]interface{})) > 0 {
		input.LustreConfiguration.LogConfiguration = expandLustreLogCreateConfiguration(v.([]interface{}))
		backupInput.LustreConfiguration.LogConfiguration = expandLustreLogCreateConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("root_squash_configuration"); ok && len(v.([]interface{})) > 0 {
		input.LustreConfiguration.RootSquashConfiguration = expandLustreRootSquashConfiguration(v.([]interface{}))
		backupInput.LustreConfiguration.RootSquashConfiguration = expandLustreRootSquashConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("backup_id"); ok {
		backupInput.BackupId = aws.String(v.(string))

		log.Printf("[DEBUG] Creating FSx Lustre File System: %s", backupInput)
		result, err := conn.CreateFileSystemFromBackupWithContext(ctx, backupInput)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "creating FSx Lustre File System from backup: %s", err)
		}

		d.SetId(aws.StringValue(result.FileSystem.FileSystemId))
	} else {
		log.Printf("[DEBUG] Creating FSx Lustre File System: %s", input)
		result, err := conn.CreateFileSystemWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "creating FSx Lustre File System: %s", err)
		}

		d.SetId(aws.StringValue(result.FileSystem.FileSystemId))
	}

	if _, err := waitFileSystemCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for FSx Lustre File System (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceLustreFileSystemRead(ctx, d, meta)...)
}

func resourceLustreFileSystemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	if d.HasChangesExcept("tags_all", "tags") {
		var waitAdminAction = false
		input := &fsx.UpdateFileSystemInput{
			ClientRequestToken:  aws.String(id.UniqueId()),
			FileSystemId:        aws.String(d.Id()),
			LustreConfiguration: &fsx.UpdateFileSystemLustreConfiguration{},
		}

		if d.HasChange("weekly_maintenance_start_time") {
			input.LustreConfiguration.WeeklyMaintenanceStartTime = aws.String(d.Get("weekly_maintenance_start_time").(string))
		}

		if d.HasChange("automatic_backup_retention_days") {
			input.LustreConfiguration.AutomaticBackupRetentionDays = aws.Int64(int64(d.Get("automatic_backup_retention_days").(int)))
		}

		if d.HasChange("daily_automatic_backup_start_time") {
			input.LustreConfiguration.DailyAutomaticBackupStartTime = aws.String(d.Get("daily_automatic_backup_start_time").(string))
		}

		if d.HasChange("auto_import_policy") {
			input.LustreConfiguration.AutoImportPolicy = aws.String(d.Get("auto_import_policy").(string))
		}

		if d.HasChange("storage_capacity") {
			input.StorageCapacity = aws.Int64(int64(d.Get("storage_capacity").(int)))
		}

		if v, ok := d.GetOk("data_compression_type"); ok {
			input.LustreConfiguration.DataCompressionType = aws.String(v.(string))
		}

		if d.HasChange("log_configuration") {
			input.LustreConfiguration.LogConfiguration = expandLustreLogCreateConfiguration(d.Get("log_configuration").([]interface{}))
			waitAdminAction = true
		}

		if d.HasChange("root_squash_configuration") {
			input.LustreConfiguration.RootSquashConfiguration = expandLustreRootSquashConfiguration(d.Get("root_squash_configuration").([]interface{}))
			waitAdminAction = true
		}

		_, err := conn.UpdateFileSystemWithContext(ctx, input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating FSX Lustre File System (%s): %s", d.Id(), err)
		}

		if _, err := waitFileSystemUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for FSx Lustre File System (%s) update: %s", d.Id(), err)
		}

		if waitAdminAction {
			if _, err := waitAdministrativeActionCompleted(ctx, conn, d.Id(), fsx.AdministrativeActionTypeFileSystemUpdate, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for FSx Lustre File System (%s) Log Configuratio to be updated: %s", d.Id(), err)
			}
		}
	}

	return append(diags, resourceLustreFileSystemRead(ctx, d, meta)...)
}

func resourceLustreFileSystemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	filesystem, err := FindFileSystemByID(ctx, conn, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] FSx Lustre File System (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading FSx Lustre File System (%s): %s", d.Id(), err)
	}

	lustreConfig := filesystem.LustreConfiguration

	if filesystem.WindowsConfiguration != nil {
		return sdkdiag.AppendErrorf(diags, "expected FSx Lustre File System, found FSx Windows File System: %s", d.Id())
	}

	if lustreConfig == nil {
		return sdkdiag.AppendErrorf(diags, "describing FSx Lustre File System (%s): empty Lustre configuration", d.Id())
	}

	if lustreConfig.DataRepositoryConfiguration == nil {
		// Initialize an empty structure to simplify d.Set() handling
		lustreConfig.DataRepositoryConfiguration = &fsx.DataRepositoryConfiguration{}
	}

	d.Set("arn", filesystem.ResourceARN)
	d.Set("dns_name", filesystem.DNSName)
	d.Set("export_path", lustreConfig.DataRepositoryConfiguration.ExportPath)
	d.Set("import_path", lustreConfig.DataRepositoryConfiguration.ImportPath)
	d.Set("auto_import_policy", lustreConfig.DataRepositoryConfiguration.AutoImportPolicy)
	d.Set("imported_file_chunk_size", lustreConfig.DataRepositoryConfiguration.ImportedFileChunkSize)
	d.Set("deployment_type", lustreConfig.DeploymentType)
	d.Set("per_unit_storage_throughput", lustreConfig.PerUnitStorageThroughput)
	d.Set("mount_name", lustreConfig.MountName)
	d.Set("storage_type", filesystem.StorageType)
	d.Set("drive_cache_type", lustreConfig.DriveCacheType)
	d.Set("kms_key_id", filesystem.KmsKeyId)

	if err := d.Set("network_interface_ids", aws.StringValueSlice(filesystem.NetworkInterfaceIds)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting network_interface_ids: %s", err)
	}

	d.Set("owner_id", filesystem.OwnerId)
	d.Set("storage_capacity", filesystem.StorageCapacity)

	if err := d.Set("subnet_ids", aws.StringValueSlice(filesystem.SubnetIds)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting subnet_ids: %s", err)
	}

	if err := d.Set("log_configuration", flattenLustreLogConfiguration(lustreConfig.LogConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting log_configuration: %s", err)
	}

	if err := d.Set("root_squash_configuration", flattenLustreRootSquashConfiguration(lustreConfig.RootSquashConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting root_squash_configuration: %s", err)
	}

	SetTagsOut(ctx, filesystem.Tags)

	d.Set("vpc_id", filesystem.VpcId)
	d.Set("weekly_maintenance_start_time", lustreConfig.WeeklyMaintenanceStartTime)
	d.Set("automatic_backup_retention_days", lustreConfig.AutomaticBackupRetentionDays)
	d.Set("daily_automatic_backup_start_time", lustreConfig.DailyAutomaticBackupStartTime)
	d.Set("copy_tags_to_backups", lustreConfig.CopyTagsToBackups)
	d.Set("data_compression_type", lustreConfig.DataCompressionType)
	d.Set("file_system_type_version", filesystem.FileSystemTypeVersion)

	return diags
}

func resourceLustreFileSystemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	request := &fsx.DeleteFileSystemInput{
		FileSystemId: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Deleting FSx Lustre File System: %s", d.Id())
	_, err := conn.DeleteFileSystemWithContext(ctx, request)

	if tfawserr.ErrCodeEquals(err, fsx.ErrCodeFileSystemNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting FSx Lustre File System (%s): %s", d.Id(), err)
	}

	if _, err := waitFileSystemDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for FSx Lustre File System (%s) to deleted: %s", d.Id(), err)
	}

	return diags
}

func expandLustreRootSquashConfiguration(l []interface{}) *fsx.LustreRootSquashConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	data := l[0].(map[string]interface{})
	req := &fsx.LustreRootSquashConfiguration{}

	if v, ok := data["root_squash"].(string); ok && v != "" {
		req.RootSquash = aws.String(v)
	}

	if v, ok := data["no_squash_nids"].(*schema.Set); ok && v.Len() > 0 {
		req.NoSquashNids = flex.ExpandStringSet(v)
	}

	return req
}

func flattenLustreRootSquashConfiguration(adopts *fsx.LustreRootSquashConfiguration) []map[string]interface{} {
	if adopts == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{}

	if adopts.RootSquash != nil {
		m["root_squash"] = aws.StringValue(adopts.RootSquash)
	}

	if adopts.NoSquashNids != nil {
		m["no_squash_nids"] = flex.FlattenStringSet(adopts.NoSquashNids)
	}

	return []map[string]interface{}{m}
}

func expandLustreLogCreateConfiguration(l []interface{}) *fsx.LustreLogCreateConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	data := l[0].(map[string]interface{})
	req := &fsx.LustreLogCreateConfiguration{
		Level: aws.String(data["level"].(string)),
	}

	if v, ok := data["destination"].(string); ok && v != "" {
		req.Destination = aws.String(logStateFunc(v))
	}

	return req
}

func flattenLustreLogConfiguration(adopts *fsx.LustreLogConfiguration) []map[string]interface{} {
	if adopts == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"level": aws.StringValue(adopts.Level),
	}

	if adopts.Destination != nil {
		m["destination"] = aws.StringValue(adopts.Destination)
	}

	return []map[string]interface{}{m}
}

func logStateFunc(v interface{}) string {
	value := v.(string)
	// API returns the specific log stream arn instead of provided log group
	logArn, _ := arn.Parse(value)
	if logArn.Service == "logs" {
		parts := strings.SplitN(logArn.Resource, ":", 3)
		if len(parts) == 3 {
			return strings.TrimSuffix(value, fmt.Sprintf(":%s", parts[2]))
		} else {
			return value
		}
	}
	return value
}
