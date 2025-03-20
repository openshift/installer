// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsClusterNetworkInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsClusterNetworkInterfaceCreate,
		ReadContext:   resourceIBMIsClusterNetworkInterfaceRead,
		UpdateContext: resourceIBMIsClusterNetworkInterfaceUpdate,
		DeleteContext: resourceIBMIsClusterNetworkInterfaceDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"cluster_network_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_interface", "cluster_network_id"),
				Description:  "The cluster network identifier.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network_interface", "name"),
				Description:  "The name for this cluster network interface. The name is unique across all interfaces in the cluster network.",
			},
			"primary_ip": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The cluster network subnet reserved IP for this cluster network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							ConflictsWith: []string{"primary_ip.0.id", "primary_ip.0.href"},
							Description:   "The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The URL for this cluster network subnet reserved IP.",
						},
						"id": &schema.Schema{
							Type:          schema.TypeString,
							ForceNew:      true,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"primary_ip.0.address", "primary_ip.0.href"},
							Description:   "The unique identifier for this cluster network subnet reserved IP.",
						},
						"name": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"primary_ip.0.id", "primary_ip.0.href"},
							Description:   "The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.",
						},
						"auto_delete": &schema.Schema{
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"primary_ip.0.id", "primary_ip.0.href"},
							Description:   "Indicates whether this cluster network subnet reserved IP member will be automatically deleted when either target is deleted, or the cluster network subnet reserved IP is unbound.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"subnet": &schema.Schema{
				Type:         schema.TypeList,
				MaxItems:     1,
				AtLeastOneOf: []string{"subnet", "primary_ip.0.id", "primary_ip.0.href"},
				Optional:     true,
				Computed:     true,
				Description:  "The associated cluster network subnet. Required if `primary_ip` does not specify a clusternetwork subnet reserved IP identity.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"subnet.0.id"},
							Computed:      true,
							Description:   "The URL for this cluster network subnet.",
						},
						"id": &schema.Schema{
							Type:          schema.TypeString,
							ForceNew:      true,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"subnet.0.href"},
							Description:   "The unique identifier for this cluster network subnet.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"allow_ip_spoofing": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether source IP spoofing is allowed on this cluster network interface. If `false`, source IP spoofing is prevented on this cluster network interface. If `true`, source IP spoofing is allowed on this cluster network interface.",
			},
			"auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this cluster network interface will be automatically deleted when `target` is deleted.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network interface was created.",
			},
			"enable_infrastructure_nat": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the virtual network interface,  allowing the workload to perform any needed NAT operations.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network interface.",
			},
			"lifecycle_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current `lifecycle_state` (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the cluster network interface.",
			},
			"mac_address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The MAC address of the cluster network interface. May be absent if`lifecycle_state` is `pending`.",
			},
			// "protocol_state_filtering_mode": &schema.Schema{
			// 	Type:        schema.TypeString,
			// 	Computed:    true,
			// 	Description: "The protocol state filtering mode used for this cluster network interface.Protocol state filtering monitors each network connection flowing over this cluster network interface, and drops any packets that are invalid based on the current connection state and protocol. See [Protocol state filtering mode](https://cloud.ibm.com/docs/vpc?topic=vpc-vni-about#protocol-state-filtering) for more information.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			// },
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The target of this cluster network interface.If absent, this cluster network interface is not attached to a target.The resources supported by this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The URL for this instance cluster network attachment.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique identifier for this instance cluster network attachment.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC this cluster network interface resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zone this cluster network interface resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this zone.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this zone.",
						},
					},
				},
			},
			"cluster_network_interface_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this cluster network interface.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIsClusterNetworkInterfaceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster_network_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_cluster_network_interface", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsClusterNetworkInterfaceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createClusterNetworkInterfaceOptions := &vpcv1.CreateClusterNetworkInterfaceOptions{}

	createClusterNetworkInterfaceOptions.SetClusterNetworkID(d.Get("cluster_network_id").(string))
	if _, ok := d.GetOk("name"); ok {
		createClusterNetworkInterfaceOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("primary_ip"); ok {
		primaryIPModel, err := ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototype(d.Get("primary_ip.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "create", "parse-primary_ip").GetDiag()
		}
		createClusterNetworkInterfaceOptions.SetPrimaryIP(primaryIPModel)
	}
	if _, ok := d.GetOk("subnet"); ok {
		subnetModel, err := ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentity(d.Get("subnet.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "create", "parse-subnet").GetDiag()
		}
		createClusterNetworkInterfaceOptions.SetSubnet(subnetModel)
	}

	clusterNetworkInterface, _, err := vpcClient.CreateClusterNetworkInterfaceWithContext(context, createClusterNetworkInterfaceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateClusterNetworkInterfaceWithContext failed: %s", err.Error()), "ibm_is_cluster_network_interface", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createClusterNetworkInterfaceOptions.ClusterNetworkID, *clusterNetworkInterface.ID))

	return resourceIBMIsClusterNetworkInterfaceRead(context, d, meta)
}

func resourceIBMIsClusterNetworkInterfaceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkInterfaceOptions := &vpcv1.GetClusterNetworkInterfaceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "sep-id-parts").GetDiag()
	}

	getClusterNetworkInterfaceOptions.SetClusterNetworkID(parts[0])
	getClusterNetworkInterfaceOptions.SetID(parts[1])

	clusterNetworkInterface, response, err := vpcClient.GetClusterNetworkInterfaceWithContext(context, getClusterNetworkInterfaceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkInterfaceWithContext failed: %s", err.Error()), "ibm_is_cluster_network_interface", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(clusterNetworkInterface.Name) {
		if err = d.Set("name", clusterNetworkInterface.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(clusterNetworkInterface.PrimaryIP) {
		primaryIPMap, err := ResourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReservedIPReferenceToMap(clusterNetworkInterface.PrimaryIP)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "primary_ip-to-map").GetDiag()
		}
		if err = d.Set("primary_ip", []map[string]interface{}{primaryIPMap}); err != nil {
			err = fmt.Errorf("Error setting primary_ip: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-primary_ip").GetDiag()
		}
	}
	if !core.IsNil(clusterNetworkInterface.Subnet) {
		subnetMap, err := ResourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReferenceToMap(clusterNetworkInterface.Subnet)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "subnet-to-map").GetDiag()
		}
		if err = d.Set("subnet", []map[string]interface{}{subnetMap}); err != nil {
			err = fmt.Errorf("Error setting subnet: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-subnet").GetDiag()
		}
	}
	if err = d.Set("allow_ip_spoofing", clusterNetworkInterface.AllowIPSpoofing); err != nil {
		err = fmt.Errorf("Error setting allow_ip_spoofing: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-allow_ip_spoofing").GetDiag()
	}
	if err = d.Set("auto_delete", clusterNetworkInterface.AutoDelete); err != nil {
		err = fmt.Errorf("Error setting auto_delete: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-auto_delete").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(clusterNetworkInterface.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("enable_infrastructure_nat", clusterNetworkInterface.EnableInfrastructureNat); err != nil {
		err = fmt.Errorf("Error setting enable_infrastructure_nat: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-enable_infrastructure_nat").GetDiag()
	}
	if err = d.Set("href", clusterNetworkInterface.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-href").GetDiag()
	}
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range clusterNetworkInterface.LifecycleReasons {
		lifecycleReasonsItemMap, err := ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		err = fmt.Errorf("Error setting lifecycle_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-lifecycle_reasons").GetDiag()
	}
	if err = d.Set("lifecycle_state", clusterNetworkInterface.LifecycleState); err != nil {
		err = fmt.Errorf("Error setting lifecycle_state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-lifecycle_state").GetDiag()
	}
	if !core.IsNil(clusterNetworkInterface.MacAddress) {
		if err = d.Set("mac_address", clusterNetworkInterface.MacAddress); err != nil {
			err = fmt.Errorf("Error setting mac_address: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-mac_address").GetDiag()
		}
	}
	// if !core.IsNil(clusterNetworkInterface.ProtocolStateFilteringMode) {
	// 	if err = d.Set("protocol_state_filtering_mode", clusterNetworkInterface.ProtocolStateFilteringMode); err != nil {
	// 		err = fmt.Errorf("Error setting protocol_state_filtering_mode: %s", err)
	// 		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-protocol_state_filtering_mode").GetDiag()
	// 	}
	// }
	if err = d.Set("resource_type", clusterNetworkInterface.ResourceType); err != nil {
		err = fmt.Errorf("Error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-resource_type").GetDiag()
	}
	targetMap := make(map[string]interface{})
	if !core.IsNil(clusterNetworkInterface.Target) {
		targetMap, err = ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetToMap(clusterNetworkInterface.Target)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "target-to-map").GetDiag()
		}
	}
	if err = d.Set("target", []map[string]interface{}{targetMap}); err != nil {
		err = fmt.Errorf("Error setting target: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-target").GetDiag()
	}
	vpcMap, err := ResourceIBMIsClusterNetworkInterfaceVPCReferenceToMap(clusterNetworkInterface.VPC)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "vpc-to-map").GetDiag()
	}
	if err = d.Set("vpc", []map[string]interface{}{vpcMap}); err != nil {
		err = fmt.Errorf("Error setting vpc: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-vpc").GetDiag()
	}
	zoneMap, err := ResourceIBMIsClusterNetworkInterfaceZoneReferenceToMap(clusterNetworkInterface.Zone)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "zone-to-map").GetDiag()
	}
	if err = d.Set("zone", []map[string]interface{}{zoneMap}); err != nil {
		err = fmt.Errorf("Error setting zone: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-zone").GetDiag()
	}
	if err = d.Set("cluster_network_interface_id", clusterNetworkInterface.ID); err != nil {
		err = fmt.Errorf("Error setting cluster_network_interface_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "read", "set-cluster_network_interface_id").GetDiag()
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_is_cluster_network_interface", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIBMIsClusterNetworkInterfaceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateClusterNetworkInterfaceOptions := &vpcv1.UpdateClusterNetworkInterfaceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "update", "sep-id-parts").GetDiag()
	}

	updateClusterNetworkInterfaceOptions.SetClusterNetworkID(parts[0])
	updateClusterNetworkInterfaceOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.ClusterNetworkInterfacePatch{}
	if d.HasChange("cluster_network_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "cluster_network_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_is_cluster_network_interface", "update", "cluster_network_id-forces-new").GetDiag()
	}
	if d.HasChange("auto_delete") {
		newAutoDelete := d.Get("auto_delete").(bool)
		patchVals.AutoDelete = &newAutoDelete
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	// updateClusterNetworkInterfaceOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateClusterNetworkInterfaceOptions.ClusterNetworkInterfacePatch = ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfacePatchAsPatch(patchVals, d)

		_, _, err = vpcClient.UpdateClusterNetworkInterfaceWithContext(context, updateClusterNetworkInterfaceOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateClusterNetworkInterfaceWithContext failed: %s", err.Error()), "ibm_is_cluster_network_interface", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMIsClusterNetworkInterfaceRead(context, d, meta)
}

func resourceIBMIsClusterNetworkInterfaceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteClusterNetworkInterfaceOptions := &vpcv1.DeleteClusterNetworkInterfaceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_cluster_network_interface", "delete", "sep-id-parts").GetDiag()
	}

	deleteClusterNetworkInterfaceOptions.SetClusterNetworkID(parts[0])
	deleteClusterNetworkInterfaceOptions.SetID(parts[1])

	_, _, err = vpcClient.DeleteClusterNetworkInterfaceWithContext(context, deleteClusterNetworkInterfaceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteClusterNetworkInterfaceWithContext failed: %s", err.Error()), "ibm_is_cluster_network_interface", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap map[string]interface{}) (vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeIntf, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextIntf, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext{}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentity(modelMap map[string]interface{}) (vpcv1.ClusterNetworkSubnetIdentityIntf, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentityByID(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetIdentityByID, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetIdentityByHref, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReservedIPReferenceToMap(model *vpcv1.ClusterNetworkSubnetReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsClusterNetworkInterfaceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkInterfaceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReferenceToMap(model *vpcv1.ClusterNetworkSubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsClusterNetworkInterfaceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceLifecycleReasonToMap(model *vpcv1.ClusterNetworkInterfaceLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetToMap(model vpcv1.ClusterNetworkInterfaceTargetIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceClusterNetworkInterfaceContext); ok {
		return ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceClusterNetworkInterfaceContextToMap(model.(*vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceClusterNetworkInterfaceContext))
	} else if _, ok := model.(*vpcv1.ClusterNetworkInterfaceTarget); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ClusterNetworkInterfaceTarget)
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.ResourceType != nil {
			modelMap["resource_type"] = *model.ResourceType
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ClusterNetworkInterfaceTargetIntf subtype encountered")
	}
}

func ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceClusterNetworkInterfaceContextToMap(model *vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceClusterNetworkInterfaceContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkInterfaceVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsClusterNetworkInterfaceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkInterfaceZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfacePatchAsPatch(patchVals *vpcv1.ClusterNetworkInterfacePatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "auto_delete"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["auto_delete"] = nil
	}
	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}

	return patch
}
