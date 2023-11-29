// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isVPCDnsResolutionBindings                = "dns_resolution_bindings"
	isVPCDnsResolutionBindingVpcId            = "vpc_id"
	isVPCDnsResolutionBindingCreatedAt        = "created_at"
	isVPCDnsResolutionBindingEndpointGateways = "endpoint_gateways"
	isVPCDnsResolutionBindingCrn              = "crn"
	isVPCDnsResolutionBindingId               = "id"
	isVPCDnsResolutionBindingHref             = "href"
	isVPCDnsResolutionBindingLifecycleState   = "lifecycle_state"
	isVPCDnsResolutionBindingName             = "name"
	isVPCDnsResolutionBindingResourceType     = "resource_type"
	isVPCDnsResolutionBindingRemote           = "remote"
	isVPCDnsResolutionBindingAccount          = "account"
	isVPCDnsResolutionBindingRegion           = "region"
)

func DataSourceIBMIsVPCDnsResolutionBindings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPCDnsResolutionBindingsRead,
		Schema: map[string]*schema.Schema{
			isVPCDnsResolutionBindings: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of VPC Dns Resolution Bindings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPCDnsResolutionBindingId: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DNS resolution binding identifier.",
						},
						isVPCDnsResolutionBindingCreatedAt: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the DNS resolution binding was created.",
						},
						isVPCDnsResolutionBindingEndpointGateways: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The endpoint gateways in the bound to VPC that are allowed to participate in this DNS resolution binding.The endpoint gateways may be remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVPCDnsResolutionBindingCrn: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this endpoint gateway.",
									},
									isVPCDnsResolutionBindingHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this endpoint gateway.",
									},
									isVPCDnsResolutionBindingId: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this endpoint gateway.",
									},
									isVPCDnsResolutionBindingName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this endpoint gateway. The name is unique across all endpoint gateways in the VPC.",
									},
									isVPCDnsResolutionBindingRemote: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVPCDnsResolutionBindingAccount: &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isVPCDnsResolutionBindingId: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this account.",
															},
															isVPCDnsResolutionBindingResourceType: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The resource type.",
															},
														},
													},
												},
												isVPCDnsResolutionBindingRegion: &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isVPCDnsResolutionBindingId: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this region.",
															},
															isVPCDnsResolutionBindingName: &schema.Schema{
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
									isVPCDnsResolutionBindingResourceType: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						isVPCDnsResolutionBindingHref: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this DNS resolution binding.",
						},
						isVPCDnsResolutionBindingLifecycleState: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the DNS resolution binding.",
						},
						isVPCDnsResolutionBindingName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this DNS resolution binding. The name is unique across all DNS resolution bindings for the VPC.",
						},
						isVPCDnsResolutionBindingResourceType: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"vpc": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC bound to for DNS resolution.The VPC may be remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVPCDnsResolutionBindingCrn: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
									},
									isVPCDnsResolutionBindingHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this VPC.",
									},
									isVPCDnsResolutionBindingId: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									isVPCDnsResolutionBindingName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this VPC. The name is unique across all VPCs in the region.",
									},
									isVPCDnsResolutionBindingRemote: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVPCDnsResolutionBindingAccount: &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isVPCDnsResolutionBindingId: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this account.",
															},
															isVPCDnsResolutionBindingResourceType: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The resource type.",
															},
														},
													},
												},
												isVPCDnsResolutionBindingRegion: &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isVPCDnsResolutionBindingId: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this region.",
															},
															isVPCDnsResolutionBindingName: &schema.Schema{
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
									isVPCDnsResolutionBindingResourceType: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
					},
				},
			},
			isVPCDnsResolutionBindingVpcId: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC identifier.",
			},
		},
	}
}

func dataSourceIBMIsVPCDnsResolutionBindingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	listVPCDnsResolutionBindingOptions := &vpcv1.ListVPCDnsResolutionBindingsOptions{}

	listVPCDnsResolutionBindingOptions.SetVPCID(d.Get(isVPCDnsResolutionBindingVpcId).(string))
	start := ""
	allrecs := []vpcv1.VpcdnsResolutionBinding{}

	for {
		if start != "" {
			listVPCDnsResolutionBindingOptions.Start = &start
		}
		vpcdnsResolutionBindingCollection, response, err := sess.ListVPCDnsResolutionBindingsWithContext(context, listVPCDnsResolutionBindingOptions)
		if err != nil {
			log.Printf("[DEBUG] ListVPCDnsResolutionBindingsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListVPCDnsResolutionBindingsWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(vpcdnsResolutionBindingCollection.Next)
		allrecs = append(allrecs, vpcdnsResolutionBindingCollection.DnsResolutionBindings...)
		if start == "" {
			break
		}
	}
	vpcdnsResolutionBindingsInfo := make([]map[string]interface{}, 0)
	if len(allrecs) != 0 {
		for _, dns := range allrecs {
			l := map[string]interface{}{}
			l[isVPCDnsResolutionBindingId] = *dns.ID

			l[isVPCDnsResolutionBindingCreatedAt] = flex.DateTimeToString(dns.CreatedAt)

			endpointGateways := []map[string]interface{}{}
			if dns.EndpointGateways != nil {
				for _, modelItem := range dns.EndpointGateways {
					modelMap, err := dataSourceIBMIsVPCDnsResolutionBindingEndpointGatewayReferenceRemoteToMap(&modelItem)
					if err != nil {
						return diag.FromErr(err)
					}
					endpointGateways = append(endpointGateways, modelMap)
				}
			}
			l[isVPCDnsResolutionBindingEndpointGateways] = endpointGateways

			l[isVPCDnsResolutionBindingId] = dns.ID

			l[isVPCDnsResolutionBindingLifecycleState] = dns.LifecycleState

			l[isVPCDnsResolutionBindingName] = dns.Name
			l[isVPCDnsResolutionBindingHref] = dns.Href

			l[isVPCDnsResolutionBindingResourceType] = dns.ResourceType

			vpc := []map[string]interface{}{}
			if dns.VPC != nil {
				modelMap, err := dataSourceIBMIsVPCDnsResolutionBindingVPCReferenceRemoteToMap(dns.VPC)
				if err != nil {
					return diag.FromErr(err)
				}
				vpc = append(vpc, modelMap)
			}
			l["vpc"] = vpc
			vpcdnsResolutionBindingsInfo = append(vpcdnsResolutionBindingsInfo, l)
		}
	}
	d.SetId(dataSourceIBMIsVPCDnsResolutionBindingsId(d))
	d.Set(isVPCDnsResolutionBindings, vpcdnsResolutionBindingsInfo)
	return nil
}

// dataSourceIBMIsVPCDnsResolutionBindingsId returns a reasonable ID for VPC Dns Resolution Bindings list.
func dataSourceIBMIsVPCDnsResolutionBindingsId(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
