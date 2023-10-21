package appstream

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	stackOperationTimeout = 4 * time.Minute
)

// @SDKResource("aws_appstream_stack", name="Stack")
// @Tags(identifierAttribute="arn")
func ResourceStack() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceStackCreate,
		ReadWithoutTimeout:   resourceStackRead,
		UpdateWithoutTimeout: resourceStackUpdate,
		DeleteWithoutTimeout: resourceStackDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"access_endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				MaxItems: 4,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(appstream.AccessEndpointType_Values(), false),
						},
						"vpce_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"application_settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"settings_group": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 100),
						},
					},
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
			},
			"embed_host_domains": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				MaxItems: 20,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(0, 128),
				},
				Set: schema.HashString,
			},
			"feedback_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"redirect_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"storage_connectors": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connector_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(appstream.StorageConnectorType_Values(), false),
						},
						"domains": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 50,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringLenBetween(1, 64),
							},
						},
						"resource_identifier": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 2048),
						},
					},
				},
			},
			"streaming_experience_settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preferred_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(appstream.PreferredProtocol_Values(), false),
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"user_settings": {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				MinItems:         1,
				DiffSuppressFunc: suppressAppsStreamStackUserSettings,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(appstream.Action_Values(), false),
						},
						"permission": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(appstream.Permission_Values(), false),
						},
					},
				},
			},
		},

		CustomizeDiff: customdiff.All(
			verify.SetTagsDiff,
			func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
				if d.Id() == "" {
					return nil
				}

				rawConfig := d.GetRawConfig()
				configApplicationSettings := rawConfig.GetAttr("application_settings")
				if configApplicationSettings.IsKnown() && !configApplicationSettings.IsNull() && configApplicationSettings.LengthInt() > 0 {
					return nil
				}

				rawState := d.GetRawState()
				stateApplicationSettings := rawState.GetAttr("application_settings")
				if stateApplicationSettings.IsKnown() && !stateApplicationSettings.IsNull() && stateApplicationSettings.LengthInt() > 0 {
					setting := stateApplicationSettings.Index(cty.NumberIntVal(0))
					if setting.IsKnown() && !setting.IsNull() {
						enabled := setting.GetAttr("enabled")
						if enabled.IsKnown() && !enabled.IsNull() && enabled.True() {
							// Trigger a diff
							return d.SetNew("application_settings", []map[string]any{
								{
									"enabled":        false,
									"settings_group": "",
								},
							})
						}
					}
				}

				return nil
			},
			func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
				_, enabled := d.GetOk("application_settings.0.enabled")
				v, sg := d.GetOk("application_settings.0.settings_group")
				log.Print(v)

				if enabled && !sg {
					return errors.New("application_settings.settings_group must be set when application_settings.enabled is true")
				} else if !enabled && sg {
					return errors.New("application_settings.settings_group must not be set when application_settings.enabled is false")
				}
				return nil
			},
		),
	}
}

func resourceStackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AppStreamConn(ctx)

	name := d.Get("name").(string)
	input := &appstream.CreateStackInput{
		Name: aws.String(name),
		Tags: GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("access_endpoints"); ok {
		input.AccessEndpoints = expandAccessEndpoints(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("application_settings"); ok {
		input.ApplicationSettings = expandApplicationSettings(v.([]interface{}))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("display_name"); ok {
		input.DisplayName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("embed_host_domains"); ok {
		input.EmbedHostDomains = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("feedback_url"); ok {
		input.FeedbackURL = aws.String(v.(string))
	}

	if v, ok := d.GetOk("redirect_url"); ok {
		input.RedirectURL = aws.String(v.(string))
	}

	if v, ok := d.GetOk("storage_connectors"); ok {
		input.StorageConnectors = expandStorageConnectors(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("streaming_experience_settings"); ok {
		input.StreamingExperienceSettings = expandStreamingExperienceSettings(v.([]interface{}))
	}

	if v, ok := d.GetOk("user_settings"); ok {
		input.UserSettings = expandUserSettings(v.(*schema.Set).List())
	}

	outputRaw, err := tfresource.RetryWhenAWSErrCodeEquals(ctx, stackOperationTimeout, func() (interface{}, error) {
		return conn.CreateStackWithContext(ctx, input)
	}, appstream.ErrCodeResourceNotFoundException)

	if err != nil {
		return diag.Errorf("creating Appstream Stack (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(outputRaw.(*appstream.CreateStackOutput).Stack.Name))

	return resourceStackRead(ctx, d, meta)
}

func resourceStackRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AppStreamConn(ctx)

	stack, err := FindStackByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Appstream Stack (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Appstream Stack (%s): %s", d.Id(), err)
	}

	if err = d.Set("access_endpoints", flattenAccessEndpoints(stack.AccessEndpoints)); err != nil {
		return diag.Errorf("setting access_endpoints: %s", err)
	}
	if err = d.Set("application_settings", flattenApplicationSettings(stack.ApplicationSettings)); err != nil {
		return diag.Errorf("setting application_settings: %s", err)
	}
	d.Set("arn", stack.Arn)
	d.Set("created_time", aws.TimeValue(stack.CreatedTime).Format(time.RFC3339))
	d.Set("description", stack.Description)
	d.Set("display_name", stack.DisplayName)
	if err = d.Set("embed_host_domains", flex.FlattenStringList(stack.EmbedHostDomains)); err != nil {
		return diag.Errorf("setting embed_host_domains: %s", err)
	}
	d.Set("feedback_url", stack.FeedbackURL)
	d.Set("name", stack.Name)
	d.Set("redirect_url", stack.RedirectURL)
	if err = d.Set("storage_connectors", flattenStorageConnectors(stack.StorageConnectors)); err != nil {
		return diag.Errorf("setting storage_connectors: %s", err)
	}
	if err = d.Set("streaming_experience_settings", flattenStreaminExperienceSettings(stack.StreamingExperienceSettings)); err != nil {
		return diag.Errorf("setting streaming_experience_settings: %s", err)
	}
	if err = d.Set("user_settings", flattenUserSettings(stack.UserSettings)); err != nil {
		return diag.Errorf("setting user_settings: %s", err)
	}

	return nil
}

func resourceStackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AppStreamConn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &appstream.UpdateStackInput{
			Name: aws.String(d.Id()),
		}

		if d.HasChange("access_endpoints") {
			input.AccessEndpoints = expandAccessEndpoints(d.Get("access_endpoints").(*schema.Set).List())
		}

		if d.HasChange("application_settings") {
			input.ApplicationSettings = expandApplicationSettings(d.Get("application_settings").([]interface{}))
		}

		if d.HasChange("description") {
			input.Description = aws.String(d.Get("description").(string))
		}

		if d.HasChange("display_name") {
			input.DisplayName = aws.String(d.Get("display_name").(string))
		}

		if d.HasChange("feedback_url") {
			input.FeedbackURL = aws.String(d.Get("feedback_url").(string))
		}

		if d.HasChange("redirect_url") {
			input.RedirectURL = aws.String(d.Get("redirect_url").(string))
		}

		if d.HasChange("streaming_experience_settings") {
			input.StreamingExperienceSettings = expandStreamingExperienceSettings(d.Get("streaming_experience_settings").([]interface{}))
		}

		if d.HasChange("user_settings") {
			input.UserSettings = expandUserSettings(d.Get("user_settings").(*schema.Set).List())
		}

		_, err := conn.UpdateStackWithContext(ctx, input)

		if err != nil {
			diag.Errorf("updating Appstream Stack (%s): %s", d.Id(), err)
		}
	}

	return resourceStackRead(ctx, d, meta)
}

func resourceStackDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AppStreamConn(ctx)

	log.Printf("[DEBUG] Deleting AppStream Stack: (%s)", d.Id())
	_, err := conn.DeleteStackWithContext(ctx, &appstream.DeleteStackInput{
		Name: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, appstream.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting Appstream Stack (%s): %s", d.Id(), err)
	}

	_, err = tfresource.RetryUntilNotFound(ctx, stackOperationTimeout, func() (interface{}, error) {
		return FindStackByName(ctx, conn, d.Id())
	})

	if err != nil {
		return diag.Errorf("waiting for Appstream Stack (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func FindStackByName(ctx context.Context, conn *appstream.AppStream, name string) (*appstream.Stack, error) {
	input := &appstream.DescribeStacksInput{
		Names: aws.StringSlice([]string{name}),
	}

	output, err := conn.DescribeStacksWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, appstream.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.Stacks) == 0 || output.Stacks[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.Stacks); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	return output.Stacks[0], nil
}

func expandAccessEndpoint(tfMap map[string]interface{}) *appstream.AccessEndpoint {
	if tfMap == nil {
		return nil
	}

	apiObject := &appstream.AccessEndpoint{
		EndpointType: aws.String(tfMap["endpoint_type"].(string)),
	}
	if v, ok := tfMap["vpce_id"]; ok {
		apiObject.VpceId = aws.String(v.(string))
	}

	return apiObject
}

func expandAccessEndpoints(tfList []interface{}) []*appstream.AccessEndpoint {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*appstream.AccessEndpoint

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandAccessEndpoint(tfMap)

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenAccessEndpoint(apiObject *appstream.AccessEndpoint) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}
	tfMap["endpoint_type"] = aws.StringValue(apiObject.EndpointType)
	tfMap["vpce_id"] = aws.StringValue(apiObject.VpceId)

	return tfMap
}

func flattenAccessEndpoints(apiObjects []*appstream.AccessEndpoint) []map[string]interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []map[string]interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenAccessEndpoint(apiObject))
	}

	return tfList
}

func expandApplicationSettings(tfList []interface{}) *appstream.ApplicationSettings {
	if len(tfList) == 0 {
		return &appstream.ApplicationSettings{
			Enabled: aws.Bool(false),
		}
	}

	tfMap := tfList[0].(map[string]interface{})

	apiObject := &appstream.ApplicationSettings{
		Enabled: aws.Bool(tfMap["enabled"].(bool)),
	}
	if v, ok := tfMap["settings_group"]; ok {
		apiObject.SettingsGroup = aws.String(v.(string))
	}

	return apiObject
}

func flattenApplicationSetting(apiObject *appstream.ApplicationSettingsResponse) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	return map[string]interface{}{
		"enabled":        aws.BoolValue(apiObject.Enabled),
		"settings_group": aws.StringValue(apiObject.SettingsGroup),
	}
}

func flattenApplicationSettings(apiObject *appstream.ApplicationSettingsResponse) []interface{} {
	if apiObject == nil {
		return nil
	}

	var tfList []interface{}

	tfList = append(tfList, flattenApplicationSetting(apiObject))

	return tfList
}

func expandStreamingExperienceSettings(tfList []interface{}) *appstream.StreamingExperienceSettings {
	if len(tfList) == 0 {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})

	apiObject := &appstream.StreamingExperienceSettings{
		PreferredProtocol: aws.String(tfMap["preferred_protocol"].(string)),
	}

	return apiObject
}

func flattenStreaminExperienceSetting(apiObject *appstream.StreamingExperienceSettings) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	return map[string]interface{}{
		"preferred_protocol": aws.StringValue(apiObject.PreferredProtocol),
	}
}

func flattenStreaminExperienceSettings(apiObject *appstream.StreamingExperienceSettings) []interface{} {
	if apiObject == nil {
		return nil
	}

	var tfList []interface{}

	tfList = append(tfList, flattenStreaminExperienceSetting(apiObject))

	return tfList
}

func expandStorageConnector(tfMap map[string]interface{}) *appstream.StorageConnector {
	if tfMap == nil {
		return nil
	}

	apiObject := &appstream.StorageConnector{
		ConnectorType: aws.String(tfMap["connector_type"].(string)),
	}
	if v, ok := tfMap["domains"]; ok && len(v.([]interface{})) > 0 {
		apiObject.Domains = flex.ExpandStringList(v.([]interface{}))
	}
	if v, ok := tfMap["resource_identifier"]; ok && v.(string) != "" {
		apiObject.ResourceIdentifier = aws.String(v.(string))
	}

	return apiObject
}

func expandStorageConnectors(tfList []interface{}) []*appstream.StorageConnector {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*appstream.StorageConnector

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandStorageConnector(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenStorageConnector(apiObject *appstream.StorageConnector) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}
	tfMap["connector_type"] = aws.StringValue(apiObject.ConnectorType)
	tfMap["domains"] = aws.StringValueSlice(apiObject.Domains)
	tfMap["resource_identifier"] = aws.StringValue(apiObject.ResourceIdentifier)

	return tfMap
}

func flattenStorageConnectors(apiObjects []*appstream.StorageConnector) []map[string]interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []map[string]interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenStorageConnector(apiObject))
	}

	return tfList
}

func expandUserSetting(tfMap map[string]interface{}) *appstream.UserSetting {
	if tfMap == nil {
		return nil
	}

	apiObject := &appstream.UserSetting{
		Action:     aws.String(tfMap["action"].(string)),
		Permission: aws.String(tfMap["permission"].(string)),
	}

	return apiObject
}

func expandUserSettings(tfList []interface{}) []*appstream.UserSetting {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*appstream.UserSetting

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandUserSetting(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenUserSetting(apiObject *appstream.UserSetting) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}
	tfMap["action"] = aws.StringValue(apiObject.Action)
	tfMap["permission"] = aws.StringValue(apiObject.Permission)

	return tfMap
}

func flattenUserSettings(apiObjects []*appstream.UserSetting) []map[string]interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []map[string]interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}
		tfList = append(tfList, flattenUserSetting(apiObject))
	}

	return tfList
}

func suppressAppsStreamStackUserSettings(k, old, new string, d *schema.ResourceData) bool {
	flagDiffUserSettings := false
	count := len(d.Get("user_settings").(*schema.Set).List())
	defaultCount := len(appstream.Action_Values())

	if count == defaultCount {
		flagDiffUserSettings = false
	}

	if count != defaultCount && (fmt.Sprintf("%d", count) == new && fmt.Sprintf("%d", defaultCount) == old) {
		flagDiffUserSettings = true
	}

	return flagDiffUserSettings
}
