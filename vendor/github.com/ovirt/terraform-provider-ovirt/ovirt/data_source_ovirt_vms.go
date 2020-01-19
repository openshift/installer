// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func dataSourceOvirtVMs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtVMsRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchemaWith(
				"max",
				"criteria",
				"case_sensitive",
			),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			// Computed
			"vms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"high_availability": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sockets": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"threads": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"reported_devices": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"comment": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ips": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"gateway": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"netmask": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"version": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceOvirtVMsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	vmsReq := conn.SystemService().VmsService().List()

	search, searchOK := d.GetOk("search")
	nameRegex, nameRegexOK := d.GetOk("name_regex")

	if searchOK {
		searchMap := search.(map[string]interface{})
		searchCriteria, searchCriteriaOK := searchMap["criteria"]
		searchMax, searchMaxOK := searchMap["max"]
		searchCaseSensitive, searchCaseSensitiveOK := searchMap["case_sensitive"]
		if !searchCriteriaOK && !searchMaxOK && !searchCaseSensitiveOK {
			return fmt.Errorf("One of criteria or max or case_sensitive in search must be assigned")
		}

		if searchCriteriaOK {
			vmsReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			vmsReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			vmsReq.CaseSensitive(csBool)
		}
	}
	vmsResp, err := vmsReq.Send()
	if err != nil {
		return err
	}
	vms, ok := vmsResp.Vms()
	if !ok || len(vms.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredVMs []*ovirtsdk4.Vm
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, vm := range vms.Slice() {
			if r.MatchString(vm.MustName()) {
				filteredVMs = append(filteredVMs, vm)
			}
		}
	} else {
		filteredVMs = vms.Slice()[:]
	}

	if len(filteredVMs) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return vmsDescriptionAttributes(d, filteredVMs, meta)
}

func vmsDescriptionAttributes(d *schema.ResourceData, vms []*ovirtsdk4.Vm, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	var s []map[string]interface{}

	for _, v := range vms {
		mapping := map[string]interface{}{
			"id":                v.MustId(),
			"name":              v.MustName(),
			"status":            v.MustStatus(),
			"template_id":       v.MustTemplate().MustId(),
			"cluster_id":        v.MustCluster().MustId(),
			"high_availability": v.MustHighAvailability().MustEnabled(),
			"memory":            v.MustMemory() / int64(math.Pow(2, 20)),
			"cores":             v.MustCpu().MustTopology().MustCores(),
			"sockets":           v.MustCpu().MustTopology().MustSockets(),
			"threads":           v.MustCpu().MustTopology().MustThreads(),
		}

		devicesResp, err := conn.SystemService().
			VmsService().
			VmService(v.MustId()).
			ReportedDevicesService().
			List().
			Send()
		if err != nil {
			return err
		}
		if devices, ok := devicesResp.ReportedDevice(); ok && len(devices.Slice()) > 0 {
			mapping["reported_devices"] = flattenOvirtNicReportedDevices(devices.Slice())
		}

		s = append(s, mapping)
	}
	d.SetId(resource.UniqueId())
	return d.Set("vms", s)
}
