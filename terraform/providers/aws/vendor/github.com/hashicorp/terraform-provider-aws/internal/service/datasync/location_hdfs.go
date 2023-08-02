package datasync

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/datasync"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

// @SDKResource("aws_datasync_location_hdfs", name="Location HDFS")
// @Tags(identifierAttribute="id")
func ResourceLocationHDFS() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLocationHDFSCreate,
		ReadWithoutTimeout:   resourceLocationHDFSRead,
		UpdateWithoutTimeout: resourceLocationHDFSUpdate,
		DeleteWithoutTimeout: resourceLocationHDFSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_arns": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidARN,
				},
			},
			"authentication_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(datasync.HdfsAuthenticationType_Values(), false),
			},
			"kerberos_keytab": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kerberos_krb5_conf": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kerberos_principal": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"kms_key_provider_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"block_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  128 * 1024 * 1024, // 128 MiB
				ValidateFunc: validation.All(
					validation.IntDivisibleBy(512),
					validation.IntBetween(1048576, 1073741824),
				),
			},
			"replication_factor": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(1, 512),
			},
			"name_node": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IsPortNumber,
						},
					},
				},
			},
			"qop_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_transfer_protection": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(datasync.HdfsDataTransferProtection_Values(), false),
						},
						"rpc_protection": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(datasync.HdfsRpcProtection_Values(), false),
						},
					},
				},
			},
			"simple_user": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"subdirectory": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/",
				ValidateFunc: validation.StringLenBetween(1, 4096),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "/" {
						return false
					}
					if strings.TrimSuffix(old, "/") == strings.TrimSuffix(new, "/") {
						return true
					}
					return false
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceLocationHDFSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	input := &datasync.CreateLocationHdfsInput{
		AgentArns:          flex.ExpandStringSet(d.Get("agent_arns").(*schema.Set)),
		NameNodes:          expandHDFSNameNodes(d.Get("name_node").(*schema.Set)),
		AuthenticationType: aws.String(d.Get("authentication_type").(string)),
		Subdirectory:       aws.String(d.Get("subdirectory").(string)),
		Tags:               GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("simple_user"); ok {
		input.SimpleUser = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kerberos_krb5_conf"); ok {
		input.KerberosKrb5Conf = []byte(v.(string))
	}

	if v, ok := d.GetOk("kerberos_keytab"); ok {
		input.KerberosKeytab = []byte(v.(string))
	}

	if v, ok := d.GetOk("kerberos_principal"); ok {
		input.KerberosPrincipal = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_provider_uri"); ok {
		input.KmsKeyProviderUri = aws.String(v.(string))
	}

	if v, ok := d.GetOk("block_size"); ok {
		input.BlockSize = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("replication_factor"); ok {
		input.ReplicationFactor = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("qop_configuration"); ok && len(v.([]interface{})) > 0 {
		input.QopConfiguration = expandHDFSQOPConfiguration(v.([]interface{}))
	}

	log.Printf("[DEBUG] Creating DataSync Location HDFS: %s", input)
	output, err := conn.CreateLocationHdfsWithContext(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating DataSync Location HDFS: %s", err)
	}

	d.SetId(aws.StringValue(output.LocationArn))

	return append(diags, resourceLocationHDFSRead(ctx, d, meta)...)
}

func resourceLocationHDFSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	output, err := FindLocationHDFSByARN(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] DataSync Location HDFS (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading DataSync Location HDFS (%s): %s", d.Id(), err)
	}

	subdirectory, err := SubdirectoryFromLocationURI(aws.StringValue(output.LocationUri))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading DataSync Location HDFS (%s): %s", d.Id(), err)
	}

	d.Set("agent_arns", flex.FlattenStringSet(output.AgentArns))
	d.Set("arn", output.LocationArn)
	d.Set("simple_user", output.SimpleUser)
	d.Set("authentication_type", output.AuthenticationType)
	d.Set("uri", output.LocationUri)
	d.Set("block_size", output.BlockSize)
	d.Set("replication_factor", output.ReplicationFactor)
	d.Set("kerberos_principal", output.KerberosPrincipal)
	d.Set("kms_key_provider_uri", output.KmsKeyProviderUri)
	d.Set("subdirectory", subdirectory)

	if err := d.Set("name_node", flattenHDFSNameNodes(output.NameNodes)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting name_node: %s", err)
	}

	if err := d.Set("qop_configuration", flattenHDFSQOPConfiguration(output.QopConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting qop_configuration: %s", err)
	}

	return diags
}

func resourceLocationHDFSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	if d.HasChangesExcept("tags_all", "tags") {
		input := &datasync.UpdateLocationHdfsInput{
			LocationArn: aws.String(d.Id()),
		}

		if d.HasChange("authentication_type") {
			input.AuthenticationType = aws.String(d.Get("authentication_type").(string))
		}

		if d.HasChange("subdirectory") {
			input.Subdirectory = aws.String(d.Get("subdirectory").(string))
		}

		if d.HasChange("simple_user") {
			input.SimpleUser = aws.String(d.Get("simple_user").(string))
		}

		if d.HasChange("kerberos_keytab") {
			input.KerberosKeytab = []byte(d.Get("kerberos_keytab").(string))
		}

		if d.HasChange("kerberos_krb5_conf") {
			input.KerberosKrb5Conf = []byte(d.Get("kerberos_krb5_conf").(string))
		}

		if d.HasChange("kerberos_principal") {
			input.KerberosPrincipal = aws.String(d.Get("kerberos_principal").(string))
		}

		if d.HasChange("kms_key_provider_uri") {
			input.KmsKeyProviderUri = aws.String(d.Get("kms_key_provider_uri").(string))
		}

		if d.HasChange("block_size") {
			input.BlockSize = aws.Int64(int64(d.Get("block_size").(int)))
		}

		if d.HasChange("replication_factor") {
			input.ReplicationFactor = aws.Int64(int64(d.Get("replication_factor").(int)))
		}

		if d.HasChange("agent_arns") {
			input.AgentArns = flex.ExpandStringSet(d.Get("agent_arns").(*schema.Set))
		}

		if d.HasChange("name_node") {
			input.NameNodes = expandHDFSNameNodes(d.Get("name_node").(*schema.Set))
		}

		if d.HasChange("qop_configuration") {
			input.QopConfiguration = expandHDFSQOPConfiguration(d.Get("qop_configuration").([]interface{}))
		}

		_, err := conn.UpdateLocationHdfsWithContext(ctx, input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating DataSync Location HDFS (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceLocationHDFSRead(ctx, d, meta)...)
}

func resourceLocationHDFSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	input := &datasync.DeleteLocationInput{
		LocationArn: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Deleting DataSync Location HDFS: %s", input)
	_, err := conn.DeleteLocationWithContext(ctx, input)

	if tfawserr.ErrMessageContains(err, datasync.ErrCodeInvalidRequestException, "not found") {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting DataSync Location HDFS (%s): %s", d.Id(), err)
	}

	return diags
}

func expandHDFSNameNodes(l *schema.Set) []*datasync.HdfsNameNode {
	nameNodes := make([]*datasync.HdfsNameNode, 0)
	for _, m := range l.List() {
		raw := m.(map[string]interface{})
		nameNode := &datasync.HdfsNameNode{
			Hostname: aws.String(raw["hostname"].(string)),
			Port:     aws.Int64(int64(raw["port"].(int))),
		}
		nameNodes = append(nameNodes, nameNode)
	}

	return nameNodes
}

func flattenHDFSNameNodes(nodes []*datasync.HdfsNameNode) []map[string]interface{} {
	dataResources := make([]map[string]interface{}, 0, len(nodes))

	for _, raw := range nodes {
		item := make(map[string]interface{})
		item["hostname"] = aws.StringValue(raw.Hostname)
		item["port"] = aws.Int64Value(raw.Port)

		dataResources = append(dataResources, item)
	}

	return dataResources
}

func expandHDFSQOPConfiguration(l []interface{}) *datasync.QopConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	qopConfig := &datasync.QopConfiguration{
		DataTransferProtection: aws.String(m["data_transfer_protection"].(string)),
		RpcProtection:          aws.String(m["rpc_protection"].(string)),
	}

	return qopConfig
}

func flattenHDFSQOPConfiguration(qopConfig *datasync.QopConfiguration) []interface{} {
	if qopConfig == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"data_transfer_protection": aws.StringValue(qopConfig.DataTransferProtection),
		"rpc_protection":           aws.StringValue(qopConfig.RpcProtection),
	}

	return []interface{}{m}
}
