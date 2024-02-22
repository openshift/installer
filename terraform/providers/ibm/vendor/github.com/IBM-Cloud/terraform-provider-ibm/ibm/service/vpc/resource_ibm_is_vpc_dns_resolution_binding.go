// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsVPCDnsResolutionBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsVPCDnsResolutionBindingCreate,
		ReadContext:   resourceIBMIsVPCDnsResolutionBindingRead,
		UpdateContext: resourceIBMIsVPCDnsResolutionBindingUpdate,
		DeleteContext: resourceIBMIsVPCDnsResolutionBindingDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC identifier.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the DNS resolution binding was created.",
			},
			"endpoint_gateways": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The endpoint gateways in the bound to VPC that are allowed to participate in this DNS resolution binding.The endpoint gateways may be remote and therefore may not be directly retrievable.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this endpoint gateway.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this endpoint gateway.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this endpoint gateway.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this endpoint gateway. The name is unique across all endpoint gateways in the VPC.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this account.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"region": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this region.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique name for this region.",
												},
											},
										},
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this DNS resolution binding.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the DNS resolution binding.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name for this DNS resolution binding. The name is unique across all DNS resolution bindings for the VPC.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				MinItems:    1,
				Description: "The VPC bound to for DNS resolution.The VPC may be remote and therefore may not be directly retrievable.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Optional:     true,
							ExactlyOneOf: []string{"vpc.0.id", "vpc.0.href", "vpc.0.crn"},
							Type:         schema.TypeString,
							Computed:     true,
							Description:  "The CRN for this VPC.",
						},
						"href": &schema.Schema{
							Type:         schema.TypeString,
							ExactlyOneOf: []string{"vpc.0.id", "vpc.0.href", "vpc.0.crn"},
							Optional:     true,
							Computed:     true,
							Description:  "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:         schema.TypeString,
							ExactlyOneOf: []string{"vpc.0.id", "vpc.0.href", "vpc.0.crn"},
							Optional:     true,
							Computed:     true,
							Description:  "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this account.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"region": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this region.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique name for this region.",
												},
											},
										},
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func resourceIBMIsVPCDnsResolutionBindingCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	spokeVPCID := d.Get("vpc_id").(string)
	createVPCDnsResolutionBindingOptions := &vpcv1.CreateVPCDnsResolutionBindingOptions{}
	vpchref := d.Get("vpc.0.href").(string)
	vpccrn := d.Get("vpc.0.crn").(string)
	vpcid := d.Get("vpc.0.id").(string)

	createVPCDnsResolutionBindingOptions.SetVPCID(spokeVPCID)
	if d.Get("name").(string) != "" {
		createVPCDnsResolutionBindingOptions.SetName(d.Get("name").(string))
	}
	if vpchref != "" {
		vPCIdentityIntf := &vpcv1.VPCIdentityByHref{
			Href: &vpchref,
		}
		createVPCDnsResolutionBindingOptions.SetVPC(vPCIdentityIntf)
	} else if vpcid != "" {
		vPCIdentityIntf := &vpcv1.VPCIdentityByID{
			ID: &vpcid,
		}
		createVPCDnsResolutionBindingOptions.SetVPC(vPCIdentityIntf)
	} else {
		vPCIdentityIntf := &vpcv1.VPCIdentityByCRN{
			CRN: &vpccrn,
		}
		createVPCDnsResolutionBindingOptions.SetVPC(vPCIdentityIntf)
	}
	vpcdnsResolutionBinding, response, err := sess.CreateVPCDnsResolutionBindingWithContext(context, createVPCDnsResolutionBindingOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateVPCDnsResolutionBindingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateVPCDnsResolutionBindingWithContext failed %s\n%s", err, response))
	}
	d.SetId(MakeTerraformVPCDNSID(spokeVPCID, *vpcdnsResolutionBinding.ID))

	err = resourceIBMIsVPCDnsResolutionBindingGet(vpcdnsResolutionBinding, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
func resourceIBMIsVPCDnsResolutionBindingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vpcId, id, err := ParseVPCDNSTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	getVPCDnsResolutionBindingOptions := &vpcv1.GetVPCDnsResolutionBindingOptions{}

	getVPCDnsResolutionBindingOptions.SetVPCID(vpcId)
	getVPCDnsResolutionBindingOptions.SetID(id)

	vpcdnsResolutionBinding, response, err := sess.GetVPCDnsResolutionBindingWithContext(context, getVPCDnsResolutionBindingOptions)
	if err != nil {
		log.Printf("[DEBUG] GetVPCDnsResolutionBindingWithContext failed %s\n%s", err, response)
		if response.StatusCode != 404 {
			return diag.FromErr(fmt.Errorf("GetVPCDnsResolutionBindingWithContext failed %s\n%s", err, response))
		} else {
			d.SetId("")
			return nil
		}
	}
	err = resourceIBMIsVPCDnsResolutionBindingGet(vpcdnsResolutionBinding, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
func resourceIBMIsVPCDnsResolutionBindingGet(vpcdnsResolutionBinding *vpcv1.VpcdnsResolutionBinding, d *schema.ResourceData) error {
	if err := d.Set("created_at", flex.DateTimeToString(vpcdnsResolutionBinding.CreatedAt)); err != nil {
		return fmt.Errorf("[ERROR] Error setting created_at: %s", err)
	}

	endpointGateways := []map[string]interface{}{}
	if vpcdnsResolutionBinding.EndpointGateways != nil {
		for _, modelItem := range vpcdnsResolutionBinding.EndpointGateways {
			modelMap, err := dataSourceIBMIsVPCDnsResolutionBindingEndpointGatewayReferenceRemoteToMap(&modelItem)
			if err != nil {
				return err
			}
			endpointGateways = append(endpointGateways, modelMap)
		}
	}
	if err := d.Set("endpoint_gateways", endpointGateways); err != nil {
		return fmt.Errorf("[ERROR] Error setting endpoint_gateways %s", err)
	}

	if err := d.Set("href", vpcdnsResolutionBinding.Href); err != nil {
		return fmt.Errorf("[ERROR] Error setting href: %s", err)
	}

	if err := d.Set("lifecycle_state", vpcdnsResolutionBinding.LifecycleState); err != nil {
		return fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err)
	}

	if err := d.Set("name", vpcdnsResolutionBinding.Name); err != nil {
		return fmt.Errorf("[ERROR] Error setting name: %s", err)
	}

	if err := d.Set("resource_type", vpcdnsResolutionBinding.ResourceType); err != nil {
		return fmt.Errorf("[ERROR] Error setting resource_type: %s", err)
	}

	vpc := []map[string]interface{}{}
	if vpcdnsResolutionBinding.VPC != nil {
		modelMap, err := dataSourceIBMIsVPCDnsResolutionBindingVPCReferenceRemoteToMap(vpcdnsResolutionBinding.VPC)
		if err != nil {
			return err
		}
		vpc = append(vpc, modelMap)
	}
	if err := d.Set("vpc", vpc); err != nil {
		return fmt.Errorf("[ERROR] Error setting vpc %s", err)
	}
	return nil
}
func resourceIBMIsVPCDnsResolutionBindingUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vpcId, id, err := ParseVPCDNSTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange(isVPCDnsResolutionBindingName) {
		nameChange := d.Get(isVPCDnsResolutionBindingName).(string)
		vpcdnsResolutionBindingPatch := &vpcv1.VpcdnsResolutionBindingPatch{
			Name: &nameChange,
		}
		vpcdnsResolutionBindingPatchAsPatch, err := vpcdnsResolutionBindingPatch.AsPatch()
		if err != nil {
			return diag.FromErr(err)
		}
		updateVPCDnsResolutionBindingOptions := &vpcv1.UpdateVPCDnsResolutionBindingOptions{}

		updateVPCDnsResolutionBindingOptions.SetVPCID(vpcId)
		updateVPCDnsResolutionBindingOptions.SetID(id)
		updateVPCDnsResolutionBindingOptions.SetVpcdnsResolutionBindingPatch(vpcdnsResolutionBindingPatchAsPatch)

		vpcdnsResolutionBinding, response, err := sess.UpdateVPCDnsResolutionBindingWithContext(context, updateVPCDnsResolutionBindingOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateVPCDnsResolutionBindingWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateVPCDnsResolutionBindingWithContext failed %s\n%s", err, response))
		}
		err = resourceIBMIsVPCDnsResolutionBindingGet(vpcdnsResolutionBinding, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
func resourceIBMIsVPCDnsResolutionBindingDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vpcId, id, err := ParseVPCDNSTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	deleteVPCDnsResolutionBindingOptions := &vpcv1.DeleteVPCDnsResolutionBindingOptions{}

	deleteVPCDnsResolutionBindingOptions.SetVPCID(vpcId)
	deleteVPCDnsResolutionBindingOptions.SetID(id)

	_, response, err := sess.DeleteVPCDnsResolutionBindingWithContext(context, deleteVPCDnsResolutionBindingOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteVPCDnsResolutionBindingWithContext failed %s\n%s", err, response)
		if response.StatusCode != 404 {
			return diag.FromErr(fmt.Errorf("DeleteVPCDnsResolutionBindingWithContext failed %s\n%s", err, response))
		}
	}
	d.SetId("")
	return nil
}
func MakeTerraformVPCDNSID(id1, id2 string) string {
	// Include both  vpc id and binding id to create a unique Terraform id.  As a bonus,
	// we can extract the bindings as needed for API calls such as READ.
	return fmt.Sprintf("%s/%s", id1, id2)
}

func ParseVPCDNSTerraformID(s string) (string, string, error) {
	segments := strings.Split(s, "/")
	if len(segments) != 2 {
		return "", "", fmt.Errorf("invalid terraform Id %s (incorrect number of segments)", s)
	}
	if segments[0] == "" || segments[1] == "" {
		return "", "", fmt.Errorf("invalid terraform Id %s (one or more empty segments)", s)
	}
	return segments[0], segments[1], nil
}
