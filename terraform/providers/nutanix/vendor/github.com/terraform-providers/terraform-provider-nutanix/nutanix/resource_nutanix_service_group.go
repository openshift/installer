package nutanix

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixServiceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixServiceGroupCreate,
		ReadContext:   resourceNutanixServiceGroupRead,
		DeleteContext: resourceNutanixServiceGroupDelete,
		UpdateContext: resourceNutanixServiceGroupUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_type_code_list": {
							Type: schema.TypeList,

							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func IsValidProtocol(category string) bool {
	switch category {
	case
		"ALL",
		"ICMP",
		"TCP",
		"UDP":
		return true
	}
	return false
}

func resourceNutanixServiceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API
	id := d.Id()
	response, err := conn.V3.GetServiceGroup(id)

	request := &v3.ServiceGroupInput{}

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return diag.Errorf("error retrieving for access control policy id (%s) :%+v", id, err)
	}

	group := response.ServiceGroup

	if d.HasChange("name") {
		group.Name = utils.StringPtr(d.Get("name").(string))
	}

	if d.HasChange("description") {
		group.Description = utils.StringPtr(d.Get("description").(string))
	}

	if d.HasChange("system_defined") {
		group.SystemDefined = utils.BoolPtr(d.Get("system_defined").(bool))
	}

	if d.HasChange("service_list") {
		serviceList, err := expandServiceEntry(d)

		if err != nil {
			return diag.FromErr(err)
		}

		group.ServiceList = serviceList
	}

	request.SystemDefined = group.SystemDefined
	request.Name = group.Name
	request.Description = group.Description
	request.ServiceList = group.ServiceList

	errUpdate := conn.V3.UpdateServiceGroup(d.Id(), request)
	if errUpdate != nil {
		return diag.Errorf("error updating service group id %s): %s", d.Id(), errUpdate)
	}

	return resourceNutanixServiceGroupRead(ctx, d, meta)
}

func flattenServiceEntry(group *v3.ServiceGroupInput) []map[string]interface{} {
	groupList := make([]map[string]interface{}, 0)

	for _, v := range group.ServiceList {
		groupItem := make(map[string]interface{})
		groupItem["protocol"] = utils.StringValue(v.Protocol)

		if v.TCPPortRangeList != nil {
			tcpprl := v.TCPPortRangeList
			tcpprList := make([]map[string]interface{}, len(tcpprl))
			for i, tcp := range tcpprl {
				tcpItem := make(map[string]interface{})
				tcpItem["end_port"] = utils.Int64Value(tcp.EndPort)
				tcpItem["start_port"] = utils.Int64Value(tcp.StartPort)
				tcpprList[i] = tcpItem
			}
			groupItem["tcp_port_range_list"] = tcpprList
		}

		if v.UDPPortRangeList != nil {
			udpprl := v.UDPPortRangeList
			udpprList := make([]map[string]interface{}, len(udpprl))
			for i, udp := range udpprl {
				udpItem := make(map[string]interface{})
				udpItem["end_port"] = utils.Int64Value(udp.EndPort)
				udpItem["start_port"] = utils.Int64Value(udp.StartPort)
				udpprList[i] = udpItem
			}
			groupItem["udp_port_range_list"] = udpprList
		}

		if v.IcmpTypeCodeList != nil {
			icmptcl := v.IcmpTypeCodeList
			icmptcList := make([]map[string]interface{}, len(icmptcl))
			for i, icmp := range icmptcl {
				icmpItem := make(map[string]interface{})
				icmpItem["code"] = strconv.FormatInt(utils.Int64Value(icmp.Code), 10)
				icmpItem["type"] = strconv.FormatInt(utils.Int64Value(icmp.Type), 10)
				icmptcList[i] = icmpItem
			}
			groupItem["icmp_type_code_list"] = icmptcList
		}

		groupList = append(groupList, groupItem)
	}
	return groupList
}

func resourceNutanixServiceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading ServiceGroup: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetServiceGroup(d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("service_list", flattenServiceEntry(resp.ServiceGroup)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", utils.StringValue(resp.ServiceGroup.Name))
	d.Set("description", utils.StringValue(resp.ServiceGroup.Description))

	return diag.FromErr(d.Set("system_defined", utils.BoolValue(resp.ServiceGroup.SystemDefined)))
}

func expandServiceEntry(d *schema.ResourceData) ([]*v3.ServiceListEntry, error) {
	if services, ok := d.GetOk("service_list"); ok {
		set := services.([]interface{})
		outbound := make([]*v3.ServiceListEntry, len(set))

		for k, v := range set {
			service := &v3.ServiceListEntry{}

			entry := v.(map[string]interface{})

			if proto, pok := entry["protocol"]; pok && proto.(string) != "" {
				if !IsValidProtocol(proto.(string)) {
					return nil, fmt.Errorf("protocol needs to be one of 'ALL', 'ICMP', 'TCP', 'UDP'")
				}
				service.Protocol = utils.StringPtr(proto.(string))
			}

			if t, tcpok := entry["tcp_port_range_list"]; tcpok {
				service.TCPPortRangeList = expandPortRangeList(t)
			}

			if u, udpok := entry["udp_port_range_list"]; udpok {
				service.UDPPortRangeList = expandPortRangeList(u)
			}

			if icmp, icmpok := entry["icmp_type_code_list"]; icmpok {
				service.IcmpTypeCodeList = expandIcmpTypeCodeList(icmp)
			}

			outbound[k] = service
		}

		return outbound, nil
	}
	return nil, nil
}

func resourceNutanixServiceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	request := &v3.ServiceGroupInput{}
	request.ServiceList = make([]*v3.ServiceListEntry, 0)

	name, nameOK := d.GetOk("name")

	// Read Arguments and set request values
	if desc, ok := d.GetOk("description"); ok {
		request.Description = utils.StringPtr(desc.(string))
	}

	// validate required fields
	if !nameOK {
		return diag.Errorf("please provide the required attribute name")
	}

	request.Name = utils.StringPtr(name.(string))

	serviceList, err := expandServiceEntry(d)

	if err != nil {
		return diag.FromErr(err)
	}

	request.ServiceList = serviceList

	requestEnc, err := json.Marshal(request)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s", requestEnc)

	resp, err := conn.V3.CreateServiceGroup(request)

	if err != nil {
		return diag.FromErr(err)
	}

	n := *resp.UUID

	// set terraform state
	d.SetId(n)

	return resourceNutanixServiceGroupRead(ctx, d, meta)
}

func resourceNutanixServiceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	log.Printf("[Debug] Destroying the service group with the ID %s", d.Id())

	if err := conn.V3.DeleteServiceGroup(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
