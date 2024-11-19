// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISSubnetRead,

		Schema: map[string]*schema.Schema{

			"identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isSubnetName, "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_subnet", "identifier"),
			},

			isSubnetIpv4CidrBlock: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isSubnetAvailableIpv4AddressCount: {
				Type:     schema.TypeInt,
				Computed: true,
			},

			isSubnetTotalIpv4AddressCount: {
				Type:     schema.TypeInt,
				Computed: true,
			},

			isSubnetName: {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{isSubnetName, "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_subnet", isSubnetName),
			},

			isSubnetTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags",
			},

			isSubnetAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
			},

			isSubnetCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			isSubnetNetworkACL: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isSubnetPublicGateway: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isSubnetStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isSubnetVPC: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{isSubnetName},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_subnet", "identifier"),
			},

			isSubnetVPCName: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isSubnetZone: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isSubnetResourceGroup: {
				Type:     schema.TypeString,
				Computed: true,
			},

			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},

			"routing_table": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The routing table for this subnet",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this routing table.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this routing table.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this routing table.",
						},
						"resource_type": {
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

func DataSourceIBMISSubnetValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "identifier",
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isSubnetName,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})

	ibmISSubnetDataSourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_subnet", Schema: validateSchema}
	return &ibmISSubnetDataSourceValidator
}

func dataSourceIBMISSubnetRead(d *schema.ResourceData, meta interface{}) error {
	err := subnetGetByNameOrID(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func subnetGetByNameOrID(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	var subnet *vpcv1.Subnet
	if v, ok := d.GetOk("identifier"); ok {
		id := v.(string)
		getSubnetOptions := &vpcv1.GetSubnetOptions{
			ID: &id,
		}
		subnetinfo, response, err := sess.GetSubnet(getSubnetOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Getting Subnet (%s): %s\n%s", id, err, response)
		}
		subnet = subnetinfo
	} else if v, ok := d.GetOk(isSubnetName); ok {
		name := v.(string)
		start := ""
		allrecs := []vpcv1.Subnet{}
		getSubnetsListOptions := &vpcv1.ListSubnetsOptions{}
		for {
			if start != "" {
				getSubnetsListOptions.Start = &start
			}
			if vpcIdOk, ok := d.GetOk(isSubnetVPC); ok {
				vpcIDOk := vpcIdOk.(string)
				getSubnetsListOptions.VPCID = &vpcIDOk
			}
			subnetsCollection, response, err := sess.ListSubnets(getSubnetsListOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Fetching subnets List %s\n%s", err, response)
			}
			start = flex.GetNext(subnetsCollection.Next)
			allrecs = append(allrecs, subnetsCollection.Subnets...)
			if start == "" {
				break
			}
		}

		for _, subnetInfo := range allrecs {
			if *subnetInfo.Name == name {
				subnet = &subnetInfo
				break
			}
		}
		if subnet == nil {
			return fmt.Errorf("[ERROR] No subnet found with name (%s)", name)
		}
	}

	d.SetId(*subnet.ID)

	if subnet.RoutingTable != nil {
		err = d.Set("routing_table", dataSourceSubnetFlattenroutingTable(*subnet.RoutingTable))
		if err != nil {
			return fmt.Errorf("Error setting routing_table %s", err)
		}
	}
	d.Set(isSubnetName, *subnet.Name)
	d.Set(isSubnetIpv4CidrBlock, *subnet.Ipv4CIDRBlock)
	d.Set(isSubnetAvailableIpv4AddressCount, *subnet.AvailableIpv4AddressCount)
	d.Set(isSubnetTotalIpv4AddressCount, *subnet.TotalIpv4AddressCount)
	if subnet.NetworkACL != nil {
		d.Set(isSubnetNetworkACL, *subnet.NetworkACL.ID)
	}
	if subnet.PublicGateway != nil {
		d.Set(isSubnetPublicGateway, *subnet.PublicGateway.ID)
	} else {
		d.Set(isSubnetPublicGateway, nil)
	}
	d.Set(isSubnetStatus, *subnet.Status)
	d.Set(isSubnetZone, *subnet.Zone.Name)
	d.Set(isSubnetVPC, *subnet.VPC.ID)
	d.Set(isSubnetVPCName, *subnet.VPC.Name)

	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}

	tags, err := flex.GetGlobalTagsUsingCRN(meta, *subnet.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"An error occured during reading of subnet (%s) tags : %s", d.Id(), err)
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *subnet.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource subnet (%s) access tags: %s", d.Id(), err)
	}

	d.Set(isSubnetTags, tags)
	d.Set(isSubnetAccessTags, accesstags)
	d.Set(isSubnetCRN, *subnet.CRN)
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/network/subnets")
	d.Set(flex.ResourceName, *subnet.Name)
	d.Set(flex.ResourceCRN, *subnet.CRN)
	d.Set(flex.ResourceStatus, *subnet.Status)
	if subnet.ResourceGroup != nil {
		d.Set(isSubnetResourceGroup, *subnet.ResourceGroup.ID)
		d.Set(flex.ResourceGroupName, *subnet.ResourceGroup.Name)
	}
	return nil
}

func dataSourcesubnetRoutingTableDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceSubnetFlattenroutingTable(result vpcv1.RoutingTableReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceSubnetRoutingTableToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceSubnetRoutingTableToMap(routingTableItem vpcv1.RoutingTableReference) (routingTableMap map[string]interface{}) {
	routingTableMap = map[string]interface{}{}

	if routingTableItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourcesubnetRoutingTableDeletedToMap(*routingTableItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		routingTableMap["deleted"] = deletedList
	}
	if routingTableItem.Href != nil {
		routingTableMap["href"] = routingTableItem.Href
	}
	if routingTableItem.ID != nil {
		routingTableMap["id"] = routingTableItem.ID
	}
	if routingTableItem.Name != nil {
		routingTableMap["name"] = routingTableItem.Name
	}
	if routingTableItem.ResourceType != nil {
		routingTableMap["resource_type"] = routingTableItem.ResourceType
	}

	return routingTableMap
}
