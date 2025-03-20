// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsPrivatePathServiceGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsPrivatePathServiceGatewayCreate,
		ReadContext:   resourceIBMIsPrivatePathServiceGatewayRead,
		UpdateContext: resourceIBMIsPrivatePathServiceGatewayUpdate,
		DeleteContext: resourceIBMIsPrivatePathServiceGatewayDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"service_endpoints": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The fully qualified domain names for this private path service gateway. ",
			},
			"default_access_policy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_private_path_service_gateway", "access_policy"),
				Description:  "The access policy for the account:- permit: access will be permitted- deny:  access will be denied- review: access will be manually reviewed.",
			},
			"load_balancer": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer for this private path service gateway. ",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of this PPSG ",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID of resource group to use.",
			},
			"zonal_affinity": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "ndicates whether this private path service gateway has zonal affinity.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the account policy was created.",
			},
			"endpoint_gateway_binding_auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether endpoint gateway bindings will be automatically deleted after endpoint_gateway_binding_auto_delete_timeout hours have passed.",
			},
			"endpoint_gateway_binding_auto_delete_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "If endpoint_gateway_binding_auto_delete is true, the hours after which endpoint gateway bindings will be automatically deleted.",
			},
			"endpoint_gateway_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of endpoint gateways using this private path service gateway.",
			},
			"published": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates the availability of this private path service gateway.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Href of this resource",
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPC this private path service gateway resides in.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "lifecycle_state of this resource",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of this resource.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the account policy was updated.",
			},
			// "remote": &schema.Schema{
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// 	Description: "If present, this property indicates that the resource associated with this reference is remote and therefore may not be directly retrievable..",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"account": &schema.Schema{
			// 				Type:        schema.TypeList,
			// 				Computed:    true,
			// 				Description: "If present, this property indicates that the referenced resource is remote to this account, and identifies the owning account.",
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"id": &schema.Schema{
			// 							Type:        schema.TypeList,
			// 							Computed:    true,
			// 							Description: "The unique identifier for this account.",
			// 						},
			// 						"resource_type": &schema.Schema{
			// 							Type:        schema.TypeString,
			// 							Computed:    true,
			// 							Description: "The resource type.",
			// 						},
			// 					},
			// 				},
			// 			},
			// 			"region": &schema.Schema{
			// 				Type:        schema.TypeList,
			// 				Computed:    true,
			// 				Description: "If present, this property indicates that the referenced resource is remote to this region, and identifies the native region.",
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"href": &schema.Schema{
			// 							Type:        schema.TypeList,
			// 							Computed:    true,
			// 							Description: "The URL for this region.",
			// 						},
			// 						"name": &schema.Schema{
			// 							Type:        schema.TypeString,
			// 							Computed:    true,
			// 							Description: "The globally unique name for this region.",
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			"private_path_service_gateway": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this account policy.",
			},
		},
	}
}

func ResourceIBMIsPrivatePathServiceGatewayValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "access_policy",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "deny, permit, review",
			Regexp:                     `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_private_path_service_gateway", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsPrivatePathServiceGatewayCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	loadBalancerId := d.Get("load_balancer").(string)

	createPrivatePathServiceGatewayOptions := &vpcv1.CreatePrivatePathServiceGatewayOptions{
		LoadBalancer: &vpcv1.LoadBalancerIdentity{
			ID: &loadBalancerId,
		},
	}
	serviceEndpoints := d.Get("service_endpoints").(*schema.Set)
	if serviceEndpoints.Len() != 0 {
		serviceEndpointsList := make([]string, serviceEndpoints.Len())
		for i, serviceEndpointsItem := range serviceEndpoints.List() {
			sEndpoint := serviceEndpointsItem.(string)
			serviceEndpointsList[i] = sEndpoint
		}
		createPrivatePathServiceGatewayOptions.ServiceEndpoints = serviceEndpointsList
	}
	if nameIntf, ok := d.GetOk("name"); ok {
		name := nameIntf.(string)
		createPrivatePathServiceGatewayOptions.Name = &name
	}
	if resGrpIntf, ok := d.GetOk("resource_group"); ok {
		resGrp := resGrpIntf.(string)
		createPrivatePathServiceGatewayOptions.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &resGrp,
		}
	}
	if defaultAccessPolicyIntf, ok := d.GetOk("default_access_policy"); ok {
		dAccessPolicy := defaultAccessPolicyIntf.(string)
		createPrivatePathServiceGatewayOptions.DefaultAccessPolicy = &dAccessPolicy
	}
	if zonalAffinityIntf, ok := d.GetOk("zonal_affinity"); ok {
		zonalAffinity := zonalAffinityIntf.(bool)
		createPrivatePathServiceGatewayOptions.ZonalAffinity = &zonalAffinity
	}

	privatePathServiceGateway, response, err := vpcClient.CreatePrivatePathServiceGatewayWithContext(context, createPrivatePathServiceGatewayOptions)
	if err != nil {
		log.Printf("[DEBUG] CreatePrivatePathServiceGatewayWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreatePrivatePathServiceGatewayWithContext failed %s\n%s", err, response))
	}

	d.SetId(*privatePathServiceGateway.ID)
	_, err = isWaitForPPSGAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMIsPrivatePathServiceGatewayUpdate(context, d, meta)
}

func resourceIBMIsPrivatePathServiceGatewayRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getPrivatePathServiceGatewayOptions := &vpcv1.GetPrivatePathServiceGatewayOptions{}

	getPrivatePathServiceGatewayOptions.SetID(d.Id())

	privatePathServiceGateway, response, err := vpcClient.GetPrivatePathServiceGatewayWithContext(context, getPrivatePathServiceGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetPrivatePathServiceGatewayWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetPrivatePathServiceGatewayWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("default_access_policy", privatePathServiceGateway.DefaultAccessPolicy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting access_policy: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(privatePathServiceGateway.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", privatePathServiceGateway.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("endpoint_gateway_count", privatePathServiceGateway.EndpointGatewayCount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoint_gateway_count: %s", err))
	}
	if err = d.Set("endpoint_gateway_binding_auto_delete", privatePathServiceGateway.EndpointGatewayBindingAutoDelete); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoint_gateway_binding_auto_delete: %s", err))
	}
	if err = d.Set("endpoint_gateway_binding_auto_delete_timeout", privatePathServiceGateway.EndpointGatewayBindingAutoDeleteTimeout); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoint_gateway_binding_auto_delete_timeout: %s", err))
	}
	if err = d.Set("published", privatePathServiceGateway.Published); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting published: %s", err))
	}
	if err = d.Set("load_balancer", *privatePathServiceGateway.LoadBalancer.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting load balancer id: %s", err))
	}
	if err = d.Set("lifecycle_state", privatePathServiceGateway.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("name", privatePathServiceGateway.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("vpc", privatePathServiceGateway.VPC.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vpc: %s", err))
	}
	if err = d.Set("zonal_affinity", privatePathServiceGateway.ZonalAffinity); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zonal_affinity: %s", err))
	}
	serviceEndpointsList := make([]string, 0)
	for i := 0; i < len(privatePathServiceGateway.ServiceEndpoints); i++ {
		serviceEndpointsList = append(serviceEndpointsList, string(privatePathServiceGateway.ServiceEndpoints[i]))
	}
	if err = d.Set("service_endpoints", serviceEndpointsList); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_endpoints: %s", err))
	}
	if err = d.Set("crn", privatePathServiceGateway.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("resource_type", privatePathServiceGateway.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("private_path_service_gateway", privatePathServiceGateway.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_path_service_gateway: %s", err))
	}

	return nil
}

func resourceIBMIsPrivatePathServiceGatewayUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updatePrivatePathServiceGatewayOptions := &vpcv1.UpdatePrivatePathServiceGatewayOptions{}
	updatePrivatePathServiceGatewayOptions.SetID(d.Id())
	hasChange := false

	patchVals := &vpcv1.PrivatePathServiceGatewayPatch{}

	if d.HasChange("default_access_policy") && !d.IsNewResource() {
		newAccessPolicy := d.Get("default_access_policy").(string)
		patchVals.DefaultAccessPolicy = &newAccessPolicy
		hasChange = true
	}

	if d.HasChange("name") && !d.IsNewResource() {
		name := d.Get("name").(string)
		patchVals.Name = &name
		hasChange = true
	}
	if d.HasChange("zonal_affinity") && !d.IsNewResource() {
		zonalAffinity := d.Get("zonal_affinity").(bool)
		patchVals.ZonalAffinity = &zonalAffinity
		hasChange = true
	}
	// if d.HasChange("published") {
	// 	published := d.Get("published").(bool)
	// 	patchVals.Published = &published
	// 	hasChange = true
	// }
	if d.HasChange("load_balancer") && !d.IsNewResource() {
		loadBalancer := d.Get("load_balancer").(string)
		patchVals.LoadBalancer = &vpcv1.LoadBalancerIdentity{
			ID: &loadBalancer,
		}
		hasChange = true
	}

	if hasChange {
		updatePrivatePathServiceGatewayOptions.PrivatePathServiceGatewayPatch, _ = patchVals.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] Error calling AsPatch for PrivatePathServiceGatewayPatch %s", err)
			return diag.FromErr(err)
		}
		_, response, err := vpcClient.UpdatePrivatePathServiceGatewayWithContext(context, updatePrivatePathServiceGatewayOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdatePrivatePathServiceGatewayWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdatePrivatePathServiceGatewayWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsPrivatePathServiceGatewayRead(context, d, meta)
}

func resourceIBMIsPrivatePathServiceGatewayDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deletePrivatePathServiceGatewayOptions := &vpcv1.DeletePrivatePathServiceGatewayOptions{}
	deletePrivatePathServiceGatewayOptions.SetID(d.Id())

	response, err := vpcClient.DeletePrivatePathServiceGatewayWithContext(context, deletePrivatePathServiceGatewayOptions)
	if err != nil {
		log.Printf("[DEBUG] DeletePrivatePathServiceGatewayWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeletePrivatePathServiceGatewayWithContext failed %s\n%s", err, response))
	}
	_, err = isWaitForPPSGDeleted(vpcClient, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func isWaitForPPSGDeleteRetry(vpcClient *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("[DEBUG] Retrying PPSG (%s) delete", id)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"ppsg_in_use"},
		Target:  []string{"deleting", "done", ""},
		Refresh: func() (interface{}, string, error) {
			deletePrivatePathServiceGatewayOptions := &vpcv1.DeletePrivatePathServiceGatewayOptions{}
			deletePrivatePathServiceGatewayOptions.SetID(id)
			log.Printf("[DEBUG] Retrying PPSG (%s) delete", id)
			response, err := vpcClient.DeletePrivatePathServiceGateway(deletePrivatePathServiceGatewayOptions)
			if err != nil {
				if response != nil && response.StatusCode == 409 {
					return response, "ppsg_in_use", nil
				} else if response != nil && response.StatusCode == 404 {
					return response, "done", nil
				}
				return response, "", fmt.Errorf("[ERROR] Error deleting ppsg: %s\n%s", err, response)
			}
			return response, "deleting", nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
func isWaitForPPSGDeleted(vpcClient *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for ppsg (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", "deleting", "stable"},
		Target:     []string{"deleted", ""},
		Refresh:    isPPSGDeleteRefreshFunc(vpcClient, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isPPSGDeleteRefreshFunc(vpcClient *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] is ppsg delete function here")
		getPPSGOptions := &vpcv1.GetPrivatePathServiceGatewayOptions{
			ID: &id,
		}
		ppsg, response, err := vpcClient.GetPrivatePathServiceGateway(getPPSGOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return ppsg, "deleted", nil
			}
			return ppsg, "", fmt.Errorf("[ERROR] The ppsg %s failed to delete: %s\n%s", id, err, response)
		}
		return ppsg, "deleting", err
	}
}
func isWaitForPPSGAvailable(vpcClient *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for ppsg (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "updating"},
		Target:     []string{"stable", "failed", "suspended"},
		Refresh:    isPPSGRefreshFunc(vpcClient, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isPPSGRefreshFunc(vpcClient *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPPSGOptions := &vpcv1.GetPrivatePathServiceGatewayOptions{
			ID: &id,
		}
		ppsg, response, err := vpcClient.GetPrivatePathServiceGateway(getPPSGOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting ppsg : %s\n%s", err, response)
		}

		return ppsg, *ppsg.LifecycleState, nil
	}
}
