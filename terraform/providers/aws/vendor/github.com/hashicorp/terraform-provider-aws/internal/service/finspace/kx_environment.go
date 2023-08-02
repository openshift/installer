package finspace

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/finspace"
	"github.com/aws/aws-sdk-go-v2/service/finspace/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_finspace_kx_environment", name="Kx Environment")
// @Tags(identifierAttribute="arn")
func ResourceKxEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceKxEnvironmentCreate,
		ReadWithoutTimeout:   resourceKxEnvironmentRead,
		UpdateWithoutTimeout: resourceKxEnvironmentUpdate,
		DeleteWithoutTimeout: resourceKxEnvironmentDelete,

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
			"availability_zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"created_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_dns_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_dns_server_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(3, 255),
						},
						"custom_dns_server_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
					},
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"infrastructure_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"last_modified_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"transit_gateway_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transit_gateway_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 32),
						},
						"routable_cidr_space": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
					},
				},
			},
		},
		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameKxEnvironment = "Kx Environment"
)

func resourceKxEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FinSpaceClient(ctx)

	in := &finspace.CreateKxEnvironmentInput{
		Name:        aws.String(d.Get("name").(string)),
		ClientToken: aws.String(id.UniqueId()),
	}

	if v, ok := d.GetOk("description"); ok {
		in.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		in.KmsKeyId = aws.String(v.(string))
	}

	out, err := conn.CreateKxEnvironment(ctx, in)
	if err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionCreating, ResNameKxEnvironment, d.Get("name").(string), err)...)
	}

	if out == nil || out.EnvironmentId == nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionCreating, ResNameKxEnvironment, d.Get("name").(string), errors.New("empty output"))...)
	}

	d.SetId(aws.ToString(out.EnvironmentId))

	if _, err := waitKxEnvironmentCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionWaitingForCreation, ResNameKxEnvironment, d.Id(), err)...)
	}

	if err := updateKxEnvironmentNetwork(ctx, d, conn); err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionCreating, ResNameKxEnvironment, d.Id(), err)...)
	}

	// The CreateKxEnvironment API currently fails to tag the environment when the
	// Tags field is set. Until the API is fixed, tag after creation instead.
	if err := createTags(ctx, conn, aws.ToString(out.EnvironmentArn), GetTagsIn(ctx)); err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionCreating, ResNameKxEnvironment, d.Id(), err)...)
	}

	return append(diags, resourceKxEnvironmentRead(ctx, d, meta)...)
}

func resourceKxEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FinSpaceClient(ctx)

	out, err := findKxEnvironmentByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] FinSpace KxEnvironment (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionReading, ResNameKxEnvironment, d.Id(), err)...)
	}

	d.Set("id", out.EnvironmentId)
	d.Set("arn", out.EnvironmentArn)
	d.Set("name", out.Name)
	d.Set("description", out.Description)
	d.Set("kms_key_id", out.KmsKeyId)
	d.Set("status", out.Status)
	d.Set("availability_zones", out.AvailabilityZoneIds)
	d.Set("infrastructure_account_id", out.DedicatedServiceAccountId)
	d.Set("created_timestamp", out.CreationTimestamp.String())
	d.Set("last_modified_timestamp", out.UpdateTimestamp.String())

	if err := d.Set("transit_gateway_configuration", flattenTransitGatewayConfiguration(out.TransitGatewayConfiguration)); err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionSetting, ResNameKxEnvironment, d.Id(), err)...)
	}

	if err := d.Set("custom_dns_configuration", flattenCustomDNSConfigurations(out.CustomDNSConfiguration)); err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionSetting, ResNameKxEnvironment, d.Id(), err)...)
	}

	return diags
}

func resourceKxEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FinSpaceClient(ctx)

	update := false

	in := &finspace.UpdateKxEnvironmentInput{
		EnvironmentId: aws.String(d.Id()),
		Name:          aws.String(d.Get("name").(string)),
	}

	if d.HasChanges("description") {
		in.Description = aws.String(d.Get("description").(string))
	}

	if d.HasChanges("name") || d.HasChanges("description") {
		update = true
		log.Printf("[DEBUG] Updating FinSpace KxEnvironment (%s): %#v", d.Id(), in)
		_, err := conn.UpdateKxEnvironment(ctx, in)
		if err != nil {
			return append(diags, create.DiagError(names.FinSpace, create.ErrActionUpdating, ResNameKxEnvironment, d.Id(), err)...)
		}
	}

	if d.HasChanges("transit_gateway_configuration") || d.HasChanges("custom_dns_configuration") {
		update = true
		if err := updateKxEnvironmentNetwork(ctx, d, conn); err != nil {
			return append(diags, create.DiagError(names.FinSpace, create.ErrActionUpdating, ResNameKxEnvironment, d.Id(), err)...)
		}
	}

	if !update {
		return diags
	}
	return append(diags, resourceKxEnvironmentRead(ctx, d, meta)...)
}

func resourceKxEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FinSpaceClient(ctx)

	log.Printf("[INFO] Deleting FinSpace KxEnvironment %s", d.Id())

	_, err := conn.DeleteKxEnvironment(ctx, &finspace.DeleteKxEnvironmentInput{
		EnvironmentId: aws.String(d.Id()),
	})
	if errs.IsA[*types.ResourceNotFoundException](err) ||
		errs.IsAErrorMessageContains[*types.ValidationException](err, "The Environment is in DELETED state") {
		log.Printf("[DEBUG] FinSpace KxEnvironment %s already deleted. Nothing to delete.", d.Id())
		return diags
	}

	if err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionDeleting, ResNameKxEnvironment, d.Id(), err)...)
	}

	if _, err := waitKxEnvironmentDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return append(diags, create.DiagError(names.FinSpace, create.ErrActionWaitingForDeletion, ResNameKxEnvironment, d.Id(), err)...)
	}

	return diags
}

// As of 2023-02-09, updating network configuration requires 2 separate requests if both DNS
// and transit gateway configurationtions are set.
func updateKxEnvironmentNetwork(ctx context.Context, d *schema.ResourceData, client *finspace.Client) error {
	transitGatewayConfigIn := &finspace.UpdateKxEnvironmentNetworkInput{
		EnvironmentId: aws.String(d.Id()),
		ClientToken:   aws.String(id.UniqueId()),
	}

	customDnsConfigIn := &finspace.UpdateKxEnvironmentNetworkInput{
		EnvironmentId: aws.String(d.Id()),
		ClientToken:   aws.String(id.UniqueId()),
	}

	updateTransitGatewayConfig := false
	updateCustomDnsConfig := false

	if v, ok := d.GetOk("transit_gateway_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil &&
		d.HasChanges("transit_gateway_configuration") {
		transitGatewayConfigIn.TransitGatewayConfiguration = expandTransitGatewayConfiguration(v.([]interface{}))
		updateTransitGatewayConfig = true
	}

	if v, ok := d.GetOk("custom_dns_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil &&
		d.HasChanges("custom_dns_configuration") {
		customDnsConfigIn.CustomDNSConfiguration = expandCustomDNSConfigurations(v.([]interface{}))
		updateCustomDnsConfig = true
	}

	if updateTransitGatewayConfig {
		if _, err := client.UpdateKxEnvironmentNetwork(ctx, transitGatewayConfigIn); err != nil {
			return err
		}

		if _, err := waitTransitGatewayConfigurationUpdated(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	if updateCustomDnsConfig {
		if _, err := client.UpdateKxEnvironmentNetwork(ctx, customDnsConfigIn); err != nil {
			return err
		}

		if _, err := waitCustomDNSConfigurationUpdated(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	return nil
}

func waitKxEnvironmentCreated(ctx context.Context, conn *finspace.Client, id string, timeout time.Duration) (*finspace.GetKxEnvironmentOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   enum.Slice(types.EnvironmentStatusCreateRequested, types.EnvironmentStatusCreating),
		Target:                    enum.Slice(types.EnvironmentStatusCreated),
		Refresh:                   statusKxEnvironment(ctx, conn, id),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*finspace.GetKxEnvironmentOutput); ok {
		return out, err
	}

	return nil, err
}

func waitTransitGatewayConfigurationUpdated(ctx context.Context, conn *finspace.Client, id string, timeout time.Duration) (*finspace.GetKxEnvironmentOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(types.TgwStatusUpdateRequested, types.TgwStatusUpdating),
		Target:  enum.Slice(types.TgwStatusSuccessfullyUpdated),
		Refresh: statusTransitGatewayConfiguration(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*finspace.GetKxEnvironmentOutput); ok {
		return out, err
	}

	return nil, err
}

func waitCustomDNSConfigurationUpdated(ctx context.Context, conn *finspace.Client, id string, timeout time.Duration) (*finspace.GetKxEnvironmentOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(types.DnsStatusUpdateRequested, types.DnsStatusUpdating),
		Target:  enum.Slice(types.DnsStatusSuccessfullyUpdated),
		Refresh: statusCustomDNSConfiguration(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*finspace.GetKxEnvironmentOutput); ok {
		return out, err
	}

	return nil, err
}

func waitKxEnvironmentDeleted(ctx context.Context, conn *finspace.Client, id string, timeout time.Duration) (*finspace.GetKxEnvironmentOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(types.EnvironmentStatusDeleteRequested, types.EnvironmentStatusDeleting),
		Target:  []string{},
		Refresh: statusKxEnvironment(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*finspace.GetKxEnvironmentOutput); ok {
		return out, err
	}

	return nil, err
}

func statusKxEnvironment(ctx context.Context, conn *finspace.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		out, err := findKxEnvironmentByID(ctx, conn, id)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, string(out.Status), nil
	}
}

func statusTransitGatewayConfiguration(ctx context.Context, conn *finspace.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		out, err := findKxEnvironmentByID(ctx, conn, id)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, string(out.TgwStatus), nil
	}
}

func statusCustomDNSConfiguration(ctx context.Context, conn *finspace.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		out, err := findKxEnvironmentByID(ctx, conn, id)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, string(out.DnsStatus), nil
	}
}

func findKxEnvironmentByID(ctx context.Context, conn *finspace.Client, id string) (*finspace.GetKxEnvironmentOutput, error) {
	in := &finspace.GetKxEnvironmentInput{
		EnvironmentId: aws.String(id),
	}
	out, err := conn.GetKxEnvironment(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}
	// Treat DELETED status as NotFound
	if out != nil && out.Status == types.EnvironmentStatusDeleted {
		return nil, &retry.NotFoundError{
			LastError:   errors.New("status is deleted"),
			LastRequest: in,
		}
	}

	if out == nil || out.EnvironmentArn == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func expandTransitGatewayConfiguration(tfList []interface{}) *types.TransitGatewayConfiguration {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})

	a := &types.TransitGatewayConfiguration{}

	if v, ok := tfMap["transit_gateway_id"].(string); ok && v != "" {
		a.TransitGatewayID = aws.String(v)
	}

	if v, ok := tfMap["routable_cidr_space"].(string); ok && v != "" {
		a.RoutableCIDRSpace = aws.String(v)
	}

	return a
}

func expandCustomDNSConfiguration(tfMap map[string]interface{}) *types.CustomDNSServer {
	if tfMap == nil {
		return nil
	}

	a := &types.CustomDNSServer{}

	if v, ok := tfMap["custom_dns_server_name"].(string); ok && v != "" {
		a.CustomDNSServerName = aws.String(v)
	}

	if v, ok := tfMap["custom_dns_server_ip"].(string); ok && v != "" {
		a.CustomDNSServerIP = aws.String(v)
	}

	return a
}

func expandCustomDNSConfigurations(tfList []interface{}) []types.CustomDNSServer {
	if len(tfList) == 0 {
		return nil
	}

	var s []types.CustomDNSServer

	for _, r := range tfList {
		m, ok := r.(map[string]interface{})

		if !ok {
			continue
		}

		a := expandCustomDNSConfiguration(m)

		if a == nil {
			continue
		}

		s = append(s, *a)
	}

	return s
}

func flattenTransitGatewayConfiguration(apiObject *types.TransitGatewayConfiguration) []interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{}

	if v := apiObject.TransitGatewayID; v != nil {
		m["transit_gateway_id"] = aws.ToString(v)
	}

	if v := apiObject.RoutableCIDRSpace; v != nil {
		m["routable_cidr_space"] = aws.ToString(v)
	}

	return []interface{}{m}
}

func flattenCustomDNSConfiguration(apiObject *types.CustomDNSServer) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{}

	if v := apiObject.CustomDNSServerName; v != nil {
		m["custom_dns_server_name"] = aws.ToString(v)
	}

	if v := apiObject.CustomDNSServerIP; v != nil {
		m["custom_dns_server_ip"] = aws.ToString(v)
	}

	return m
}

func flattenCustomDNSConfigurations(apiObjects []types.CustomDNSServer) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var l []interface{}

	for _, apiObject := range apiObjects {
		l = append(l, flattenCustomDNSConfiguration(&apiObject))
	}

	return l
}
