// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isInstanceProfiles = "profiles"
)

func dataSourceIBMISInstanceProfiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceProfilesRead,

		Schema: map[string]*schema.Schema{

			isInstanceProfiles: {
				Type:        schema.TypeList,
				Description: "List of instance profile maps",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this virtual server instance profile belongs to.",
						},
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default OS architecture for an instance with this profile.",
						},
						"architecture_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for the OS architecture.",
						},

						"architecture_values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The supported OS architecture(s) for an instance with this profile.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"bandwidth": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"disks": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of the instance profile's disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"quantity": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The value for this profile field.",
												},
												"default": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The default value for this profile field.",
												},
												"max": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The maximum value for this profile field.",
												},
												"min": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The minimum value for this profile field.",
												},
												"step": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The increment step value for this profile field.",
												},
												"values": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The permitted values for this profile field.",
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
											},
										},
									},
									"size": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The value for this profile field.",
												},
												"default": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The default value for this profile field.",
												},
												"max": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The maximum value for this profile field.",
												},
												"min": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The minimum value for this profile field.",
												},
												"step": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The increment step value for this profile field.",
												},
												"values": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The permitted values for this profile field.",
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
											},
										},
									},
									"supported_interface_types": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"default": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"values": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The supported disk interfaces used for attaching the disk.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance profile.",
						},
						"memory": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"port_speed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
								},
							},
						},
						"vcpu_architecture": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The default VCPU architecture for an instance with this profile.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The VCPU architecture for an instance with this profile.",
									},
								},
							},
						},
						"vcpu_count": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
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

func dataSourceIBMISInstanceProfilesRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		err := classicInstanceProfilesList(d, meta)
		if err != nil {
			return err
		}
	} else {
		err := instanceProfilesList(d, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicInstanceProfilesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcclassicv1.InstanceProfile{}
	for {
		listInstanceProfilesOptions := &vpcclassicv1.ListInstanceProfilesOptions{}
		if start != "" {
			listInstanceProfilesOptions.Start = &start
		}
		availableProfiles, response, err := sess.ListInstanceProfiles(listInstanceProfilesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Instance Profiles %s\n%s", err, response)
		}
		start = GetNext(availableProfiles.Next)
		allrecs = append(allrecs, availableProfiles.Profiles...)
		if start == "" {
			break
		}
	}
	profilesInfo := make([]map[string]interface{}, 0)
	for _, profile := range allrecs {

		l := map[string]interface{}{
			"name":   *profile.Name,
			"family": *profile.Family,
		}
		profilesInfo = append(profilesInfo, l)
	}
	d.SetId(dataSourceIBMISInstanceProfilesID(d))
	d.Set(isInstanceProfiles, profilesInfo)
	return nil
}

func instanceProfilesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listInstanceProfilesOptions := &vpcv1.ListInstanceProfilesOptions{}
	availableProfiles, response, err := sess.ListInstanceProfiles(listInstanceProfilesOptions)
	if err != nil {
		return fmt.Errorf("Error Fetching Instance Profiles %s\n%s", err, response)
	}
	profilesInfo := make([]map[string]interface{}, 0)
	for _, profile := range availableProfiles.Profiles {

		l := map[string]interface{}{
			"name":   *profile.Name,
			"family": *profile.Family,
		}
		if profile.OsArchitecture != nil {
			if profile.OsArchitecture.Default != nil {
				l["architecture"] = *profile.OsArchitecture.Default
			}
			if profile.OsArchitecture.Type != nil {
				l["architecture_type"] = profile.OsArchitecture.Type
			}
			if profile.OsArchitecture.Values != nil {
				l["architecture_values"] = profile.OsArchitecture.Values
			}
		}
		if profile.Bandwidth != nil {
			bandwidthList := []map[string]interface{}{}
			bandwidthMap := dataSourceInstanceProfileBandwidthToMap(*profile.Bandwidth.(*vpcv1.InstanceProfileBandwidth))
			bandwidthList = append(bandwidthList, bandwidthMap)
			l["bandwidth"] = bandwidthList
		}
		if profile.Disks != nil {
			disksList := []map[string]interface{}{}
			for _, disksItem := range profile.Disks {
				disksList = append(disksList, dataSourceInstanceProfileDisksToMap(disksItem))
			}
			l["disks"] = disksList
		}
		if profile.Href != nil {
			l["href"] = profile.Href
		}
		if profile.Memory != nil {
			memoryList := []map[string]interface{}{}
			memoryMap := dataSourceInstanceProfileMemoryToMap(*profile.Memory.(*vpcv1.InstanceProfileMemory))
			memoryList = append(memoryList, memoryMap)
			l["memory"] = memoryList
		}
		if profile.PortSpeed != nil {
			portSpeedList := []map[string]interface{}{}
			portSpeedMap := dataSourceInstanceProfilePortSpeedToMap(*profile.PortSpeed.(*vpcv1.InstanceProfilePortSpeed))
			portSpeedList = append(portSpeedList, portSpeedMap)
			l["port_speed"] = portSpeedList
		}
		if profile.VcpuArchitecture != nil {
			vcpuArchitectureList := []map[string]interface{}{}
			vcpuArchitectureMap := dataSourceInstanceProfileVcpuArchitectureToMap(*profile.VcpuArchitecture)
			vcpuArchitectureList = append(vcpuArchitectureList, vcpuArchitectureMap)
			l["vcpu_architecture"] = vcpuArchitectureList
		}
		if profile.VcpuCount != nil {
			vcpuCountList := []map[string]interface{}{}
			vcpuCountMap := dataSourceInstanceProfileVcpuCountToMap(*profile.VcpuCount.(*vpcv1.InstanceProfileVcpu))
			vcpuCountList = append(vcpuCountList, vcpuCountMap)
			l["vcpu_count"] = vcpuCountList
		}
		if profile.Disks != nil {
			l[isInstanceDisks] = dataSourceInstanceProfileFlattenDisks(profile.Disks)
			if err != nil {
				return fmt.Errorf("Error setting disks %s", err)
			}
		}
		profilesInfo = append(profilesInfo, l)
	}
	d.SetId(dataSourceIBMISInstanceProfilesID(d))
	d.Set(isInstanceProfiles, profilesInfo)
	return nil
}

// dataSourceIBMISInstanceProfilesID returns a reasonable ID for a Instance Profile list.
func dataSourceIBMISInstanceProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
