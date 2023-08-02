package lightsail

import (
	"context"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_lightsail_container_service", name="Container Service")
// @Tags(identifierAttribute="id")
func ResourceContainerService() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceContainerServiceCreate,
		ReadWithoutTimeout:   resourceContainerServiceRead,
		UpdateWithoutTimeout: resourceContainerServiceUpdate,
		DeleteWithoutTimeout: resourceContainerServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 63),
					validation.StringMatch(regexp.MustCompile(`^[a-z0-9]{1,2}|[a-z0-9][a-z0-9-]+[a-z0-9]$`), ""),
				),
			},
			"power": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(flattenContainerServicePowerValues(types.ContainerServicePowerName("").Values()), false),
			},
			"power_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"principal_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_domain_names": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"domain_names": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"private_registry_access": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ecr_image_puller_role": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_active": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"principal_arn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scale": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 20),
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceContainerServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailClient(ctx)

	serviceName := d.Get("name").(string)
	input := &lightsail.CreateContainerServiceInput{
		ServiceName: aws.String(serviceName),
		Power:       types.ContainerServicePowerName(d.Get("power").(string)),
		Scale:       aws.Int32(int32(d.Get("scale").(int))),
		Tags:        GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("public_domain_names"); ok {
		input.PublicDomainNames = expandContainerServicePublicDomainNames(v.([]interface{}))
	}

	if v, ok := d.GetOk("private_registry_access"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.PrivateRegistryAccess = expandPrivateRegistryAccess(v.([]interface{})[0].(map[string]interface{}))
	}

	_, err := conn.CreateContainerService(ctx, input)
	if err != nil {
		return diag.Errorf("creating Lightsail Container Service (%s): %s", serviceName, err)
	}

	d.SetId(serviceName)

	if err := waitContainerServiceCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("waiting for Lightsail Container Service (%s) creation: %s", d.Id(), err)
	}

	// once container service creation and/or deployment successful (now enabled by default), disable it if "is_disabled" is true
	if v, ok := d.GetOk("is_disabled"); ok && v.(bool) {
		input := &lightsail.UpdateContainerServiceInput{
			ServiceName: aws.String(d.Id()),
			IsDisabled:  aws.Bool(true),
		}

		_, err := conn.UpdateContainerService(ctx, input)
		if err != nil {
			return diag.Errorf("disabling Lightsail Container Service (%s): %s", d.Id(), err)
		}

		if err := waitContainerServiceDisabled(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("waiting for Lightsail Container Service (%s) to be disabled: %s", d.Id(), err)
		}
	}

	return resourceContainerServiceRead(ctx, d, meta)
}

func resourceContainerServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailClient(ctx)

	cs, err := FindContainerServiceByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Lightsail Container Service (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Lightsail Container Service (%s): %s", d.Id(), err)
	}

	d.Set("name", cs.ContainerServiceName)
	d.Set("power", cs.Power)
	d.Set("scale", cs.Scale)
	d.Set("is_disabled", cs.IsDisabled)

	if err := d.Set("public_domain_names", flattenContainerServicePublicDomainNames(cs.PublicDomainNames)); err != nil {
		return diag.Errorf("setting public_domain_names for Lightsail Container Service (%s): %s", d.Id(), err)
	}
	if err := d.Set("private_registry_access", []interface{}{flattenPrivateRegistryAccess(cs.PrivateRegistryAccess)}); err != nil {
		return diag.Errorf("setting private_registry_access for Lightsail Container Service (%s): %s", d.Id(), err)
	}
	d.Set("arn", cs.Arn)
	d.Set("availability_zone", cs.Location.AvailabilityZone)
	d.Set("created_at", aws.ToTime(cs.CreatedAt).Format(time.RFC3339))
	d.Set("power_id", cs.PowerId)
	d.Set("principal_arn", cs.PrincipalArn)
	d.Set("private_domain_name", cs.PrivateDomainName)
	d.Set("resource_type", cs.ResourceType)
	d.Set("state", cs.State)
	d.Set("url", cs.Url)

	SetTagsOut(ctx, cs.Tags)

	return nil
}

func resourceContainerServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailClient(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		publicDomainNames, _ := containerServicePublicDomainNamesChanged(d)

		input := &lightsail.UpdateContainerServiceInput{
			ServiceName:       aws.String(d.Id()),
			IsDisabled:        aws.Bool(d.Get("is_disabled").(bool)),
			Power:             types.ContainerServicePowerName(d.Get("power").(string)),
			PublicDomainNames: publicDomainNames,
			Scale:             aws.Int32(int32(d.Get("scale").(int))),
		}

		_, err := conn.UpdateContainerService(ctx, input)
		if err != nil {
			return diag.Errorf("updating Lightsail Container Service (%s): %s", d.Id(), err)
		}

		if d.HasChange("is_disabled") && d.Get("is_disabled").(bool) {
			if err := waitContainerServiceDisabled(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
				return diag.Errorf("waiting for Lightsail Container Service (%s) update: %s", d.Id(), err)
			}
		} else {
			if err := waitContainerServiceUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
				return diag.Errorf("waiting for Lightsail Container Service (%s) update: %s", d.Id(), err)
			}
		}
	}

	return resourceContainerServiceRead(ctx, d, meta)
}

func resourceContainerServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailClient(ctx)

	input := &lightsail.DeleteContainerServiceInput{
		ServiceName: aws.String(d.Id()),
	}

	_, err := conn.DeleteContainerService(ctx, input)

	if IsANotFoundError(err) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting Lightsail Container Service (%s): %s", d.Id(), err)
	}

	if err := waitContainerServiceDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("waiting for Lightsail Container Service (%s) deletion: %s", d.Id(), err)
	}

	return nil
}

func expandContainerServicePublicDomainNames(rawPublicDomainNames []interface{}) map[string][]string {
	if len(rawPublicDomainNames) == 0 {
		return nil
	}

	resultMap := make(map[string][]string)

	for _, rpdn := range rawPublicDomainNames {
		rpdnMap := rpdn.(map[string]interface{})

		rawCertificates := rpdnMap["certificate"].(*schema.Set).List()

		for _, rc := range rawCertificates {
			rcMap := rc.(map[string]interface{})

			var domainNames []string
			for _, rawDomainName := range rcMap["domain_names"].([]interface{}) {
				domainNames = append(domainNames, rawDomainName.(string))
			}

			certificateName := rcMap["certificate_name"].(string)

			resultMap[certificateName] = domainNames
		}
	}

	return resultMap
}

func expandPrivateRegistryAccess(tfMap map[string]interface{}) *types.PrivateRegistryAccessRequest {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.PrivateRegistryAccessRequest{}

	if v, ok := tfMap["ecr_image_puller_role"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.EcrImagePullerRole = expandECRImagePullerRole(v[0].(map[string]interface{}))
	}

	return apiObject
}

func expandECRImagePullerRole(tfMap map[string]interface{}) *types.ContainerServiceECRImagePullerRoleRequest {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.ContainerServiceECRImagePullerRoleRequest{}

	if v, ok := tfMap["is_active"].(bool); ok {
		apiObject.IsActive = aws.Bool(v)
	}

	return apiObject
}

func flattenPrivateRegistryAccess(apiObject *types.PrivateRegistryAccess) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.EcrImagePullerRole; v != nil {
		tfMap["ecr_image_puller_role"] = []interface{}{flattenECRImagePullerRole(v)}
	}

	return tfMap
}

func flattenECRImagePullerRole(apiObject *types.ContainerServiceECRImagePullerRole) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.IsActive; v != nil {
		tfMap["is_active"] = aws.ToBool(v)
	}

	if v := apiObject.PrincipalArn; v != nil {
		tfMap["principal_arn"] = aws.ToString(v)
	}

	return tfMap
}

func flattenContainerServicePublicDomainNames(domainNames map[string][]string) []interface{} {
	if domainNames == nil {
		return []interface{}{}
	}

	var rawCertificates []interface{}

	for certName, domains := range domainNames {
		rawCertificate := map[string]interface{}{
			"certificate_name": certName,
			"domain_names":     domains,
		}

		rawCertificates = append(rawCertificates, rawCertificate)
	}

	return []interface{}{
		map[string]interface{}{
			"certificate": rawCertificates,
		},
	}
}

func containerServicePublicDomainNamesChanged(d *schema.ResourceData) (map[string][]string, bool) {
	o, n := d.GetChange("public_domain_names")
	oldPublicDomainNames := expandContainerServicePublicDomainNames(o.([]interface{}))
	newPublicDomainNames := expandContainerServicePublicDomainNames(n.([]interface{}))

	changed := !reflect.DeepEqual(oldPublicDomainNames, newPublicDomainNames)
	if changed {
		if newPublicDomainNames == nil {
			newPublicDomainNames = map[string][]string{}
		}

		// if the change is to detach a certificate, in .tf, a certificate block is removed
		// however, an empty []*string entry must be added to tell Lightsail that we want none of the domain names
		// under the certificate, effectively detaching the certificate
		for certificateName := range oldPublicDomainNames {
			if _, ok := newPublicDomainNames[certificateName]; !ok {
				newPublicDomainNames[certificateName] = []string{}
			}
		}
	}

	return newPublicDomainNames, changed
}

func flattenContainerServicePowerValues(t []types.ContainerServicePowerName) []string {
	var out []string

	for _, v := range t {
		out = append(out, string(v))
	}

	return out
}

func FindContainerServiceByName(ctx context.Context, conn *lightsail.Client, serviceName string) (*types.ContainerService, error) {
	input := &lightsail.GetContainerServicesInput{
		ServiceName: aws.String(serviceName),
	}

	output, err := conn.GetContainerServices(ctx, input)

	if IsANotFoundError(err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.ContainerServices) == 0 {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.ContainerServices); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	return &output.ContainerServices[0], nil
}
