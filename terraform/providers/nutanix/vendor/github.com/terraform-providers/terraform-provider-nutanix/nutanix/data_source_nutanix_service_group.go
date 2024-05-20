package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// dataSourceNutanixServiceGroup returns schema for datasource nutanix service group
// /v3/service_groups/{uuid}
// https://www.nutanix.dev/api_references/prism-central-v3/#/b3A6MjU1ODc2OTk-get-a-existing-service-group
func dataSourceNutanixServiceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixServiceGroupRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"is_system_defined": {
				Type:        schema.TypeBool,
				Description: "specifying whether it is a system defined service group",
				Computed:    true,
			},
		},
	}
}

func dataSourceNutanixServiceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	if uuid, uuidOk := d.GetOk("uuid"); uuidOk {
		svGroup, reqErr := conn.V3.GetServiceGroup(uuid.(string))

		if reqErr != nil {
			return diag.Errorf("error reading user with error %s", reqErr)
		}

		if err := d.Set("name", utils.StringValue(svGroup.ServiceGroup.Name)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("description", utils.StringValue(svGroup.ServiceGroup.Description)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("is_system_defined", utils.BoolValue(svGroup.ServiceGroup.SystemDefined)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("service_list", flattenServiceEntry(svGroup.ServiceGroup)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(uuid.(string))
	} else {
		return diag.Errorf("please provide `uuid`")
	}
	return nil
}
