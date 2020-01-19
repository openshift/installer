package ovirt

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func dataSourceOvirtNics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtNicsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"vm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nics": {
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
						"boot_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interface": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"linked": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"on_boot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"plugged": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
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

func dataSourceOvirtNicsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	vmID := d.Get("vm_id").(string)
	nameRegex, nameRegexOK := d.GetOk("name_regex")

	vmNicsServiceReq := conn.SystemService().VmsService().VmService(vmID).NicsService().List()

	nicsResp, err := vmNicsServiceReq.Send()
	if err != nil {
		return err
	}

	nics, ok := nicsResp.Nics()
	if !ok || len(nics.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredNics []*ovirtsdk4.Nic

	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, nic := range nics.Slice() {
			if r.MatchString(nic.MustName()) {
				filteredNics = append(filteredNics, nic)
			}
		}
	} else {
		filteredNics = nics.Slice()[:]
	}

	if len(filteredNics) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return nicsDecriptionAttributes(d, filteredNics, meta)
}

func nicsDecriptionAttributes(d *schema.ResourceData, nics []*ovirtsdk4.Nic, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range nics {
		mapping := map[string]interface{}{
			"id":        v.MustId(),
			"name":      v.MustName(),
			"interface": string(v.MustInterface()),
			"plugged":   v.MustPlugged(),
			"linked":    v.MustLinked(),
		}

		if bootProtocol, ok := v.BootProtocol(); ok {
			mapping["boot_protocol"] = string(bootProtocol)
		}
		if comment, ok := v.Comment(); ok {
			mapping["comment"] = comment
		}
		if desc, ok := v.Description(); ok {
			mapping["description"] = desc
		}
		if onBoot, ok := v.OnBoot(); ok {
			mapping["on_boot"] = onBoot
		}
		if mac, ok := v.Mac(); ok {
			mapping["mac_address"] = mac.MustAddress()
		}

		if reportedDevices, ok := v.ReportedDevices(); ok && len(reportedDevices.Slice()) > 0 {
			mapping["reported_devices"] = flattenOvirtNicReportedDevices(reportedDevices.Slice())
		}

		s = append(s, mapping)
	}
	d.SetId(resource.UniqueId())
	return d.Set("nics", s)
}

func flattenOvirtNicReportedDevices(reportedDevices []*ovirtsdk4.ReportedDevice) []map[string]interface{} {
	var s []map[string]interface{}

	for _, v := range reportedDevices {
		mapping := map[string]interface{}{
			"id":   v.MustId(),
			"name": v.MustName(),
			"type": string(v.MustType()),
		}

		if desc, ok := v.Description(); ok {
			mapping["description"] = desc
		}
		if comment, ok := v.Comment(); ok {
			mapping["comment"] = comment
		}

		if mac, ok := v.Mac(); ok {
			mapping["mac_address"] = mac.MustAddress()
		}

		if ips, ok := v.Ips(); ok && len(ips.Slice()) > 0 {
			mapping["ips"] = flattenOvirtNicIps(ips.Slice())
		}

		s = append(s, mapping)
	}

	return s
}

func flattenOvirtNicIps(ips []*ovirtsdk4.Ip) []map[string]interface{} {
	var s []map[string]interface{}

	for _, v := range ips {
		mapping := map[string]interface{}{}

		if address, ok := v.Address(); ok {
			mapping["address"] = address
		}
		if gateway, ok := v.Gateway(); ok {
			mapping["gateway"] = gateway
		}
		if netmask, ok := v.Netmask(); ok {
			mapping["netmask"] = netmask
		}
		if version, ok := v.Version(); ok {
			mapping["version"] = version
		}
		s = append(s, mapping)
	}

	return s
}
