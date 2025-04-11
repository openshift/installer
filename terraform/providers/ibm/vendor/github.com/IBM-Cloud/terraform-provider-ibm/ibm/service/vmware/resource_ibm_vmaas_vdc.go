// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.98.0-8be2046a-20241205-162752
 */

package vmware

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vmware-go-sdk/vmwarev1"
)

func ResourceIbmVmaasVdc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmVmaasVdcCreate,
		ReadContext:   resourceIbmVmaasVdcRead,
		UpdateContext: resourceIbmVmaasVdcUpdate,
		DeleteContext: resourceIbmVmaasVdcDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"accept_language": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_vmaas_vdc", "accept_language"),
				Description:  "Language.",
			},
			"cpu": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_vmaas_vdc", "cpu"),
				Description:  "The vCPU usage limit on the virtual data center (VDC). Supported for VDCs deployed on a multitenant Cloud Director site. This property is applicable when the resource pool type is reserved.",
			},
			"director_site": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The Cloud Director site in which to deploy the virtual data center (VDC).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "A unique ID for the Cloud Director site.",
						},
						"pvdc": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "The resource pool within the Director Site in which to deploy the virtual data center (VDC).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "A unique ID for the resource pool.",
									},
									"provider_type": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Determines how resources are made available to the virtual data center (VDC). Required for VDCs deployed on a multitenant Cloud Director site.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of the resource pool type.",
												},
											},
										},
									},
								},
							},
						},
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the VMware Cloud Director tenant portal where this virtual data center (VDC) can be managed.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_vmaas_vdc", "name"),
				Description:  "A human readable ID for the virtual data center (VDC).",
			},
			"ram": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_vmaas_vdc", "ram"),
				Description:  "The RAM usage limit on the virtual data center (VDC) in GB (1024^3 bytes). Supported for VDCs deployed on a multitenant Cloud Director site. This property is applicable when the resource pool type is reserved.",
			},
			"fast_provisioning_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether this virtual data center has fast provisioning enabled or not.",
			},
			"rhel_byol": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if the RHEL VMs will be using the license from IBM or the customer will use their own license (BYOL).",
			},
			"windows_byol": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if the Microsoft Windows VMs will be using the license from IBM or the customer will use their own license (BYOL).",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of this virtual data center (VDC).",
			},
			"provisioned_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time that the virtual data center (VDC) is provisioned and available to use.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique ID for the virtual data center (VDC) in IBM Cloud.",
			},
			"deleted_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time that the virtual data center (VDC) is deleted.",
			},
			"edges": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VMware NSX-T networking edges deployed on the virtual data center (VDC). NSX-T edges are used for bridging virtualization networking to the physical public-internet and IBM private networking.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique ID for the edge.",
						},
						"public_ips": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The public IP addresses assigned to the edge.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"private_ips": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The private IP addresses assigned to the edge.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"private_only": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the edge is private only. The default value is True for a private Cloud Director site and False for a public Cloud Director site.",
						},
						"size": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The size of the edge.The size can be specified only for performance edges. Larger sizes require more capacity from the Cloud Director site in which the virtual data center (VDC) was created to be deployed.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Determines the state of the edge.",
						},
						"transit_gateways": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connected IBM Transit Gateways.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique ID for an IBM Transit Gateway.",
									},
									"connections": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IBM Transit Gateway connections.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The autogenerated name for this connection.",
												},
												"transit_gateway_connection_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name of the connection created on the IBM Transit Gateway.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Determines the state of the connection.",
												},
												"local_gateway_ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Local gateway IP address for the connection.",
												},
												"remote_gateway_ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Remote gateway IP address for the connection.",
												},
												"local_tunnel_ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Local tunnel IP address for the connection.",
												},
												"remote_tunnel_ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Remote tunnel IP address for the connection.",
												},
												"local_bgp_asn": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Local network BGP ASN for the connection.",
												},
												"remote_bgp_asn": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Remote network BGP ASN for the connection.",
												},
												"network_account_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the account that owns the connected network.",
												},
												"network_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of the network that is connected through this connection. Only \"unbound_gre_tunnel\" is supported.",
												},
												"base_network_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of the network that the unbound GRE tunnel is targeting. Only \"classic\" is supported.",
												},
												"zone": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The location of the connection.",
												},
											},
										},
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Determines the state of the IBM Transit Gateway based on its connections.",
									},
									"region": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The region where the IBM Transit Gateway is deployed.",
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of edge to be deployed.Efficiency edges allow for multiple VDCs to share some edge resources. Performance edges do not share resources between VDCs.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The edge version.",
						},
					},
				},
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about why the request to create the virtual data center (VDC) cannot be completed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An error code specific to the error encountered.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A message that describes why the error ocurred.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL that links to a page with more information about this error.",
						},
					},
				},
			},
			"ordered_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time that the virtual data center (VDC) is ordered.",
			},
			"org_href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the organization that owns the VDC.",
			},
			"org_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the VMware Cloud Director organization that contains this virtual data center (VDC). VMware Cloud Director organizations are used to create strong boundaries between VDCs. There is a complete isolation of user administration, networking, workloads, and VMware Cloud Director catalogs between different Director organizations.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines the state of the virtual data center.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines whether this virtual data center is in a single-tenant or multitenant Cloud Director site.",
			},
		},
	}
}

func ResourceIbmVmaasVdcValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "accept_language",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9-,;=\.\*\s]{1,256}$`,
			MinValueLength:             1,
			MaxValueLength:             256,
		},
		validate.ValidateSchema{
			Identifier:                 "cpu",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "0",
			MaxValue:                   "2000",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z][A-Za-z0-9_\-]{1,128}$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "ram",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "0",
			MaxValue:                   "40960",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_vmaas_vdc", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmVmaasVdcCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vmwareClient, err := meta.(conns.ClientSession).VmwareV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createVdcOptions := &vmwarev1.CreateVdcOptions{}

	createVdcOptions.SetName(d.Get("name").(string))
	directorSiteModel, err := ResourceIbmVmaasVdcMapToVDCDirectorSitePrototype(d.Get("director_site.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "create", "parse-director_site").GetDiag()
	}
	createVdcOptions.SetDirectorSite(directorSiteModel)
	if _, ok := d.GetOk("edge"); ok {
		edgeModel, err := ResourceIbmVmaasVdcMapToVDCEdgePrototype(d.Get("edge.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "create", "parse-edge").GetDiag()
		}
		createVdcOptions.SetEdge(edgeModel)
	}
	if _, ok := d.GetOk("fast_provisioning_enabled"); ok {
		createVdcOptions.SetFastProvisioningEnabled(d.Get("fast_provisioning_enabled").(bool))
	}
	if _, ok := d.GetOk("resource_group"); ok {
		resourceGroupModel, err := ResourceIbmVmaasVdcMapToResourceGroupIdentity(d.Get("resource_group.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "create", "parse-resource_group").GetDiag()
		}
		createVdcOptions.SetResourceGroup(resourceGroupModel)
	}
	if _, ok := d.GetOk("cpu"); ok {
		createVdcOptions.SetCpu(int64(d.Get("cpu").(int)))
	}
	if _, ok := d.GetOk("ram"); ok {
		createVdcOptions.SetRam(int64(d.Get("ram").(int)))
	}
	if _, ok := d.GetOk("rhel_byol"); ok {
		createVdcOptions.SetRhelByol(d.Get("rhel_byol").(bool))
	}
	if _, ok := d.GetOk("windows_byol"); ok {
		createVdcOptions.SetWindowsByol(d.Get("windows_byol").(bool))
	}
	if _, ok := d.GetOk("accept_language"); ok {
		createVdcOptions.SetAcceptLanguage(d.Get("accept_language").(string))
	}

	vdc, _, err := vmwareClient.CreateVdcWithContext(context, createVdcOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateVdcWithContext failed: %s", err.Error()), "ibm_vmaas_vdc", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*vdc.ID)

	if waitForVdcStatus {
		_, err = waitForVdcStatusUpdate(context, d, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for VDC (%s) to be reeady: %s", *vdc.ID, err))
		}
	}

	return resourceIbmVmaasVdcRead(context, d, meta)
}

func resourceIbmVmaasVdcRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vmwareClient, err := meta.(conns.ClientSession).VmwareV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getVdcOptions := &vmwarev1.GetVdcOptions{}

	getVdcOptions.SetID(d.Id())
	if _, ok := d.GetOk("accept_language"); ok {
		getVdcOptions.SetAcceptLanguage(d.Get("accept_language").(string))
	}

	vdc, response, err := vmwareClient.GetVdcWithContext(context, getVdcOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVdcWithContext failed: %s", err.Error()), "ibm_vmaas_vdc", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(vdc.Cpu) {
		if err = d.Set("cpu", flex.IntValue(vdc.Cpu)); err != nil {
			err = fmt.Errorf("Error setting cpu: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-cpu").GetDiag()
		}
	}
	directorSiteMap, err := ResourceIbmVmaasVdcVDCDirectorSiteToMap(vdc.DirectorSite)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "director_site-to-map").GetDiag()
	}
	if err = d.Set("director_site", []map[string]interface{}{directorSiteMap}); err != nil {
		err = fmt.Errorf("Error setting director_site: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-director_site").GetDiag()
	}
	if err = d.Set("name", vdc.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-name").GetDiag()
	}
	if !core.IsNil(vdc.Ram) {
		if err = d.Set("ram", flex.IntValue(vdc.Ram)); err != nil {
			err = fmt.Errorf("Error setting ram: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-ram").GetDiag()
		}
	}
	if !core.IsNil(vdc.FastProvisioningEnabled) {
		if err = d.Set("fast_provisioning_enabled", vdc.FastProvisioningEnabled); err != nil {
			err = fmt.Errorf("Error setting fast_provisioning_enabled: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-fast_provisioning_enabled").GetDiag()
		}
	}
	if !core.IsNil(vdc.RhelByol) {
		if err = d.Set("rhel_byol", vdc.RhelByol); err != nil {
			err = fmt.Errorf("Error setting rhel_byol: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-rhel_byol").GetDiag()
		}
	}
	if !core.IsNil(vdc.WindowsByol) {
		if err = d.Set("windows_byol", vdc.WindowsByol); err != nil {
			err = fmt.Errorf("Error setting windows_byol: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-windows_byol").GetDiag()
		}
	}
	if err = d.Set("href", vdc.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-href").GetDiag()
	}
	if !core.IsNil(vdc.ProvisionedAt) {
		if err = d.Set("provisioned_at", flex.DateTimeToString(vdc.ProvisionedAt)); err != nil {
			err = fmt.Errorf("Error setting provisioned_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-provisioned_at").GetDiag()
		}
	}
	if err = d.Set("crn", vdc.Crn); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-crn").GetDiag()
	}
	if !core.IsNil(vdc.DeletedAt) {
		if err = d.Set("deleted_at", flex.DateTimeToString(vdc.DeletedAt)); err != nil {
			err = fmt.Errorf("Error setting deleted_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-deleted_at").GetDiag()
		}
	}
	edges := []map[string]interface{}{}
	for _, edgesItem := range vdc.Edges {
		edgesItemMap, err := ResourceIbmVmaasVdcEdgeToMap(&edgesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "edges-to-map").GetDiag()
		}
		edges = append(edges, edgesItemMap)
	}
	if err = d.Set("edges", edges); err != nil {
		err = fmt.Errorf("Error setting edges: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-edges").GetDiag()
	}
	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range vdc.StatusReasons {
		statusReasonsItemMap, err := ResourceIbmVmaasVdcStatusReasonToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "status_reasons-to-map").GetDiag()
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		err = fmt.Errorf("Error setting status_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-status_reasons").GetDiag()
	}
	if err = d.Set("ordered_at", flex.DateTimeToString(vdc.OrderedAt)); err != nil {
		err = fmt.Errorf("Error setting ordered_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-ordered_at").GetDiag()
	}
	if err = d.Set("org_href", vdc.OrgHref); err != nil {
		err = fmt.Errorf("Error setting org_href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-org_href").GetDiag()
	}
	if err = d.Set("org_name", vdc.OrgName); err != nil {
		err = fmt.Errorf("Error setting org_name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-org_name").GetDiag()
	}
	if err = d.Set("status", vdc.Status); err != nil {
		err = fmt.Errorf("Error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-status").GetDiag()
	}
	if err = d.Set("type", vdc.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "read", "set-type").GetDiag()
	}

	return nil
}

func resourceIbmVmaasVdcUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vmwareClient, err := meta.(conns.ClientSession).VmwareV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateVdcOptions := &vmwarev1.UpdateVdcOptions{}

	updateVdcOptions.SetID(d.Id())

	hasChange := false

	patchVals := &vmwarev1.VDCPatch{}
	if d.HasChange("accept_language") {
		updateVdcOptions.SetAcceptLanguage(d.Get("accept_language").(string))
		hasChange = true
	}
	if d.HasChange("cpu") {
		newCpu := int64(d.Get("cpu").(int))
		patchVals.Cpu = &newCpu
		hasChange = true
	}
	if d.HasChange("fast_provisioning_enabled") {
		newFastProvisioningEnabled := d.Get("fast_provisioning_enabled").(bool)
		patchVals.FastProvisioningEnabled = &newFastProvisioningEnabled
		hasChange = true
	}
	if d.HasChange("ram") {
		newRam := int64(d.Get("ram").(int))
		patchVals.Ram = &newRam
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateVdcOptions.VDCPatch = ResourceIbmVmaasVdcVDCPatchAsPatch(patchVals, d)

		_, _, err = vmwareClient.UpdateVdcWithContext(context, updateVdcOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateVdcWithContext failed: %s", err.Error()), "ibm_vmaas_vdc", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if waitForVdcStatus {
		_, err = waitForVdcStatusUpdate(context, d, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for VDC to be reeady: %s", err))
		}
	}

	return resourceIbmVmaasVdcRead(context, d, meta)
}

func resourceIbmVmaasVdcDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vmwareClient, err := meta.(conns.ClientSession).VmwareV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_vmaas_vdc", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteVdcOptions := &vmwarev1.DeleteVdcOptions{}

	deleteVdcOptions.SetID(d.Id())
	if _, ok := d.GetOk("accept_language"); ok {
		deleteVdcOptions.SetAcceptLanguage(d.Get("accept_language").(string))
	}

	vdcDelete, _, err := vmwareClient.DeleteVdcWithContext(context, deleteVdcOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteVdcWithContext failed: %s", err.Error()), "ibm_vmaas_vdc", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if waitForVdcStatus {
		_, err = waitForVdcToDelete(context, d, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for VDC (%s) to be deleted: %s", *vdcDelete.ID, err))
		}
	}

	d.SetId("")

	return nil
}

func ResourceIbmVmaasVdcMapToVDCDirectorSitePrototype(modelMap map[string]interface{}) (*vmwarev1.VDCDirectorSitePrototype, error) {
	model := &vmwarev1.VDCDirectorSitePrototype{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	PvdcModel, err := ResourceIbmVmaasVdcMapToDirectorSitePVDC(modelMap["pvdc"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Pvdc = PvdcModel
	return model, nil
}

func ResourceIbmVmaasVdcMapToDirectorSitePVDC(modelMap map[string]interface{}) (*vmwarev1.DirectorSitePVDC, error) {
	model := &vmwarev1.DirectorSitePVDC{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	if modelMap["provider_type"] != nil && len(modelMap["provider_type"].([]interface{})) > 0 {
		ProviderTypeModel, err := ResourceIbmVmaasVdcMapToVDCProviderType(modelMap["provider_type"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ProviderType = ProviderTypeModel
	}
	return model, nil
}

func ResourceIbmVmaasVdcMapToVDCProviderType(modelMap map[string]interface{}) (*vmwarev1.VDCProviderType, error) {
	model := &vmwarev1.VDCProviderType{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIbmVmaasVdcMapToVDCEdgePrototype(modelMap map[string]interface{}) (*vmwarev1.VDCEdgePrototype, error) {
	model := &vmwarev1.VDCEdgePrototype{}
	if modelMap["size"] != nil && modelMap["size"].(string) != "" {
		model.Size = core.StringPtr(modelMap["size"].(string))
	}
	model.Type = core.StringPtr(modelMap["type"].(string))
	if modelMap["private_only"] != nil {
		model.PrivateOnly = core.BoolPtr(modelMap["private_only"].(bool))
	}
	return model, nil
}

func ResourceIbmVmaasVdcMapToResourceGroupIdentity(modelMap map[string]interface{}) (*vmwarev1.ResourceGroupIdentity, error) {
	model := &vmwarev1.ResourceGroupIdentity{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIbmVmaasVdcVDCDirectorSiteToMap(model *vmwarev1.VDCDirectorSite) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	pvdcMap, err := ResourceIbmVmaasVdcDirectorSitePVDCToMap(model.Pvdc)
	if err != nil {
		return modelMap, err
	}
	modelMap["pvdc"] = []map[string]interface{}{pvdcMap}
	modelMap["url"] = *model.URL
	return modelMap, nil
}

func ResourceIbmVmaasVdcDirectorSitePVDCToMap(model *vmwarev1.DirectorSitePVDC) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	if model.ProviderType != nil {
		providerTypeMap, err := ResourceIbmVmaasVdcVDCProviderTypeToMap(model.ProviderType)
		if err != nil {
			return modelMap, err
		}
		modelMap["provider_type"] = []map[string]interface{}{providerTypeMap}
	}
	return modelMap, nil
}

func ResourceIbmVmaasVdcVDCProviderTypeToMap(model *vmwarev1.VDCProviderType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIbmVmaasVdcEdgeToMap(model *vmwarev1.Edge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["public_ips"] = model.PublicIps
	modelMap["private_ips"] = model.PrivateIps
	if model.PrivateOnly != nil {
		modelMap["private_only"] = *model.PrivateOnly
	}
	modelMap["size"] = *model.Size
	modelMap["status"] = *model.Status
	transitGateways := []map[string]interface{}{}
	for _, transitGatewaysItem := range model.TransitGateways {
		transitGatewaysItemMap, err := ResourceIbmVmaasVdcTransitGatewayToMap(&transitGatewaysItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		transitGateways = append(transitGateways, transitGatewaysItemMap)
	}
	modelMap["transit_gateways"] = transitGateways
	modelMap["type"] = *model.Type
	modelMap["version"] = *model.Version
	return modelMap, nil
}

func ResourceIbmVmaasVdcTransitGatewayToMap(model *vmwarev1.TransitGateway) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	connections := []map[string]interface{}{}
	for _, connectionsItem := range model.Connections {
		connectionsItemMap, err := ResourceIbmVmaasVdcTransitGatewayConnectionToMap(&connectionsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		connections = append(connections, connectionsItemMap)
	}
	modelMap["connections"] = connections
	modelMap["status"] = *model.Status
	modelMap["region"] = *model.Region
	return modelMap, nil
}

func ResourceIbmVmaasVdcTransitGatewayConnectionToMap(model *vmwarev1.TransitGatewayConnection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	if model.TransitGatewayConnectionName != nil {
		modelMap["transit_gateway_connection_name"] = *model.TransitGatewayConnectionName
	}
	modelMap["status"] = *model.Status
	if model.LocalGatewayIp != nil {
		modelMap["local_gateway_ip"] = *model.LocalGatewayIp
	}
	if model.RemoteGatewayIp != nil {
		modelMap["remote_gateway_ip"] = *model.RemoteGatewayIp
	}
	if model.LocalTunnelIp != nil {
		modelMap["local_tunnel_ip"] = *model.LocalTunnelIp
	}
	if model.RemoteTunnelIp != nil {
		modelMap["remote_tunnel_ip"] = *model.RemoteTunnelIp
	}
	if model.LocalBgpAsn != nil {
		modelMap["local_bgp_asn"] = flex.IntValue(model.LocalBgpAsn)
	}
	if model.RemoteBgpAsn != nil {
		modelMap["remote_bgp_asn"] = flex.IntValue(model.RemoteBgpAsn)
	}
	modelMap["network_account_id"] = *model.NetworkAccountID
	modelMap["network_type"] = *model.NetworkType
	modelMap["base_network_type"] = *model.BaseNetworkType
	modelMap["zone"] = *model.Zone
	return modelMap, nil
}

func ResourceIbmVmaasVdcStatusReasonToMap(model *vmwarev1.StatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIbmVmaasVdcVDCPatchAsPatch(patchVals *vmwarev1.VDCPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "cpu"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["cpu"] = nil
	} else if !exists {
		delete(patch, "cpu")
	}
	path = "fast_provisioning_enabled"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["fast_provisioning_enabled"] = nil
	} else if !exists {
		delete(patch, "fast_provisioning_enabled")
	}
	path = "ram"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["ram"] = nil
	} else if !exists {
		delete(patch, "ram")
	}

	return patch
}
