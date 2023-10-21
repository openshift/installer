package ssm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_ssm_maintenance_window_target")
func ResourceMaintenanceWindowTarget() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMaintenanceWindowTargetCreate,
		ReadWithoutTimeout:   resourceMaintenanceWindowTargetRead,
		UpdateWithoutTimeout: resourceMaintenanceWindowTargetUpdate,
		DeleteWithoutTimeout: resourceMaintenanceWindowTargetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected WINDOW_ID/WINDOW_TARGET_ID", d.Id())
				}
				d.Set("window_id", idParts[0])
				d.SetId(idParts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"window_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(ssm.MaintenanceWindowResourceType_Values(), true),
			},

			"targets": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile(`^[\p{L}\p{Z}\p{N}_.:/=\-@]*$|resource-groups:ResourceTypeFilters|resource-groups:Name$`), ""),
								validation.StringLenBetween(1, 163),
							),
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 50,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_\-.]{3,128}$`), "Only alphanumeric characters, hyphens, dots & underscores allowed"),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"owner_information": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
		},
	}
}

func resourceMaintenanceWindowTargetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMConn(ctx)

	log.Printf("[INFO] Registering SSM Maintenance Window Target")

	params := &ssm.RegisterTargetWithMaintenanceWindowInput{
		WindowId:     aws.String(d.Get("window_id").(string)),
		ResourceType: aws.String(d.Get("resource_type").(string)),
		Targets:      expandTargets(d.Get("targets").([]interface{})),
	}

	if v, ok := d.GetOk("name"); ok {
		params.Name = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		params.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("owner_information"); ok {
		params.OwnerInformation = aws.String(v.(string))
	}

	resp, err := conn.RegisterTargetWithMaintenanceWindowWithContext(ctx, params)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating SSM Maintenance Window Target: %s", err)
	}

	d.SetId(aws.StringValue(resp.WindowTargetId))

	return append(diags, resourceMaintenanceWindowTargetRead(ctx, d, meta)...)
}

func resourceMaintenanceWindowTargetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMConn(ctx)

	windowID := d.Get("window_id").(string)
	params := &ssm.DescribeMaintenanceWindowTargetsInput{
		WindowId: aws.String(windowID),
		Filters: []*ssm.MaintenanceWindowFilter{
			{
				Key:    aws.String("WindowTargetId"),
				Values: []*string{aws.String(d.Id())},
			},
		},
	}

	resp, err := conn.DescribeMaintenanceWindowTargetsWithContext(ctx, params)
	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, ssm.ErrCodeDoesNotExistException) {
		log.Printf("[WARN] Maintenance Window (%s) Target (%s) not found, removing from state", windowID, d.Id())
		d.SetId("")
		return diags
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "getting Maintenance Window (%s) Target (%s): %s", windowID, d.Id(), err)
	}

	found := false
	for _, t := range resp.Targets {
		if aws.StringValue(t.WindowTargetId) == d.Id() {
			found = true

			d.Set("owner_information", t.OwnerInformation)
			d.Set("window_id", t.WindowId)
			d.Set("resource_type", t.ResourceType)
			d.Set("name", t.Name)
			d.Set("description", t.Description)

			if err := d.Set("targets", flattenTargets(t.Targets)); err != nil {
				return sdkdiag.AppendErrorf(diags, "setting targets: %s", err)
			}

			break
		}
	}

	if !d.IsNewResource() && !found {
		log.Printf("[INFO] Maintenance Window Target not found. Removing from state")
		d.SetId("")
		return diags
	}

	return diags
}

func resourceMaintenanceWindowTargetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMConn(ctx)

	log.Printf("[INFO] Updating SSM Maintenance Window Target: %s", d.Id())

	params := &ssm.UpdateMaintenanceWindowTargetInput{
		Targets:        expandTargets(d.Get("targets").([]interface{})),
		WindowId:       aws.String(d.Get("window_id").(string)),
		WindowTargetId: aws.String(d.Id()),
	}

	if d.HasChange("name") {
		params.Name = aws.String(d.Get("name").(string))
	}

	if d.HasChange("description") {
		params.Description = aws.String(d.Get("description").(string))
	}

	if d.HasChange("owner_information") {
		params.OwnerInformation = aws.String(d.Get("owner_information").(string))
	}

	_, err := conn.UpdateMaintenanceWindowTargetWithContext(ctx, params)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating SSM Maintenance Window Target (%s): %s", d.Id(), err)
	}

	return diags
}

func resourceMaintenanceWindowTargetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMConn(ctx)

	log.Printf("[INFO] Deregistering SSM Maintenance Window Target: %s", d.Id())

	params := &ssm.DeregisterTargetFromMaintenanceWindowInput{
		WindowId:       aws.String(d.Get("window_id").(string)),
		WindowTargetId: aws.String(d.Id()),
	}

	_, err := conn.DeregisterTargetFromMaintenanceWindowWithContext(ctx, params)
	if tfawserr.ErrCodeEquals(err, ssm.ErrCodeDoesNotExistException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deregistering SSM Maintenance Window Target (%s): %s", d.Id(), err)
	}

	return diags
}
